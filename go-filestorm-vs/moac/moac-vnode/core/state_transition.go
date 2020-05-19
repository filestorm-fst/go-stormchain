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

package core

import (
	"errors"
	"math/big"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/common/math"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/contracts"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/vm"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
)

var (
	Big0                         = big.NewInt(0)
	errInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")
	precompiledContracts         = contracts.GetInstance()
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==
  4a) Attempt to run transaction data
  4b) If valid, use result as code for the new state object
== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	gp           *GasPool
	msg          Message
	gasRemaining uint64
	gasPrice     *big.Int
	initialGas   *big.Int
	value        *big.Int
	data         []byte
	state        vm.StateDB
	evm          *vm.EVM
}

// Message represents a message sent to a contract.
type Message interface {
	From() common.Address
	//FromFrontier() (common.Address, error)
	To() *common.Address

	//Whether contract will flush automatically.
	AutoFlush() bool
	//Whether contract will return in the same block.
	// Removed after QueryFlag is not used any more 2018/
	// SyncFlag() bool
	//Wait for number of blocks. Will not work when SyncFlag is true.
	WaitBlockNumber() *big.Int

	//Identity consensus smart contract addr, removed
	Via() *common.Address
	ShardFlag() uint64

	GasPrice() *big.Int
	GasLimit() *big.Int
	Value() *big.Int

	Nonce() uint64
	CheckNonce() bool
	Data() []byte
	GetSystem() uint64
	GetMsgHash() *common.Hash
}

// IntrinsicGas computes the 'intrinsic gas' for a message
// with the given data.
//
// TODO convert to uint64
func IntrinsicGas(data []byte, contractCreation, pangu bool) *big.Int {
	igas := new(big.Int)
	if contractCreation && pangu {
		igas.SetUint64(params.TxGasContractCreation)
	} else {
		igas.SetUint64(params.TxGas)
	}
	if len(data) > 0 {
		var nz int64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		m := big.NewInt(nz)
		m.Mul(m, new(big.Int).SetUint64(params.TxDataNonZeroGas))
		igas.Add(igas, m)
		m.SetInt64(int64(len(data)) - nz)
		m.Mul(m, new(big.Int).SetUint64(params.TxDataZeroGas))
		igas.Add(igas, m)
	}
	return igas
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:         gp,
		evm:        evm,
		msg:        msg,
		gasPrice:   msg.GasPrice(),
		initialGas: new(big.Int),
		value:      msg.Value(),
		data:       msg.Data(),
		state:      evm.StateDB,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg Message, gp *GasPool) ([]byte, *big.Int, bool, error) {
	if msg.GetMsgHash() == nil {
		log.Debugf("[core/state_transition.go->ApplyMessage] msgHash %v", msg.GetMsgHash())
	} else {
		log.Debugf("[core/state_transition.go->ApplyMessage] msgHash %v", msg.GetMsgHash().Hex())
	}

	st := NewStateTransition(evm, msg, gp)

	ret, _, gasUsed, failed, err := st.TransitionDb()
	return ret, gasUsed, failed, err
}

func ApplyMessageForCalculateGas(evm *vm.EVM, msg Message, gp *GasPool) ([]byte, *big.Int, bool, error) {
	if msg.GetMsgHash() == nil {
		log.Debugf("[core/state_transition.go->ApplyMessage] msgHash %v", msg.GetMsgHash())
	} else {
		log.Debugf("[core/state_transition.go->ApplyMessage] msgHash %v", msg.GetMsgHash().Hex())
	}

	st := NewStateTransition(evm, msg, gp)

	ret, _, gasUsed, failed, err := st.TransitionDbForCalculateGas()
	return ret, gasUsed, failed, err
}

func (st *StateTransition) from() vm.AccountRef {
	f := st.msg.From()
	if !st.state.Exist(f) {
		st.state.CreateAccount(f)
	}
	return vm.AccountRef(f)
}

func (st *StateTransition) to() vm.AccountRef {
	if st.msg == nil {
		return vm.AccountRef{}
	}
	to := st.msg.To()
	if to == nil {
		return vm.AccountRef{} // contract creation
	}

	reference := vm.AccountRef(*to)
	if !st.state.Exist(*to) {
		st.state.CreateAccount(*to)
	}
	return reference
}

func (st *StateTransition) useGas(amount uint64) error {
	if st.gasRemaining < amount {
		return vm.ErrOutOfGas
	}
	st.gasRemaining -= amount

	return nil
}

func (st *StateTransition) buyGas() error {
	mgas := st.msg.GasLimit()
	if mgas.BitLen() > 64 {
		return vm.ErrOutOfGas
	}

	mgval := new(big.Int).Mul(mgas, st.gasPrice)

	var (
		state  = st.state
		sender = st.from()
	)
	if state.GetBalance(sender.Address()).Cmp(mgval) < 0 {
		return errInsufficientBalanceForGas
	}
	if err := st.gp.SubGas(mgas); err != nil {
		return err
	}
	st.gasRemaining += mgas.Uint64()

	st.initialGas.Set(mgas)
	state.SubBalance(sender.Address(), mgval)
	return nil
}

func (st *StateTransition) preCheck() error {
	msg := st.msg
	sender := st.from()

	// Make sure this transaction's nonce is correct
	if msg.CheckNonce() {
		nonce := st.state.GetNonce(sender.Address())
		if nonce < msg.Nonce() {
			return ErrNonceTooHigh
		} else if nonce > msg.Nonce() {
			return ErrNonceTooLow
		}
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and returning the result
// including the required gas for the operation as well as the used gas. It returns an error if it
// failed. An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, requiredGas, usedGas *big.Int, failed bool, err error) {
	log.Debugf("[core/state_transition.go->TransitionDb in]")
	if err = st.preCheck(); err != nil {
		return nil, nil, nil, false, err
	}
	msg := st.msg
	sender := st.from() // err checked in preCheck

	pangu := st.evm.ChainConfig().IsPangu(st.evm.BlockNumber)
	contractCreation := msg.To() == nil

	// Pay intrinsic gas
	// TODO convert to uint64
	intrinsicGas := IntrinsicGas(st.data, contractCreation, pangu)
	if intrinsicGas.BitLen() > 64 {
		return nil, nil, nil, false, vm.ErrOutOfGas
	}

	//if system msg, do not use gasRemaining
	if st.msg.GetSystem() == 0 {
		log.Debugf("[core/state_transition.go->TransitionDb] not system contract, useGas")
		if err = st.useGas(intrinsicGas.Uint64()); err != nil {
			return nil, nil, nil, false, err
		}
	}

	var (
		evm     = st.evm
		msgHash = msg.GetMsgHash()
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr        error
		contractAddr common.Address
	)
	if contractCreation {
		if msg.GetSystem() > 0 {
			addr := common.BytesToAddress([]byte{byte(msg.GetSystem())})
			ret, vmerr = evm.CreateSystemContract(addr, st.data, precompiledContracts, msgHash)
			contractAddr = addr
		} else {
			ret, contractAddr, st.gasRemaining, vmerr = evm.Create(sender, st.data, st.gasRemaining, st.value,
				msg.ShardFlag(), precompiledContracts, msgHash)
		}
		log.Debugf("created contract address:%v, leftOverGas:%v, codeHash:%v, ret:%v", contractAddr.String(), st.gasRemaining,
			evm.StateDB.GetCodeHash(contractAddr).String(), common.Bytes2Hex(ret))

	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(sender.Address(), st.state.GetNonce(sender.Address())+1)
		to := st.to().Address()
		data := st.data
		ret, st.gasRemaining, vmerr = evm.Call(sender, to, data, st.gasRemaining, st.value, msg.WaitBlockNumber().Cmp(big.NewInt(0)) == 0, msg.ShardFlag(), precompiledContracts, msgHash)
		if vmerr == nil && evm.Nr != nil && len(data) >= 4 {
			if contracts.IsWhiteListCall(to, data[:4]) {
				evm.Nr.UpdateWhiteState(0 /*evm.BlockNumber.Uint64()*/)
			}
		}
	}
	if vmerr != nil {
		log.Debug("VM returned with error", "err", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		// The above statement was wrong. If it's out of gas, the
		// entire transaction should revert.

		if vmerr == vm.ErrInsufficientBalance {
			return nil, nil, nil, false, vmerr
		}
	}
	requiredGas = new(big.Int).Set(st.gasUsed())

	st.refundGas()
	st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(st.gasUsed(), st.gasPrice))

	return ret, requiredGas, st.gasUsed(), vmerr != nil, err
}

func (st *StateTransition) TransitionDbForCalculateGas() (ret []byte, requiredGas, usedGas *big.Int, failed bool, err error) {
	log.Debugf("[core/state_transition.go->TransitionDb in]")
	if err = st.buyGas(); err != nil {
		return nil, nil, nil, false, err
	}
	msg := st.msg
	sender := st.from() // err checked in preCheck

	pangu := st.evm.ChainConfig().IsPangu(st.evm.BlockNumber)
	contractCreation := msg.To() == nil

	// Pay intrinsic gas
	// TODO convert to uint64
	intrinsicGas := IntrinsicGas(st.data, contractCreation, pangu)
	if intrinsicGas.BitLen() > 64 {
		return nil, nil, nil, false, vm.ErrOutOfGas
	}

	//if system msg, do not use gasRemaining
	if st.msg.GetSystem() == 0 {
		log.Debugf("[core/state_transition.go->TransitionDb] not system contract, useGas")
		if err = st.useGas(intrinsicGas.Uint64()); err != nil {
			return nil, nil, nil, false, err
		}
	}

	var (
		evm     = st.evm
		msgHash = msg.GetMsgHash()
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr        error
		contractAddr common.Address
	)
	if contractCreation {
		if msg.GetSystem() > 0 {
			addr := common.BytesToAddress([]byte{byte(msg.GetSystem())})
			ret, vmerr = evm.CreateSystemContract(addr, st.data, precompiledContracts, msgHash)
			contractAddr = addr
		} else {
			ret, contractAddr, st.gasRemaining, vmerr = evm.Create(sender, st.data, st.gasRemaining, st.value,
				msg.ShardFlag(), precompiledContracts, msgHash)
		}
		log.Debugf("created contract address:%v, codeHash:%v, ret:%v", contractAddr.String(),
			evm.StateDB.GetCodeHash(contractAddr).String(), common.Bytes2Hex(ret))

	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(sender.Address(), st.state.GetNonce(sender.Address())+1)
		ret, st.gasRemaining, vmerr = evm.Call(sender, st.to().Address(), st.data, st.gasRemaining, st.value, msg.WaitBlockNumber().Cmp(big.NewInt(0)) == 0, msg.ShardFlag(), precompiledContracts, msgHash)
	}
	if vmerr != nil {
		log.Debug("VM returned with error", "err", vmerr)
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		// The above statement was wrong. If it's out of gas, the
		// entire transaction should revert.

		if vmerr == vm.ErrInsufficientBalance {
			return nil, nil, nil, false, vmerr
		}
	}
	requiredGas = new(big.Int).Set(st.gasUsed())

	st.refundGas()
	st.state.AddBalance(st.evm.Coinbase, new(big.Int).Mul(st.gasUsed(), st.gasPrice))

	return ret, requiredGas, st.gasUsed(), vmerr != nil, err
}

func (st *StateTransition) refundGas() {
	// Return mc for remaining gas to the sender account,
	// exchanged at the original rate.
	sender := st.from() // err already checked
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gasRemaining), st.gasPrice)
	if st.initialGas.Cmp(new(big.Int).SetUint64(st.gasRemaining)) > 0 {
		st.state.AddBalance(sender.Address(), remaining)
	}
	// Apply refund counter, capped to half of the used gas.
	uhalf := remaining.Div(st.gasUsed(), common.Big2)
	refund := math.BigMin(uhalf, st.state.GetRefund())
	st.gasRemaining += refund.Uint64()

	st.state.AddBalance(sender.Address(), refund.Mul(refund, st.gasPrice))

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	if st.initialGas.Cmp(new(big.Int).SetUint64(st.gasRemaining)) > 0 {
		st.gp.AddGas(new(big.Int).SetUint64(st.gasRemaining))
	}
}

func (st *StateTransition) gasUsed() *big.Int {
	if st.initialGas.Cmp(new(big.Int).SetUint64(st.gasRemaining)) > 0 {
		return new(big.Int).Sub(st.initialGas, new(big.Int).SetUint64(st.gasRemaining))
	}
	return big.NewInt(0)
}
