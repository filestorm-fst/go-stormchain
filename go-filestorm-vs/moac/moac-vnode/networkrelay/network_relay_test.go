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
	"testing"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	libtypes "github.com/filestorm/go-filestorm/moac/moac-lib/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
	"golang.org/x/net/context"

	pb "github.com/filestorm/go-filestorm/moac/moac-lib/proto"
)

type testScsHandler struct {
	ScsPushMsgChan chan pb.ScsPushMsg
	ScsPushResChan chan pb.ScsPushMsg
	ScsRegChan     chan pb.ScsPushMsg
}

var (
	scsMsgChan = make(chan *pb.ScsPushMsg, 1000)
	scsResChan = make(chan *pb.ScsPushMsg)
	scsHandler = new(testScsHandler)
)

func TestRoleCache(t *testing.T) {

	networkRelay := NewNetworkRelay(scsMsgChan, scsResChan, scsHandler)
	scsId1 := common.HexToAddress("0x1000000000000000000000000000000000000001")
	scsId2 := common.HexToAddress("0x1000000000000000000000000000000000000002")
	subchain1 := common.HexToAddress("0x2000000000000000000000000000000000000001")
	subchain2 := common.HexToAddress("0x2000000000000000000000000000000000000002")
	networkRelay.UpdateScsRoleCache(scsId1, subchain1, 11)
	networkRelay.UpdateScsRoleCache(scsId1, subchain2, 12)
	networkRelay.UpdateScsRoleCache(scsId2, subchain1, 21)
	networkRelay.UpdateScsRoleCache(scsId2, subchain2, 22)
	fmt.Println(networkRelay.ScsRoleCache)
	r1 := networkRelay.GetScsRoleCache(scsId1, subchain1)
	r2 := networkRelay.GetScsRoleCache(scsId1, subchain2)
	r3 := networkRelay.GetScsRoleCache(scsId2, subchain1)
	r4 := networkRelay.GetScsRoleCache(scsId2, subchain2)
	fmt.Printf("r1:%v, r2:%v, r3:%v, r4:%v\n", r1, r2, r3, r4)
	networkRelay.DelScsRoleCache(scsId1, subchain2)
	r1 = networkRelay.GetScsRoleCache(scsId1, subchain1)
	r2 = networkRelay.GetScsRoleCache(scsId1, subchain2)
	r3 = networkRelay.GetScsRoleCache(scsId2, subchain1)
	r4 = networkRelay.GetScsRoleCache(scsId2, subchain2)
	fmt.Printf("r1:%v, r2:%v, r3:%v, r4:%v\n", r1, r2, r3, r4)
	ScsServerList := make(map[string]*libtypes.ScsServerConnection)

	ScsServerList[scsId1.String()] = &libtypes.ScsServerConnection{
		ScsId:      scsId1.String(),
		LiveFlag:   true,
		Stream:     nil,
		Req:        make(chan *pb.ScsPushMsg, 32768),
		Cancel:     make(chan bool),
		RetryCount: 0,
	}
	ScsServerList[scsId2.String()] = &libtypes.ScsServerConnection{
		ScsId:      scsId2.String(),
		LiveFlag:   true,
		Stream:     nil,
		Req:        make(chan *pb.ScsPushMsg, 32768),
		Cancel:     make(chan bool),
		RetryCount: 0,
	}
	networkRelay.ScsServerList = &ScsServerList
	networkRelay.SetScsRoleCacheOld(subchain1)
	r1 = networkRelay.GetScsRoleCache(scsId1, subchain1)
	r2 = networkRelay.GetScsRoleCache(scsId1, subchain2)
	r3 = networkRelay.GetScsRoleCache(scsId2, subchain1)
	r4 = networkRelay.GetScsRoleCache(scsId2, subchain2)
	fmt.Printf("r1:%v, r2:%v, r3:%v, r4:%v\n", r1, r2, r3, r4)
}

func (s *testScsHandler) ScsPush(stream pb.Vnode_ScsPushServer) error {
	return nil
}
func (s *testScsHandler) AccountInfo(ctx context.Context, in *pb.AccountInfoRequest) (*pb.AccountInfoReply, error) {
	return nil, nil
}
func (s *testScsHandler) ChainInfo(ctx context.Context, in *pb.ChainInfoRequest) (*pb.ChainInfoReply, error) {
	return nil, nil
}
func (s *testScsHandler) RemoteCall(ctx context.Context, in *pb.RemoteCallRequest) (*pb.RemoteCallReply, error) {
	return nil, nil
}
func (s *testScsHandler) GetBlockNumber() uint64 {
	return 0
}
func (s *testScsHandler) GetContractAddressMappingMember(addr common.Address, pos uint32, key common.Address) ([]common.Address, error) {
	return nil, nil
}
func (s *testScsHandler) NotifyMsgRunState(hash common.Hash) bool {
	return false
}
func (sh *testScsHandler) GetSCSRole(contractAddress common.Address, nodeAddress common.Address) params.ScsKind {
	return params.None
}
func (s *testScsHandler) GetScsPushMsgChan() chan pb.ScsPushMsg {
	return s.ScsPushMsgChan

}
func (s *testScsHandler) GetScsRegChan() chan pb.ScsPushMsg {
	return s.ScsPushResChan
}
func (s *testScsHandler) GetScsPushResChan() chan pb.ScsPushMsg {
	return s.ScsRegChan
}
