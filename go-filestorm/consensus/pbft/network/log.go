// Copyright 2019 The go-filestorm Authors
// This file is part of the go-filestorm library.
//
// The go-filestorm library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-filestorm library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-filestorm library. If not, see <http://www.gnu.org/licenses/>.

// Package pbft implements the practical Byzantine fault tolerance consensus engine.
package pbft

import (
	"fmt"
	"github.com/filestorm/go-filestorm/pbft"
)

func LogMsg(msg interface{}) {
	switch msg.(type) {
	case *consensus.RequestMsg:
		reqMsg := msg.(*consensus.RequestMsg)
		fmt.Printf("[REQUEST] ClientID: %s, Timestamp: %d, Operation: %s\n", reqMsg.ClientID, reqMsg.Timestamp, reqMsg.Operation)
	case *consensus.PrePrepareMsg:
		prePrepareMsg := msg.(*consensus.PrePrepareMsg)
		fmt.Printf("[PREPREPARE] ClientID: %s, Operation: %s, SequenceID: %d\n", prePrepareMsg.RequestMsg.ClientID, prePrepareMsg.RequestMsg.Operation, prePrepareMsg.SequenceID)
	case *consensus.VoteMsg:
		voteMsg := msg.(*consensus.VoteMsg)
		if voteMsg.MsgType == consensus.PrepareMsg {
			fmt.Printf("[PREPARE] NodeID: %s\n", voteMsg.NodeID)
		} else if voteMsg.MsgType == consensus.CommitMsg {
			fmt.Printf("[COMMIT] NodeID: %s\n", voteMsg.NodeID)
		}
	}
}

func LogStage(stage string, isDone bool) {
	if isDone {
		fmt.Printf("[STAGE-DONE] %s\n", stage)
	} else {
		fmt.Printf("[STAGE-BEGIN] %s\n", stage)
	}
}