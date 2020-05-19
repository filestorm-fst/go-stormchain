// Copyright 2015 The MOAC-core Authors
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

package mc

import (
	"context"
	"math/big"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/common/math"
	"github.com/filestorm/go-filestorm/moac/moac-lib/mcdb"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/bloombits"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/state"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/vm"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/event"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc/downloader"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc/gasprice"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/rpc"
)

// MoacApiBackend implements mcapi.Backend for full nodes
type MoacApiBackend struct {
	mc  *MoacService
	gpo *gasprice.Oracle
}

func (b *MoacApiBackend) ChainConfig() *params.ChainConfig {
	return b.mc.chainConfig
}

func (b *MoacApiBackend) CurrentBlock() *types.Block {
	return b.mc.blockchain.CurrentBlock()
}

func (b *MoacApiBackend) SetHead(number uint64) {
	b.mc.ProtocolManager.downloader.Cancel()
	b.mc.blockchain.SetHead(number)
}

func (b *MoacApiBackend) HeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Header, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.mc.miner.PendingBlock()
		return block.Header(), nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.mc.blockchain.CurrentBlock().Header(), nil
	}
	return b.mc.blockchain.GetHeaderByNumber(uint64(blockNr)), nil
}

func (b *MoacApiBackend) BlockByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block := b.mc.miner.PendingBlock()
		return block, nil
	}
	// Otherwise resolve and return the block
	if blockNr == rpc.LatestBlockNumber {
		return b.mc.blockchain.CurrentBlock(), nil
	}
	return b.mc.blockchain.GetBlockByNumber(uint64(blockNr)), nil
}

func (b *MoacApiBackend) StateAndHeaderByNumber(ctx context.Context, blockNr rpc.BlockNumber) (*state.StateDB, *types.Header, error) {
	// Pending state is only known by the miner
	if blockNr == rpc.PendingBlockNumber {
		block, state := b.mc.miner.Pending()
		return state, block.Header(), nil
	}
	// Otherwise resolve the block number and return its state
	header, err := b.HeaderByNumber(ctx, blockNr)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := b.mc.BlockChain().StateAt(header.Root)
	return stateDb, header, err
}

func (b *MoacApiBackend) GetBlock(ctx context.Context, blockHash common.Hash) (*types.Block, error) {
	return b.mc.blockchain.GetBlockByHash(blockHash), nil
}

func (b *MoacApiBackend) GetReceipts(ctx context.Context, blockHash common.Hash) (types.Receipts, error) {
	return core.GetBlockReceipts(b.mc.chainDb, blockHash, core.GetBlockNumber(b.mc.chainDb, blockHash)), nil
}

func (b *MoacApiBackend) GetTd(blockHash common.Hash) *big.Int {
	return b.mc.blockchain.GetTdByHash(blockHash)
}

func (b *MoacApiBackend) GetEVM(ctx context.Context, msg core.Message, state *state.StateDB, header *types.Header, vmCfg vm.Config) (*vm.EVM, func() error, error) {
	state.SetBalance(msg.From(), math.MaxBig256)
	vmError := func() error { return nil }

	context := core.NewEVMContext(msg, header, b.mc.BlockChain(), nil, nil)
	//nr := b.mc.ProtocolManager.NetworkRelay
	return vm.NewEVM(context, state, b.mc.chainConfig, vmCfg, /*nr*/nil), vmError, nil
}

func (b *MoacApiBackend) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return b.mc.BlockChain().SubscribeRemovedLogsEvent(ch)
}

func (b *MoacApiBackend) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return b.mc.BlockChain().SubscribeChainEvent(ch)
}

func (b *MoacApiBackend) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return b.mc.BlockChain().SubscribeChainHeadEvent(ch)
}

func (b *MoacApiBackend) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return b.mc.BlockChain().SubscribeChainSideEvent(ch)
}

func (b *MoacApiBackend) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return b.mc.BlockChain().SubscribeLogsEvent(ch)
}

func (b *MoacApiBackend) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	// log.Info("[mc/api_backend.go->MoacApiBackend.SendTx]")
	return b.mc.txPool.AddLocal(signedTx)
}

func (b *MoacApiBackend) GetPoolTransactions() (types.Transactions, error) {
	pending, err := b.mc.txPool.Pending()
	if err != nil {
		return nil, err
	}
	var txs types.Transactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (b *MoacApiBackend) GetPoolTransaction(hash common.Hash) *types.Transaction {
	return b.mc.txPool.Get(hash)
}

func (b *MoacApiBackend) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return b.mc.txPool.State().GetNonce(addr), nil
}

func (b *MoacApiBackend) Stats() (pending int, queued int) {
	return b.mc.txPool.Stats()
}

func (b *MoacApiBackend) TxPoolContent() (map[common.Address]types.Transactions, map[common.Address]types.Transactions) {
	return b.mc.TxPool().Content()
}

func (b *MoacApiBackend) SubscribeTxPreEvent(ch chan<- core.TxPreEvent) event.Subscription {
	return b.mc.TxPool().SubscribeTxPreEvent(ch)
}

func (b *MoacApiBackend) Downloader() *downloader.Downloader {
	return b.mc.Downloader()
}

func (b *MoacApiBackend) ProtocolVersion() int {
	return b.mc.McVersion()
}

func (b *MoacApiBackend) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return b.gpo.SuggestPrice(ctx)
}

func (b *MoacApiBackend) ChainDb() mcdb.Database {
	return b.mc.ChainDb()
}

func (b *MoacApiBackend) EventMux() *event.TypeMux {
	return b.mc.EventMux()
}

func (b *MoacApiBackend) AccountManager() *accounts.Manager {
	return b.mc.AccountManager()
}

func (b *MoacApiBackend) BloomStatus() (uint64, uint64) {
	sections, _, _ := b.mc.bloomIndexer.Sections()
	return params.BloomBitsBlocks, sections
}

func (b *MoacApiBackend) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	for i := 0; i < bloomFilterThreads; i++ {
		go session.Multiplex(bloomRetrievalBatch, bloomRetrievalWait, b.mc.bloomRequests)
	}
}
