package types

import (
	"math/big"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	pb "github.com/filestorm/go-filestorm/moac/moac-lib/proto"
)

type ShakeInfo struct {
	Pbhs       string
	Scsid      string
	Capability uint32
	Stream     *pb.Vnode_ScsPushServer
	ChainId    int64
}

func (s *ShakeInfo) GetScsid() string { return s.Scsid }

type AccountInfo struct {
	Addr                common.Address
	Balance             *big.Int
	Nonce               uint64
	CodeHash            common.Hash
	Query               uint64
	Shard               uint64
	CreationBlockNumber *big.Int
	WaitBlockNumber     *big.Int
}

type VnodeInfo struct {
}

type ContractInfo struct {
	Balance  *big.Int
	Nonce    uint64
	Root     common.Hash
	CodeHash []byte
	Code     []byte
	Storage  map[string]string
}

type LiveInfo struct {
	CurrentBlockNum *big.Int
}
