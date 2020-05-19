// Copyright 2017 The MOAC-core Authors
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

package networkrelay

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"gopkg.in/fatih/set.v0"

	"bytes"
	"math/big"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	pb "github.com/filestorm/go-filestorm/moac/moac-lib/proto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
	libtypes "github.com/filestorm/go-filestorm/moac/moac-lib/types"
	"golang.org/x/net/context"

	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/contracts"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/node"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
)

var nr *NetworkRelay
var mu sync.Mutex
var scsMsgLimit = uint(params.ScsMsgLimit)
var subChainMsgLimit = uint(params.SubchainMsgLimit)

const maxKnownMsgs = 32768 // Maximum Msg to keep in the known list
const maxRoleCache = 20

type RoleCache struct {
	//SubchainId common.Address
	Role   params.ScsKind
	Number uint64
}

type NotifyScsMsg struct {
	Address   *common.Address
	Msg       *[]byte
	MsgHash   *common.Hash
	Block     *big.Int
	TimeStamp time.Time
}

type ContractRef interface {
	Address() common.Address
}

type AccountRef common.Address

// Address casts AccountRef to a Address
func (ar AccountRef) Address() common.Address { return (common.Address)(ar) }

type VmInterface interface {
	StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64) (ret []byte, leftOverGas uint64, err error)
}

type VnodeServerInterface interface {
	//ConnectScsServers(scsid string, blknum *big.Int) error
	ScsPush(stream pb.Vnode_ScsPushServer) error
	AccountInfo(ctx context.Context, in *pb.AccountInfoRequest) (*pb.AccountInfoReply, error)
	ChainInfo(ctx context.Context, in *pb.ChainInfoRequest) (*pb.ChainInfoReply, error)
	RemoteCall(ctx context.Context, in *pb.RemoteCallRequest) (*pb.RemoteCallReply, error)
	GetBlockNumber() uint64
	GetContractAddressMappingMember(addr common.Address, pos uint32, key common.Address) ([]common.Address, error)
	NotifyMsgRunState(hash common.Hash) bool
	GetSCSRole(contractAddress common.Address, nodeAddress common.Address) params.ScsKind
	GetScsPushMsgChan() chan pb.ScsPushMsg
	GetScsRegChan() chan pb.ScsPushMsg
	GetScsPushResChan() chan pb.ScsPushMsg
}

type NetworkRelay struct {
	ChainId           uint64
	lock              sync.RWMutex
	ScsServers        map[string]*ScsServerConfig
	scsMsgChan        chan *pb.ScsPushMsg
	scsResChan        chan *pb.ScsPushMsg
	scsMsgChanMainnet chan *pb.ScsPushMsg
	scsResChanMainnet chan *pb.ScsPushMsg
	scsHandler        VnodeServerInterface
	vm                VmInterface
	_vnodePushChan    chan pb.ScsPushMsg
	currReqId         uint32
	NotifyScsCache    map[string]*NotifyScsMsg
	ScsServerList     *map[string]*libtypes.ScsServerConnection
	knownMsgs         *set.Set // Set of msg request ids known to be sent to the scs
	scsMsgCounts      map[string]uint
	req               chan *pb.ScsPushMsg
	stream            pb.Vnode_ScsPushServer
	blockNumber       uint64
	// For role cache
	lockRole     sync.RWMutex
	ScsRoleCache map[common.Address]map[common.Address]*RoleCache
	lockWhite    sync.RWMutex
	WhiteFlag    uint64
	WhiteList    []common.Hash
	WhiteFunc    []common.Hash
}

type ScsServerConfig struct {
	Stream    pb.Vnode_ScsPushServer
	PublicKey string
	ScsId     string
}

func NewNetworkRelay(
	scsMsgChan chan *pb.ScsPushMsg,
	scsResChan chan *pb.ScsPushMsg,
	scsMsgChanMainnet chan *pb.ScsPushMsg,
	scsResChanMainnet chan *pb.ScsPushMsg,
	scsHandler VnodeServerInterface,
	chainId uint64,
) *NetworkRelay {
	log.Debugf("[mc/network_relay] NewNetworkRelay")
	nr = &NetworkRelay{
		scsMsgChan:        scsMsgChan,
		scsResChan:        scsResChan,
		scsMsgChanMainnet: scsMsgChanMainnet,
		scsResChanMainnet: scsResChanMainnet,
		scsHandler:        scsHandler,
		ChainId:           chainId,
	}

	go nr.runScsServersChannels()
	nr.ScsServers = make(map[string]*ScsServerConfig)
	nr._vnodePushChan = make(chan pb.ScsPushMsg) //TODO
	nr.NotifyScsCache = make(map[string]*NotifyScsMsg)
	nr.ScsRoleCache = make(map[common.Address]map[common.Address]*RoleCache)
	nr.currReqId = 0
	nr.knownMsgs = set.New()
	nr.blockNumber = 0
	nr.WhiteFlag = 0

	nr.resetMsgCounts()

	return nr
}

func GetInstance() *NetworkRelay {
	return nr
}

func (self *NetworkRelay) SetScsServerList(scsServerList *map[string]*libtypes.ScsServerConnection) {
	self.lock.RLock()
	defer self.lock.RUnlock()

	self.ScsServerList = scsServerList
}

func (self *NetworkRelay) GetNewRequestId() uint32 {
	if self != nil {
		self.currReqId++
		return self.currReqId
	}
	log.Debugf("self is nil", self)
	return 0
}

func (self *NetworkRelay) UpdateWhiteState(bn uint64) {
	if !params.PriorityChain(self.ChainId) {
		return
	}
	log.Debugf("[mc/network_relay] UpdateWhiteState WhiteFlag:%d, bn:%d", self.WhiteFlag, bn)
	if bn == 0 {
		self.WhiteFlag = bn
		return
	}
	if self.WhiteFlag == 0 {
		self.WhiteFlag = bn
	} else if self.WhiteFlag+2 < bn {
		return
	}

	self.updateWhiteInfo()
}

func (self *NetworkRelay) inWhiteFunc(contractAddr common.Address, code []byte) bool {
	self.lockWhite.RLock()
	defer self.lockWhite.RUnlock()

	if self.WhiteList != nil {
		codeHash := contracts.GetWhiteListCommonHash(contractAddr)
		if codeHash != (common.Hash{}) {
			for _, h := range self.WhiteList {
				if h == codeHash {
					return true
				}
			}
		}
	}

	if self.WhiteFunc != nil {
		hash := common.HexToHash(contractAddr.String() + common.Bytes2Hex(code) + "0000000000000000")
		for _, h := range self.WhiteFunc {
			if h == hash {
				return true
			}
		}
	}

	return false
}

func (self *NetworkRelay) inWhiteList(contractAddr common.Address) bool {
	if self.WhiteList == nil {
		return false
	}

	self.lockWhite.RLock()
	defer self.lockWhite.RUnlock()

	codeHash := contracts.GetWhiteListCommonHash(contractAddr)
	if codeHash == (common.Hash{}) {
		return false
	}
	hash := common.HexToHash(contractAddr.String() + "000000000000000000000000")
	for _, h := range self.WhiteList {
		if h == codeHash || h == hash {
			return true
		}
	}

	return false
}

func (self *NetworkRelay) updateWhiteInfo() {
	self.lockWhite.Lock()
	defer self.lockWhite.Unlock()

	self.WhiteList, self.WhiteFunc = contracts.GetWhiteInfo()
}

func (self *NetworkRelay) Priority(contractAddr, scsId common.Address, code []byte) bool {
	if params.ClearanceChain(self.ChainId) {
		return true
	}
	if !params.PriorityChain(self.ChainId) {
		return false
	}

	if self.inWhiteList(contractAddr) && self.inWhiteFunc(contractAddr, code) {
		role := self.GetSCSRole(contractAddr, scsId)
		return role != params.None && role != params.LockScs

	}

	return false
}

func (self *NetworkRelay) SetBlockNumber(bn uint64) {
	if self.blockNumber != bn {
		self.resetMsgCounts()
	}
	self.blockNumber = bn
}

func (self *NetworkRelay) NotifyScsFinalize(curblock *big.Int, liveFlag bool) {
	self.resetMsgCounts()
	if !liveFlag {
		self.NotifyScsCache = make(map[string]*NotifyScsMsg)
		return
	}
	var (
		notifyScsMsg []*NotifyScsMsg
		realMsg      []byte
		scsType      []byte
		scsStatus    []byte
	)

	//iterate through maps
	self.lock.Lock()
	for k, v := range self.NotifyScsCache {
		log.Debugf("[networkrelay/network_relay.go->NotifyScsFinalize] k:%v", k)
		log.Debugf("[networkrelay/network_relay.go->NotifyScsFinalize] v.Block:%v", v.Block)
		log.Debugf("[networkrelay/network_relay.go->NotifyScsFinalize] v.Msg:%v", common.Bytes2Hex(*v.Msg))
		if v.Block.Cmp(curblock) < 0 {
			log.Debugf("[networkrelay/network_relay.go->NotifyScsFinalize] need to notify v.Msg:%v", v.Msg)
			notifyScsMsg = append(notifyScsMsg, v)
			delete(self.NotifyScsCache, k)
		}
	}
	self.lock.Unlock()

	for _, notify := range notifyScsMsg {
		msg := *notify.Msg
		msgHash := *notify.MsgHash
		address := *notify.Address
		lenMsg := len(msg)
		if lenMsg > 31 {
			realMsg = msg[(lenMsg - 32):]
		}

		//convert to bigInt
		//TODO: change to bigInt
		msgInBigInt := big.NewInt(0)
		msgInBigInt.SetBytes(realMsg)
		msgInHex := int(msgInBigInt.Int64())
		scsStatus = common.IntToBytes(msgInHex)
		log.Debugf("NetworkRelay NotifyScsFinalize msgInHex %v", msgInHex)

		//TODO 1:proposal created, 2: proposal disputed.
		scsType = common.IntToBytes(params.ControlMsg)
		switch msgInHex {
		case params.RegOpen, params.RegClose, params.RegAdd, params.RegAsMonitor, params.RegAsBackup, params.ApproveProposal:
			log.Debugf("NetworkRelay NotifyScsFinalize clear cache of subchain[%v] %v", address.String(), msgInHex)
			self.SetScsRoleCacheOld(address)
		default:
			// reqType = ""
		}

		requestflag := true
		scsid := []byte("")
		nodeObj := node.GetInstance()
		var requestId []byte
		if nodeObj != nil {
			requestId = nodeObj.GetRequestId(false)
			id := common.RlpHash([]interface{}{requestId, scsid, requestflag, msgHash.Bytes()})
			if self.knownMsgs != nil && self.knownMsgs.Has(id) {
				requestId = nodeObj.GetRequestId(true)
				log.Debugf("NetworkRelay NotifyScsFinalize update requestId %s", string(requestId))
			}
		} else {
			requestId = []byte("")
		}
		//TODO: get a unified requestId
		conReq := &pb.ScsPushMsg{
			Requestid:   requestId,
			Timestamp:   common.Int64ToBytes(time.Now().Unix()),
			Requestflag: requestflag,
			Type:        scsType,
			Status:      scsStatus,
			Scsid:       scsid,
			Sender:      []byte(""), //don't need sender for notify
			Receiver:    address.Bytes(),
			Subchainid:  address.Bytes(),
			Msghash:     msgHash.Bytes(),
		}
		self.VnodePushMsg(conReq)
	}
}

func (self *NetworkRelay) NotifyScs(address common.Address, msg []byte, hash common.Hash, blocknumber *big.Int) {
	log.Debugf("[mc/network_relay] NotifyScs address %v msg %v hash %v", address, msg, hash)

	notifyScsMsg := NotifyScsMsg{
		Address:   &address,
		Msg:       &msg,
		MsgHash:   &hash,
		Block:     blocknumber,
		TimeStamp: time.Now(), //add a timestamp so that too old message could be cleaned up at some point.
	}
	self.lock.Lock()
	defer self.lock.Unlock()
	self.NotifyScsCache[hash.String()] = &notifyScsMsg
}

func (self *NetworkRelay) runScsServersChannels() error {
	log.Debugf("[mc/network_relay] runScsServersChannels")
	scsPushMsgChan := self.scsHandler.GetScsPushMsgChan()
	scsRegChan := self.scsHandler.GetScsRegChan()
	for {
		select {
		case scsPushMsg := <-scsPushMsgChan:
			self.OnReceiveMsg(&scsPushMsg)
		case scsPushMsg := <-scsRegChan:
			self.OnReceiveRegisterMsg(&scsPushMsg)
		}
	}

	return nil
}

func (self *NetworkRelay) BroadcastRes(msg *pb.ScsPushMsg, forceToMainnet bool) {
	log.Debugf("[mc/network_relay] BroadcastRes")
	if forceToMainnet {
		self.scsResChanMainnet <- msg
	} else {
		self.scsResChan <- msg
	}
}

// BroadcastMsg will
// 1) forward the msg to other vnode
// 2) forward the msg to local scs node
func (self *NetworkRelay) BroadcastMsg(scsmsg *pb.ScsPushMsg, forceToMainnet bool) {
	log.Debugf("[mc/network_relay] BroadcastMsg Requestid: %s, Subchainid: %s, forceToMainnet: %t", string(scsmsg.Requestid), common.ToHex(scsmsg.Subchainid[:]), forceToMainnet)
	if !self.ShouldBroadcast(scsmsg) {
		log.Debugf("[mc/network_relay] Too many msgs type: %v, status: %v", common.BytesToInt(scsmsg.Type), common.BytesToInt(scsmsg.Status))
		return
	}
	if forceToMainnet {
		self.scsMsgChanMainnet <- scsmsg
	} else {
		// msg pushed into scsmsgchan will be send out by pm's shardingBroadcastLoop() function
		self.scsMsgChan <- scsmsg
	}

	// broadcast messages to connected scs
	scsPushMsgChan := self.scsHandler.GetScsPushMsgChan()
	scsPushMsgChan <- *scsmsg
}

func (self *NetworkRelay) resetMsgCounts() {
	mu.Lock()
	defer mu.Unlock()

	// log.Infof("resetMsgCounts")
	self.scsMsgCounts = make(map[string]uint)
}

func (self *NetworkRelay) ShouldBroadcast(msg *pb.ScsPushMsg) bool {
	mu.Lock()
	defer mu.Unlock()

	msgType := common.BytesToInt(msg.Type)
	msgStatus := common.BytesToInt(msg.Status)
	// Pass NewBlock msg and flush msg
	if (msgType == params.BroadCast && msgStatus == params.NewBlock) || (msgType == params.ControlMsg &&
		(msgStatus == params.CreateProposal || msgStatus == params.DisputeProposal ||
			msgStatus == params.ApproveProposal || msgStatus == params.UpdateLastFlushBlk ||
			msgStatus == params.DistributeProposal || msgStatus == params.UploadRedeemData ||
			msgStatus == params.EnterAndRedeem)) {
		return true
	}

	scsId := common.ToHex(msg.Scsid)
	subChainId := common.ToHex(msg.Subchainid)
	if self.scsMsgCounts[scsId] == 0 {
		self.scsMsgCounts[scsId] = 1
	} else {
		self.scsMsgCounts[scsId]++
	}
	if self.scsMsgCounts[subChainId] == 0 {
		self.scsMsgCounts[subChainId] = 1
	} else {
		self.scsMsgCounts[subChainId]++
	}

	if self.scsMsgCounts[scsId] > scsMsgLimit || self.scsMsgCounts[subChainId] > subChainMsgLimit {
		return false
	}
	return true
}

type DecodedScsRegisterMsg struct {
	ScsId      string
	RequestId  string
	Operation  string
	PublicKey  string
	Capability uint32
}

type DecodedScsPushMsg struct {
	ScsId      string
	RequestId  string
	Sender     common.Address
	Receiver   common.Address
	MsgType    string
	SubChainId string
	Status     string
	MsgHash    []byte
}

type registerPayload struct {
	Operation  string `json:"operation"         gencodec:"required"`
	PublicKey  string `json:"publickey"         gencodec:"required"`
	Capability uint32 `json:"capability"        gencodec:"required"`
}

func decodeScsRegisterMsg(raw *pb.ScsPushMsg) DecodedScsRegisterMsg {
	log.Debugf("[mc/network_relay.go->decodeScsRegisterMsg] %v", raw)
	var (
		ScsId     string
		RequestId string
		RegMsg    registerPayload
	)

	rlp.Decode(bytes.NewReader(raw.Scsid), &ScsId)
	rlp.Decode(bytes.NewReader(raw.Requestid), &RequestId)
	rlp.Decode(bytes.NewReader(raw.Msghash), &RegMsg)

	decodedScsRegisterMsg := DecodedScsRegisterMsg{
		ScsId:      ScsId,
		RequestId:  RequestId,
		Operation:  RegMsg.Operation,
		PublicKey:  RegMsg.PublicKey,
		Capability: RegMsg.Capability,
	}

	return decodedScsRegisterMsg
}

func decodeScsPushMsg(raw *pb.ScsPushMsg, scsId string) DecodedScsPushMsg {
	log.Debugf("[mc/network_relay] decodeScsPushMsg %v", raw)
	var (
		RequestId  string
		ScsId      string
		Sender     common.Address
		Receiver   common.Address
		MsgType    string
		SubChainId string
		Status     string
		//RegMsg     []byte
	)

	rlp.Decode(bytes.NewReader(raw.Sender), &Sender)
	rlp.Decode(bytes.NewReader(raw.Receiver), &Receiver)
	rlp.Decode(bytes.NewReader(raw.Type), &MsgType)
	rlp.Decode(bytes.NewReader(raw.Subchainid), &SubChainId)
	rlp.Decode(bytes.NewReader(raw.Status), &Status)
	rlp.Decode(bytes.NewReader(raw.Scsid), &ScsId)
	rlp.Decode(bytes.NewReader(raw.Requestid), &RequestId)

	decodedScsPushMsg := DecodedScsPushMsg{
		ScsId:      ScsId,
		RequestId:  RequestId,
		Sender:     Sender,
		Receiver:   Receiver,
		MsgType:    MsgType,
		SubChainId: SubChainId,
		Status:     Status,
		MsgHash:    raw.Msghash,
	}

	return decodedScsPushMsg
}

func (self *NetworkRelay) OnReceiveMsg(msg *pb.ScsPushMsg) {
	log.Debugf("[mc/network_relay] OnReceiveMsg %s", string(msg.Requestid))
	self.VnodePushMsg(msg)
}

func (self *NetworkRelay) OnReceiveRes(res *pb.ScsPushMsg) {
	//Doing nothing for now
}

func (self *NetworkRelay) OnReceiveRegisterMsg(msg *pb.ScsPushMsg) {
	//Doing nothing for now
}

func (self *NetworkRelay) GetBlockNumber() uint64 {
	return self.scsHandler.GetBlockNumber()
}

func (self *NetworkRelay) NotifyMsgRunState(hash common.Hash) bool {
	return self.scsHandler.NotifyMsgRunState(hash)
}

func (self *NetworkRelay) GetSCSRole(subchainId, scsId common.Address) params.ScsKind {
	var scsRole params.ScsKind
	curNumber := self.GetBlockNumber()
	roleInfo := self.GetScsRoleCache(scsId, subchainId)
	if roleInfo != nil && roleInfo.Role == params.LockScs {
		return params.LockScs
	}
	//Three questions need to get scs role: 1.Cleared; 2.Old Cache; 3. smaller than maxRoleCache
	condition := (roleInfo == nil || roleInfo.Number == 0 || curNumber > roleInfo.Number+maxRoleCache || curNumber < maxRoleCache)
	if condition {
		scsRole = self.scsHandler.GetSCSRole(subchainId, scsId)
		self.UpdateScsRoleCache(scsId, subchainId, scsRole)
	} else {
		scsRole = roleInfo.Role
	}
	if roleInfo != nil {
		log.Debugf("GetSCSRole cache info scsRole:%v, number:%v, curNumber:%v", scsRole, roleInfo.Number, curNumber)
	} else {
		log.Debugf("GetSCSRole cache info scsRole:%v, number:0, curNumber:%v, roleInfo is nil", scsRole, curNumber)
	}

	return scsRole
}

func (self *NetworkRelay) ScsMsgCheck(role params.ScsKind, scs common.Address, conReq *pb.ScsPushMsg) bool {
	types, status := common.BytesToInt(conReq.Type), common.BytesToInt(conReq.Status)
	toScs := common.BytesToAddress(conReq.GetReceiver())
	log.Debugf("[mc/network_relay.go->ScsMsgCheck] role:%v, types:%v, status:%v", role, types, status)
	if types == params.ScsPing || types == params.ScsShakeHand {
		return true
	} else if types == params.BroadCast && (status == params.SyncRequest || status == params.SyncComplete) {
		return (toScs == scs)
	} else {
		switch role {
		case params.None:
			return false

		case params.ConsensusScs:
			if types == params.ControlMsg && (status == params.RegOpen || status == params.RegAdd || status == params.RegAsMonitor || status == params.RegAsBackup) {
				return false
			} else {
				return true
			}
		case params.MonitorScs:
			if types == params.BroadCast || types == params.DirectCall || (types == params.ControlMsg && (status == params.DistributeProposal || status == params.RegAsMonitor ||
				status == params.ResetAll)) {
				return true
			} else {
				return false
			}
		case params.BackupScs:
			if types == params.BroadCast || (types == params.ControlMsg && (status == params.DistributeProposal || status == params.RegAsBackup || status == params.ResetAll)) {
				return true
			} else {
				return false
			}
		case params.MatchSelTarget:
			if types == params.ControlMsg && (status == params.RegOpen || status == params.RegAdd || status == params.DistributeProposal || status == params.RequestRelease) {
				return true
			} else {
				return false
			}
		case params.LockScs:
			return false
		default:
			return false
		}
	}
}

// send push msg to all connected scs servers
func (self *NetworkRelay) VnodePushMsg(conReq *pb.ScsPushMsg) (map[int]*pb.ScsPushMsg, error) {
	consensusAddress := common.BytesToAddress(conReq.GetSubchainid())
	log.Debugf("[mc/network_relay.go->VnodePushMsg] called: %v", consensusAddress.String())
	serverList := self.ServerWithoutMsg(conReq)
	if self != nil && serverList != nil {
		for _, scsConnection := range serverList {
			if len(scsConnection.ScsId) < 5 || scsConnection.ScsId == "0x0000000000000000000000000000000000000000" {
				log.Debugf("skipping invalid scsid: %v", scsConnection.ScsId)
				continue
			}

			requestID := fmt.Sprintf("%s", conReq.Requestid)
			if strings.Contains(requestID, scsConnection.ScsId) {
				log.Debugf("skipping self message to scsid: %v", scsConnection.ScsId)
				self.MarkMsg(scsConnection.ScsId, conReq)
				continue
			}

			scsRole := self.GetSCSRole(consensusAddress, common.HexToAddress(scsConnection.ScsId))
			checkResult := self.ScsMsgCheck(scsRole, common.HexToAddress(scsConnection.ScsId), conReq)
			log.Debugf("VnodePushMsg scsRole:%v, checkResult:%v", scsRole, checkResult)
			if checkResult {
				if len(requestID) > common.AddressLength {
					self.MarkMsg(scsConnection.ScsId, conReq)
				}
				self.sendPushMsg(scsConnection, conReq)
				subchainid := common.Bytes2Hex(conReq.Subchainid)
				// take last 40 bytes
				subchainid = fmt.Sprintf("0x%s", subchainid[len(subchainid)-40:])
				log.Debugf(
					"!!! sending to scs: %v, Requestid: %s, type: %v, status: %v, subchainID: %s",
					scsConnection.ScsId, requestID[:16], conReq.Type, conReq.Status, subchainid,
				)
			}
		}
	}

	return nil, nil
}

func (self *NetworkRelay) MarkMsg(scsId string, msg *pb.ScsPushMsg) {
	if self.knownMsgs == nil {
		self.knownMsgs = set.New()
	}
	id := common.RlpHash([]interface{}{msg.Requestid, msg.Scsid, msg.Requestflag, msg.Msghash})
	// If we reached the memory allowance, drop a previously known transaction hash
	for self.knownMsgs.Size() >= maxKnownMsgs {
		self.knownMsgs.Pop()
	}
	self.knownMsgs.Add(id)
}

func (self *NetworkRelay) ServerWithoutMsg(msg *pb.ScsPushMsg) []*libtypes.ScsServerConnection {
	if self.ScsServerList == nil {
		return nil
	}

	id := common.RlpHash([]interface{}{msg.Requestid, msg.Scsid, msg.Requestflag, msg.Msghash})
	list := make([]*libtypes.ScsServerConnection, 0, len(*self.ScsServerList))
	for _, svr := range *self.ScsServerList {
		if self.knownMsgs == nil || !self.knownMsgs.Has(id) {
			list = append(list, svr)
		}
	}
	return list
}

func (self *NetworkRelay) UpdateScsRoleCache(scsId common.Address, subchainId common.Address, role params.ScsKind) {
	self.lockRole.Lock()
	defer self.lockRole.Unlock()

	if _, ok := self.ScsRoleCache[scsId]; !ok {
		self.ScsRoleCache[scsId] = make(map[common.Address]*RoleCache)
	}
	self.ScsRoleCache[scsId][subchainId] = &RoleCache{Role: role, Number: self.GetBlockNumber()}
}

func (self *NetworkRelay) DelScsRoleCache(scsId common.Address, subchainId common.Address) {
	self.lockRole.Lock()
	defer self.lockRole.Unlock()

	if roleList, ok := self.ScsRoleCache[scsId]; ok {
		if _, geted := roleList[subchainId]; geted {
			delete(self.ScsRoleCache[scsId], subchainId)
		}
	}
}

func (self *NetworkRelay) GetScsRoleCache(scsId common.Address, subchainId common.Address) *RoleCache {
	self.lockRole.RLock()
	defer self.lockRole.RUnlock()

	if roleList, ok := self.ScsRoleCache[scsId]; ok {
		if roleInfo, geted := roleList[subchainId]; geted {
			return roleInfo
		}
	}

	return nil
}

func (self *NetworkRelay) SetScsRoleCacheOld(subchainId common.Address) {
	self.lockRole.Lock()
	defer self.lockRole.Unlock()

	if self.ScsServerList == nil {
		return
	}

	for _, svr := range *self.ScsServerList {
		scsId := common.HexToAddress(svr.ScsId)
		if roleList, ok := self.ScsRoleCache[scsId]; ok {
			if roleInfo, geted := roleList[subchainId]; geted && roleInfo.Role != params.LockScs {
				roleInfo.Number = 0
			}
		}
	}
}

type DecodedPayload struct {
	Addr      common.Address
	TimeStamp int64
	MsgHash   []byte
}

type DecodedVnodePushMsg struct {
	RequestId  uint32
	Sender     common.Address
	Receiver   common.Address
	Status     string
	SubChainId string
	TypeRaw    string
	Payload    DecodedPayload
}

func (self *NetworkRelay) sendPushMsg(scsConnection *libtypes.ScsServerConnection, req *pb.ScsPushMsg) {
	select {
	case scsConnection.Req <- req:
	default:
		log.Errorf("[mc/network_relay.go->sendPushMsg] channel is full, disregard %v", req)
	}
}

func (self *NetworkRelay) SetupMsgSender(scsId string) {
	scsServer := (*self.ScsServerList)[scsId]
	go func() {
		for {
			select {
			case req := <-scsServer.Req:
				stream := *scsServer.Stream
				if err := stream.Send(req); err != nil {
					log.Errorf("Failed to send push msg: %v", err)
				}
			case <-scsServer.Cancel:
				return
			}
		}
	}()
}
