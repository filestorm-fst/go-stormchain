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

package vm

import (
	"math/big"
	"sync"
	"sync/atomic"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/crypto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"

	"time"

	pb "github.com/filestorm/go-filestorm/moac/moac-lib/proto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts"
)

var (
	// emptyCodeHash is used by create to ensure deployment is disallowed to already
	// deployed contract addresses (relevant after the account abstraction).
	emptyCodeHash = crypto.Keccak256Hash(nil)
	requestId     = uint32(0)
	// the latest EVM instance.
	evmCache *EVM
	evmMu    sync.RWMutex
)

const (
	DirectCall   = 1
	BroadCast    = 2
	ControlMsg   = 3
	ScsShakeHand = 4
	ScsPing      = 5

	None            = -1
	RegOpen         = 0
	RegClose        = 1
	CreateProposal  = 2
	DisputeProposal = 3
	ApproveProposal = 4
)

type ContractsInterface interface {
	PrecompiledContractsPangu() map[common.Address]PrecompiledContract
	PrecompiledContractsByzantium() map[common.Address]PrecompiledContract
	RunPrecompiledContract(evm *EVM, snapshot int, p PrecompiledContract, input []byte, contract *Contract, hash *common.Hash) (ret []byte, err error)
	SystemContractCallAddr() common.Address
	SystemContractEntryAddr(num *big.Int) common.Address
	IsSystemCaller(caller ContractRef) (ret bool)
	WhiteListCallAddr() common.Address
}

type PrecompiledContract interface {
	RequiredGas(input []byte) uint64                                                                 // RequiredPrice calculates the contract gas use
	Run(evm *EVM, snapshot int, contract *Contract, input []byte, hash *common.Hash) ([]byte, error) // Run runs the precompiled contract
}

type NetworkRelayInterface interface {
	VnodePushMsg(conReq *pb.ScsPushMsg) (map[int]*pb.ScsPushMsg, error)
	NotifyScs(address common.Address, msg []byte, hash common.Hash, block *big.Int)
	UpdateWhiteState(block uint64)
}

type (
	CanTransferFunc func(StateDB, common.Address, *big.Int) bool
	TransferFunc    func(StateDB, common.Address, common.Address, *big.Int)
	// GetHashFunc returns the nth block hash in the blockchain
	// and is used by the BLOCKHASH EVM op code.
	GetHashFunc func(uint64) common.Hash
)

// Run runs the given contract and takes care of running precompiles with a fallback to the byte code interpreter.
func Run(evm *EVM, snapshot int, contract *Contract, input []byte, precompiledContracts ContractsInterface, msgHash *common.Hash) ([]byte, error) {
	//log.Debugf("Run contract %v, input %v", contract, common.Bytes2Hex(input))
	if msgHash != nil {
		log.Debugf("[core/vm/evm.go->run] msgHash %v ", msgHash.Hex())
	} else {
		log.Debugf("[core/vm/evm.go->run] msgHash %v", msgHash)
	}
	if contract.CodeAddr != nil {
		precompiles := precompiledContracts.PrecompiledContractsPangu()
		// if evm.ChainConfig().IsByzantium(evm.BlockNumber) {
		// 	precompiles = precompiledContracts.PrecompiledContractsByzantium()
		// }
		if p := precompiles[*contract.CodeAddr]; p != nil {
			return precompiledContracts.RunPrecompiledContract(evm, snapshot, p, input, contract, msgHash)
		}
	}
	ret, err := evm.interpreter.Run(snapshot, contract, input, precompiledContracts, msgHash)
	//log.Infof("Run %v %v", ret, err)
	return ret, err
}

type CTXPool interface {
	AddShardJoinCall(tx *types.Transaction) error
}

// Context provides the EVM with auxiliary information. Once provided
// it shouldn't be modified.
type Context struct {
	// CanTransfer returns whether the account contains
	// sufficient ether to transfer the value
	CanTransfer CanTransferFunc
	// Transfer transfers ether from one account to the other
	Transfer TransferFunc
	// GetHash returns the hash corresponding to n
	GetHash GetHashFunc

	// Message information
	Origin   common.Address // Provides information for ORIGIN
	GasPrice *big.Int       // Provides information for GASPRICE

	// Block information
	Coinbase    common.Address // Provides information for COINBASE
	GasLimit    *big.Int       // Provides information for GASLIMIT
	BlockNumber *big.Int       // Provides information for NUMBER
	Time        *big.Int       // Provides information for TIME
	Difficulty  *big.Int       // Provides information for DIFFICULTY
	CtxTxpool   CTXPool
}

// EVM is the MoacNode Virtual Machine base object and provides
// the necessary tools to run a contract on the given state with
// the provided context. It should be noted that any error
// generated through any of the calls should be considered a
// revert-state-and-consume-all-gas operation, no checks on
// specific errors should ever be performed. The interpreter makes
// sure that any errors generated are to be considered faulty code.
//
// The EVM should never be reused and is not thread safe.
type EVM struct {
	// Context provides auxiliary blockchain related information
	Context
	// StateDB gives access to the underlying state
	StateDB StateDB
	// Depth is the current call stack
	depth int

	// chainConfig contains information about the current chain
	chainConfig *params.ChainConfig
	// chain rules contains the chain rules for the current epoch
	chainRules params.Rules
	// virtual machine configuration options used to initialise the
	// evm.
	VmConfig Config
	// global (to this context) moac virtual machine
	// used throughout the execution of the tx.
	interpreter *Interpreter
	// abort is used to abort the EVM calling operations
	// NOTE: must be set atomically
	abort int32
	Nr    NetworkRelayInterface
}

type ContractCallStruct struct {
	Sender       []byte
	Contractaddr []byte
	Input        []byte
	Gas          uint64
	Val          []byte
	Sync         bool
	Code         []byte
	Codehash     []byte
}

type QueryResult struct {
	address common.Address
}

// NewEVM retutrns a new EVM . The returned EVM is not thread safe and should
// only ever be used *once*.
func NewEVM(ctx Context, statedb StateDB, chainConfig *params.ChainConfig, vmConfig Config, nr NetworkRelayInterface) *EVM {
	//log.Info("NewEVM()")
	evm := &EVM{
		Context:     ctx,
		StateDB:     statedb,
		VmConfig:    vmConfig,
		chainConfig: chainConfig,
		chainRules:  chainConfig.Rules(ctx.BlockNumber),
		Nr:          nr,
	}

	evm.interpreter = NewInterpreter(evm, vmConfig)
	return evm
}

// SetEVM the latest EVM instance.
func SetEVM(evm *EVM) {
	evmMu.RLock()
	defer evmMu.RUnlock()

	evmCache = evm
}

// GetEVM the latest EVM instance.
func GetEVM() *EVM {
	if evmCache == nil {
		return nil
	}

	evmMu.RLock()
	defer evmMu.RUnlock()

	evmCache.StateDB = evmCache.StateDB.Copy2().(StateDB)

	return evmCache
}

// Cancel cancels any running EVM operation. This may be called concurrently and
// it's safe to be called multiple times.
func (evm *EVM) Cancel() {
	atomic.StoreInt32(&evm.abort, 1)
}

// Call executes the contract associated with the addr with the given input as
// parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
func (evm *EVM) Call(
	caller ContractRef,
	addr common.Address,
	input []byte,
	gasLimit uint64,
	value *big.Int,
	synccall bool,
	shardingFlag uint64,
	precompiledContracts ContractsInterface,
	msgHash *common.Hash) (
	ret []byte,
	leftOverGas uint64,
	err error) {
	log.Debugf("[core/vm/evm.go->Call] caller=%v addr=%v input=%v code=%v &evm.StateDB=%v, gasLimit=%v", caller.Address().Hex(), addr.String(), common.Bytes2Hex(input), len(evm.StateDB.GetCode(addr)), &evm.StateDB, gasLimit)

	leftOverGas = gasLimit
	if evm.VmConfig.NoRecursion && evm.depth > 0 {
		// log.Info("here 1")
		return nil, leftOverGas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		// log.Info("here 2")
		return nil, leftOverGas, ErrDepth
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)
	// replace with entry address
	curcaller := caller
	if caller.Address() == precompiledContracts.SystemContractEntryAddr(big.NewInt(0)) && addr == common.BytesToAddress([]byte{12}) {
		//get address from last of value
		if len(input) >= 100 {
			curcaller = AccountRef(common.BytesToAddress(input[16:36]))
			to = AccountRef(common.BytesToAddress(input[48:68]))
			value.SetBytes(input[80:])
			//gas allos
			log.Info("delegate tx", "from", curcaller, "to", to, "value", value)
		}
	}
	// log.Info("here 3")

	// Fail if we're trying to transfer more than the available balance
	if !evm.Context.CanTransfer(evm.StateDB, curcaller.Address(), value) {
		// log.Info("here 4")
		return nil, leftOverGas, ErrInsufficientBalance
	}

	if !evm.StateDB.Exist(to.Address()) {
		// log.Info("here 5")
		//log.Info("Cannot find addr:" + to.Address().String())
		precompiles := precompiledContracts.PrecompiledContractsPangu()
		// if evm.ChainConfig().IsByzantium(evm.BlockNumber) {
		// 	precompiles = precompiledContracts.PrecompiledContractsByzantium()
		// }
		// if precompiles[to.Address()] == nil && evm.ChainConfig().IsEIP158(evm.BlockNumber) && value.Sign() == 0 {
		// 	return nil, gas, nil
		// }
		if precompiles[to.Address()] == nil && value.Sign() == 0 {
			// Calling a non existing account, don't do anything, but ping the tracer
			if evm.VmConfig.Debug && evm.depth == 0 {
				evm.VmConfig.Tracer.CaptureStart(curcaller.Address(), addr, false, input, leftOverGas, value)
				evm.VmConfig.Tracer.CaptureEnd(ret, 0, 0, nil)
			}
			return nil, leftOverGas, nil
		}
		evm.StateDB.CreateAccount(to.Address())
	} else {
		//log.Info("Found addr:" + to.Address().String())
	}
	evm.Transfer(evm.StateDB, curcaller.Address(), to.Address(), value)
	// log.Info("here 6")

	// initialise a new contract and set the code that is to be used by the
	// E The contract is a scoped environment for this execution context
	// only.
	contract := NewContract(curcaller, to, value, leftOverGas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	//skip delegate call for contract call for now
	if evm.StateDB.GetCodeSize(to.Address()) <= 0 || addr == common.BytesToAddress([]byte{12}) {
		// log.Info("only transfer transaction:" + to.Address().String())
		return nil, leftOverGas, nil
	}

	ret, err = Run(evm, snapshot, contract, input, precompiledContracts, msgHash)

	// Capture the tracer start/end events in debug mode
	if evm.VmConfig.Debug && evm.depth == 0 {
		// Even if the account has no code, we need to continue because it might be a precompile
		start := time.Now()
		evm.VmConfig.Tracer.CaptureStart(caller.Address(), addr, false, input, leftOverGas, value)

		defer func() { // Lazy evaluation of the parameters
			evm.VmConfig.Tracer.CaptureEnd(ret, leftOverGas-contract.GasRemaining, time.Since(start), err)
		}()
	}

	leftOverGas = contract.GasRemaining

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in pangu this also counts for code storage gas errors.
	if err != nil {
		log.Debug("error running evm", "err", err)
		evm.StateDB.RevertToSnapshot(snapshot)
		// if err != errExecutionReverted {
		// 	contract.UseGas(contract.GasRemaining)
		// }
	}
	log.Debugf("input %v", common.Bytes2Hex(input))
	log.Debugf("Call returning ret %v leftOverGas %v err %v", common.Bytes2Hex(ret), leftOverGas, err)
	return ret, leftOverGas, err
}

func (evm *EVM) decodeRespond(resp map[int]*pb.ScsPushMsg) (ret []byte, leftOverGas uint64, err error) {
	//TODO write real code to decode respond
	ret = []byte("test")
	leftOverGas = 0
	err = nil
	return ret, leftOverGas, err
}

func (evm *EVM) decodeQueryResult(result []byte) (queryResult QueryResult, err error) {
	//TODO put the real decode in it.
	ret := QueryResult{}
	return ret, nil
}

// CallCode executes the contract associated with the addr with the given input
// as parameters. It also handles any necessary value transfer required and takes
// the necessary steps to create accounts and reverses the state in case of an
// execution error or failed value transfer.
//
// CallCode differs from Call in the sense that it executes the given address'
// code with the caller as context.
func (evm *EVM) CallCode(caller ContractRef, addr common.Address, input []byte, gas uint64, value *big.Int, precompiledContracts ContractsInterface, msgHash *common.Hash) (ret []byte, leftOverGas uint64, err error) {
	log.Debug("evm.CallCode")
	if evm.VmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}

	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Fail if we're trying to transfer more than the available balance
	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, gas, ErrInsufficientBalance
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)
	// initialise a new contract and set the code that is to be used by the
	// E The contract is a scoped evmironment for this execution context
	// only.
	contract := NewContract(caller, to, value, gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = Run(evm, snapshot, contract, input, precompiledContracts, msgHash)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.GasRemaining)
		}
	}
	return ret, contract.GasRemaining, err
}

// DelegateCall executes the contract associated with the addr with the given input
// as parameters. It reverses the state in case of an execution error.
//
// DelegateCall differs from CallCode in the sense that it executes the given address'
// code with the caller as context and the caller is set to the caller of the caller.
func (evm *EVM) DelegateCall(caller ContractRef, addr common.Address, input []byte, gas uint64, precompiledContracts ContractsInterface, msgHash *common.Hash) (ret []byte, leftOverGas uint64, err error) {
	if evm.VmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}

	var (
		snapshot = evm.StateDB.Snapshot()
		to       = AccountRef(caller.Address())
	)

	// Initialise a new contract and make initialise the delegate values
	contract := NewContract(caller, to, nil, gas)
	if precompiledContracts.IsSystemCaller(caller) {
		contract.AsDelegateWithData(input)
	} else {
		contract.AsDelegate()
	}
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	ret, err = Run(evm, snapshot, contract, input, precompiledContracts, msgHash)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.GasRemaining)
		}
	}
	return ret, contract.GasRemaining, err
}

// StaticCall executes the contract associated with the addr with the given input
// as parameters while disallowing any modifications to the state during the call.
// Opcodes that attempt to perform such modifications will result in exceptions
// instead of performing the modifications.
func (evm *EVM) StaticCall(caller ContractRef, addr common.Address, input []byte, gas uint64, precompiledContracts ContractsInterface, msgHash *common.Hash) (ret []byte, leftOverGas uint64, err error) {
	if evm.VmConfig.NoRecursion && evm.depth > 0 {
		return nil, gas, nil
	}
	// Fail if we're trying to execute above the call depth limit
	if evm.depth > int(params.CallCreateDepth) {
		return nil, gas, ErrDepth
	}
	// Make sure the readonly is only set if we aren't in readonly yet
	// this makes also sure that the readonly flag isn't removed for
	// child calls.
	if !evm.interpreter.readOnly {
		evm.interpreter.readOnly = true
		defer func() { evm.interpreter.readOnly = false }()
	}

	var (
		to       = AccountRef(addr)
		snapshot = evm.StateDB.Snapshot()
	)
	// Initialise a new contract and set the code that is to be used by the
	// EVM. The contract is a scoped environment for this execution context
	// only.
	contract := NewContract(caller, to, new(big.Int), gas)
	contract.SetCallCode(&addr, evm.StateDB.GetCodeHash(addr), evm.StateDB.GetCode(addr))

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in Pangu this also counts for code storage gas errors.
	ret, err = Run(evm, snapshot, contract, input, precompiledContracts, msgHash)
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			contract.UseGas(contract.GasRemaining)
		}
	}
	return ret, contract.GasRemaining, err
}

// Create creates a new contract using code as deployment code.
func (evm *EVM) Create(caller ContractRef, code []byte, gasRemaining uint64, value *big.Int,
	shardflag uint64, precompiledContracts ContractsInterface, msgHash *common.Hash) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	log.Debugf("[core/vm/evm.go->Create] caller=%v len(code)=%v gasRemaining=%v value=%v &StateDB=%v", caller.Address(), len(code), gasRemaining, value, &evm.StateDB)
	log.Debugf("[core/vm/evm.go->Create] shardflag=%v", shardflag)

	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if evm.depth > int(params.CallCreateDepth) {
		return nil, common.Address{}, gasRemaining, ErrDepth
	}
	if !evm.CanTransfer(evm.StateDB, caller.Address(), value) {
		return nil, common.Address{}, gasRemaining, ErrInsufficientBalance
	}
	// Ensure there's no existing contract already at the designated address
	nonce := evm.StateDB.GetNonce(caller.Address())
	evm.StateDB.SetNonce(caller.Address(), nonce+1)

	contractAddr = crypto.CreateAddress(caller.Address(), nonce)
	contractHash := evm.StateDB.GetCodeHash(contractAddr)
	if evm.StateDB.GetNonce(contractAddr) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, common.Address{}, 0, ErrContractAddressCollision
	}

	/* 	if len(code) <= 0 {
	   		return nil, common.Address{}, gasRemaining, ErrEmptyCode
	   	}
	*/
	// Create a new account on the state
	snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(contractAddr)

	// MOAC, take EIP158 as default, remove empty account
	// if evm.ChainConfig().IsEIP158(evm.BlockNumber) {
	evm.StateDB.SetNonce(contractAddr, 1)
	// }

	evm.Transfer(evm.StateDB, caller.Address(), contractAddr, value)

	manager := accounts.GetManager()
	manager.AddContract(contractAddr)

	if evm.VmConfig.NoRecursion && evm.depth > 0 {
		return nil, contractAddr, gasRemaining, nil
	}

	if evm.VmConfig.Debug && evm.depth == 0 {
		evm.VmConfig.Tracer.CaptureStart(caller.Address(), contractAddr, true, code, gasRemaining, value)
	}
	start := time.Now()

	log.Debugf("create contract address %v, nonce %v", contractAddr.String(), nonce)
	contract := NewContract(caller, AccountRef(contractAddr), value, gasRemaining)
	codeHash := crypto.Keccak256Hash(code)
	contract.SetCallCode(&contractAddr, codeHash, code)
	//cGasRemaining = contract.GasRemaining
	log.Debugf("contract code hash %v", codeHash.Hex())
	log.Debugf("contract gasRemaining %v", contract.GasRemaining)
	ret, err = Run(evm, snapshot, contract, nil, precompiledContracts, msgHash)

	log.Debugf("after contract creating Run retLen: %d err: %v", len(ret), err)
	// check whether the max code size has been exceeded
	// MOAC set EIP158 as default
	// maxCodeSizeExceeded := evm.ChainConfig().IsEIP158(evm.BlockNumber) && len(ret) > params.MaxCodeSize
	maxCodeSizeExceeded := len(ret) > params.MaxCodeSize

	if evm.ChainConfig().IsNuwa(evm.BlockNumber) {
		maxCodeSizeExceeded = len(ret) > params.NuwaMaxCodeSize
	}

	// lg.Print("maxCodeSizeExceeded: %v", maxCodeSizeExceeded)
	// lg.Print("err == nil: %v", err == nil)
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		gasRequiredToCreateData := uint64(len(ret)) * params.CreateDataGas
		log.Debugf("contract.GasRemaining: %v gasRequiredToCreateData: %v", contract.GasRemaining, gasRequiredToCreateData)
		if contract.GasRemaining >= gasRequiredToCreateData && contract.UseGas(gasRequiredToCreateData) {
			log.Debugf("created code")
			evm.StateDB.SetCode(contractAddr, ret)
		} else {
			contract.UseGas(contract.GasRemaining)
			log.Debugf("did not create code, out of gas. ")
			//evm.StateDB.SetCode(contractAddr, ret) //TODO delete it
			err = ErrCodeStoreOutOfGas
		}
	} else {
		log.Debugf("maxCodeSizeExceeded %v, error %v", maxCodeSizeExceeded, err)
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in pangu this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (evm.ChainConfig().IsPangu(evm.BlockNumber) || err != ErrCodeStoreOutOfGas)) {
		log.Debugf("Create reverting")
		evm.StateDB.RevertToSnapshot(snapshot)
		if err != errExecutionReverted {
			log.Debugf("Create using all gas %v", contract.GasRemaining)
			contract.UseGas(contract.GasRemaining)
		}
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	if evm.VmConfig.Debug && evm.depth == 0 {
		evm.VmConfig.Tracer.CaptureEnd(ret, gasRemaining-contract.GasRemaining, time.Since(start), err)
	}
	return ret, contractAddr, contract.GasRemaining, err
}

// Create creates a system contract using code as deployment code.
func (evm *EVM) CreateSystemContract(sysaddr common.Address, code []byte, precompiledContracts ContractsInterface, msgHash *common.Hash) (ret []byte, err error) {
	log.Info("[core/vm/evm.go->EVM.CreateSystemContract] sysaddr=" + sysaddr.String())
	// Depth check execution. Fail if we're trying to execute above the
	// limit.
	if evm.depth > int(params.CallCreateDepth) {
		return nil, ErrDepth
	}

	// Ensure there's no existing contract already at the designated address
	nonce := evm.StateDB.GetNonce(precompiledContracts.SystemContractCallAddr())
	evm.StateDB.SetNonce(precompiledContracts.SystemContractCallAddr(), nonce+1)

	contractAddr := sysaddr
	contractHash := evm.StateDB.GetCodeHash(contractAddr)
	if evm.StateDB.GetNonce(contractAddr) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		return nil, ErrContractAddressCollision
	}
	// Create a new account on the state
	snapshot := evm.StateDB.Snapshot()
	evm.StateDB.CreateAccount(contractAddr)

	// MOAC use EIP158 check
	// if evm.ChainConfig().IsEIP158(evm.BlockNumber) {
	evm.StateDB.SetNonce(contractAddr, 1)
	// }

	// initialise a new contract and set the code that is to be used by the
	// E The contract is a scoped evmironment for this execution context
	// only.
	value := big.NewInt(0)
	contract := NewContract(AccountRef(precompiledContracts.SystemContractCallAddr()), AccountRef(contractAddr), value, 0)
	contract.SetCallCode(&contractAddr, crypto.Keccak256Hash(code), code)

	if evm.VmConfig.NoRecursion && evm.depth > 0 {
		return nil, nil
	}
	evm.interpreter.Cfg.DisableGasMetering = true
	ret, err = Run(evm, snapshot, contract, nil, precompiledContracts, msgHash)
	// check whether the max code size has been exceeded
	// MOAC use EIP158 check
	// maxCodeSizeExceeded := evm.ChainConfig().IsEIP158(evm.BlockNumber) && len(ret) > params.MaxCodeSize
	maxCodeSizeExceeded := len(ret) > params.MaxCodeSize

	if evm.ChainConfig().IsNuwa(evm.BlockNumber) {
		maxCodeSizeExceeded = len(ret) > params.NuwaMaxCodeSize
	}

	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		evm.StateDB.SetCode(contractAddr, ret)
	}

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in pangu this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && (evm.ChainConfig().IsPangu(evm.BlockNumber) || err != ErrCodeStoreOutOfGas)) {
		evm.StateDB.RevertToSnapshot(snapshot)
	}
	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}
	return ret, err
}

func (evm *EVM) Query(addr common.Address) (ret []byte, leftOverGas uint64, err error) {
	log.Debugf("[core/vm/evm.go->Query in] addr=%v", addr.String())
	var (
		snapshot = evm.StateDB.Snapshot()
	)

	senderBytes, _ := rlp.EncodeToBytes(addr)
	subChainIdBytes, _ := rlp.EncodeToBytes("testchainid")
	contractAddrBytes, _ := rlp.EncodeToBytes(addr)
	valBytes, _ := rlp.EncodeToBytes(0)

	//requestId := vnode.GetInstance().GetRequestId(true)
	//TODO: make a unified requestId
	requestId := []byte("test")

	msgStruct := &ContractCallStruct{
		Sender:       senderBytes,
		Contractaddr: contractAddrBytes,
		Input:        nil,
		Gas:          0,
		Val:          valBytes,
		Sync:         true,
		Code:         nil,
		Codehash:     nil}
	msgHash, _ := rlp.EncodeToBytes(msgStruct)

	conReq := &pb.ScsPushMsg{
		Requestid:   requestId,
		Timestamp:   common.Int64ToBytes(time.Now().Unix()),
		Requestflag: true,
		Type:        common.IntToBytes(DirectCall),
		Status:      common.IntToBytes(None),
		Scsid:       []byte(""),
		Subchainid:  subChainIdBytes,
		Sender:      senderBytes,
		Receiver:    contractAddrBytes,
		Msghash:     msgHash,
	}

	r, err := evm.Nr.VnodePushMsg(conReq)
	if err != nil {
		log.Debug("could not call query. requestid:" + string(conReq.Requestid) + " err:" + err.Error())
	} else {
		log.Debug("call returned successfully.")
	}

	ret, cgas, err := evm.decodeRespond(r)
	//add by frank end

	// When an error was returned by the EVM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in pangu this also counts for code storage gas errors.
	if err != nil {
		evm.StateDB.RevertToSnapshot(snapshot)
		// if err != errExecutionReverted {
		// 	contract.UseGas(contract.GasRemaining)
		// }
	}
	evm.ApplyFlushChanges(r)
	return ret, cgas, err

}

func (evm *EVM) ApplyFlushChanges(flushResult map[int]*pb.ScsPushMsg) (err error) {
	log.Debugf("[core/vm/evm.go->ApplyFlushChanges in] flushResult=%v", flushResult)

	//0) read through the detail
	ret, _, err := evm.decodeRespond(flushResult)
	queryResult, err := evm.decodeQueryResult(ret)
	address := queryResult.address
	ChgObj, err := evm.StateDB.ScsStateChgStringToObj(address, ret)
	if err != nil {
		return err
	}

	accounts := ChgObj.Accounts
	for _, account := range accounts {
		//1) apply data changes
		for key, value := range account.Storage {
			evm.StateDB.SetState(address, common.HexToHash(key), common.HexToHash(value))
		}
		//2) apply transaction
		log.Debug("ApplyFlushChanges balanceChg:" + account.BalanceChg)
		balanceChg := new(big.Int)
		balanceChg.SetString(account.BalanceChg, 10)
		if balanceChg.Cmp(big.NewInt(0)) > 0 {
			evm.StateDB.AddBalance(address, balanceChg)
		}
	}

	return err
}

// ChainConfig returns the evmironment's chain configuration
func (evm *EVM) ChainConfig() *params.ChainConfig { return evm.chainConfig }

// Interpreter returns the EVM interpreter
func (evm *EVM) Interpreter() *Interpreter { return evm.interpreter }
