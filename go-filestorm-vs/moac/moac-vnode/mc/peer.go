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
	"errors"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"time"

	"gopkg.in/fatih/set.v0"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	pb "github.com/filestorm/go-filestorm/moac/moac-lib/proto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/p2p"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/p2p/discover"
)

var (
	errClosed            = errors.New("Peer set is closed")
	errAlreadyRegistered = errors.New("Peer is already registered")
	errNotRegistered     = errors.New("Peer is not registered")
)

const (
	maxKnownMsgs     = 3276800 // Maximum directcall msgs requestid to keep in the known list (prevent DOS)
	maxKnownTxs      = 32768   // Maximum transactions hashes to keep in the known list (prevent DOS)
	maxKnownBlocks   = 1024    // Maximum block hashes to keep in the known list (prevent DOS)
	handshakeTimeout = 5 * time.Second

	maxQueuedTxs   = 512  // maxQueuedTxs is the maximum number of transaction lists to queue up
	maxQueuedProps = 8    // maxQueuedProps is the maximum number of block propagations to queue up
	maxQueuedAnns  = 8    // maxQueuedAnns is the maximum number of block announcements to queue up
	maxQueuedMsgs  = 1000 // maxQueuedProps is the maximum number of scs msgs to queue up
	maxQueuedRes   = 8    // maxQueuedAnns is the maximum number of scs  res to queue up
)

// PeerInfo represents a short summary of the MoacNode sub-protocol metadata known
// about a connected Peer.
type PeerInfo struct {
	Version    int      `json:"version"`    // MoacNode protocol version negotiated
	Difficulty *big.Int `json:"difficulty"` // Total difficulty of the Peer's blockchain
	Head       string   `json:"head"`       // SHA3 hash of the Peer's best owned block
}

// propEvent is a block propagation, waiting for its turn in the broadcast queue.
type propEvent struct {
	block *types.Block
	td    *big.Int
}

type Peer struct {
	id string

	*p2p.Peer
	rw          p2p.MsgReadWriter
	version     int         // Protocol version negotiated
	forkDrop    *time.Timer // Timed connection dropper if forks aren't validated in time
	head        common.Hash
	td          *big.Int
	lock        sync.RWMutex
	subnet      string
	knownTxs    *set.Set                // Set of transaction hashes known to be known by this Peer
	knownBlocks *set.Set                // Set of block hashes known to be known by this Peer
	knownMsgs   *set.Set                // Set of msg request ids known to be sent to the Peer
	queuedTxs   chan types.Transactions // Queue of transactions to broadcast to the peer
	queuedProps chan *propEvent         // Queue of blocks to broadcast to the peer
	queuedAnns  chan *types.Block       // Queue of blocks to announce to the peer
	queuedMsgs  chan *pb.ScsPushMsg     // Queue of scs msgs to broadcast to the peer
	queuedRes   chan *pb.ScsPushMsg     // Queue of scs Res to broadcast to the peer
	term        chan struct{}           // Termination channel to stop the broadcaster
}

func newPeer(version int, p *p2p.Peer, rw p2p.MsgReadWriter) *Peer {
	id := p.ID()

	return &Peer{
		Peer:        p,
		rw:          rw,
		version:     version,
		id:          fmt.Sprintf("%x", id[:8]),
		knownTxs:    set.New(),
		knownBlocks: set.New(),
		knownMsgs:   set.New(),
		queuedTxs:   make(chan types.Transactions, maxQueuedTxs),
		queuedProps: make(chan *propEvent, maxQueuedProps),
		queuedAnns:  make(chan *types.Block, maxQueuedAnns),
		queuedMsgs:  make(chan *pb.ScsPushMsg, maxQueuedMsgs),
		queuedRes:   make(chan *pb.ScsPushMsg, maxQueuedRes),
		term:        make(chan struct{}),
		subnet:      p.Subnet(),
	}
}

// broadcast is a write loop that multiplexes block propagations, announcements
// and transaction broadcasts into the remote peer. The goal is to have an async
// writer that does not lock up node internals.
func (p *Peer) broadcastloop() {
	for {
		select {
		case txs := <-p.queuedTxs:
			if err := p.SendTransactions(txs); err != nil {
				return
			}
		case prop := <-p.queuedProps:
			if err := p.SendNewBlock(prop.block, prop.td); err != nil {
				return
			}
		case block := <-p.queuedAnns:
			if err := p.SendNewBlockHashes([]common.Hash{block.Hash()}, []uint64{block.NumberU64()}); err != nil {
				return
			}
		case msg := <-p.queuedMsgs:
			if err := p.SendScsMsg(msg); err != nil {
				return
			}
		case res := <-p.queuedRes:
			if err := p.SendScsRes(res); err != nil {
				return
			}
		case <-p.term:
			return
		}
	}
}

func (p *Peer) IsMainnet() bool {
	return p.Peer.IsMainnet()
}

func (p *Peer) SetNodeType(id discover.NodeID, nodeType int) {
	// call p2p.peer to set node type
	p.Peer.SetNodeType(id, nodeType)
}

func (p *Peer) Id() string {
	return p.id
}

// Info gathers and returns a collection of metadata known about a Peer.
func (p *Peer) Info() *PeerInfo {
	hash, td := p.Head()

	return &PeerInfo{
		Version:    p.version,
		Difficulty: td,
		Head:       hash.Hex(),
	}
}

// Head retrieves a copy of the current head hash and total difficulty of the
// Peer.
func (p *Peer) Head() (hash common.Hash, td *big.Int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	copy(hash[:], p.head[:])
	return hash, new(big.Int).Set(p.td)
}

// SetHead updates the head hash and total difficulty of the Peer.
func (p *Peer) SetHead(hash common.Hash, td *big.Int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	copy(p.head[:], hash[:])
	p.td.Set(td)
}

// MarkBlock marks a block as known for the Peer, ensuring that the block will
// never be propagated to this particular Peer.
func (p *Peer) MarkBlock(hash common.Hash) {
	// If we reached the memory allowance, drop a previously known block hash
	for p.knownBlocks.Size() >= maxKnownBlocks {
		p.knownBlocks.Pop()
	}
	p.knownBlocks.Add(hash)
}

func (p *Peer) MarkMsg(id common.Hash) {
	// If we reached the memory allowance, drop a previously known transaction hash
	for p.knownMsgs.Size() >= maxKnownMsgs {
		p.knownMsgs.Pop()
	}
	p.knownMsgs.Add(id)
}

// MarkTransaction marks a transaction as known for the Peer, ensuring that it
// will never be propagated to this particular Peer.
func (p *Peer) MarkTransaction(hash common.Hash) {
	// If we reached the memory allowance, drop a previously known transaction hash
	for p.knownTxs.Size() >= maxKnownTxs {
		p.knownTxs.Pop()
	}
	p.knownTxs.Add(hash)
}

// SendTransactions sends transactions to the Peer and includes the hashes
// in its transaction hash set for future reference.
func (p *Peer) SendTransactions(txs types.Transactions) error {
	for _, tx := range txs {
		p.knownTxs.Add(tx.Hash())
	}

	log.Debugf("p2p send txmsg for %s", p.subnet)
	return p2p.Send(p.rw, TxMsg, txs)
}

// AsyncSendTransactions queues list of transactions propagation to a remote
// peer. If the peer's broadcast queue is full, the event is silently dropped.
func (p *Peer) AsyncSendTransactions(txs types.Transactions) {
	select {
	case p.queuedTxs <- txs:
		for _, tx := range txs {
			p.knownTxs.Add(tx.Hash())
		}
	default:
		p.Log().Debug("Dropping transaction propagation", "count", len(txs))
	}
}

// SendNewBlockHashes announces the availability of a number of blocks through
// a hash notification.
func (p *Peer) SendNewBlockHashes(hashes []common.Hash, numbers []uint64) error {
	for _, hash := range hashes {
		p.knownBlocks.Add(hash)
	}

	request := make(NewBlockHashesData, len(hashes))
	for i := 0; i < len(hashes); i++ {
		request[i].Hash = hashes[i]
		request[i].Number = numbers[i]
	}
	log.Debugf("p2p send newblockhashesmsg for %s", p.subnet)
	return p2p.Send(p.rw, NewBlockHashesMsg, request)
}

// AsyncSendNewBlockHash queues the availability of a block for propagation to a
// remote peer. If the peer's broadcast queue is full, the event is silently
// dropped.
func (p *Peer) AsyncSendNewBlockHash(block *types.Block) {
	select {
	case p.queuedAnns <- block:
		p.knownBlocks.Add(block.Hash())
	default:
		p.Log().Debug("Dropping block announcement", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// SendNewBlock propagates an entire block to a remote Peer.
func (p *Peer) SendNewBlock(block *types.Block, td *big.Int) error {
	p.knownBlocks.Add(block.Hash())
	log.Debugf("p2p send newblockmsg for %s", p.subnet)
	return p2p.Send(p.rw, NewBlockMsg, []interface{}{block, td})
}

// AsyncSendNewBlock queues an entire block for propagation to a remote peer. If
// the peer's broadcast queue is full, the event is silently dropped.
func (p *Peer) AsyncSendNewBlock(block *types.Block, td *big.Int) {
	select {
	case p.queuedProps <- &propEvent{block: block, td: td}:
		p.knownBlocks.Add(block.Hash())
	default:
		p.Log().Debug("Dropping block propagation", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// SendBlockHeaders sends a batch of block headers to the remote Peer.
func (p *Peer) SendBlockHeaders(headers []*types.Header) error {
	log.Debugf("p2p send blockheadersmsg for %s", p.subnet)
	return p2p.Send(p.rw, BlockHeadersMsg, headers)
}

// SendBlockBodies sends a batch of block contents to the remote Peer.
func (p *Peer) SendBlockBodies(bodies []*BlockBody) error {
	log.Debugf("p2p send blockbodiesmsg for %s", p.subnet)
	return p2p.Send(p.rw, BlockBodiesMsg, BlockBodiesData(bodies))
}

// SendBlockBodiesRLP sends a batch of block contents to the remote Peer from
// an already RLP encoded format.
func (p *Peer) SendBlockBodiesRLP(bodies []rlp.RawValue) error {
	log.Debugf("p2p send blockbodiesmsgRLP for %s", p.subnet)
	return p2p.Send(p.rw, BlockBodiesMsg, bodies)
}

// SendNodeDataRLP sends a batch of arbitrary internal data, corresponding to the
// hashes requested.
func (p *Peer) SendNodeData(data [][]byte) error {
	log.Debugf("p2p send nodedatamsg for %s", p.subnet)
	return p2p.Send(p.rw, NodeDataMsg, data)
}

// SendReceiptsRLP sends a batch of transaction receipts, corresponding to the
// ones requested from an already RLP encoded format.
func (p *Peer) SendReceiptsRLP(receipts []rlp.RawValue) error {
	log.Debugf("p2p send receiptsRLP for %s", p.subnet)
	return p2p.Send(p.rw, ReceiptsMsg, receipts)
}

func (p *Peer) SendScsRes(msg *pb.ScsPushMsg) error {
	id := GetScsPushMsgHash(msg)
	p.MarkMsg(id)
	size, _, _ := rlp.EncodeToReader(*msg)
	log.Debugf("p2p send scsres[%d] for %s", size, p.subnet)
	return p2p.Send(p.rw, ScsRes, *msg)
}

// AsyncSendScsRes queues an scs push msg for propagation to a remote peer. If
// the peer's broadcast queue is full, the event is silently dropped.
func (p *Peer) AsyncSendScsRes(msg *pb.ScsPushMsg) {
	select {
	case p.queuedRes <- msg:
		id := GetScsPushMsgHash(msg)
		p.MarkMsg(id)
	default:
		p.Log().Debug("Dropping push res", "requestId", string(msg.Requestid))
	}
}

func (p *Peer) SendScsMsg(msg *pb.ScsPushMsg) error {
	id := GetScsPushMsgHash(msg)
	p.MarkMsg(id)
	size, _, _ := rlp.EncodeToReader(*msg)
	log.Debugf("p2p send scsmsg[%d] to network: %s", size, p.subnet)
	return p2p.Send(p.rw, ScsMsg, *msg)
}

// AsyncSendScsMsg queues an scs push msg for propagation to a remote peer. If
// the peer's broadcast queue is full, the event is silently dropped.
func (p *Peer) AsyncSendScsMsg(msg *pb.ScsPushMsg) {
	select {
	case p.queuedMsgs <- msg:
		id := GetScsPushMsgHash(msg)
		p.MarkMsg(id)
	default:
		p.Log().Debug("Dropping push msg", "requestId", string(msg.Requestid))
	}
}

// RequestOneHeader is a wrapper around the header query functions to fetch a
// single header. It is used solely by the fetcher.
func (p *Peer) RequestOneHeader(hash common.Hash) error {
	p.Log().Debug("Fetching single header", "hash", hash)
	log.Debugf("p2p send getblockheadersmsg one header for %s", p.subnet)
	return p2p.Send(p.rw, GetBlockHeadersMsg, &GetBlockHeadersData{Origin: HashOrNumber{Hash: hash}, Amount: uint64(1), Skip: uint64(0), Reverse: false})
}

// RequestHeadersByHash fetches a batch of blocks' headers corresponding to the
// specified header query, based on the hash of an origin block.
func (p *Peer) RequestHeadersByHash(origin common.Hash, amount int, skip int, reverse bool) error {
	p.Log().Debug("Fetching batch of headers", "count", amount, "fromhash", origin, "skip", skip, "reverse", reverse)
	log.Debugf("p2p send getblockheadersmsg by hash for %s", p.subnet)
	return p2p.Send(p.rw, GetBlockHeadersMsg, &GetBlockHeadersData{Origin: HashOrNumber{Hash: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestHeadersByNumber fetches a batch of blocks' headers corresponding to the
// specified header query, based on the number of an origin block.
func (p *Peer) RequestHeadersByNumber(origin uint64, amount int, skip int, reverse bool) error {
	p.Log().Debug("Fetching batch of headers", "count", amount, "fromnum", origin, "skip", skip, "reverse", reverse)
	log.Debugf("p2p send getblockheadersmsg by number for %s", p.subnet)
	return p2p.Send(p.rw, GetBlockHeadersMsg, &GetBlockHeadersData{Origin: HashOrNumber{Number: origin}, Amount: uint64(amount), Skip: uint64(skip), Reverse: reverse})
}

// RequestBodies fetches a batch of blocks' bodies corresponding to the hashes
// specified.
func (p *Peer) RequestBodies(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of block bodies", "count", len(hashes))
	log.Debugf("p2p send getblockbodiesmsg for %s", p.subnet)
	return p2p.Send(p.rw, GetBlockBodiesMsg, hashes)
}

// RequestNodeData fetches a batch of arbitrary data from a node's known state
// data, corresponding to the specified hashes.
func (p *Peer) RequestNodeData(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of state data", "count", len(hashes))
	log.Debugf("p2p send getnodedatamsg for %s", p.subnet)
	return p2p.Send(p.rw, GetNodeDataMsg, hashes)
}

// RequestReceipts fetches a batch of transaction receipts from a remote node.
func (p *Peer) RequestReceipts(hashes []common.Hash) error {
	p.Log().Debug("Fetching batch of receipts", "count", len(hashes))
	log.Debugf("p2p send getreceiptsmsg for %s", p.subnet)
	return p2p.Send(p.rw, GetReceiptsMsg, hashes)
}

// Handshake executes the mc protocol handshake, negotiating version number,
// network IDs, difficulties, head and genesis blocks.
func (p *Peer) Handshake(network uint64, td *big.Int, head common.Hash, genesis common.Hash) error {
	// Send out own handshake in a new thread
	errc := make(chan error, 2)
	var status StatusData // safe to read after two values have been received from errc

	go func() {
		log.Debugf("p2p send statusmsg for %s", p.subnet)
		errc <- p2p.Send(p.rw, StatusMsg, &StatusData{
			ProtocolVersion: uint32(p.version),
			NetworkId:       network,
			TD:              td,
			CurrentBlock:    head,
			GenesisBlock:    genesis,
		})
	}()
	go func() {
		err := p.checkRemoteStatus(network, &status, genesis)
		id := p.Peer.ID()
		if err != nil {
			Err := err.Error()
			if strings.Contains(Err, "Genesis block mismatch") ||
				strings.Contains(Err, "Protocol version mismatch") ||
				strings.Contains(Err, "NetworkId mismatch") {
				p.SetNodeType(id, discover.AlienNode)
				log.Debugf(
					"Node with mismatch status [%v] = %v, will be blacklisted",
					Err,
					id.String()[:16],
				)
			}
		} else {
			// if peer passes check remote status, it's at least an uncle node
			p.SetNodeType(id, discover.UncleNode)
		}
		errc <- err
	}()
	timeout := time.NewTimer(handshakeTimeout)
	defer timeout.Stop()
	for i := 0; i < 2; i++ {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return p2p.DiscReadTimeout
		}
	}
	p.td, p.head = status.TD, status.CurrentBlock
	return nil
}

func (p *Peer) checkRemoteStatus(network uint64, status *StatusData, genesis common.Hash) (err error) {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Code != StatusMsg {
		return ErrResp(ErrNoStatusMsg, "first msg has code %x (!= %x)", msg.Code, StatusMsg)
	}
	if msg.Size > ProtocolMaxMsgSize {
		return ErrResp(ErrMsgTooLarge, "%v > %v", msg.Size, ProtocolMaxMsgSize)
	}
	// Decode the handshake and make sure everything matches
	if err := msg.Decode(&status); err != nil {
		return ErrResp(ErrDecode, "msg %v: %v", msg, err)
	}
	if status.GenesisBlock != genesis {
		return ErrResp(ErrGenesisBlockMismatch, "%x (!= %x)", status.GenesisBlock[:8], genesis[:8])
	}
	if status.NetworkId != network {
		return ErrResp(ErrNetworkIdMismatch, "%d (!= %d)", status.NetworkId, network)
	}
	if int(status.ProtocolVersion) != p.version {
		return ErrResp(ErrProtocolVersionMismatch, "%d (!= %d)", status.ProtocolVersion, p.version)
	}
	return nil
}

// String implements fmt.Stringer.
func (p *Peer) String() string {
	return fmt.Sprintf("Peer %s [%s]", p.id,
		fmt.Sprintf("mc/%2d", p.version),
	)
}

// send sharding message
func (p *Peer) SendShardingMsg() {
}

// close signals the broadcast goroutine to terminate.
func (p *Peer) close() {
	close(p.term)
}

// PeerSet represents the collection of active peers currently participating in
// the MoacNode sub-protocol.
type PeerSet struct {
	peers  map[string]*Peer
	lock   sync.RWMutex
	closed bool
}

// newPeerSet creates a new Peer set to track the active participants.
func newPeerSet() *PeerSet {
	return &PeerSet{
		peers: make(map[string]*Peer),
	}
}

// Register injects a new Peer into the working set, or returns an error if the
// Peer is already known.
func (ps *PeerSet) Register(p *Peer) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if ps.closed {
		return errClosed
	}
	if _, ok := ps.peers[p.id]; ok {
		return errAlreadyRegistered
	}
	ps.peers[p.Id()] = p
	go p.broadcastloop()
	return nil
}

// Unregister removes a remote Peer from the active set, disabling any further
// actions to/from that particular entity.
func (ps *PeerSet) Unregister(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	p, ok := ps.peers[id]
	if !ok {
		return errNotRegistered
	}
	delete(ps.peers, id)
	p.close()

	return nil
}

// Peer retrieves the registered Peer with the given id.
func (ps *PeerSet) Peer(id string) *Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return ps.peers[id]
}

// Len returns if the current number of peers in the set.
func (ps *PeerSet) Len() int {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	return len(ps.peers)
}

// PeersWithoutBlock retrieves a list of peers that do not have a given block in
// their set of known hashes.
func (ps *PeerSet) PeersWithoutBlock(hash common.Hash) []*Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownBlocks.Has(hash) {
			list = append(list, p)
		}
	}
	return list
}

// PeersWithoutTx retrieves a list of peers that do not have a given transaction
// in their set of known hashes.
func (ps *PeerSet) PeersWithoutTx(hash common.Hash) []*Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownTxs.Has(hash) {
			list = append(list, p)
		}
	}
	return list
}

func GetScsPushMsgHash(msg *pb.ScsPushMsg) common.Hash {
	return common.RlpHash([]interface{}{
		msg.Requestid, msg.Scsid, msg.Requestflag, msg.Msghash},
	)
}

func (ps *PeerSet) PeersWithoutSCSMsg(msg *pb.ScsPushMsg) []*Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	id := GetScsPushMsgHash(msg)
	list := make([]*Peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		if !p.knownMsgs.Has(id) {
			list = append(list, p)
		}
	}
	return list

}

// BestPeer retrieves the known Peer with the currently highest total difficulty.
func (ps *PeerSet) BestPeer() *Peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	var (
		bestPeer *Peer
		bestTd   *big.Int
	)
	for _, p := range ps.peers {
		if _, td := p.Head(); bestPeer == nil || td.Cmp(bestTd) > 0 {
			bestPeer, bestTd = p, td
		}
	}
	return bestPeer
}

// Close disconnects all peers.
// No new peers can be registered after Close has returned.
func (ps *PeerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}
