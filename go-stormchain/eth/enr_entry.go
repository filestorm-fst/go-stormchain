// Copyright 2019 The go-stormchain Authors
// This file is part of the go-stormchain library.
//
// The go-stormchain library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-stormchain library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-stormchain library. If not, see <http://www.gnu.org/licenses/>.

package storm

import (
	"github.com/filestorm-fst/go-stormchain/core"
	"github.com/filestorm-fst/go-stormchain/core/forkid"
	"github.com/filestorm-fst/go-stormchain/p2p/enode"
	"github.com/filestorm-fst/go-stormchain/rlp"
)

// ethEntry is the "storm" ENR entry which advertises storm protocol
// on the discovery network.
type ethEntry struct {
	ForkID forkid.ID // Fork identifier per EIP-2124

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e ethEntry) ENRKey() string {
	return "storm"
}

func (storm *StormChain) startEthEntryUpdate(ln *enode.LocalNode) {
	var newHead = make(chan core.ChainHeadEvent, 10)
	sub := storm.blockchain.SubscribeChainHeadEvent(newHead)

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-newHead:
				ln.Set(storm.currentEthEntry())
			case <-sub.Err():
				// Would be nice to sync with storm.Stop, but there is no
				// good way to do that.
				return
			}
		}
	}()
}

func (storm *StormChain) currentEthEntry() *ethEntry {
	return &ethEntry{ForkID: forkid.NewID(storm.blockchain)}
}
