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

package contracts

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"math/big"

	"sync"

	"strings"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/common/math"
	"github.com/filestorm/go-filestorm/moac/moac-lib/crypto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/crypto/bn256"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/vm"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
	"golang.org/x/crypto/ripemd160"
)

var (
	systemCntEntry = 101

	systemContractCallAddr    = common.BytesToAddress([]byte{100}) // call addr
	systemContractEntryAddrV1 = common.BytesToAddress([]byte{101}) // entry addr
	systemContractEntryAddrV2 = common.BytesToAddress([]byte{102}) // entry addr
	systemContractEntryAddrV3 = common.BytesToAddress([]byte{103}) // entry addr
	systemContractEntryAddrV4 = common.BytesToAddress([]byte{104}) // entry addr
	systemContractEntryAddrV5 = common.BytesToAddress([]byte{105}) // entry addr
	systemContractEntryAddrV6 = common.BytesToAddress([]byte{106}) // entry addr

	//needs to be changed
	//Updated Nov 14th, 2018 from vnode-shanghai-l for 99
	whiteListContractCallAddr105 = common.HexToAddress("0xaafec39f9eab730ccfb9e8d94391d671c6495355") // entry addr
	//Updated Jan 15th, 2019 from vnode-shanghai-l for 99
	//whiteListContractCallAddr = common.HexToAddress("0x5c451e808a625ef0c2b07fb4b4c731dafe3a7097")
	//Updated Jan 15th, 2019 from vnode-shanghai-l for 188
	//whiteListContractCallAddr = common.HexToAddress("0xc07d377de9821b7642d8cd2f97f89e53f87447da")
	//Updated Mar 28th, 2019 from vnode-shanghai-l for 99
	whiteListContractCallAddr = common.HexToAddress("0x406f4c2d1ca8be4709b6a3cc66066cf176aa01ca")
	isValidHash               = "10a0c274"
	getWhiteListHash          = "0fa38fb7"
	getWhiteFuncHash          = "93a670fb"
	getWhiteInfoHash          = "6fb7d382"
	registerHash              = "e1fa8e84"
	registerFuncHash          = "dccbd014"
	removeHash                = "95bc2673"
	removeFuncHash            = "21b04c93"
)

var mu sync.Mutex

type PrecompiledContracts struct {
}

var instance *PrecompiledContracts

//Singleton
func GetInstance() *PrecompiledContracts {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		instance = &PrecompiledContracts{}
	}
	return instance
}

// RunPrecompiledContract runs and evaluates the output of a precompiled contract.
func (pc *PrecompiledContracts) IsSystemCaller(caller vm.ContractRef) (ret bool) {
	//for future system contract update
	return caller.Address() == systemContractCallAddr
}

func IsSystemCaller(caller vm.ContractRef) (ret bool) {
	return caller.Address() == systemContractCallAddr
}

func (pc *PrecompiledContracts) SystemCntEntry() int {
	return systemCntEntry
}

// PrecompiledContract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
//type PrecompiledContract interface {
//	RequiredGas(input []byte) uint64                                              // RequiredPrice calculates the contract gas use
//	Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte) ([]byte, error) // Run runs the precompiled contract
//}

// PrecompiledContractsPangu contains the default set of pre-compiled MoacNode
// contracts used in the Frontier and Pangu releases.
var precompiledContractsPangu = map[common.Address]vm.PrecompiledContract{
	common.BytesToAddress([]byte{1}):  &ecrecover{},
	common.BytesToAddress([]byte{2}):  &sha256hash{},
	common.BytesToAddress([]byte{3}):  &ripemd160hash{},
	common.BytesToAddress([]byte{4}):  &dataCopy{},
	common.BytesToAddress([]byte{9}):  &localShardCheckAndEnroll{},
	common.BytesToAddress([]byte{10}): &checkShardValid{},
	common.BytesToAddress([]byte{11}): &queryContract{},
	common.BytesToAddress([]byte{12}): &delegateSend{},
	common.BytesToAddress([]byte{13}): &notifySCS{},
	common.BytesToAddress([]byte{14}): &spendGas{},
	//system contract
	systemContractEntryAddrV1: &systemContract{},
}

func (pc *PrecompiledContracts) PrecompiledContractsPangu() map[common.Address]vm.PrecompiledContract {
	return precompiledContractsPangu
}

// PrecompiledContractsByzantium contains the default set of pre-compiled MoacNode
// contracts used in the Byzantium release.
var precompiledContractsByzantium = map[common.Address]vm.PrecompiledContract{
	common.BytesToAddress([]byte{1}):  &ecrecover{},
	common.BytesToAddress([]byte{2}):  &sha256hash{},
	common.BytesToAddress([]byte{3}):  &ripemd160hash{},
	common.BytesToAddress([]byte{4}):  &dataCopy{},
	common.BytesToAddress([]byte{5}):  &bigModExp{},
	common.BytesToAddress([]byte{6}):  &bn256Add{},
	common.BytesToAddress([]byte{7}):  &bn256ScalarMul{},
	common.BytesToAddress([]byte{8}):  &bn256Pairing{},
	common.BytesToAddress([]byte{9}):  &localShardCheckAndEnroll{},
	common.BytesToAddress([]byte{10}): &checkShardValid{},
	common.BytesToAddress([]byte{11}): &queryContract{},
	common.BytesToAddress([]byte{12}): &delegateSend{},
	common.BytesToAddress([]byte{13}): &notifySCS{},
	common.BytesToAddress([]byte{14}): &spendGas{},
	//system contract
	systemContractEntryAddrV1: &systemContract{},
}

func (pc *PrecompiledContracts) PrecompiledContractsByzantium() map[common.Address]vm.PrecompiledContract {
	return precompiledContractsByzantium
}

func (pc *PrecompiledContracts) SystemContractCallAddr() common.Address {
	return systemContractCallAddr
}

func (pc *PrecompiledContracts) WhiteListCallAddr() common.Address {
	return whiteListContractCallAddr
}

func (pc *PrecompiledContracts) SystemContractEntryAddr(num *big.Int) common.Address {
	// for upgrade
	return systemContractEntryAddrV1
}

// RunPrecompiledContract runs and evaluates the output of a precompiled contract.
func (pc *PrecompiledContracts) RunPrecompiledContract(evm *vm.EVM, snapshot int, p vm.PrecompiledContract, input []byte, contract *vm.Contract, hash *common.Hash) (ret []byte, err error) {
	log.Debugf("RunPrecompiledContract")
	gas := p.RequiredGas(input)
	if contract.UseGas(gas) {
		log.Debugf("RunPrecompiledContract contract.UseGas(gas)")
		return p.Run(evm, snapshot, contract, input, hash)
	}
	log.Debugf("RunPrecompiledContract ErrOutOfGas")
	return nil, vm.ErrOutOfGas
}

// ECRECOVER implemented as a native contract.
type ecrecover struct{}

func (c *ecrecover) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (c *ecrecover) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	const ecRecoverInputLength = 128

	input = common.RightPadBytes(input, ecRecoverInputLength)
	// "input" is (hash, v, r, s), each 32 bytes
	// but for ecrecover we want (r, s, v)

	r := new(big.Int).SetBytes(input[64:96])
	s := new(big.Int).SetBytes(input[96:128])
	v := input[63] - 27

	// tighter sig s values input pangu only apply to tx sigs
	if !vm.AllZero(input[32:63]) || !crypto.ValidateSignatureValues(v, r, s, false) {
		return nil, nil
		//for test only
		//return common.LeftPadBytes([]byte("0101020202010101020202"), 32), nil
	}
	// v needs to be at the end for libsecp256k1
	pubKey, err := crypto.Ecrecover(input[:32], append(input[64:128], v))
	// make sure the public key is a valid one
	if err != nil {
		return nil, nil
	}

	// the first byte of pubkey is bitcoin heritage
	return common.LeftPadBytes(crypto.Keccak256(pubKey[1:])[12:], 32), nil
}

// SHA256 implemented as a native contract.
type sha256hash struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *sha256hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*params.Sha256PerWordGas + params.Sha256BaseGas
}
func (c *sha256hash) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	h := sha256.Sum256(input)
	return h[:], nil
}

// RIPMED160 implemented as a native contract.
type ripemd160hash struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *ripemd160hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*params.Ripemd160PerWordGas + params.Ripemd160BaseGas
}
func (c *ripemd160hash) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	ripemd := ripemd160.New()
	ripemd.Write(input)
	return common.LeftPadBytes(ripemd.Sum(nil), 32), nil
}

// data copy implemented as a native contract.
type dataCopy struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *dataCopy) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*params.IdentityPerWordGas + params.IdentityBaseGas
}
func (c *dataCopy) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, in []byte, hash *common.Hash) ([]byte, error) {
	return in, nil
}

// bigModExp implements a native big integer exponential modular operation.
type bigModExp struct{}

var (
	big1      = big.NewInt(1)
	big4      = big.NewInt(4)
	big8      = big.NewInt(8)
	big16     = big.NewInt(16)
	big32     = big.NewInt(32)
	big64     = big.NewInt(64)
	big96     = big.NewInt(96)
	big480    = big.NewInt(480)
	big1024   = big.NewInt(1024)
	big3072   = big.NewInt(3072)
	big199680 = big.NewInt(199680)
)

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bigModExp) RequiredGas(input []byte) uint64 {
	var (
		baseLen = new(big.Int).SetBytes(vm.GetData(input, 0, 32))
		expLen  = new(big.Int).SetBytes(vm.GetData(input, 32, 32))
		modLen  = new(big.Int).SetBytes(vm.GetData(input, 64, 32))
	)
	if len(input) > 96 {
		input = input[96:]
	} else {
		input = input[:0]
	}
	// Retrieve the head 32 bytes of exp for the adjusted exponent length
	var expHead *big.Int
	if big.NewInt(int64(len(input))).Cmp(baseLen) <= 0 {
		expHead = new(big.Int)
	} else {
		if expLen.Cmp(big32) > 0 {
			expHead = new(big.Int).SetBytes(vm.GetData(input, baseLen.Uint64(), 32))
		} else {
			expHead = new(big.Int).SetBytes(vm.GetData(input, baseLen.Uint64(), expLen.Uint64()))
		}
	}
	// Calculate the adjusted exponent length
	var msb int
	if bitlen := expHead.BitLen(); bitlen > 0 {
		msb = bitlen - 1
	}
	adjExpLen := new(big.Int)
	if expLen.Cmp(big32) > 0 {
		adjExpLen.Sub(expLen, big32)
		adjExpLen.Mul(big8, adjExpLen)
	}
	adjExpLen.Add(adjExpLen, big.NewInt(int64(msb)))

	// Calculate the gas cost of the operation
	gas := new(big.Int).Set(math.BigMax(modLen, baseLen))
	switch {
	case gas.Cmp(big64) <= 0:
		gas.Mul(gas, gas)
	case gas.Cmp(big1024) <= 0:
		gas = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(gas, gas), big4),
			new(big.Int).Sub(new(big.Int).Mul(big96, gas), big3072),
		)
	default:
		gas = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(gas, gas), big16),
			new(big.Int).Sub(new(big.Int).Mul(big480, gas), big199680),
		)
	}
	gas.Mul(gas, math.BigMax(adjExpLen, big1))
	gas.Div(gas, new(big.Int).SetUint64(params.ModExpQuadCoeffDiv))

	if gas.BitLen() > 64 {
		return math.MaxUint64
	}
	return gas.Uint64()
}

func (c *bigModExp) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	var (
		baseLen = new(big.Int).SetBytes(vm.GetData(input, 0, 32)).Uint64()
		expLen  = new(big.Int).SetBytes(vm.GetData(input, 32, 32)).Uint64()
		modLen  = new(big.Int).SetBytes(vm.GetData(input, 64, 32)).Uint64()
	)
	if len(input) > 96 {
		input = input[96:]
	} else {
		input = input[:0]
	}
	// Handle a special case when both the base and mod length is zero
	if baseLen == 0 && modLen == 0 {
		return []byte{}, nil
	}
	// Retrieve the operands and execute the exponentiation
	var (
		base = new(big.Int).SetBytes(vm.GetData(input, 0, baseLen))
		exp  = new(big.Int).SetBytes(vm.GetData(input, baseLen, expLen))
		mod  = new(big.Int).SetBytes(vm.GetData(input, baseLen+expLen, modLen))
	)
	if mod.BitLen() == 0 {
		// Modulo 0 is undefined, return zero
		return common.LeftPadBytes([]byte{}, int(modLen)), nil
	}
	return common.LeftPadBytes(base.Exp(base, exp, mod).Bytes(), int(modLen)), nil
}

var (
	// errNotOnCurve is returned if a point being unmarshalled as a bn256 elliptic
	// curve point is not on the curve.
	errNotOnCurve = errors.New("point not on elliptic curve")

	// errInvalidCurvePoint is returned if a point being unmarshalled as a bn256
	// elliptic curve point is invalid.
	errInvalidCurvePoint = errors.New("invalid elliptic curve point")
)

// newCurvePoint unmarshals a binary blob into a bn256 elliptic curve point,
// returning it, or an error if the point is invalid.
func newCurvePoint(blob []byte) (*bn256.G1, error) {
	p, onCurve := new(bn256.G1).Unmarshal(blob)
	if !onCurve {
		return nil, errNotOnCurve
	}
	gx, gy, _, _ := p.CurvePoints()
	if gx.Cmp(bn256.P) >= 0 || gy.Cmp(bn256.P) >= 0 {
		return nil, errInvalidCurvePoint
	}
	return p, nil
}

// newTwistPoint unmarshals a binary blob into a bn256 elliptic curve point,
// returning it, or an error if the point is invalid.
func newTwistPoint(blob []byte) (*bn256.G2, error) {
	p, onCurve := new(bn256.G2).Unmarshal(blob)
	if !onCurve {
		return nil, errNotOnCurve
	}
	x2, y2, _, _ := p.CurvePoints()
	if x2.Real().Cmp(bn256.P) >= 0 || x2.Imag().Cmp(bn256.P) >= 0 ||
		y2.Real().Cmp(bn256.P) >= 0 || y2.Imag().Cmp(bn256.P) >= 0 {
		return nil, errInvalidCurvePoint
	}
	return p, nil
}

// bn256Add implements a native elliptic curve point addition.
type bn256Add struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256Add) RequiredGas(input []byte) uint64 {
	return params.Bn256AddGas
}

func (c *bn256Add) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	x, err := newCurvePoint(vm.GetData(input, 0, 64))
	if err != nil {
		return nil, err
	}
	y, err := newCurvePoint(vm.GetData(input, 64, 64))
	if err != nil {
		return nil, err
	}
	res := new(bn256.G1)
	res.Add(x, y)
	return res.Marshal(), nil
}

// bn256ScalarMul implements a native elliptic curve scalar multiplication.
type bn256ScalarMul struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256ScalarMul) RequiredGas(input []byte) uint64 {
	return params.Bn256ScalarMulGas
}

func (c *bn256ScalarMul) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	p, err := newCurvePoint(vm.GetData(input, 0, 64))
	if err != nil {
		return nil, err
	}
	res := new(bn256.G1)
	res.ScalarMult(p, new(big.Int).SetBytes(vm.GetData(input, 64, 32)))
	return res.Marshal(), nil
}

var (
	// true32Byte is returned if the bn256 pairing check succeeds.
	true32Byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// false32Byte is returned if the bn256 pairing check fails.
	false32Byte = make([]byte, 32)

	// errBadPairingInput is returned if the bn256 pairing input is invalid.
	errBadPairingInput = errors.New("bad elliptic curve pairing size")
)

// bn256Pairing implements a pairing pre-compile for the bn256 curve
type bn256Pairing struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256Pairing) RequiredGas(input []byte) uint64 {
	return params.Bn256PairingBaseGas + uint64(len(input)/192)*params.Bn256PairingPerPointGas
}

func (c *bn256Pairing) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	// Handle some corner cases cheaply
	if len(input)%192 > 0 {
		return nil, errBadPairingInput
	}
	// Convert the input into a set of coordinates
	var (
		cs []*bn256.G1
		ts []*bn256.G2
	)
	for i := 0; i < len(input); i += 192 {
		c, err := newCurvePoint(input[i : i+64])
		if err != nil {
			return nil, err
		}
		t, err := newTwistPoint(input[i+64 : i+192])
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
		ts = append(ts, t)
	}
	// Execute the pairing checks and return the results
	if bn256.PairingCheck(cs, ts) {
		return true32Byte, nil
	}
	return false32Byte, nil
}

var (
	// errBadPairingInput is returned if the bn256 pairing input is invalid.
	errBadEnrollCheckArgs = errors.New("bad check enroll args")
)

// localShardCheckAndEnroll implements sharding enrollment with check
type checkShardValid struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *checkShardValid) RequiredGas(input []byte) uint64 {
	return 100000
}

func (c *checkShardValid) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	//log.Info("[core/vm/contracts.go->checkShardValid.Run]")
	// Handle some corner cases cheaply
	if len(input) != 100 {
		return false32Byte, errBadEnrollCheckArgs
	}

	return false32Byte, nil
}

var (
	// errBadPairingInput is returned if the bn256 pairing input is invalid.
	errBadEnrollArgs = errors.New("bad check and enroll args")
)

// localShardCheckAndEnroll implements sharding enrollment with check
type localShardCheckAndEnroll struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *localShardCheckAndEnroll) RequiredGas(input []byte) uint64 {
	return 100000
}

func (c *localShardCheckAndEnroll) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	// Handle some corner cases cheaply
	if len(input) != 132 {
		return false32Byte, errBadEnrollArgs
	}

	return false32Byte, nil
}

// systemContract implements the system contract features
type systemContract struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *systemContract) RequiredGas(input []byte) uint64 {
	//for each function
	return params.SystemContractGas
}

func (c *systemContract) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, msgHash *common.Hash) ([]byte, error) {
	// log.Info("systemcontract run", "from", contract.CallerAddress, "input", input)
	caller := vm.AccountRef(contract.CallerAddress)

	if IsSystemCaller(caller) {
		evm.Interpreter().Cfg.DisableGasMetering = true
	} else {
		//log.Info("user call")

	}
	//statedb := evm.StateDB.(*state.StateDB)
	//log.Info("statedb", "statedb", string(statedb.DumpAccountStorage(SystemContractEntryAddr)))

	ret, err := evm.Interpreter().Run(snapshot, contract, input, GetInstance(), msgHash)
	return ret, err
}

// queryContract implements the query contract features
type queryContract struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *queryContract) RequiredGas(input []byte) uint64 {
	//for each function
	return 100000
}

func (c *queryContract) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	log.Info("[core/vm/contracts.go->queryContract.Run]", "from", contract.CallerAddress, "input", input)

	//getting first element
	if len(input) < 132 {
		return false32Byte, vm.ErrInputFormat
	}

	return false32Byte, nil
}

// delegateSend implements the query contract features
type delegateSend struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *delegateSend) RequiredGas(input []byte) uint64 {
	//for each function
	return 100000
}

func (c *delegateSend) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	//log.Debugf("[core/vm/contracts.go->delegateSend.Run] from:%v input:%v", contract.CallerAddress, "input", input)

	return false32Byte, nil
}

type notifySCS struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *notifySCS) RequiredGas(input []byte) uint64 {
	log.Debugf("[core/vm/contracts.go->notifySCS.RequiredGas] input:%v output:%v", common.Bytes2Hex(input), params.NotifyScsGas)
	return params.NotifyScsGas
}

func (c *notifySCS) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	log.Debugf("[core/vm/contracts.go->notifySCS.Run] from:%v input:%v", contract.CallerAddress.String(), common.Bytes2Hex(input))

	//Use networkRelay to notify SCS about the transaction.
	if evm.Nr != nil && hash != nil && IsInWhiteList(evm, contract.CallerAddress) {
		log.Debugf("[core/vm/contracts.go->notifySCS.Run] hash:%v", hash.String())
		evm.Nr.NotifyScs(contract.CallerAddress, input, *hash, evm.BlockNumber)
		if evm.ChainConfig().IsNuwa(evm.BlockNumber) {
			return true32Byte, nil
		}
	}

	// call v-node check and relay functions.
	if evm.ChainConfig().IsNuwa(evm.BlockNumber) {
		return false32Byte, nil
	} else {
		return nil, nil
	}
}

func IsWhiteListCall(addr common.Address, data []byte) bool {
	if addr != whiteListContractCallAddr {
		return false
	}

	if bytes.Equal(common.Hex2Bytes(registerFuncHash), data) || bytes.Equal(common.Hex2Bytes(removeFuncHash), data) ||
		bytes.Equal(common.Hex2Bytes(registerHash), data) || bytes.Equal(common.Hex2Bytes(removeHash), data) {
		return true
	}

	return false
}

func IsInWhiteList(evm *vm.EVM, callerAddress common.Address) bool {
	if evm == nil {
		evm = vm.GetEVM()
		if evm == nil {
			log.Debug("[core/vm/contracts.go->IsInWhiteList] evm nil")
			return false
		}
	}

	networkId := evm.ChainConfig().ChainId.Uint64()
	if !params.PriorityChain(networkId) {
		return true
	}
	snapshot := evm.StateDB.Snapshot()
	//log.Infof("IsInWhiteList callerAddress %v, whiteListContractCallAddr %v", callerAddress.Hex(), whiteListContractCallAddr.Hex())
	whiteListContract := vm.NewContract(vm.AccountRef(callerAddress), vm.AccountRef(whiteListContractCallAddr), big.NewInt(0), 1000000)
	whiteListContract.SetCallCode(&whiteListContractCallAddr, evm.StateDB.GetCodeHash(whiteListContractCallAddr), evm.StateDB.GetCode(whiteListContractCallAddr))

	//code := evm.StateDB.GetCode(callerAddress)
	//h := crypto.Keccak256(code);; h==callContractHashcode?
	callContractHashcode := evm.StateDB.GetCodeHash(callerAddress)

	codeHashHex1 := strings.ToLower(callContractHashcode.Hex())
	codeHashHex1 = codeHashHex1[2:len(codeHashHex1)]
	codeHashHex2 := strings.ToLower(callerAddress.Hex())
	codeHashHex2 = codeHashHex2[2:len(codeHashHex2)]
	//"10a0c274" isValid(bytes32)
	callingHash := isValidHash + codeHashHex2 + "000000000000000000000000" + codeHashHex1
	log.Debugf("callingHash %v", callingHash)
	input := common.Hex2Bytes(callingHash)

	precompiledContracts := GetInstance()
	ret, err := vm.Run(evm, snapshot, whiteListContract, input, precompiledContracts, nil)
	if err != nil {
		log.Errorf("IsInWhiteList error %v", err.Error())
		return false
	}
	retValue := common.Bytes2Hex(ret)
	log.Debugf("IsInWhiteList retValue %v, ret %v", retValue, ret)
	if retValue != "0000000000000000000000000000000000000000000000000000000000000000" {
		log.Debugf("IsInWhiteList returning true")
		return true
	} else {
		log.Debugf("IsInWhiteList returning false")
		return false
	}
}

// GetWhiteListCommonHash returns white list contract code hash
func GetWhiteListCommonHash(callerAddress common.Address) common.Hash {
	evm := vm.GetEVM()
	if evm == nil {
		log.Debug("[core/vm/contracts.go->GetWhiteListCommonHash] evm nil")
		return common.Hash{}
	}

	return evm.StateDB.GetCodeHash(callerAddress)
}

// GetWhiteInfo return White Functions
func GetWhiteInfo() (wList, wFunc []common.Hash) {
	evm := vm.GetEVM()
	if evm == nil {
		log.Debug("[core/vm/contracts.go->GetWhiteInfo] evm nil")
		return
	}

	snapshot := evm.StateDB.Snapshot()
	//log.Infof("IsInWhiteList callerAddress %v, whiteListContractCallAddr %v", callerAddress.Hex(), whiteListContractCallAddr.Hex())
	whiteListContract := vm.NewContract(vm.AccountRef(common.Address{}), vm.AccountRef(whiteListContractCallAddr), big.NewInt(0), 1000000)
	whiteListContract.SetCallCode(&whiteListContractCallAddr, evm.StateDB.GetCodeHash(whiteListContractCallAddr), evm.StateDB.GetCode(whiteListContractCallAddr))

	callingHash := getWhiteInfoHash + "0000000000000000000000000000000000000000000000000000000000000000"
	callingHash += "00000000000000000000000000000000000000000000000000000000000003e8"
	log.Debugf("callingHash %v", callingHash)
	input := common.Hex2Bytes(callingHash)

	precompiledContracts := GetInstance()
	ret, err := vm.Run(evm, snapshot, whiteListContract, input, precompiledContracts, nil)
	if err != nil {
		log.Errorf("GetWhiteInfo error %v", err.Error())
		return
	}
	log.Debugf("GetWhiteInfo retValue %v", common.Bytes2Hex(ret))
	if len(ret) > 128 {
		temp := ret[:]
		index := common.BytesToInt(temp[28:32])
		if index > len(ret) {
			return
		}
		temp = ret[index:]
		l := common.BytesToInt(temp[28:32])
		temp = ret[index+32:]
		if l*32 > len(temp) {
			return
		}
		for i := 0; i < l && len(temp) >= 32; i++ {
			wList = append(wList, common.BytesToHash(temp[:32]))
			temp = temp[32:]
		}

		temp = ret[:]
		index = common.BytesToInt(temp[60:64])
		if index > len(ret) {
			return
		}
		temp = ret[index:]
		l = common.BytesToInt(temp[28:32])
		temp = ret[index+32:]
		if l*32 > len(temp) {
			return
		}
		for i := 0; i < l && len(temp) >= 32; i++ {
			wFunc = append(wFunc, common.BytesToHash(temp[:32]))
			temp = temp[32:]
		}
	}

	return
}

// PrintWhiteList prints WhiteList
func PrintWhiteList(evm *vm.EVM, callerAddress common.Address) {
	if evm == nil {
		evm = vm.GetEVM()
		if evm == nil {
			log.Debug("[core/vm/contracts.go->PrintWhiteList] evm nil")
			return
		}
	}

	snapshot := evm.StateDB.Snapshot()
	//log.Infof("IsInWhiteList callerAddress %v, whiteListContractCallAddr %v", callerAddress.Hex(), whiteListContractCallAddr.Hex())
	whiteListContract := vm.NewContract(vm.AccountRef(callerAddress), vm.AccountRef(whiteListContractCallAddr), big.NewInt(0), 1000000)
	whiteListContract.SetCallCode(&whiteListContractCallAddr, evm.StateDB.GetCodeHash(whiteListContractCallAddr), evm.StateDB.GetCode(whiteListContractCallAddr))

	callingHash := getWhiteListHash + "0000000000000000000000000000000000000000000000000000000000000000"
	callingHash += "00000000000000000000000000000000000000000000000000000000000003e8"
	log.Debugf("callingHash %v", callingHash)
	input := common.Hex2Bytes(callingHash)

	precompiledContracts := GetInstance()
	ret, err := vm.Run(evm, snapshot, whiteListContract, input, precompiledContracts, nil)
	if err != nil {
		log.Errorf("getWhiteList error %v", err.Error())
		return
	}
	retValue := common.Bytes2Hex(ret)
	log.Debugf("IsInWhiteList retValue %v, ret %v", retValue, ret)
}

type spendGas struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *spendGas) RequiredGas(input []byte) uint64 {

	//getting first element
	if len(input) != 36 {
		return uint64(8000000)
	}

	//getting number
	bnum := big.NewInt(0)
	bnum.SetBytes(input[4:])
	num := bnum.Int64()

	//if pangu version
	gas := int64(1000)

	if num > 10 {
		gas += 2000 * (num - 10)
	}
	if num > 15 {
		gas += 4000 * (num - 15)
	}
	if num > 20 {
		gas += 8000 * (num - 20)
	}
	if num > 25 {
		gas += 16000 * (num - 25)
	}
	if num > 30 {
		gas += 32000 * (num - 30)
	}
	if num > 35 {
		gas += 64000 * (num - 35)
	}
	if num > 40 {
		gas += 128000 * (num - 40)
	}
	if num > 45 {
		gas += 256000 * (num - 45)
	}

	log.Debugf("[core/vm/contracts.go->RequiredGas.RequiredGas] input:%v output:%v", input, gas)

	return uint64(gas)
}

func (c *spendGas) Run(evm *vm.EVM, snapshot int, contract *vm.Contract, input []byte, hash *common.Hash) ([]byte, error) {
	log.Debugf("[core/vm/contracts.go->notifySCS.Run] from:%v input:%v", contract.CallerAddress, input)

	if evm.ChainConfig().IsNuwa(evm.BlockNumber) {
		return false32Byte, nil
	} else {
		return nil, nil
	}
}
