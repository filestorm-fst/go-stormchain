// Copyright 2014 The MOAC-core Authors
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

// Package mc implements the MoacNode protocol.
package mc

import (
	"errors"
	"fmt"
	"math/big"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/common/hexutil"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-lib/mcdb"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/consensus"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/consensus/ethash"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/bloombits"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/vm"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/event"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/internal/mcapi"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc/downloader"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc/filters"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc/gasprice"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/miner"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/node"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/p2p"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/rpc"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/vnode"
)

// Add a instance of MoacService handling SCS
var Instance *MoacService

type LesServer interface {
	Start(srvr *p2p.Server)
	Stop()
	Protocols() []p2p.Protocol
}

// MoacService implements the MoacService full node service.
// 2018/07/06 Added scs interface to scsHandler
type MoacService struct {
	config      *Config
	chainConfig *params.ChainConfig

	// Channel for shutting down the service
	shutdownChan  chan bool    // Channel for shutting down the moacnode
	stopDbUpgrade func() error // stop chain db sequential key upgrade

	// Handlers
	txPool          *core.TxPool
	blockchain      *core.BlockChain
	ProtocolManager *ProtocolManager
	lesServer       LesServer
	scsHandler      *vnode.VnodeServer

	// DB interfaces
	chainDb        mcdb.Database // Block chain database
	eventMux       *event.TypeMux
	engine         consensus.Engine
	accountManager *accounts.Manager
	bloomRequests  chan chan *bloombits.Retrieval // Channel receiving bloom data retrieval requests
	bloomIndexer   *core.ChainIndexer             // Bloom indexer operating during block imports
	ApiBackend     *MoacApiBackend
	miner          *miner.Miner
	gasPrice       *big.Int
	moacbase       common.Address
	networkId      uint64
	netRPCService  *mcapi.PublicNetAPI
	lock           sync.RWMutex // Protects the variadic fields (e.g. gas price and moacbase)
}

func (s *MoacService) AddLesServer(ls LesServer) {
	s.lesServer = ls
}

func GetInstance() *MoacService {
	return Instance
}

// New creates a new MoacService object (including the
// initialisation of the common MoacService object)
func New(ctx *node.ServiceContext, config *Config) (*MoacService, error) {
	if config.SyncMode == downloader.LightSync {
		return nil, errors.New("can't run mc.MoacService in light sync mode, use les.LightMoacService")
	}
	if !config.SyncMode.IsValid() {
		return nil, fmt.Errorf("invalid sync mode %d", config.SyncMode)
	}
	chainDb, err := CreateDB(ctx, config, "chaindata")
	if err != nil {
		return nil, err
	}

	syncDb, err := CreateSyncDB(ctx, config, "syncdata")
	if err != nil {
		return nil, err
	}

	stopDbUpgrade := upgradeDeduplicateData(chainDb)
	chainConfig, genesisHash, genesisErr := core.SetupGenesisBlock(chainDb, config.Genesis)

	if _, ok := genesisErr.(*params.ConfigCompatError); genesisErr != nil && !ok {
		return nil, genesisErr
	}
	log.Debugf("Initialized chain configuration %v", chainConfig)

	mcSrv := &MoacService{
		config:         config,
		chainDb:        chainDb,
		chainConfig:    chainConfig,
		eventMux:       ctx.EventMux,
		accountManager: ctx.AccountManager,
		engine:         CreateConsensusEngine(ctx, config, chainConfig, chainDb),
		shutdownChan:   make(chan bool),
		stopDbUpgrade:  stopDbUpgrade,
		networkId:      config.NetworkId,
		gasPrice:       config.GasPrice,
		moacbase:       config.Moacbase,
		bloomRequests:  make(chan chan *bloombits.Retrieval),
		bloomIndexer:   NewBloomIndexer(chainDb, params.BloomBitsBlocks),
	}

	log.Infof("Initializing MoacService protocol versions=%v network=%v", ProtocolVersions, config.NetworkId)

	if !config.SkipBcVersionCheck {
		bcVersion := core.GetBlockChainVersion(chainDb)
		if bcVersion != core.BlockChainVersion && bcVersion != 0 {
			return nil, fmt.Errorf("Blockchain DB version mismatch (%d / %d). Run MOAC upgradedb.\n", bcVersion, core.BlockChainVersion)
		}
		core.WriteBlockChainVersion(chainDb, core.BlockChainVersion)
	}

	vmConfig := vm.Config{EnablePreimageRecording: config.EnablePreimageRecording}

	mcSrv.blockchain, err = core.NewBlockChain(chainDb, mcSrv.chainConfig, mcSrv.engine, vmConfig)
	if err != nil {
		return nil, err
	}
	//add by frank
	mcSrv.scsHandler = vnode.NewScsService(chainDb, syncDb, mcSrv.blockchain, mcSrv, ctx)

	// Rewind the chain in case of an incompatible config upgrade.
	if compat, ok := genesisErr.(*params.ConfigCompatError); ok {
		log.Warn("Rewinding chain to upgrade configuration", "err", compat)
		mcSrv.blockchain.SetHead(compat.RewindTo)
		core.WriteChainConfig(chainDb, genesisHash, chainConfig)
	}
	mcSrv.bloomIndexer.Start(mcSrv.blockchain.CurrentHeader(), mcSrv.blockchain.SubscribeChainEvent)

	if config.TxPool.Journal != "" {
		config.TxPool.Journal = ctx.ResolvePath(config.TxPool.Journal)
	}

	mcSrv.txPool = core.NewTxPool(config.TxPool, mcSrv.chainConfig, mcSrv.blockchain)

	log.Debugf("create new protocol manager")
	if mcSrv.ProtocolManager, err = NewProtocolManager(
		mcSrv.chainConfig,
		config.SyncMode,
		config.NetworkId,
		mcSrv.eventMux,
		mcSrv.txPool,
		mcSrv.engine,
		mcSrv.blockchain,
		chainDb,
		mcSrv.scsHandler); err != nil {
		return nil, err
	}

	mcSrv.miner = miner.New(mcSrv, mcSrv.chainConfig, mcSrv.EventMux(), mcSrv.engine, mcSrv.ProtocolManager.NetworkRelay)
	mcSrv.miner.SetExtra(makeExtraData(config.ExtraData))

	mcSrv.ApiBackend = &MoacApiBackend{mcSrv, nil}
	gpoParams := config.GPO
	if gpoParams.Default == nil {
		gpoParams.Default = config.GasPrice
	}
	mcSrv.ApiBackend.gpo = gasprice.NewOracle(mcSrv.ApiBackend, gpoParams)

	//Set the instance with the current MOAC NODE
	Instance = mcSrv

	log.Infof("******Current Block %v, isNuwa %v", mcSrv.blockchain.CurrentBlock().Number(), chainConfig.IsNuwa(mcSrv.blockchain.CurrentBlock().Number()))

	return mcSrv, nil
}

/*
 * Update the extra data from
 * 	uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
 * to
 *  pangu 0.8.x
 */
func makeExtraData(extra []byte) []byte {
	if len(extra) == 0 {
		// create default extradata
		extra, _ = rlp.EncodeToBytes([]interface{}{
			"MOAC-", params.VersionNum, "-",
			runtime.Version(),
			runtime.GOOS,
		})
	}
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}
	return extra
}

// CreateDB creates the chain database.
func CreateDB(ctx *node.ServiceContext, config *Config, name string) (mcdb.Database, error) {
	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*mcdb.LDBDatabase); ok {
		db.Meter("mc/db/chaindata/")
	}
	return db, nil
}

// CreateSyncDB creates the sync database.
func CreateSyncDB(ctx *node.ServiceContext, config *Config, name string) (mcdb.Database, error) {
	dbdir := ctx.ResolvePath(name)
	if common.FileExist(dbdir) {
		start := time.Now()
		os.RemoveAll(dbdir)
		log.Info("Database successfully deleted", "elapsed", common.PrettyDuration(time.Since(start)))
	}

	db, err := ctx.OpenDatabase(name, config.DatabaseCache, config.DatabaseHandles)
	if err != nil {
		return nil, err
	}
	if db, ok := db.(*mcdb.LDBDatabase); ok {
		db.Meter("mc/db/syncdata/")
	}
	return db, nil
}

// CreateConsensusEngine creates the required type of consensus engine instance for an MoacService service
func CreateConsensusEngine(ctx *node.ServiceContext, config *Config, chainConfig *params.ChainConfig, db mcdb.Database) consensus.Engine {
	switch {
	case config.PowFake:
		log.Warn("Ethash used in fake mode")
		return ethash.NewFaker()
	case config.PowTest:
		log.Warn("Ethash used in test mode")
		return ethash.NewTester()
	case config.PowShared:
		log.Warn("Ethash used in shared mode")
		return ethash.NewShared()
	default:
		engine := ethash.New(ctx.ResolvePath(config.EthashCacheDir), config.EthashCachesInMem, config.EthashCachesOnDisk,
			config.EthashDatasetDir, config.EthashDatasetsInMem, config.EthashDatasetsOnDisk)
		engine.SetThreads(-1) // Disable CPU mining
		return engine
	}
}

// APIs returns the collection of RPC services the moaccore package offers.
// NOTE, some of these services probably need to be moved to somewhere else.
// 2018/07/06 added the scsHandler as scs provider
func (s *MoacService) APIs() []rpc.API {

	apis := mcapi.GetAPIs(s.ApiBackend)
	fmt.Printf("GetAPI list: %d ......\n", len(apis))
	// Append any APIs exposed explicitly by the consensus engine
	apis = append(apis, s.engine.APIs(s.BlockChain())...)
	log.Infof("GetAPI list after engine.APIS: %d ......\n", len(apis))
	// Append all the local APIs and return
	return append(apis, []rpc.API{
		{
			Namespace: "mc",
			Version:   "1.0",
			Service:   NewPublicMoacAPI(s),
			Public:    true,
		}, {
			Namespace: "mc",
			Version:   "1.0",
			Service:   NewPublicMinerAPI(s),
			Public:    true,
		}, {
			Namespace: "mc",
			Version:   "1.0",
			Service:   downloader.NewPublicDownloaderAPI(s.ProtocolManager.downloader, s.eventMux),
			Public:    true,
		}, {
			Namespace: "miner",
			Version:   "1.0",
			Service:   NewPrivateMinerAPI(s),
			Public:    false,
		}, {
			Namespace: "mc",
			Version:   "1.0",
			Service:   filters.NewPublicFilterAPI(s.ApiBackend, false),
			Public:    true,
		}, {
			Namespace: "admin",
			Version:   "1.0",
			Service:   NewPrivateAdminAPI(s),
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPublicDebugAPI(s),
			Public:    true,
		}, {
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(s.chainConfig, s),
		}, {
			Namespace: "net",
			Version:   "1.0",
			Service:   s.netRPCService,
			Public:    true,
		}, {
			Namespace: "vnode",
			Version:   "1.0",
			Service:   s.scsHandler,
			Public:    true,
		},
	}...)
}

func (s *MoacService) ResetWithGenesisBlock(gb *types.Block) {
	s.blockchain.ResetWithGenesisBlock(gb)
}

/*
 * Return the base58 encoded address instead of the HEX address.
 */
func (s *MoacService) Moacbase() (eb common.Address, err error) {
	log.Debug("[mc/backend.go->MoacService.Moacbase]")
	s.lock.RLock()
	moacbase := s.moacbase
	s.lock.RUnlock()

	if moacbase != (common.Address{}) {
		return moacbase, nil
	}
	if wallets := s.AccountManager().Wallets(); len(wallets) > 0 {
		if accounts := wallets[0].Accounts(); len(accounts) > 0 {
			return accounts[0].Address, nil
		}
	}
	return common.Address{}, fmt.Errorf("MOAC node must have a specified base address!")
}

// set in js console via admin interface or wrapper from cli flags
func (self *MoacService) SetMoacbase(moacbase common.Address) {
	self.lock.Lock()
	self.moacbase = moacbase
	log.Debugf("Set moacbase to %s\n", moacbase.Hex())
	self.lock.Unlock()

	self.miner.SetMoacbase(moacbase)
}

func (s *MoacService) StartMining(local bool) error {
	eb, err := s.Moacbase()
	if err != nil {
		log.Error("Cannot start mining without moacbase", "err", err)
		return fmt.Errorf("moacbase missing: %v", err)
	}
	if local {
		// If local (CPU) mining is started, we can disable the transaction rejection
		// mechanism introduced to speed sync times. CPU mining on mainnet is ludicrous
		// so noone will ever hit this path, whereas marking sync done on CPU mining
		// will ensure that private networks work in single miner mode too.
		atomic.StoreUint32(&s.ProtocolManager.acceptTxs, 1)
	}
	go s.miner.Start(eb)
	return nil
}

func (s *MoacService) StopMining()         { s.miner.Stop() }
func (s *MoacService) IsMining() bool      { return s.miner.Mining() }
func (s *MoacService) Miner() *miner.Miner { return s.miner }

func (s *MoacService) AccountManager() *accounts.Manager  { return s.accountManager }
func (s *MoacService) BlockChain() *core.BlockChain       { return s.blockchain }
func (s *MoacService) TxPool() *core.TxPool               { return s.txPool }
func (s *MoacService) EventMux() *event.TypeMux           { return s.eventMux }
func (s *MoacService) Engine() consensus.Engine           { return s.engine }
func (s *MoacService) ChainDb() mcdb.Database             { return s.chainDb }
func (s *MoacService) IsListening() bool                  { return true } // Always listening
func (s *MoacService) McVersion() int                     { return int(s.ProtocolManager.SubProtocols[0].Version) }
func (s *MoacService) NetVersion() uint64                 { return s.networkId }
func (s *MoacService) Downloader() *downloader.Downloader { return s.ProtocolManager.downloader }
func (s *MoacService) IsSubnetP2PEnabled(contractAddress common.Address, where string) bool {
	return s.ProtocolManager.isSubnetP2PEnabled(contractAddress, where)
}

// Protocols implements node.Service, returning all the currently configured
// network protocols to start.
func (s *MoacService) Protocols() []p2p.Protocol {
	if s.lesServer == nil {
		return s.ProtocolManager.SubProtocols
	}
	return append(s.ProtocolManager.SubProtocols, s.lesServer.Protocols()...)
}

// Start implements node.Service, starting all internal goroutines needed by the
// MoacService protocol implementation.
func (s *MoacService) Start(server *p2p.Server) error {
	// Start the bloom bits servicing goroutines
	s.startBloomHandlers()

	// Start the RPC service
	s.netRPCService = mcapi.NewPublicNetAPI(server, s.NetVersion())

	// Figure out a max peers count based on the server limits
	maxPeers := server.MaxPeers
	if s.config.LightServ > 0 {
		maxPeers -= s.config.LightPeers
		if maxPeers < server.MaxPeers/2 {
			maxPeers = server.MaxPeers / 2
		}
	}
	// Start the networking layer and the light server if requested
	s.ProtocolManager.Start(maxPeers)
	if s.lesServer != nil {
		s.lesServer.Start(server)
	}
	return nil
}

// Stop implements node.Service, terminating all internal goroutines used by the
// MoacService protocol.
func (s *MoacService) Stop() error {
	if s.stopDbUpgrade != nil {
		s.stopDbUpgrade()
	}
	s.bloomIndexer.Close()
	s.blockchain.Stop()
	s.ProtocolManager.Stop()
	if s.lesServer != nil {
		s.lesServer.Stop()
	}
	s.txPool.Stop()
	s.miner.Stop()
	s.eventMux.Stop()

	s.chainDb.Close()
	s.scsHandler.Close()
	log.Info("MOAC NODE shutting down....")
	close(s.shutdownChan)

	return nil
}
