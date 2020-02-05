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

package fst

import (
	"github.com/filestorm/go-filestorm/core"
	"github.com/filestorm/go-filestorm/core/forkid"
	"github.com/filestorm/go-filestorm/p2p/enode"
	"github.com/filestorm/go-filestorm/rlp"
)

// ethEntry is the "fst" ENR entry which advertises fst protocol
// on the discovery network.
type ethEntry struct {
	ForkID forkid.ID // Fork identifier per EIP-2124

	// Ignore additional fields (for forward compatibility).
	Rest []rlp.RawValue `rlp:"tail"`
}

// ENRKey implements enr.Entry.
func (e ethEntry) ENRKey() string {
	return "fst"
}

func (fst *Filestorm) startEthEntryUpdate(ln *enode.LocalNode) {
	var newHead = make(chan core.ChainHeadEvent, 10)
	sub := fst.blockchain.SubscribeChainHeadEvent(newHead)

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case <-newHead:
				ln.Set(fst.currentEthEntry())
			case <-sub.Err():
				// Would be nice to sync with fst.Stop, but there is no
				// good way to do that.
				return
			}
		}
	}()
}

func (fst *Filestorm) currentEthEntry() *ethEntry {
	return &ethEntry{ForkID: forkid.NewID(fst.blockchain)}
}
