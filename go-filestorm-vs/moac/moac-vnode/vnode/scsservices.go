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

package vnode

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"errors"
	"fmt"
	"hash"
	"io"
	"math/big"
	"net"
	"os"
	"sync"
	"time"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/crypto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-lib/mcdb"
	pb "github.com/filestorm/go-filestorm/moac/moac-lib/proto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
	libtypes "github.com/filestorm/go-filestorm/moac/moac-lib/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/contracts"
	ctypes "github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/vm"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/node"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/vnode/config"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/vnode/types"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	Server             = &VnodeServer{}
	Grequestid         uint32
	mu                 sync.Mutex
	scsPushMsgChanSize uint32 = 1000
	scsRetryLimit      uint   = 10
)

type Backend interface {
	BlockChain() *core.BlockChain
	TxPool() *core.TxPool
	ChainDb() mcdb.Database
	IsSubnetP2PEnabled(contractAddress common.Address, where string) bool
}

// VNODE server data structure
type VnodeServer struct {
	Config         *config.Configuration
	chaindb        mcdb.Database
	blockchain     *core.BlockChain
	moacnode       Backend
	ctx            *node.ServiceContext
	ScsPushMsgChan chan pb.ScsPushMsg
	ScsPushResChan chan pb.ScsPushMsg
	ScsRegChan     chan pb.ScsPushMsg
	ScsServerList  *map[string]*libtypes.ScsServerConnection
	PingTickers    map[string]*time.Ticker
	syncdb         mcdb.Database
}

// SCSService, API functions, returns if the SCSService is open (= true) or not (=false), indication if the Vnode is providing SCS service or not.
func (s *VnodeServer) ScsService() bool {
	return s.Config.SCSService
}

func (s *VnodeServer) GetScsPushMsgChan() chan pb.ScsPushMsg { return s.ScsPushMsgChan }
func (s *VnodeServer) GetScsRegChan() chan pb.ScsPushMsg     { return s.ScsRegChan }
func (s *VnodeServer) GetScsPushResChan() chan pb.ScsPushMsg { return s.ScsPushResChan }

func NewScsService(db mcdb.Database, syncdb mcdb.Database, bc *core.BlockChain, mn Backend, ctx *node.ServiceContext) *VnodeServer {
	log.Debugf("[vnode/scsservices.go->NewScsService]")
	Server.chaindb = db
	Server.syncdb = syncdb
	Server.blockchain = bc
	Server.moacnode = mn
	Server.ctx = ctx
	Server.ScsRegChan = make(chan pb.ScsPushMsg)
	Server.ScsPushMsgChan = make(chan pb.ScsPushMsg, scsPushMsgChanSize) // TODO
	Server.ScsPushResChan = make(chan pb.ScsPushMsg)
	ScsServerList := make(map[string]*libtypes.ScsServerConnection)
	Server.ScsServerList = &ScsServerList
	Server.PingTickers = make(map[string]*time.Ticker)

	return Server
}

//Close close the VNODE scs service
func (s *VnodeServer) Close() {
}

//Singleton
func GetInstance() *VnodeServer {
	return Server
}

func getPubkeyHash(pb []byte) (pbhs string) {
	log.Debugf("[vnode/scsservices.go->getPubkeyHash]")
	var hs hash.Hash
	hs = md5.New()
	io.WriteString(hs, string(pb))
	return fmt.Sprintf("%x", hs.Sum(nil))
}

func (s *VnodeServer) ConnectScsServers(scsid string, blknum *big.Int) error {
	log.Debugf("[vnode/scsservices.go->Server.ConnectScsServers] scsid %v", scsid)
	scsServerList := *s.ScsServerList
	if scsServerList[scsid] != nil {
		scsServer := scsServerList[scsid]
		log.Debugf("scsId %v, scsServer %v", scsid, scsServer)

		liveinfo := types.LiveInfo{CurrentBlockNum: blknum}
		liveBytes, _ := rlp.EncodeToBytes(liveinfo)

		nodeObj := node.GetInstance()
		var requestId []byte
		if nodeObj != nil {
			requestId = nodeObj.GetRequestId(true)
		} else {
			requestId = []byte("")
		}
		conReq := &pb.ScsPushMsg{
			Requestid:   requestId,
			Timestamp:   common.Int64ToBytes(time.Now().Unix()),
			Requestflag: true,
			Type:        common.IntToBytes(params.ScsPing),
			Status:      []byte(""),
			Scsid:       []byte(scsid),
			Subchainid:  []byte(""),
			Sender:      []byte(""),
			Receiver:    []byte(""),
			Msghash:     liveBytes,
		}
		stream := *scsServer.Stream
		if err := stream.Send(conReq); err != nil {
			log.Errorf("do Ping error!!! err:" + err.Error())
			scsServer.LiveFlag = false
			scsServer.RetryCount++
			if scsServer.RetryCount > scsRetryLimit {
				s.CloseScsServer(scsid)
			}
		} else {
			log.Debugf("Ping reply from scs: %v LiveFlag true", conReq.Timestamp)
			scsServer.LiveFlag = true
			scsServer.RetryCount = 0
		}

	} else {
		err := errors.New("Scs server not exists")
		log.Debugf("ConnectScsServers %v %v", scsid, err)
		return err
	}

	return nil
}

type DecodedScsRegisterMsg struct {
	ScsId      string
	RequestId  string
	Operation  string
	PublicKey  string
	Capability uint32
	ChainId    *big.Int
}

type registerPayload struct {
	Operation  string   `json:"operation"         gencodec:"required"`
	PublicKey  string   `json:"publickey"         gencodec:"required"`
	Capability uint32   `json:"capability"        gencodec:"required"`
	ChainId    *big.Int `json:"ChainId"           gencodec:"required"`
}

func decodeScsRegisterMsg(raw *pb.ScsPushMsg) DecodedScsRegisterMsg {
	log.Debugf("[vnode/scsservices.go->decodeScs`RegisterMsg] %v", raw)
	var (
		ScsId     string
		RequestId string
		RegMsg    registerPayload
	)

	ScsId = string(raw.Scsid)
	RequestId = string(raw.Requestid)
	rlp.Decode(bytes.NewReader(raw.Msghash), &RegMsg)

	decodedScsRegisterMsg := DecodedScsRegisterMsg{
		ScsId:      ScsId,
		RequestId:  RequestId,
		Operation:  RegMsg.Operation,
		PublicKey:  RegMsg.PublicKey,
		Capability: RegMsg.Capability,
		ChainId:    RegMsg.ChainId,
	}

	return decodedScsRegisterMsg
}

func (s *VnodeServer) UpdateScsServerList(shakeInfo *types.ShakeInfo) error {
	log.Debugf("[vnode/scsservices.go->Server.UpdateScsServerList] shakeInfo %v scsId %v LiveFlag true", shakeInfo, shakeInfo.Scsid)
	if len(shakeInfo.Scsid) < 5 || shakeInfo.Scsid == "0x0000000000000000000000000000000000000000" {
		return fmt.Errorf("Invalid scsId %v", shakeInfo.Scsid)
	}
	scsServerList := *s.ScsServerList
	if scsServerList[shakeInfo.Scsid] != nil {
		s.CloseScsServer(shakeInfo.Scsid)
	}
	scsServerList[shakeInfo.Scsid] = &libtypes.ScsServerConnection{
		ScsId:      shakeInfo.Scsid,
		LiveFlag:   true,
		Stream:     shakeInfo.Stream,
		Req:        make(chan *pb.ScsPushMsg, 32768),
		Cancel:     make(chan bool),
		RetryCount: 0,
	}

	log.Debugf("[vnode/scsservices.go->Server.UpdateScsServerList] list %v", s.ScsServerList)
	return nil
}

func (s *VnodeServer) CloseScsServer(scsId string) {
	log.Infof("closing scs %v", scsId)
	scsServerList := *s.ScsServerList
	scsServer := scsServerList[scsId]
	if scsServer != nil {
		scsServer.Cancel <- true
		close(scsServer.Cancel)
		close(scsServer.Req)
		delete(scsServerList, scsId)
	}
}

func (s *VnodeServer) ScsPush(stream pb.Vnode_ScsPushServer) error {
	log.Debugf("[vnode/scsservices.go->ScsPush]")
	nr := core.Nr

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Errorf("[vnode/scsservices.go->ScsPush] %v", err)
			break
		}
		log.Debugf("ScsPush got msg: %v", string(in.GetRequestid()))

		if s != nil &&
			s.moacnode != nil &&
			s.moacnode.BlockChain() != nil &&
			s.moacnode.BlockChain().Config() != nil &&
			s.moacnode.BlockChain().CurrentBlock() != nil &&
			s.moacnode.BlockChain().CurrentBlock().Header() != nil &&
			!s.moacnode.BlockChain().Config().IsNuwa(s.moacnode.BlockChain().CurrentBlock().Header().Number) {
			log.Error("not supported until Nuwa block#")
			break
		}

		log.Debugf("here")
		var operationReplyBytes []byte
		typeInt := common.BytesToInt(in.Type)
		switch typeInt {
		case params.ScsShakeHand:
			go func() {
				s.ScsRegChan <- *in
			}()

			// Not initialize
			if nr == nil || Server == nil || Server.chaindb == nil {
				log.Debug("[vnode/scsservices.go->ScsPush]vnode hasn't been initialized")
				operationReplyBytes, _ = rlp.EncodeToBytes("Not initialize")
				msgs := []*pb.ScsPushMsg{
					{
						Requestid:   in.Requestid,
						Timestamp:   common.Int64ToBytes(time.Now().Unix()),
						Requestflag: false,
						Type:        in.Type,
						Status:      in.Status,
						Scsid:       in.Scsid,
						Subchainid:  nil,
						Sender:      common.FromHex(params.VnodeBeneficialAddress),
						Receiver:    nil,
						Msghash:     operationReplyBytes,
					},
				}
				if err := stream.Send(msgs[0]); err != nil {
					log.Errorf("send error %v", err)
					//return err
				}
				continue
			}

			decodedIn := decodeScsRegisterMsg(in)
			pbhs := getPubkeyHash([]byte(decodedIn.PublicKey))
			log.Debugf("Convert pbhs: %v", pbhs)

			//save scs shake info
			oldinfo := core.GetScsShakeInfo(Server.chaindb, decodedIn.ScsId)
			log.Debugf("got shake info from db: %v", oldinfo)
			chainId := params.MainnetChainConfig.ChainId.Int64()
			shakeinfo := types.ShakeInfo{
				Pbhs:    pbhs,
				Scsid:   decodedIn.ScsId,
				Stream:  &stream,
				ChainId: chainId,
			}
			if oldinfo == shakeinfo {
				log.Debugf("shake info is old")
			}

			validChainid := decodedIn.ChainId.Cmp(s.moacnode.BlockChain().Config().ChainId) == 0
			if validChainid {
				log.Debugf("save shake info to db: %v", shakeinfo)
				core.WriteScsShakeInfo(Server.chaindb, decodedIn.ScsId, shakeinfo)
				err = s.UpdateScsServerList(&shakeinfo)
				nr.ScsServerList = s.ScsServerList
				if err != nil {
					errMsg := fmt.Sprintf("shakereply: error: invalid scsId %v", shakeinfo.Scsid)
					operationReplyBytes, _ = rlp.EncodeToBytes(errMsg)
				} else {
					nr.SetupMsgSender(shakeinfo.Scsid)
					operationReplyBytes, _ = rlp.EncodeToBytes("shakereply")
				}
			} else {
				log.Errorf("ChainIdError: Local ChainId: %v, SCS ChainId: %v ", s.moacnode.BlockChain().ChainId(), decodedIn.ChainId)
				operationReplyBytes, _ = rlp.EncodeToBytes("chainiderror")
			}

			msgs := []*pb.ScsPushMsg{
				{
					Requestid:   in.Requestid,
					Timestamp:   common.Int64ToBytes(time.Now().Unix()),
					Requestflag: false,
					Type:        in.Type,
					Status:      in.Status,
					Scsid:       in.Scsid,
					Subchainid:  nil,
					Sender:      common.FromHex(params.VnodeBeneficialAddress),
					Receiver:    nil,
					Msghash:     operationReplyBytes,
				},
			}
			if err := stream.Send(msgs[0]); err != nil {
				log.Errorf("send error %v", err)
				continue
			}

			//connect and ping scs
			log.Debugf("before PingTicker")
			if validChainid && s.PingTickers[decodedIn.ScsId] == nil {
				log.Debugf("PingTicker for %v is nil", decodedIn.ScsId)
				s.PingTickers[decodedIn.ScsId] = time.NewTicker(params.TimerPingInterval * time.Second)
				CurPingBlkNum := new(big.Int)
				PingTickerRepeat := 0
				for {
					select {
					case <-s.PingTickers[decodedIn.ScsId].C:
						if CurPingBlkNum.Cmp(s.blockchain.CurrentBlock().Header().Number) < 0 {
							CurPingBlkNum = s.blockchain.CurrentBlock().Header().Number
							PingTickerRepeat = 0
							err := s.ConnectScsServers(decodedIn.ScsId, s.blockchain.CurrentBlock().Header().Number)
							if err != nil {
								log.Debugf("failed to connect to scs server error %v", err)
							}

						} else {
							PingTickerRepeat++
						}
					}

				}
			}

		case params.ScsPing:
			operationReplyBytes, _ := rlp.EncodeToBytes("pingreply")
			msgs := []*pb.ScsPushMsg{
				{
					Requestid:   in.Requestid,
					Timestamp:   common.Int64ToBytes(time.Now().Unix()),
					Requestflag: false,
					Type:        in.Type,
					Status:      in.Status,
					Scsid:       in.Scsid,
					Subchainid:  []byte(""),
					Sender:      []byte(""),
					Receiver:    []byte(""),
					Msghash:     operationReplyBytes},
			}

			go func() {
				mu.Lock()
				defer mu.Unlock()
				if err := stream.Send(msgs[0]); err != nil {
					log.Debugf("pinging: %v", err)
				}
			}()

		case params.DirectCall:
			//TODO: these are response from scs, need to handle them.
		case params.ControlMsg:
		case params.BroadCast:
			forceToMainnet := false
			contractAddress := common.BytesToAddress(in.Subchainid)
			if !Server.moacnode.IsSubnetP2PEnabled(contractAddress, "scsservices.go/scspush()") {
				forceToMainnet = true
			}
			nr.BroadcastMsg(in, forceToMainnet)
		default:
			go func() {
				s.ScsPushMsgChan <- *in
			}()

			log.Errorf("steam operation error: typeInt %v", typeInt)
		}
	}
	return nil
}

//AccountInfo:
func (s *VnodeServer) AccountInfo(ctx context.Context, in *pb.AccountInfoRequest) (*pb.AccountInfoReply, error) {
	if !s.moacnode.BlockChain().Config().IsNuwa(s.moacnode.BlockChain().CurrentBlock().Header().Number) {
		return nil, errors.New("not supported until Nuwa block#")
	}

	// var account common.Address
	account := common.BytesToAddress(in.Addr)

	st, _ := Server.blockchain.State()
	var bls = st.GetBalance(account)
	var nus = s.moacnode.TxPool().State().GetNonce(account)
	var cdh = st.GetCodeHash(account)
	qry := uint64(0) //Query flag, set to 0 and may remove later
	var shd, _ = st.GetFlag(account)
	var cbk, wbk, _ = st.GetFlushInfo(account)
	log.Debugf("[vnode/scsservices.go->AccountInfo] account address:%v", account.String())
	log.Debugf("[vnode/scsservices.go->AccountInfo] account nonce:%v", nus)
	accountinfo := types.AccountInfo{Addr: account, Balance: bls, Nonce: nus, CodeHash: cdh, Query: qry,
		Shard: shd, CreationBlockNumber: cbk, WaitBlockNumber: wbk}

	replybd, _ := rlp.EncodeToBytes(accountinfo)
	ret := &pb.AccountInfoReply{Requestid: in.Requestid, Replybody: replybd}
	return ret, nil
}

//ChainInfo: MicroChain info?
func (s *VnodeServer) ChainInfo(ctx context.Context, in *pb.ChainInfoRequest) (*pb.ChainInfoReply, error) {
	conaddr := common.BytesToAddress(in.Consensusaddr)
	if !s.moacnode.BlockChain().Config().IsNuwa(s.moacnode.BlockChain().CurrentBlock().Header().Number) {
		return nil, errors.New("not supported until Nuwa block#")
	}

	st, _ := Server.blockchain.State()
	var res = st.DumpContractStorage(conaddr, in.Request)
	log.Debugf("[vnode/scsservices.go->ChainInfo] consensus address:%v", conaddr.String())

	ret := &pb.ChainInfoReply{Requestid: in.Requestid, Replybody: res}
	return ret, nil
}

func (s *VnodeServer) RemoteCall(ctx context.Context, in *pb.RemoteCallRequest) (*pb.RemoteCallReply, error) {
	if !s.moacnode.BlockChain().Config().IsNuwa(s.moacnode.BlockChain().CurrentBlock().Header().Number) {
		return nil, errors.New("not supported until Nuwa block#")
	}

	var sender common.Address
	sender = common.BytesToAddress(in.Sender)
	log.Debugf("RemoteCall: sender address:%v", sender.String())

	var conaddr common.Address
	conaddr = common.BytesToAddress(in.Contractaddr)
	log.Debugf("RemoteCall: conaddr address:%v", conaddr.String())

	var signedTx ctypes.Transaction
	if err := rlp.DecodeBytes(in.Data, &signedTx); err != nil {
		log.Errorf("failed to DecodeBytes body: %v", err)
	}

	log.Info("[RemoteCall->SendTx]")
	if err := s.moacnode.TxPool().AddLocal(&signedTx); err != nil {
		log.Debugf("[RemoteCall->SendTxErr]%v", err)
		return nil, err
	}

	// print left balance
	st, _ := Server.blockchain.State()
	log.Debugf("RemoteCall GetBalance('%s'): %d", sender.String(), st.GetBalance(sender))

	replybd, _ := rlp.EncodeToBytes("success")
	ret := &pb.RemoteCallReply{Requestid: in.Requestid, Replybody: replybd}
	return ret, nil
}

func (s *VnodeServer) GetBlockNumber() uint64 {
	return s.blockchain.CurrentBlock().NumberU64()
}

//GetContractAddressMappingMember:
func (s *VnodeServer) GetContractAddressMappingMember(addr common.Address, pos uint32, key common.Address) ([]common.Address, error) {
	// first get state db
	st, _ := Server.blockchain.State()

	// getstate member object by position (ignore constant member)
	posbytes := make([]byte, 32)
	binary.LittleEndian.PutUint32(posbytes, pos)
	lenhash := st.GetState(addr, common.BytesToHash(posbytes))

	//allocate dynamic array
	len := int64(binary.BigEndian.Uint64(lenhash[:]))
	var arry []common.Address

	indexbigint := new(big.Int)

	indexbigint.SetBytes(crypto.Keccak256(posbytes))

	//get each element
	for i := int64(0); i < len; i++ {
		indexbigint = indexbigint.Add(indexbigint, big.NewInt(i))
		indexhash := st.GetState(addr, common.BytesToHash(indexbigint.Bytes()))
		arry = append(arry, common.BytesToAddress(indexhash.Bytes()[12:]))
	}

	return arry, nil
}

func (s *VnodeServer) NotifyMsgRunState(hash common.Hash) bool {
	if receipt, _, _, _ := core.GetReceipt(Server.chaindb, hash); receipt != nil {
		return !receipt.Failed
	}
	return false
}

//GetSCSRole: return the connected SCS info.
func (s *VnodeServer) GetSCSRole(contractAddress common.Address, nodeAddress common.Address) params.ScsKind {
	isMonitorHash := "0x50859fd9"
	data := common.FromHex(isMonitorHash)
	nodeAddressBytes := nodeAddress.Bytes()
	data = append(data, common.LeftPadBytes(nodeAddressBytes, 32)...)
	st, _ := Server.blockchain.State()
	if codeHash := st.GetCodeHash(contractAddress); codeHash == (common.Hash{}) {
		return params.None
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

	context := core.NewEVMContext(msg, Server.blockchain.CurrentBlock().Header(), Server.blockchain, nil, nil)
	evm := vm.NewEVM(context, st, params.TestChainConfig,
		vm.Config{EnableJit: false, ForceJit: false}, nil)

	contractRef := vm.AccountRef(contractAddress)
	log.Debugf("GetSCSRole contractAddress=%v nodeAddress=%v", contractAddress.String(), nodeAddress.String())

	precompiledContracts := contracts.GetInstance()

	ret, leftGas, err := evm.Call(contractRef, contractAddress, data, params.GenesisGasLimit.Uint64(), big.NewInt(0), false, uint64(0), precompiledContracts, msg.GetMsgHash())
	if err != nil {
		log.Errorf("GetSCSRole error %v", err)
		return params.LockScs
	}
	log.Debugf("GetSCSRole leftGas %v, ret:%v", leftGas, ret)

	var retString string
	retString = common.Bytes2Hex(ret)
	log.Debugf("GetSCSRole ret %v", retString)
	switch retString {
	case "0000000000000000000000000000000000000000000000000000000000000000":
		return params.None
	case "0000000000000000000000000000000000000000000000000000000000000001":
		return params.ConsensusScs
	case "0000000000000000000000000000000000000000000000000000000000000002":
		return params.MonitorScs
	case "0000000000000000000000000000000000000000000000000000000000000003":
		return params.BackupScs
	case "0000000000000000000000000000000000000000000000000000000000000004":
		return params.MatchSelTarget
	default:
		return params.None
	}
}

// func (s *Server) GetNodeList(addr common.Address) []common.Address {
// 	var nodeList []common.Address
// 	st, _ := ServerObj.blockchain.State()
// 	v := "000000000000000000000000000000000000000000000000000000000000000f"
// 	request := make([]*pb.StorageRequest, 0)
// 	request = append(request, &pb.StorageRequest{1, common.Hex2Bytes(v), []byte{}, []byte{}})
// 	res := st.DumpContractStorage(addr, request)
// 	var contractInfo types.ContractInfo
// 	if err := json.Unmarshal(res, &contractInfo); err != nil {
// 		log.Errorf("failed to Unmarshal body: %v", err)
// 	}
// 	log.Debugf("GetNodeList contractInfo.Storage:%v", contractInfo.Storage)
// 	num, _ := strconv.ParseInt(common.GetValue(contractInfo.Storage, v), 16, 64)
// 	key := common.KeytoKey(v)
// 	for i := int64(0); i < num; i++ {
// 		if str := common.GetValue(contractInfo.Storage, key); str != "" {
// 			nodeList = append(nodeList, common.HexToAddress(str))
// 		}
// 		key = common.IncreaseHexByOne(key)
// 	}
// 	return nodeList
// }

// func (s *Server) IsMonitor(contractAddress common.Address, nodeAddress common.Address) bool {
// 	monitorFlag := false

// 	isMonitorHash := "0x562aa201"
// 	data := common.FromHex(isMonitorHash)
// 	nodeAddressBytes := nodeAddress.Bytes()
// 	data = append(data, common.LeftPadBytes(nodeAddressBytes, 32)...)

// 	nr := core.Nr

// 	st, _ := ServerObj.blockchain.State()

// 	to := common.Address{}
// 	to.SetString("")
// 	from := common.Address{}
// 	viaaddress := common.Address{}
// 	viaaddress.SetString(*ServerObj.UCfg.VnodeBeneficialAddress)
// 	msgHash := common.Hash{}
// 	msgHash.SetString("")
// 	msg := Message{
// 		to:              &to,
// 		from:            from,
// 		nonce:           0,
// 		amount:          big.NewInt(0),
// 		price:           big.NewInt(0),
// 		gasLimit:        big.NewInt(0),
// 		data:            []byte{},
// 		checkNonce:      false,
// 		system:          0,
// 		autoFlush:       false,
// 		waitBlockNumber: big.NewInt(0),
// 		via:             &viaaddress,
// 		shardFlag:       0,
// 		msgHash:         &msgHash,
// 	}

// 	context := core.NewEVMContext(msg, ServerObj.blockchain.CurrentBlock().Header(), ServerObj.blockchain, nil, nil)
// 	evm := vm.NewEVM(context, st, params.TestChainConfig,
// 		vm.Config{EnableJit: false, ForceJit: false}, nil)

// 	contractRef := vm.AccountRef(contractAddress)
// 	log.Debugf("IsMonitor contractAddress=%v contractRef=%v", contractAddress, contractRef)

// 	precompiledContracts := contracts.GetInstance()

// 	ret, leftGas, err := evm.Call(contractRef, contractAddress, data, params.GenesisGasLimit.Uint64(), big.NewInt(0), false, uint64(0), precompiledContracts, msg.GetMsgHash())
// 	if err != nil {
// 		log.Errorf("IsMonitor error %v", err)
// 	}
// 	log.Infof("IsMonitor leftGas %v, ret:%v", leftGas, ret)

// 	var retString string
// 	retString = common.Bytes2Hex(ret)
// 	log.Infof("IsMonitor ret %v", retString)
// 	if retString != "0000000000000000000000000000000000000000000000000000000000000000" {
// 		monitorFlag = true
// 	}
// 	return monitorFlag
// }

// func (s *Server) MatchSelTarget(contractAddress common.Address, nodeAddress common.Address) bool {
// 	//TODO test this function
// 	selected := false

// 	matchSelTargetHash := "0xb8598896"
// 	data := common.FromHex(matchSelTargetHash)
// 	nodeAddressBytes := nodeAddress.Bytes()
// 	data = append(data, common.LeftPadBytes(nodeAddressBytes, 32)...)
// 	data = append(data, common.LeftPadBytes([]byte{0}, 32)...)
// 	data = append(data, common.LeftPadBytes([]byte{1}, 32)...)

// 	nr := core.Nr

// 	st, _ := ServerObj.blockchain.State()

// 	to := common.Address{}
// 	to.SetString("")
// 	from := common.Address{}
// 	viaaddress := common.Address{}
// 	vnodeAddress := *getConfig().VnodeBeneficialAddress
// 	viaaddress.SetString(vnodeAddress)
// 	msgHash := common.Hash{}
// 	msgHash.SetString("")
// 	msg := Message{
// 		to:         &to,
// 		from:       from,
// 		nonce:      0,
// 		amount:     big.NewInt(0),
// 		price:      big.NewInt(0),
// 		gasLimit:   big.NewInt(0),
// 		data:       []byte{},
// 		checkNonce: false,
// 		system:     0,
// 		//	syncFlag                bool
// 		autoFlush:       false,
// 		waitBlockNumber: big.NewInt(0),
// 		via:             &viaaddress,
// 		shardFlag:       0,
// 		msgHash:         &msgHash,
// 	}

// 	context := core.NewEVMContext(msg, ServerObj.blockchain.CurrentBlock().Header(), ServerObj.blockchain, nil, nil)
// 	evm := vm.NewEVM(context, st, params.TestChainConfig,
// 		vm.Config{EnableJit: false, ForceJit: false}, nil)

// 	contractRef := vm.AccountRef(contractAddress)

// 	precompiledContracts := contracts.GetInstance()

// 	ret, leftGas, err := evm.Call(contractRef, contractAddress, data, params.GenesisGasLimit.Uint64(), big.NewInt(0), false, uint64(0), precompiledContracts, msg.GetMsgHash())
// 	if err != nil {
// 		log.Errorf("MatchSelTarget error %v", err)
// 	}
// 	log.Debugf("MatchSelTarget leftGas %v", leftGas)

// 	var retString string
// 	retString = common.Bytes2Hex(ret)
// 	log.Debugf("MatchSelTarget ret %v", retString)
// 	if retString != "0000000000000000000000000000000000000000000000000000000000000000" {
// 		selected = true
// 	}

// 	return selected
// }

// VnodeServiceStart -- run grpc server
func VnodeServiceStart(inconfigfile string) {
	log.Info("[vnode/scsservices.go->VnodeServiceStart]")
	if err := loadSCSConfig(inconfigfile); err != nil {
		log.Errorf("Load SCS config error:%v with %v\n", err, inconfigfile)
		os.Exit(1)
	}

	if params.SCSService == true {
		lis, err := net.Listen("tcp", params.VnodeServiceCfg)
		log.Debugf("try to listen to %v lis %v", params.VnodeServiceCfg, lis)
		if err != nil {
			log.Errorf("failed to listen: %v %v", err, params.VnodeServiceCfg)
		} else {
			log.Debugf("listen to %v", params.VnodeServiceCfg)
		}
		s := grpc.NewServer()

		pb.RegisterVnodeServer(s, Server)
		// Register reflection service on gRPC server.
		reflection.Register(s)
		if err := s.Serve(lis); err != nil {
			log.Errorf("failed to serve: %v", err)
		}
	}

}

/*
 * Load in scs configuration
 * for vnode and saved the info in
 * params.SCSService
 *
 */
func loadSCSConfig(configFilePath string) error {
	vnodeConfig, err := config.GetConfiguration(configFilePath)

	if err != nil {
		log.Debugf("Error reading: %v", configFilePath)
		return err
	}

	// No error, continue to setup
	// init config to its default values from params
	Server.Config = &config.Configuration{}
	Server.Config.ShowToPublic = params.ShowToPublic
	Server.Config.VnodeBeneficialAddress = &params.VnodeBeneficialAddress
	Server.Config.VnodeServiceCfg = &params.VnodeServiceCfg
	Server.Config.SCSService = params.SCSService
	Server.Config.VnodeIp = &params.VnodeIp
	Server.Config.ForceSubnetP2P = params.ForceSubnetP2P

	if vnodeConfig != nil {
		//Update the default values with valid inputs
		Server.Config.ShowToPublic = vnodeConfig.ShowToPublic
		Server.Config.SCSService = vnodeConfig.SCSService
		Server.Config.ForceSubnetP2P = vnodeConfig.ForceSubnetP2P

		params.ShowToPublic = Server.Config.ShowToPublic
		params.SCSService = Server.Config.SCSService
		params.ForceSubnetP2P = Server.Config.ForceSubnetP2P

		if vnodeConfig.VnodeServiceCfg != nil {

			Server.Config.VnodeServiceCfg = vnodeConfig.VnodeServiceCfg
			params.VnodeServiceCfg = *vnodeConfig.VnodeServiceCfg
		}
		log.Debugf("Config.VnodeBeneficialAddress is: %v", Server.Config.VnodeBeneficialAddress)

		if vnodeConfig.VnodeBeneficialAddress != nil {
			if common.IsHexAddress(*vnodeConfig.VnodeBeneficialAddress) {
				Server.Config.VnodeBeneficialAddress = vnodeConfig.VnodeBeneficialAddress
				params.VnodeBeneficialAddress = *Server.Config.VnodeBeneficialAddress
			} else {
				//invalid input address, return error
				return errors.New("Invalid VnodeBeneficialAddress")
			}

		}

		if vnodeConfig.VnodeIp != nil {
			Server.Config.VnodeIp = vnodeConfig.VnodeIp
			params.VnodeIp = *Server.Config.VnodeIp
		}

		//Check the input parameters
		log.Debugf("setting params.VnodeServiceCfg to %v", params.VnodeServiceCfg)
		log.Debugf("setting params.VnodeBeneficialAddress to %v", params.VnodeBeneficialAddress)
		log.Debugf("setting params.VnodeIp to %v", params.VnodeIp)

	}

	node.VnodeServiceCfg = Server.Config.VnodeServiceCfg
	node.ShowToPublic = Server.Config.ShowToPublic
	node.Ip = Server.Config.VnodeIp

	//Need to check for a valid MOAC HEX address if SCSService is true
	if Server.Config.VnodeBeneficialAddress != nil {
		if Server.Config.SCSService == true {
			if common.IsHexAddress(*Server.Config.VnodeBeneficialAddress) {
				vnodeBeneficialAddress := common.HexToAddress(*Server.Config.VnodeBeneficialAddress)
				node.VnodeBeneficialAddress = &vnodeBeneficialAddress
				return nil
			} else {
				//
				fmt.Printf("Error in vnodeBeneficialAddress: %v is not a valid HEX address\n", *Server.Config.VnodeBeneficialAddress)
				return errors.New("invalid vnodeBeneficialAddress when SCSService is true")
			}

		}

	} else {
		//Create a new ADDRESS structure and assign the pointer if no VnodeBeneficialAddress info
		emptyAdd := common.Address{}
		node.VnodeBeneficialAddress = &emptyAdd
	}

	return nil
}

func (s *VnodeServer) ScbPublicCall(ctx context.Context, in *pb.ScbPublicCallRequest) (*pb.ScbPublicCallReply, error) {
	if !s.moacnode.BlockChain().Config().IsNuwa(s.moacnode.BlockChain().CurrentBlock().Header().Number) {
		return nil, errors.New("not supported until Nuwa block#")
	}

	contractAddress := common.BytesToAddress(in.Contractaddr)
	st, _ := Server.blockchain.State()

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

	context := core.NewEVMContext(msg, Server.blockchain.CurrentBlock().Header(), Server.blockchain, nil, nil)
	evm := vm.NewEVM(context, st, params.TestChainConfig,
		vm.Config{EnableJit: false, ForceJit: false, DisableGasMetering: true}, nil)

	contractRef := vm.AccountRef(contractAddress)
	log.Debugf("ScbPublicCall contractAddress=%v contractRef=%v input=%s", contractAddress, contractRef, common.Bytes2Hex(in.Data))

	precompiledContracts := contracts.GetInstance()

	ret, leftGas, err := evm.Call(contractRef, contractAddress, in.Data, params.GenesisGasLimit.Uint64(), big.NewInt(0), false, uint64(0), precompiledContracts, msg.GetMsgHash())
	if err != nil {
		log.Errorf("ScbPublicCall error %v", err)
	}
	log.Debugf("ScbPublicCall leftGas: %v, retLen: %v", leftGas, len(ret))
	return &pb.ScbPublicCallReply{Requestid: in.Requestid, Replybody: ret}, nil
}

//UploadBlock - upload receiving the block from SCS, save it to local disk for persistency
func (s *VnodeServer) UploadBlock(ctx context.Context, in *pb.UploadBlockRequest) (*pb.UploadBlockReply, error) {
	log.Debugf("Got Upload Block request msg from scs Requestid:%v Subchainid:%v Sender:%v Blocknumber:%v Blockhash:%v",
		in.Requestid, common.BytesToAddress(in.Subchainid).String(), common.BytesToAddress(in.Sender).String(), in.Blocknumber,
		common.BytesToHash(in.Blockhash).String())
	if !s.moacnode.BlockChain().Config().IsNuwa(s.moacnode.BlockChain().CurrentBlock().Header().Number) {
		return nil, errors.New("not supported until Nuwa block#")
	}

	hash := common.BytesToHash(in.Blockhash)
	if err := core.WriteSyncBlock(Server.syncdb, in.Sender, in.Subchainid, hash, in.Blocknumber, in.Blockdata); err != nil {
		log.Infof("Failed to write the block")
		return nil, err
	}

	cnt := core.GetSyncBlockCnt(Server.syncdb, in.Sender, in.Subchainid, hash) + 1
	core.WriteSyncBlockCnt(Server.syncdb, in.Sender, in.Subchainid, hash, cnt)

	log.Debugf("Upload the block %v times", cnt)
	replybd, _ := rlp.EncodeToBytes("success")
	return &pb.UploadBlockReply{Requestid: in.Requestid, Replybody: replybd}, nil
}

//DownloadBlock - Download the block from vnode so that the SCS are able to retrieve the missing block data
func (s *VnodeServer) DownloadBlock(ctx context.Context, in *pb.DownloadBlockRequest) (*pb.DownloadBlockReply, error) {
	log.Debugf("Got Downoad Block request msg from scs Requestid:%v Subchainid:%v Sender:%v Blocknumber:%v Blockhash:%v",
		in.Requestid, common.BytesToAddress(in.Subchainid).String(), common.BytesToAddress(in.Sender).String(), in.Blocknumber,
		common.BytesToHash(in.Blockhash).String())
	if !s.moacnode.BlockChain().Config().IsNuwa(s.moacnode.BlockChain().CurrentBlock().Header().Number) {
		return nil, errors.New("not supported until Nuwa block#")
	}

	hash := common.Hash{}
	if len(in.Blockhash) > 0 {
		hash = common.BytesToHash(in.Blockhash)
	}

	if (hash == common.Hash{}) {
		hash = core.GetSyncBlockHash(Server.syncdb, in.Sender, in.Subchainid, in.Blocknumber)
	}

	data := []byte{}
	if (hash != common.Hash{}) {
		data = core.GetSyncBlock(Server.syncdb, in.Sender, in.Subchainid, hash, in.Blocknumber)

		//TODO: currently just delete the sync block after download
		//We need create another DB to log the same key but the value as timestamp, so it can delete old by time
		if cnt := core.GetSyncBlockCnt(Server.syncdb, in.Sender, in.Subchainid, hash); cnt > 1 {
			cnt--
			core.WriteSyncBlockCnt(Server.syncdb, in.Sender, in.Subchainid, hash, cnt)
			log.Debugf("Debase Sync block storeed count: %v", cnt)
		} else {
			log.Debugf("Delete Sync record!")
			core.DeleteSyncBlockHash(Server.syncdb, in.Sender, in.Subchainid, in.Blocknumber)
			core.DeleteSyncBlock(Server.syncdb, in.Sender, in.Subchainid, hash, in.Blocknumber)
		}
	}

	return &pb.DownloadBlockReply{Requestid: in.Requestid, Replybody: data}, nil
}
