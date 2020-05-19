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

package types

import (
	"fmt"
	"math/big"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/common/hexutil"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
)

//go:generate gencodec -type txdata -field-override txdataMarshaling -out gen_tx_json.go

var ()

type ShardingMessage struct {
	ContractAddress *common.Address `json:"contractaddr"      gencodec:"required"`
	MessageType     uint64          `json:"MessageType" 	gencodec:"required"`
	BlockStart      *big.Int        `json:"BlockStart"    gencodec:"required"`
	BlockExpire     *big.Int        `json:"BlockExpire"    gencodec:"required"`
	Payload         []byte          `json:"input"    gencodec:"required"`
}

type ShardingMessageMarshaling struct {
	ContractAddress hexutil.Bytes
	MessageType     hexutil.Uint64
	BlockStart      *hexutil.Big
	BlockExpire     *hexutil.Big
	Payload         hexutil.Bytes
}

func NewShardingMessage(cnt *common.Address, msgtype uint64, blockstart, blockexpire *big.Int, data []byte) *ShardingMessage {
	if len(data) > 0 {
		data = common.CopyBytes(data)
	}
	d := ShardingMessage{
		ContractAddress: cnt,
		MessageType:     msgtype,
		BlockStart:      blockstart,
		BlockExpire:     blockexpire,
		Payload:         data,
	}
	return &d
}

func (msg *ShardingMessage) String() string {

	enc, _ := rlp.EncodeToBytes(msg)
	return fmt.Sprintf(`
	ShardingMessage
	Contract: %v
	type:     %#x
	BlockStart:       %#x
	BlockExpire:    %#x
	Data:     0x%x
	enc: %x
`,
		msg.ContractAddress,
		msg.MessageType,
		msg.BlockStart,
		msg.BlockExpire,
		msg.Payload,
		enc,
	)
}
