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

package core

import (
	"errors"
	"math/big"
	"sync"
	"time"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/state"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/event"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
)

const ()

var (
	ErrInvalidPeer = errors.New("invalid sharding peer")
)

var (
	shardingExpireInterval = 8 * time.Second
)

var ()

// ShardingMonitorConfig are the configuration parameters of the sharding monitor.
type ShardingMonitorConfig struct {
	Lifetime time.Duration // Maximum amount of time msg queued
}

var DefaultShardingMonitorConfig = ShardingMonitorConfig{
	Lifetime: 3 * time.Minute,
}

type ShardingMonitor struct {
	config       ShardingMonitorConfig
	chainconfig  *params.ChainConfig
	chain        blockChain
	gasPrice     *big.Int
	txFeed       event.Feed
	scope        event.SubscriptionScope
	chainHeadCh  chan ChainHeadEvent
	chainHeadSub event.Subscription
	signer       types.Signer
	mu           sync.RWMutex

	currentState  *state.StateDB      // Current state in the blockchain head
	pendingState  *state.ManagedState // Pending state tracking virtual nonces
	currentMaxGas *big.Int            // Current gasRemaining limit for transaction caps

	locals  *accountSet // Set of local transaction to exepmt from evicion rules
	journal *txJournal  // Journal of local transaction to back up to disk

	pending map[common.Address]*txList         // All currently processable transactions
	queue   map[common.Address]*txList         // Queued but non-processable transactions
	beats   map[common.Address]time.Time       // Last heartbeat from each known account
	all     map[common.Hash]*types.Transaction // All transactions to allow lookups
	priced  *txPricedList                      // All transactions sorted by price

	wg sync.WaitGroup // for shutdown sync

	pangu bool
}

func NewShardingMonitor(config ShardingMonitorConfig, chainconfig *params.ChainConfig, chain blockChain) *ShardingMonitor {

	// Create the transaction pool with its initial settings
	monitor := &ShardingMonitor{
		config:      config,
		chainconfig: chainconfig,
		chain:       chain,
		signer:      types.NewPanguSigner(chainconfig.ChainId),
		pending:     make(map[common.Address]*txList),
		queue:       make(map[common.Address]*txList),
		beats:       make(map[common.Address]time.Time),
		all:         make(map[common.Hash]*types.Transaction),
		chainHeadCh: make(chan ChainHeadEvent, chainHeadChanSize),
		//gasPrice:    new(big.Int).SetUint64(config.PriceLimit),
	}

	monitor.locals = newAccountSet(monitor.signer)
	monitor.priced = newTxPricedList(&monitor.all)
	monitor.reset(nil, chain.CurrentBlock().Header())

	// Start the event loop and return
	monitor.wg.Add(1)
	go monitor.loop()

	return monitor
}

func (monitor *ShardingMonitor) loop() {
	defer monitor.wg.Done()

	// Start the stats reporting and transaction eviction tickers
	var prevPending, prevQueued, prevStales int

	report := time.NewTicker(statsReportInterval)
	defer report.Stop()

	evict := time.NewTicker(evictionInterval)
	defer evict.Stop()

	//journal := time.NewTicker(monitor.config.Rejournal)
	//defer journal.Stop()

	// Track the previous head headers for transaction reorgs
	head := monitor.chain.CurrentBlock()

	// Keep waiting for and reacting to the various events
	for {
		select {
		// Handle ChainHeadEvent
		case ev := <-monitor.chainHeadCh:
			if ev.Block != nil {
				monitor.mu.Lock()
				if monitor.chainconfig.IsPangu(ev.Block.Number()) {
					monitor.pangu = true
				}
				monitor.reset(head.Header(), ev.Block.Header())
				head = ev.Block

				monitor.mu.Unlock()
			}
		// Be unsubscribed due to system stopped
		case <-monitor.chainHeadSub.Err():
			return

		// Handle stats reporting ticks
		case <-report.C:
			monitor.mu.RLock()
			pending, queued := monitor.stats()
			stales := monitor.priced.stales
			//log.Info("Txmonitor.loop", "report", pending)
			monitor.mu.RUnlock()

			if pending != prevPending || queued != prevQueued || stales != prevStales {
				log.Debug("Transaction pool status report", "executable", pending, "queued", queued, "stales", stales)
				prevPending, prevQueued, prevStales = pending, queued, stales
			}

		}
	}
}

// reset retrieves the current state of the blockchain and ensures the content
// of the monitor is valid with regard to the chain state.
func (monitor *ShardingMonitor) reset(oldHead, newHead *types.Header) {

}

// Stop
func (monitor *ShardingMonitor) Stop() {
	// Unsubscribe all subscriptions registered from txpool
	monitor.scope.Close()

	// Unsubscribe subscriptions registered from blockchain
	monitor.chainHeadSub.Unsubscribe()
	monitor.wg.Wait()

	if monitor.journal != nil {
		monitor.journal.close()
	}
	log.Info("monitor stopped")
}

// State returns the virtual managed state of the monitor
func (monitor *ShardingMonitor) State() *state.ManagedState {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()

	return monitor.pendingState
}

// Stats retrieves the current monitor stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (monitor *ShardingMonitor) Stats() (int, int) {
	monitor.mu.RLock()
	defer monitor.mu.RUnlock()

	return monitor.stats()
}

// stats retrieves the current pool stats, namely the number of pending and the
// number of queued (non-executable) transactions.
func (monitor *ShardingMonitor) stats() (int, int) {
	pending := 0
	for _, list := range monitor.pending {
		pending += list.Len()
	}
	queued := 0
	for _, list := range monitor.queue {
		queued += list.Len()
	}
	return pending, queued
}
