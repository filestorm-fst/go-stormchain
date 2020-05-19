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

package params

import "time"

// These are network parameters that need to be constant between clients, but
// aren't necesarilly consensus related.

const (
	// BloomBitsBlocks is the number of blocks a single bloom bit section vector
	// contains.
	BloomBitsBlocks             uint64        = 4096
	TimerPingInterval           time.Duration = 1
	DirectCallLimitPerBlock                   = 2048
	DirectCallGasLimit                        = 4000000
	SubchainMsgLimit                          = 1000
	ScsMsgLimit                               = 1000
	SubnetListenPortMin                       = 40333 // port range min for subnet p2p listen addr
	SubnetListenPortMax                       = 40999 // port range max for subnet p2p listen addr
	SubnetP2PConnectionFraction               = 3     // fraction of subnet tcp connections to the mainnet
	SubnetP2PConnectionMin                    = 5     // min subnet p2p tpc connections
	SubnetBootNodeLimits                      = 5     // max number of bootnodes return by a findvalue query for subnet

	// scspushmsg type
	DirectCall   = 1
	BroadCast    = 2
	ControlMsg   = 3
	ScsShakeHand = 4
	ScsPing      = 5

	// None            = -1
	RegOpen            = 0
	RegClose           = 1
	CreateProposal     = 2
	DisputeProposal    = 3
	ApproveProposal    = 4
	RegAdd             = 5
	RegAsMonitor       = 6
	RegAsBackup        = 7
	UpdateLastFlushBlk = 8
	DistributeProposal = 9
	ResetAll           = 10
	UploadRedeemData   = 11
	EnterAndRedeem     = 12
	RequestRelease     = 13

	// scspushmsg status
	NewBlock     = 0
	SyncRequest  = 1
	SyncComplete = 2
)

type ScsKind int64

const (
	None ScsKind = iota
	ConsensusScs
	MonitorScs
	BackupScs
	MatchSelTarget
	LockScs
)
