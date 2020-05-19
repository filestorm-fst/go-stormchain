// Copyright 2015 The MOAC-core Authors
// This file is part of the MOAC-core library.
//
// The MOAC-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The MOAC-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the MOAC-core library. If not, see <http://www.gnu.org/licenses/>.

package mc

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-lib/mcdb"
	pb "github.com/filestorm/go-filestorm/moac/moac-lib/proto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/consensus"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/contracts"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	ctypes "github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/vm"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/event"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc/downloader"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc/fetcher"
	nr "github.com/filestorm/go-filestorm/moac/moac-vnode/networkrelay"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/node"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/p2p"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/p2p/discover"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/vnode"
	gocache "github.com/patrickmn/go-cache"
)

const (
	softResponseLimit = 2 * 1024 * 1024 // Target maximum size of returned blocks, headers or node data.
	estHeaderRlpSize  = 500             // Approximate size of an RLP encoded block header

	// txChanSize is the size of channel listening to TxPreEvent.
	// The number is referenced from the size of tx pool.
	txChanSize       = 4096
	scsMsgChanSize   = 1000
	txSyncChanSize   = 4096
	subnetsConfigKey = "subnets config"
	maxBlockNumber   = uint64(10000000000000000000)
	blockBuffer      = uint64(10)
)

var (
	daoChallengeTimeout       = 15 * time.Second      // Time allowance for a node to reply to the DAO handshake challenge
	subnetP2PServerStarting   = make(map[string]bool) // Used together with the mutex for singleton of p2p server start
	subnetP2PServerStartingMu sync.Mutex              // Used together with the bool mapping for singleton of p2p server start

	// hash of subchainbase var "isSubnetP2PEnabled"
	isSubnetP2PEnabledHex = "0x90348ed6"

	// Keep the value cached for 10 minute and clean up old ones every 5 mins
	// key is contractaddress + blocknumber, so it will be queried with different
	// keys after blockchain grows
	isSubnetP2PEnabledContractSettingCache = gocache.New(10*time.Minute, 5*time.Minute)

	removeZombieSubnetInterval      = 3 * time.Minute // check for zombie subnet every 30 mins
	zombieSubnetTTL                 = 2 * time.Hour   // remove subnets that do not have traffic for 2 hours
	subnetOpInAdvance               = uint64(100)     // create subnet p2p n mainnet blocks in advance
	checkCurrentBlockNumberInterval = 3 * time.Second // interval for check if new block arrives
)

// errIncompatibleConfig is returned if the requested protocols and configs are
// not compatible (low protocol version restrictions and high requirements).
var errIncompatibleConfig = errors.New("incompatible configuration")

func ErrResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

const (
	SubnetBegin  = 0
	SubnetEnd    = 1
	SubnetRemove = 2
)

type SubnetNotify struct {
	subnetid    string
	blockNumber uint64
	op          int // enum defined as SubnetBegin, SubnetEnd
}

type SubnetConfig struct {
	Subnetid   string
	BlockStart uint64
	BlockEnd   uint64
}

type ProtocolManager struct {
	networkId       uint64
	fastSync        uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	acceptTxs       uint32 // Flag whether we're considered synchronised (enables transaction processing)
	txpool          txPool
	blockchain      *core.BlockChain
	chaindb         mcdb.Database
	chainconfig     *params.ChainConfig
	subnetsConfigdb mcdb.Database
	subnetsConfig   map[string]SubnetConfig
	downloader      *downloader.Downloader
	NetworkRelay    *nr.NetworkRelay
	fetcher         *fetcher.Fetcher
	SubProtocols    []p2p.Protocol
	eventMux        *event.TypeMux
	txCh            chan core.TxPreEvent
	txSub           event.Subscription
	minedBlockSub   *event.TypeMuxSubscription

	// channels for fetcher, syncer, txsyncLoop
	newPeerCh   chan *Peer
	txsyncCh    chan *txsync
	quitSync    chan struct{}
	noMorePeers chan struct{}

	// wait group is used for graceful shutdowns during downloading
	// and processing
	wg                sync.WaitGroup
	sharding          core.ShardingMonitor // shardingmonitor
	scsMsgChan        chan *pb.ScsPushMsg
	scsResChan        chan *pb.ScsPushMsg
	scsMsgChanMainnet chan *pb.ScsPushMsg
	scsResChanMainnet chan *pb.ScsPushMsg
	p2pManager        *P2PManager
}

// about the host Peer.
type MoacNodeInfo struct {
	Network    uint64      `json:"network"`    // MoacNode network ID (1=Frontier, 2=Morden, Ropsten=3)
	Difficulty *big.Int    `json:"difficulty"` // Total difficulty of the host's blockchain
	Genesis    common.Hash `json:"genesis"`    // SHA3 hash of the host's genesis block
	Head       common.Hash `json:"head"`       // SHA3 hash of the host's best owned block
}

type P2PManager struct {
	maxPeers             int
	peerSet              *PeerSet                    // this is for mainnet p2p network
	subnetPeerSets       map[string]*PeerSet         // mapping from subnet name to peerset for the subnet p2p network
	subnetPushMsgBuffers map[string][]*pb.ScsPushMsg // mapping from subnet to its buffered msg
	subnetResMsgBuffers  map[string][]*pb.ScsPushMsg // mapping from subnet to its buffered msg
	subnetLastMsg        map[string]time.Time        //mapping from subnet to its last msg timestamp
}

func (p *P2PManager) String() string {
	availableSubnets := []string{}
	for subnetName, peerset := range p.subnetPeerSets {
		availableSubnets = append(
			availableSubnets,
			fmt.Sprintf("%s(peers=%d)", subnetName, len(peerset.peers)),
		)
	}

	return fmt.Sprintf("%s", availableSubnets)
}

// NewProtocolManager returns a new moaccore sub protocol manager. The MoacNode sub protocol manages peers capable
// with the moaccore network.
func NewProtocolManager(
	config *params.ChainConfig,
	mode downloader.SyncMode,
	networkId uint64,
	mux *event.TypeMux,
	txpool txPool,
	engine consensus.Engine,
	blockchain *core.BlockChain,
	chaindb mcdb.Database,
	scsHandler *vnode.VnodeServer,
) (*ProtocolManager, error) {
	// Create the protocol manager with the base fields
	manager := &ProtocolManager{
		networkId:         networkId,
		eventMux:          mux,
		txpool:            txpool,
		blockchain:        blockchain,
		chaindb:           chaindb,
		chainconfig:       config,
		subnetsConfigdb:   NewSubnetsConfigDB(chaindb),
		subnetsConfig:     make(map[string]SubnetConfig),
		newPeerCh:         make(chan *Peer),
		noMorePeers:       make(chan struct{}),
		txsyncCh:          make(chan *txsync, txSyncChanSize),
		quitSync:          make(chan struct{}),
		scsMsgChan:        make(chan *pb.ScsPushMsg, scsMsgChanSize),
		scsResChan:        make(chan *pb.ScsPushMsg, scsMsgChanSize),
		scsMsgChanMainnet: make(chan *pb.ScsPushMsg, scsMsgChanSize),
		scsResChanMainnet: make(chan *pb.ScsPushMsg, scsMsgChanSize),
		p2pManager: &P2PManager{
			subnetPeerSets:       make(map[string]*PeerSet),
			peerSet:              newPeerSet(),
			subnetPushMsgBuffers: make(map[string][]*pb.ScsPushMsg),
			subnetResMsgBuffers:  make(map[string][]*pb.ScsPushMsg),
			subnetLastMsg:        make(map[string]time.Time),
		},
	}
	// Figure out whether to allow fast sync or not
	if mode == downloader.FastSync && blockchain.CurrentBlock().NumberU64() > 0 {
		log.Warn("Blockchain not empty, fast sync disabled")
		mode = downloader.FullSync
	}
	if mode == downloader.FastSync {
		manager.fastSync = uint32(1)
	}
	// Initiate a sub-protocol for every implemented version we can handle
	manager.SubProtocols = make([]p2p.Protocol, 0, len(ProtocolVersions))
	for i, version := range ProtocolVersions {
		// Skip protocol version if incompatible with the mode of operation
		if mode == downloader.FastSync && version < mc63 {
			continue
		}

		// Compatible; initialise the sub-protocol
		version := version // Closure for the run
		manager.SubProtocols = append(manager.SubProtocols, p2p.Protocol{
			Name:    ProtocolName,
			Version: version,
			Length:  ProtocolLengths[i],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := manager.newPeer(int(version), p, rw)
				log.Debugf(
					"subnet run func for peer %v, local: %v, remote: %v, subnet name=%v, protocol=%v, version=%v",
					peer.id,
					peer.Peer.LocalAddr(),
					peer.Peer.RemoteAddr(),
					string(peer.Peer.Subnet()),
					ProtocolName,
					version,
				)
				select {
				case manager.newPeerCh <- peer:
					manager.wg.Add(1)
					defer manager.wg.Done()
					return manager.runPeer(peer)
				case <-manager.quitSync:
					return p2p.DiscQuitting
				}
			},
			NodeInfo: func() interface{} {
				return manager.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p := manager.p2pManager.peerSet.Peer(fmt.Sprintf("%x", id[:8])); p != nil {
					return p.Info()
				}
				return nil
			},
		})
	}

	if len(manager.SubProtocols) == 0 {
		return nil, errIncompatibleConfig
	}

	// Construct the different synchronisation mechanisms
	manager.downloader = downloader.New(mode, chaindb, manager.eventMux, blockchain, nil, manager.removePeer)

	validator := func(header *types.Header) error {
		return engine.VerifyHeader(blockchain, header, true)
	}
	heighter := func() uint64 {
		return blockchain.CurrentBlock().NumberU64()
	}
	inserter := func(blocks types.Blocks) (int, error) {
		// If fast sync is running, deny importing weird blocks
		if atomic.LoadUint32(&manager.fastSync) == 1 {
			log.Warn("Discarded bad propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		atomic.StoreUint32(&manager.acceptTxs, 1) // Mark initial sync done on any fetcher import
		liveFlag := false
		return manager.blockchain.InsertChain(blocks, liveFlag)
	}
	manager.fetcher = fetcher.New(blockchain.GetBlockByHash, validator, manager.BroadcastBlock, heighter, inserter, manager.removePeer)

	manager.NetworkRelay = nr.NewNetworkRelay(
		manager.scsMsgChan,
		manager.scsResChan,
		manager.scsMsgChanMainnet,
		manager.scsResChanMainnet,
		scsHandler,
		blockchain.ChainId().Uint64(),
	)

	//Add scs Servers to NetworkRelay
	core.Nr = manager.NetworkRelay

	return manager, nil
}

func (pm *ProtocolManager) removePeer(id string, subnet string) {
	// Short circuit if the Peer was already removed
	var peerSet *PeerSet
	var ok bool
	if subnet == "mainnet" {
		peerSet = pm.p2pManager.peerSet
	} else {
		peerSet, ok = pm.p2pManager.subnetPeerSets[subnet]
		// if can not find the subnet peerset
		if !ok || peerSet == nil {
			log.Errorf("subnet can not find peerset in removepeer %s", subnet)
			return
		}
	}
	peer := peerSet.Peer(id)
	if peer == nil {
		return
	}
	log.Debug("Removing MoacNode Peer", "Peer", id)

	// Unregister the Peer from the downloader and MoacNode Peer set
	pm.downloader.UnregisterPeer(id)
	if err := peerSet.Unregister(id); err != nil {
		log.Error("Peer removal failed", "Peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

func (pm *ProtocolManager) Start(maxPeers int) {
	pm.p2pManager.maxPeers = maxPeers

	// broadcast transactions
	pm.txCh = make(chan core.TxPreEvent, txChanSize)
	pm.txSub = pm.txpool.SubscribeTxPreEvent(pm.txCh)
	go pm.txBroadcastLoop()

	// broadcast mined blocks
	pm.minedBlockSub = pm.eventMux.Subscribe(core.NewMinedBlockEvent{})
	go pm.minedBroadcastLoop()

	// start sync handlers
	go pm.syncer()
	go pm.txsyncLoop()

	// start sharding monitor loop
	go pm.shardingBroadcastLoop()

	// start buffered subnet msg broadcast loop
	go pm.subnetBufferedMsgBroadcastLoop()

	// start subnet notify loop
	go pm.subnetCheckLoop()

	// initialize all persisted subnets from db [Disabled]
	if subnetsConfig, err := pm.loadSubnetsConfig(); err == nil {
		pm.subnetsConfig = subnetsConfig
	}

	// start zombie subnet removal loop
	go pm.removeZombieSubnetsLoop()
}

func (pm *ProtocolManager) Stop() {
	log.Info("Stopping MoacNode protocol")

	pm.txSub.Unsubscribe()         // quits txBroadcastLoop
	pm.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop

	// Quit the sync loop.
	// After this send has completed, no new peers will be accepted.
	pm.noMorePeers <- struct{}{}

	// Quit fetcher, txsyncLoop.
	close(pm.quitSync)

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the Peer set.
	// sessions which are already established but not added to pm.peers yet
	// will exit when they try to register.
	peerSet := pm.p2pManager.peerSet
	log.Debugf("subnet close mainnet peerset")
	peerSet.Close()

	for subnetid, peerSet := range pm.p2pManager.subnetPeerSets {
		log.Debugf("subnet close subnet:%s peerset", subnetid)
		peerSet.Close()
	}

	// Wait for all Peer handler goroutines and the loops to come down.
	pm.wg.Wait()

	log.Info("MoacNode protocol stopped")
}

func (pm *ProtocolManager) newPeer(protocolVersion int, p *p2p.Peer, rw p2p.MsgReadWriter) *Peer {
	if !p.IsMainnet() {
		subnet := p.Subnet()
		// if peerset does not exist for the given peer
		if _, ok := pm.p2pManager.subnetPeerSets[subnet]; !ok {
			pm.p2pManager.subnetPeerSets[subnet] = newPeerSet()
			log.Debugf("subnet create new subnet peerset for %s", subnet)
		}
	}

	return newPeer(protocolVersion, p, newMeteredMsgWriter(rw))
}

// runpeer is the callback invoked to manage the life cycle of an mc Peer. When
// this function terminates, the Peer is disconnected.
func (pm *ProtocolManager) runPeer(p *Peer) error {
	log.Debugf("subnet runpeer %s", p.subnet)
	var peerSet *PeerSet
	if p.IsMainnet() {
		peerSet = pm.p2pManager.peerSet
	} else {
		peerSet = pm.p2pManager.subnetPeerSets[p.subnet]
		if peerSet == nil {
			return fmt.Errorf("subnet can not found peerset in runpeer")
		}
	}

	if peerSet.Len() >= pm.p2pManager.maxPeers {
		log.Debugf("[MC handle] err = %s, [%d/%d]", p2p.DiscTooManyPeers, peerSet.Len(), pm.p2pManager.maxPeers)
		return p2p.DiscTooManyPeers
	}
	p.Log().Debug("MoacNode Peer connected", "name", p.Name())

	// Execute the MoacNode handshake
	td, head, genesis := pm.blockchain.Status()
	if err := p.Handshake(pm.networkId, td, head, genesis); err != nil {
		p.Log().Debug("MoacNode handshake failed", "err", err)
		return err
	}

	// init read and writer, sets the protocol version
	if rw, ok := p.rw.(*meteredMsgReadWriter); ok {
		rw.Init(p.version)
	}

	// Register the Peer to the peerset, whatever it is the mainnet one or subnet one
	if err := peerSet.Register(p); err != nil {
		p.Log().Error("MoacNode Peer registration failed", "err", err)
		return err
	}
	defer pm.removePeer(p.id, p.subnet)

	// only do this for mainnet p2p peers
	if p.IsMainnet() {
		// Register the Peer in the downloader. If the downloader considers it banned, we disconnect
		if err := pm.downloader.RegisterPeer(p.id, p.version, p); err != nil {
			return err
		}

		// Propagate existing transactions. new transactions appearing
		// after this will be sent via broadcasts.
		pm.syncTransactions(p)
	}

	for {
		if err := pm.handleMsg(p); err != nil {
			p.Log().Debug("MoacNode message handling failed", "err", err)
			return err
		}
	}
}

func (pm *ProtocolManager) subnetCheckLoop() {
	// check every 3 seconds
	checkCurrentBlockNumberTicker := time.NewTicker(checkCurrentBlockNumberInterval)

	for {
		select {
		case <-checkCurrentBlockNumberTicker.C:
			currentBlockNumber := pm.blockchain.CurrentHeader().Number.Uint64()
			subnetsToRemove := []string{}
			for subnetid, subnetConfig := range pm.subnetsConfig {

				// case a: create the subnet p2p server
				createAtBlockNumber := subnetConfig.BlockStart - subnetOpInAdvance
				if currentBlockNumber >= createAtBlockNumber {
					log.Debugf("subnet get new start subnet p2p server notify for %s", subnetid)
					n := node.GetInstance()
					if !n.HasSubnetServer(subnetid) {
						pm.maybeStartSubnetP2PServer(subnetid)
					}
				}

				// case b: remove the subnet p2p server, note that subnet forward
				// should already be disabled by the logic in issubnetp2penabled()
				endAtBlockNumber := subnetConfig.BlockEnd
				if currentBlockNumber > endAtBlockNumber {
					subnetsToRemove = append(subnetsToRemove, subnetid)
				}
			}

			if len(subnetsToRemove) > 0 {
				for _, subnetid := range subnetsToRemove {
					// should updating the db is successful, remove the server as well
					if err := pm.updateSubnetsConfigDB(subnetid, 0, SubnetRemove); err == nil {
						n := node.GetInstance()
						pm.RemoveSubnetServer(subnetid, n)
					}
				}
			}
		}
	}
}

func (pm *ProtocolManager) subnetBufferedMsgBroadcastLoop() {
	// run every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	for _ = range ticker.C {
		// broadcast buffered push msg
		for subnet, msgBuffer := range pm.p2pManager.subnetPushMsgBuffers {
			peerSet, _ := pm.p2pManager.subnetPeerSets[subnet]
			// if peerset is available
			if peerSet != nil && len(peerSet.peers) > 0 {
				for _, msg := range msgBuffer {
					pm.BroadcastPushMsgToPeerSet(msg, peerSet)
					log.Debugf("subnet broadcast buffered msg for %s", subnet)
				}
				// set the array back to empty
				pm.p2pManager.subnetPushMsgBuffers[subnet] = nil
			}
		}

		// broadcast buffered res msg
		for subnet, msgBuffer := range pm.p2pManager.subnetResMsgBuffers {
			peerSet, _ := pm.p2pManager.subnetPeerSets[subnet]
			// if peerset is available
			if peerSet != nil && len(peerSet.peers) > 0 {
				for _, msg := range msgBuffer {
					pm.BroadcastPushResToPeerSet(msg, peerSet)
					log.Debugf("subnet broadcast buffered msg for %s", subnet)
				}
				// set the array back to empty
				pm.p2pManager.subnetResMsgBuffers[subnet] = nil
			}
		}
	}
}

// updateSubnetsConfigDB update the subnet config db
func (pm *ProtocolManager) updateSubnetsConfigDB(subnetid string, blockNumber uint64, op int) error {
	// validate inputs
	currentBlockNumber := pm.blockchain.CurrentHeader().Number.Uint64()
	// for all ops except remove, check blockNumber
	if op != SubnetRemove && blockNumber < currentBlockNumber+blockBuffer {
		return fmt.Errorf("blockNumber needs to be at least %d blocks in the future", blockBuffer)
	}

	// load existing configs
	key := []byte(subnetsConfigKey)
	value, _ := pm.subnetsConfigdb.Get(key)
	log.Debugf("subnet config current config \"%s\"", string(value))
	subnetsConfig := make(map[string]SubnetConfig)
	if len(value) > 0 {
		if err := json.Unmarshal(value, &subnetsConfig); err != nil {
			log.Errorf("subnet config unmarshal failed: %v", err)
			return err
		}
	}

	// update config
	var subnetConfig SubnetConfig
	var found bool
	if subnetConfig, found = subnetsConfig[subnetid]; !found {
		subnetConfig = SubnetConfig{subnetid, maxBlockNumber, maxBlockNumber}
	}

	switch op {
	case SubnetBegin:
		if blockNumber >= subnetConfig.BlockEnd {
			return fmt.Errorf(
				"block start(%d) can not be larger than block end(%d)",
				subnetConfig.BlockStart, subnetConfig.BlockEnd,
			)
		}
		subnetConfig.BlockStart = blockNumber
		subnetsConfig[subnetid] = subnetConfig
	case SubnetEnd:
		if blockNumber <= subnetConfig.BlockStart {
			return fmt.Errorf(
				"block end(%d) can not be smaller than block start(%d)",
				subnetConfig.BlockStart, subnetConfig.BlockEnd,
			)
		}
		subnetConfig.BlockEnd = blockNumber
		subnetsConfig[subnetid] = subnetConfig
	case SubnetRemove:
		delete(subnetsConfig, subnetid)
	}

	// write it back
	if scBytes, err := json.Marshal(subnetsConfig); err != nil {
		log.Debugf(
			"subnet p2p config for %s updated failed with error: %v",
			subnetid, err,
		)
		return err
	} else {
		// persist new config back
		if errPut := pm.subnetsConfigdb.Put(key, scBytes); errPut != nil {
			return errPut
		}
		log.Debugf(
			"subnet p2p config for %s updated to %d/%d",
			subnetid,
			subnetConfig.BlockStart,
			subnetConfig.BlockEnd,
		)
	}
	pm.subnetsConfig = subnetsConfig
	return nil
}

func (pm *ProtocolManager) loadSubnetsConfig() (map[string]SubnetConfig, error) {
	key := []byte(subnetsConfigKey)
	currentConfig, _ := pm.subnetsConfigdb.Get(key)
	sc := make(map[string]SubnetConfig)
	if err := json.Unmarshal(currentConfig, &sc); err != nil {
		return sc, err
	}

	return sc, nil
}

func (pm *ProtocolManager) maybeStartSubnetP2PServer(subnetid string) {
	// use mutex to make sure we only have one instance
	// of goroutine running for each subnet
	log.Debugf("subnet in maybeStartSubnetP2PServer %s", subnetid)
	subnetP2PServerStartingMu.Lock()
	start := !subnetP2PServerStarting[subnetid]
	subnetP2PServerStarting[subnetid] = true
	subnetP2PServerStartingMu.Unlock()
	log.Debugf("subnet in maybeStartSubnetP2PServer after lock, start = %t, %s", start, subnetid)
	if start {
		go func(_subnetid string) {
			n := node.GetInstance()
			if !n.HasSubnetServer(_subnetid) {
				n.StartSubnetP2PServer(_subnetid)
			}
			subnetP2PServerStartingMu.Lock()
			subnetP2PServerStarting[_subnetid] = false
			subnetP2PServerStartingMu.Unlock()
		}(subnetid)
	}
}

func (pm *ProtocolManager) getContractParam(contractAddress common.Address, data []byte) []byte {
	nr := pm.NetworkRelay
	st, _ := pm.blockchain.State()
	if codeHash := st.GetCodeHash(contractAddress); codeHash == (common.Hash{}) {
		return []byte{}
	}
	to := common.Address{}
	to.SetString("")
	from := common.Address{}
	viaaddress := common.Address{}
	viaaddress.SetString(params.VnodeBeneficialAddress)
	msgHash := common.Hash{}
	msgHash.SetString("")
	msg := ctypes.NewMessage(
		from,
		&to,
		0,
		big.NewInt(0),
		big.NewInt(0),
		big.NewInt(0),
		[]byte{},
		false,
		false,
		false,
		big.NewInt(0),
		0,
		&viaaddress,
		&msgHash,
	)
	context := core.NewEVMContext(msg, pm.blockchain.CurrentBlock().Header(), pm.blockchain, nil, nil)
	evm := vm.NewEVM(context, st, params.TestChainConfig,
		vm.Config{EnableJit: false, ForceJit: false}, nr)
	contractRef := vm.AccountRef(contractAddress)
	precompiledContracts := contracts.GetInstance()
	ret, leftGas, err := evm.Call(
		contractRef,
		contractAddress,
		data,
		params.GenesisGasLimit.Uint64(),
		big.NewInt(0),
		false,
		uint64(0),
		precompiledContracts,
		msg.GetMsgHash(),
	)
	log.Debugf(" %v, %v, %v", ret, leftGas, err)

	return ret
}

func (pm *ProtocolManager) isSubnetP2PEnabled(contractAddress common.Address, where string) bool {
	// check vnode local setting first
	subnetid := strings.ToLower(contractAddress.String())
	currentBlockNumber := pm.blockchain.CurrentHeader().Number.Uint64()
	subnetConfig, found := pm.subnetsConfig[subnetid]
	subnetOnBlockNumber := subnetConfig.BlockStart
	subnetOffBlockNumber := subnetConfig.BlockEnd
	// check if current block number is in the active range of subnet server
	if found && subnetOnBlockNumber <= currentBlockNumber && currentBlockNumber < subnetOffBlockNumber {
		log.Debugf("subnet isSubnetP2PEnabled %s: %t in [%s] %d / %d / %d", subnetid, true, where, subnetOnBlockNumber, currentBlockNumber, subnetOffBlockNumber)
		return true
	}

	// check subchainbase setting
	key := fmt.Sprintf("isSubnetP2PEnabled:%s:%d", contractAddress.String(), currentBlockNumber)
	var ret = false
	if value, found := isSubnetP2PEnabledContractSettingCache.Get(key); !found {
		data := common.FromHex(isSubnetP2PEnabledHex)
		_ret := pm.getContractParam(contractAddress, data)
		_ret_int := 0
		if len(_ret) > 8 {
			_ret_int = int(binary.BigEndian.Uint32(_ret[(len(_ret) - 4):]))
		}
		ret = _ret_int != 0
		isSubnetP2PEnabledContractSettingCache.Set(key, ret, gocache.DefaultExpiration)
	} else {
		ret = value.(bool)
	}

	log.Debugf("subnet isSubnetP2PEnabled %s, %t in [%s] %d / %d / %d", key, ret, where, subnetOnBlockNumber, currentBlockNumber, subnetOffBlockNumber)

	return ret
}

func (pm *ProtocolManager) removeZombieSubnetsLoop() {
	ticker := time.NewTicker(removeZombieSubnetInterval)
	for _ = range ticker.C {
		n := node.GetInstance()
		for subnetid, lastMsgTime := range pm.p2pManager.subnetLastMsg {
			deadline := lastMsgTime.Add(zombieSubnetTTL)
			// if we passed deadline
			if deadline.Before(time.Now()) {
				log.Debugf(
					"subnet zombie %s should be removed: lastMsg=%v, current=%v",
					subnetid, lastMsgTime, time.Now(),
				)

				pm.RemoveSubnetServer(subnetid, n)
			}
		}
	}
}

func (pm *ProtocolManager) RemoveSubnetServer(subnetid string, n *node.Node) {
	// 1. stop p2p discovery
	n.StopSubnetServer(subnetid)

	// 2. disconnect all active peer connections
	if peerset, found := pm.p2pManager.subnetPeerSets[subnetid]; found {
		peerset.Close()
		delete(pm.p2pManager.subnetPeerSets, subnetid)
	}
}

func (pm *ProtocolManager) shardingBroadcastLoop() {
	// 1) if the msg is received from mainnet, broadcast back to the mainnet
	// 2) if the msg is received from subnet, broadcast back to the subnet
	for {
		select {
		case msg := <-pm.scsMsgChan:
			// replace the first two bytes with "0x"
			subnetid := fmt.Sprintf("0x%s", common.Bytes2Hex(msg.Subchainid)[2:])

			// get subnet peerset, if peerset is not ready,
			// notify to create the p2p server for the subnet and buffer the message
			peerSet, _ := pm.p2pManager.subnetPeerSets[subnetid]
			if peerSet != nil {
				// update subnets last msg timestamp
				pm.p2pManager.subnetLastMsg[subnetid] = time.Now()

				pm.BroadcastPushMsgToPeerSet(msg, peerSet)
			} else {
				// buffer the subnet pusuh msg
				if pm.p2pManager.subnetPushMsgBuffers[subnetid] == nil {
					pm.p2pManager.subnetPushMsgBuffers[subnetid] = make([]*pb.ScsPushMsg, 0)
				}
				pm.p2pManager.subnetPushMsgBuffers[subnetid] = append(
					pm.p2pManager.subnetPushMsgBuffers[subnetid],
					msg,
				)
				log.Errorf("subnet pushmsg peerset not found for subchain %s, msg buffered", subnetid)
			}

		case msg := <-pm.scsResChan:
			// replace the first two bytes with "0x"
			subnetid := fmt.Sprintf("0x%s", common.Bytes2Hex(msg.Subchainid)[2:])
			log.Debugf("subnet resmsg avaiable subnets: %s, msg subnet: %s", pm.p2pManager, subnetid)
			peerSet, _ := pm.p2pManager.subnetPeerSets[subnetid]
			if peerSet != nil {
				pm.BroadcastPushResToPeerSet(msg, peerSet)
			} else {
				// buffer the subnet res msg
				if pm.p2pManager.subnetResMsgBuffers[subnetid] == nil {
					pm.p2pManager.subnetResMsgBuffers[subnetid] = make([]*pb.ScsPushMsg, 0)
				}
				pm.p2pManager.subnetResMsgBuffers[subnetid] = append(
					pm.p2pManager.subnetResMsgBuffers[subnetid],
					msg,
				)
				log.Errorf("subnet pushres peerset not found for subchain %s", subnetid)
			}

		case msg := <-pm.scsMsgChanMainnet:
			// broadcast to mainnet
			pm.BroadcastPushMsgToPeerSet(msg, pm.p2pManager.peerSet)

		case msg := <-pm.scsResChanMainnet:
			// broadcast to mainnet
			pm.BroadcastPushResToPeerSet(msg, pm.p2pManager.peerSet)
		}
	}
}

// handleMsg is invoked whenever an inbound message is received from a remote
// Peer. The remote connection is torn down upon returning any error.
func (pm *ProtocolManager) handleMsg(p *Peer) error {
	log.Debugf("[mc/handler.go->pm.handleMsg]")
	// Read the next message from the remote Peer, and ensure it's fully consumed
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Size > ProtocolMaxMsgSize {
		return ErrResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
	}
	defer msg.Discard()

	if p.IsMainnet() {
		err = pm.handleMainnetMsg(p, msg)
	} else {
		err = pm.handleSubnetMsg(p, msg)
	}

	return err
}

func (pm *ProtocolManager) handleSubnetMsg(p *Peer, msg p2p.Msg) error {
	switch {
	case msg.Code == ScsMsg:
		var scsMsg pb.ScsPushMsg
		if err := msg.Decode(&scsMsg); err != nil {
			return ErrResp(ErrDecode, "msg %v: %v", msg, err)
		}
		id := GetScsPushMsgHash(&scsMsg)
		p.MarkMsg(id)
		// broadcastmsg will
		// 1) forward msg to other vnode (end up in shardingBroadcastLoop())
		// 2) forward msg to local scs
		forceToMainnet := false
		subnetid := fmt.Sprintf("0x%s", common.Bytes2Hex(scsMsg.Subchainid)[2:])
		contractAddress := common.HexToAddress(subnetid)
		if !pm.isSubnetP2PEnabled(contractAddress, "handleSubnetMsg") {
			forceToMainnet = true
		}
		pm.NetworkRelay.BroadcastMsg(&scsMsg, forceToMainnet)

	case msg.Code == ShardingCreateMsg:

	case msg.Code == ShardingFlushMsg:

	case msg.Code == ScsRes:

	case msg.Code == ScsReg:

	}

	return nil
}

func (pm *ProtocolManager) handleMainnetMsg(p *Peer, msg p2p.Msg) error {
	// Handle the message depending on its contents
	switch {
	case msg.Code == StatusMsg:
		// Status messages should never arrive after the handshake
		return ErrResp(ErrExtraStatusMsg, "uncontrolled status message")

	// Block header query, collect the requested headers and reply
	case msg.Code == GetBlockHeadersMsg:
		// Decode the complex header query
		var query GetBlockHeadersData
		if err := msg.Decode(&query); err != nil {
			return ErrResp(ErrDecode, "%v: %v", msg, err)
		}
		hashMode := query.Origin.Hash != (common.Hash{})

		// Gather headers until the fetch or network limits is reached
		var (
			bytes   common.StorageSize
			headers []*types.Header
			unknown bool
		)

		for !unknown && len(headers) < int(query.Amount) && bytes < softResponseLimit && len(headers) < downloader.MaxHeaderFetch {
			// Retrieve the next header satisfying the query
			var origin *types.Header
			if hashMode {
				origin = pm.blockchain.GetHeaderByHash(query.Origin.Hash)
			} else {
				origin = pm.blockchain.GetHeaderByNumber(query.Origin.Number)
			}
			if origin == nil {
				break
			}
			number := origin.Number.Uint64()
			headers = append(headers, origin)
			bytes += estHeaderRlpSize

			// Advance to the next header of the query
			switch {
			case query.Origin.Hash != (common.Hash{}) && query.Reverse:
				// Hash based traversal towards the genesis block
				for i := 0; i < int(query.Skip)+1; i++ {
					if header := pm.blockchain.GetHeader(query.Origin.Hash, number); header != nil {
						query.Origin.Hash = header.ParentHash
						number--
					} else {
						unknown = true
						break
					}
				}
			case query.Origin.Hash != (common.Hash{}) && !query.Reverse:
				// Hash based traversal towards the leaf block
				var (
					current = origin.Number.Uint64()
					next    = current + query.Skip + 1
				)
				if next <= current {
					infos, _ := json.MarshalIndent(p.Peer.Info(), "", "  ")
					p.Log().Warn("GetBlockHeaders skip overflow attack", "current", current, "skip", query.Skip, "next", next, "attacker", infos)
					unknown = true
				} else {
					if header := pm.blockchain.GetHeaderByNumber(next); header != nil {
						if pm.blockchain.GetBlockHashesFromHash(header.Hash(), query.Skip+1)[query.Skip] == query.Origin.Hash {
							query.Origin.Hash = header.Hash()
						} else {
							unknown = true
						}
					} else {
						unknown = true
					}
				}
			case query.Reverse:
				// Number based traversal towards the genesis block
				if query.Origin.Number >= query.Skip+1 {
					query.Origin.Number -= (query.Skip + 1)
				} else {
					unknown = true
				}

			case !query.Reverse:
				// Number based traversal towards the leaf block
				query.Origin.Number += (query.Skip + 1)
			}
		}
		return p.SendBlockHeaders(headers)

	case msg.Code == BlockHeadersMsg:
		// A batch of headers arrived to one of our previous requests
		var headers []*types.Header
		if err := msg.Decode(&headers); err != nil {
			return ErrResp(ErrDecode, "msg %v: %v", msg, err)
		}

		filter := len(headers) == 1
		if filter {
			headers = pm.fetcher.FilterHeaders(headers, time.Now())
		}
		if len(headers) > 0 || !filter {
			err := pm.downloader.DeliverHeaders(p.id, headers)
			if err != nil {
				log.Debug("Failed to deliver headers", "err", err)
			}
		}

	case msg.Code == GetBlockBodiesMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather blocks until the fetch or network limits is reached
		var (
			hash   common.Hash
			bytes  int
			bodies []rlp.RawValue
		)
		for bytes < softResponseLimit && len(bodies) < downloader.MaxBlockFetch {
			// Retrieve the hash of the next block
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return ErrResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested block body, stopping if enough was found
			if data := pm.blockchain.GetBodyRLP(hash); len(data) != 0 {
				bodies = append(bodies, data)
				bytes += len(data)
			}
		}
		return p.SendBlockBodiesRLP(bodies)

	case msg.Code == BlockBodiesMsg:
		// A batch of block bodies arrived to one of our previous requests
		var request BlockBodiesData
		if err := msg.Decode(&request); err != nil {
			return ErrResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver them all to the downloader for queuing
		trasactions := make([][]*types.Transaction, len(request))
		uncles := make([][]*types.Header, len(request))

		for i, body := range request {
			trasactions[i] = body.Transactions
			uncles[i] = body.Uncles
		}
		// Filter out any explicitly requested bodies, deliver the rest to the downloader
		filter := len(trasactions) > 0 || len(uncles) > 0
		if filter {
			trasactions, uncles = pm.fetcher.FilterBodies(trasactions, uncles, time.Now())
		}
		if len(trasactions) > 0 || len(uncles) > 0 || !filter {
			err := pm.downloader.DeliverBodies(p.id, trasactions, uncles)
			if err != nil {
				log.Debug("Failed to deliver bodies", "err", err)
			}
		}

	case p.version >= mc63 && msg.Code == GetNodeDataMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather state data until the fetch or network limits is reached
		var (
			hash  common.Hash
			bytes int
			data  [][]byte
		)
		for bytes < softResponseLimit && len(data) < downloader.MaxStateFetch {
			// Retrieve the hash of the next state entry
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return ErrResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested state entry, stopping if enough was found
			if entry, err := pm.chaindb.Get(hash.Bytes()); err == nil {
				data = append(data, entry)
				bytes += len(entry)
			}
		}
		return p.SendNodeData(data)

	case p.version >= mc63 && msg.Code == NodeDataMsg:
		// A batch of node state data arrived to one of our previous requests
		var data [][]byte
		if err := msg.Decode(&data); err != nil {
			return ErrResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverNodeData(p.id, data); err != nil {
			log.Debug("Failed to deliver node state data", "err", err)
		}

	case p.version >= mc63 && msg.Code == GetReceiptsMsg:
		// Decode the retrieval message
		msgStream := rlp.NewStream(msg.Payload, uint64(msg.Size))
		if _, err := msgStream.List(); err != nil {
			return err
		}
		// Gather state data until the fetch or network limits is reached
		var (
			hash     common.Hash
			bytes    int
			receipts []rlp.RawValue
		)
		for bytes < softResponseLimit && len(receipts) < downloader.MaxReceiptFetch {
			// Retrieve the hash of the next block
			if err := msgStream.Decode(&hash); err == rlp.EOL {
				break
			} else if err != nil {
				return ErrResp(ErrDecode, "msg %v: %v", msg, err)
			}
			// Retrieve the requested block's receipts, skipping if unknown to us
			results := core.GetBlockReceipts(pm.chaindb, hash, core.GetBlockNumber(pm.chaindb, hash))
			if results == nil {
				if header := pm.blockchain.GetHeaderByHash(hash); header == nil || header.ReceiptHash != types.EmptyRootHash {
					continue
				}
			}
			// If known, encode and queue for response packet
			if encoded, err := rlp.EncodeToBytes(results); err != nil {
				log.Error("Failed to encode receipt", "err", err)
			} else {
				receipts = append(receipts, encoded)
				bytes += len(encoded)
			}
		}
		return p.SendReceiptsRLP(receipts)

	case p.version >= mc63 && msg.Code == ReceiptsMsg:
		// A batch of receipts arrived to one of our previous requests
		var receipts [][]*types.Receipt
		if err := msg.Decode(&receipts); err != nil {
			return ErrResp(ErrDecode, "msg %v: %v", msg, err)
		}
		// Deliver all to the downloader
		if err := pm.downloader.DeliverReceipts(p.id, receipts); err != nil {
			log.Debug("Failed to deliver receipts", "err", err)
		}

	case msg.Code == NewBlockHashesMsg:
		var announces NewBlockHashesData
		if err := msg.Decode(&announces); err != nil {
			return ErrResp(ErrDecode, "%v: %v", msg, err)
		}
		// Mark the hashes as present at the remote node
		for _, block := range announces {
			pm.NetworkRelay.SetBlockNumber(block.Number)
			p.MarkBlock(block.Hash)
		}
		// Schedule all the unknown hashes for retrieval
		unknown := make(NewBlockHashesData, 0, len(announces))
		for _, block := range announces {
			if !pm.blockchain.HasBlock(block.Hash, block.Number) {
				unknown = append(unknown, block)
			}
		}
		for _, block := range unknown {
			pm.fetcher.Notify(p.id, block.Hash, block.Number, time.Now(), p.RequestOneHeader, p.RequestBodies)
		}

	case msg.Code == NewBlockMsg:
		// Retrieve and decode the propagated block
		var request newBlockData
		if err := msg.Decode(&request); err != nil {
			return ErrResp(ErrDecode, "%v: %v", msg, err)
		}
		request.Block.ReceivedAt = msg.ReceivedAt
		request.Block.ReceivedFrom = p

		// Mark the Peer as owning the block and schedule it for import
		p.MarkBlock(request.Block.Hash())
		pm.fetcher.Enqueue(p.id, request.Block)

		// Assuming the block is importable by the Peer, but possibly not yet done so,
		// calculate the head hash and TD that the Peer truly must have.
		var (
			trueHead = request.Block.ParentHash()
			trueTD   = new(big.Int).Sub(request.TD, request.Block.Difficulty())
		)
		// Update the peers total difficulty if better than the previous
		if _, td := p.Head(); trueTD.Cmp(td) > 0 {
			p.SetHead(trueHead, trueTD)

			// Schedule a sync if above ours. Note, this will not fire a sync for a gap of
			// a singe block (as the true TD is below the propagated block), however this
			// scenario should easily be covered by the fetcher.
			currentBlock := pm.blockchain.CurrentBlock()
			if trueTD.Cmp(pm.blockchain.GetTd(currentBlock.Hash(), currentBlock.NumberU64())) > 0 {
				go pm.synchronise(p)
			}
		}

	case msg.Code == TxMsg:
		var txs, txsShard []*types.Transaction
		// Transactions can be processed, parse all of them and deliver to the pool
		if err := msg.Decode(&txs); err != nil {
			return ErrResp(ErrDecode, "msg %v: %v", msg, err)
		}

		// Transactions arrived, make sure we have a valid and fresh chain to handle them
		if atomic.LoadUint32(&pm.acceptTxs) == 0 {
			log.Debug("no chain available")
			for _, tx := range txs {
				if tx.ShardingFlag() != 0 {
					txsShard = append(txsShard, tx)
				}
			}
			if len(txsShard) == 0 {
				break
			}
			txs = txsShard
		}

		for i, tx := range txs {
			// Validate and mark the remote transaction
			if tx == nil {
				return ErrResp(ErrDecode, "transaction %d is nil", i)
			}
			p.MarkTransaction(tx.Hash())
		}
		pm.txpool.AddRemotes(txs)

	case msg.Code == ScsMsg:
		var scsMsg pb.ScsPushMsg
		// msg.payload is a ScsPushMsg
		// this code decode the payload and put its content into scsMsg
		if err := msg.Decode(&scsMsg); err != nil {
			return ErrResp(ErrDecode, "msg %v: %v", msg, err)
		}
		id := GetScsPushMsgHash(&scsMsg)
		p.MarkMsg(id)
		// broadcastmsg will
		// 1) forward msg to other vnode (end up in shardingBroadcastLoop())
		// 2) forward msg to local scs
		pm.NetworkRelay.BroadcastMsg(&scsMsg, true)

	case msg.Code == ShardingCreateMsg:

	case msg.Code == ShardingFlushMsg:

	case msg.Code == ScsRes:

	case msg.Code == ScsReg:

	default:
		return ErrResp(ErrInvalidMsgCode, "%v", msg.Code)
	}

	return nil
}

// BroadcastBlock will either propagate a block to a subset of it's peers, or
// will only announce it's availability (depending what's requested).
func (pm *ProtocolManager) BroadcastBlock(block *types.Block, propagate bool) {
	hash := block.Hash()
	peerSet := pm.p2pManager.peerSet
	peers := peerSet.PeersWithoutBlock(hash)

	// If propagation is requested, send to a subset of the Peer
	if propagate {
		// Calculate the TD of the block (it's not imported yet, so block.Td is not valid)
		var td *big.Int
		if parent := pm.blockchain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent != nil {
			td = new(big.Int).Add(block.Difficulty(), pm.blockchain.GetTd(block.ParentHash(), block.NumberU64()-1))
		} else {
			log.Error("Propagating dangling block", "number", block.Number(), "hash", hash)
			return
		}
		// Send the block to a subset of our peers
		transfer := peers[:int(math.Sqrt(float64(len(peers))))]
		for _, peer := range transfer {
			peer.AsyncSendNewBlock(block, td)
		}
		log.Trace("Propagated block", "hash", hash, "recipients", len(transfer), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
		return
	}
	// Otherwise if the block is indeed in out own chain, announce it
	if pm.blockchain.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			peer.AsyncSendNewBlockHash(block)
		}
		log.Trace("Announced block", "hash", hash, "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
	}
}

// BroadcastTx will propagate a transaction to all peers which are not known to
// already have the given transaction.
func (pm *ProtocolManager) BroadcastTx(hash common.Hash, tx *types.Transaction) {
	// Broadcast transaction to a batch of peers not knowing about it
	peerSet := pm.p2pManager.peerSet
	peers := peerSet.PeersWithoutTx(hash)
	//FIXME include this again: peers = peers[:int(math.Sqrt(float64(len(peers))))]
	for _, peer := range peers {
		peer.AsyncSendTransactions(types.Transactions{tx})
	}
}

func (pm *ProtocolManager) BroadcastPushResToPeerSet(msg *pb.ScsPushMsg, peerSet *PeerSet) {
	peers := peerSet.PeersWithoutSCSMsg(msg)
	for _, peer := range peers {
		peer.AsyncSendScsRes(msg)
	}
}

func (pm *ProtocolManager) BroadcastPushMsgToPeerSet(msg *pb.ScsPushMsg, peerSet *PeerSet) {
	peers := peerSet.PeersWithoutSCSMsg(msg)
	for _, peer := range peers {
		peer.AsyncSendScsMsg(msg)
	}
	log.Debugf("!!!  peer SendScsMsg Requestid: %s, peerNum: %d", string(msg.Requestid), len(peers))
}

// Mined broadcast loop
func (pm *ProtocolManager) minedBroadcastLoop() {
	// automatically stops if unsubscribe
	for obj := range pm.minedBlockSub.Chan() {
		switch ev := obj.Data.(type) {
		case core.NewMinedBlockEvent:
			pm.BroadcastBlock(ev.Block, true)  // First propagate block to peers
			pm.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}
}

func (pm *ProtocolManager) txBroadcastLoop() {
	for {
		select {
		case event := <-pm.txCh:
			pm.BroadcastTx(event.Tx.Hash(), event.Tx)
		case <-pm.txSub.Err():
			return
		}
	}
}

// NodeInfo retrieves some protocol metadata about the running host node.
func (pm *ProtocolManager) NodeInfo() *MoacNodeInfo {
	currentBlock := pm.blockchain.CurrentBlock()
	return &MoacNodeInfo{
		Network:    pm.networkId,
		Difficulty: pm.blockchain.GetTd(currentBlock.Hash(), currentBlock.NumberU64()),
		Genesis:    pm.blockchain.Genesis().Hash(),
		Head:       currentBlock.Hash(),
	}
}
