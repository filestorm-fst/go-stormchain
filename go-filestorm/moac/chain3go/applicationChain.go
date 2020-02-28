// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain3go

import (
	"math/big"
	"strings"

	filestorm "github.com/filestorm/go-filestorm"
	"github.com/filestorm/go-filestorm/accounts/abi"
	"github.com/filestorm/go-filestorm/accounts/abi/bind"
	"github.com/filestorm/go-filestorm/common"
	"github.com/filestorm/go-filestorm/core/types"
	"github.com/filestorm/go-filestorm/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = filestorm.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ApplicationChainABI is the input ABI used to generate the binding from.
const ApplicationChainABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"FLUSH_AMOUNT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"flushEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushList\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"admins\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushMapping\",\"outputs\":[{\"name\":\"flushId\",\"type\":\"uint256\"},{\"name\":\"validator\",\"type\":\"address\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"FOUNDATION_BLACK_HOLE_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushValidatorList\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chainId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"addFund\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"balance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"current_validators\",\"type\":\"address[]\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"string\"}],\"name\":\"flush\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"genesis\",\"type\":\"string\"}],\"name\":\"setGenesisInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"FOUNDATION_MOAC_REQUIRED_AMOUNT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"period\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recv\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGenesisInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"blockSec\",\"type\":\"uint256\"},{\"name\":\"flushNumber\",\"type\":\"uint256\"},{\"name\":\"initial_validators\",\"type\":\"address[]\"},{\"name\":\"exchangeRate\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"}]"

// ApplicationChainFuncSigs maps the 4-byte function signature to its string representation.
var ApplicationChainFuncSigs = map[string]string{
	"08925f61": "FLUSH_AMOUNT()",
	"5a35b2fc": "FOUNDATION_BLACK_HOLE_ADDRESS()",
	"dfe89af2": "FOUNDATION_MOAC_REQUIRED_AMOUNT()",
	"70480275": "addAdmin(address)",
	"a2f09dfa": "addFund()",
	"429b62e5": "admins(address)",
	"b69ef8a8": "balance()",
	"9a8a0592": "chainId()",
	"d66cfad5": "flush(address[],uint256,string)",
	"090aaae5": "flushEpoch()",
	"18c15211": "flushList(uint256)",
	"4c2f8619": "flushMapping(uint256)",
	"635cf62a": "flushValidatorList(uint256,uint256)",
	"fd8735a5": "getGenesisInfo()",
	"ef78d4fd": "period()",
	"1785f53c": "removeAdmin(address)",
	"d9b11dc8": "setGenesisInfo(string)",
	"f7c8d221": "withdrawFund(address,uint256)",
}

// ApplicationChainBin is the compiled bytecode used for deploying new contracts.
var ApplicationChainBin = "0x6080604081905260006002556801158e460913d00000600d5567016345785d8a0000600e55600f8054600160a060020a0319167348328afc8dd45c1c252e7e883fc89bd17ddee7c017905561100f3881900390819083398101604090815281516020830151918301516060840151600d5492949190910191600090819034101561011057604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f4e6f7420456e6f756768204d4f414320746f20437265617465204170706c696360448201527f6174696f6e20436861696e000000000000000000000000000000000000000000606482015290519081900360840190fd5b50506000805433600160a060020a0319918216811783554360035560048790556005869055346002556007805460ff1916905582805260086020527f5eff886ea0ce6ca488a3d6e336d6c0f75f46d19b42c06ce5ee98e42c96d256c78390557f5eff886ea0ce6ca488a3d6e336d6c0f75f46d19b42c06ce5ee98e42c96d256c88054909216179055805b8351811015610204576000828152600a6020526040902084518590839081106101bf57fe5b6020908102919091018101518254600180820185556000948552929093209092018054600160a060020a031916600160a060020a03909316929092179091550161019a565b600082815260086020818152604080842060016002820155815180840192839052858152948790529290915291516102429260039092019190610296565b50506009805460018181019092557f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af019190915534909102600c553360009081526020829052604090205550610331915050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106102d757805160ff1916838001178555610304565b82800160010185558215610304579182015b828111156103045782518255916020019190600101906102e9565b50610310929150610314565b5090565b61032e91905b80821115610310576000815560010161031a565b90565b610ccf806103406000396000f3006080604052600436106100fb5763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166308925f618114610100578063090aaae5146101275780631785f53c1461013c57806318c152111461015f578063429b62e5146101775780634c2f8619146101985780635a35b2fc1461024f578063635cf62a14610280578063704802751461029b5780639a8a0592146102bc578063a2f09dfa146102d1578063b69ef8a8146102d9578063d66cfad5146102ee578063d9b11dc814610387578063dfe89af2146103e0578063ef78d4fd146103f5578063f7c8d2211461040a578063fd8735a51461042e575b600080fd5b34801561010c57600080fd5b506101156104b8565b60408051918252519081900360200190f35b34801561013357600080fd5b506101156104be565b34801561014857600080fd5b5061015d600160a060020a03600435166104c4565b005b34801561016b57600080fd5b506101156004356105cd565b34801561018357600080fd5b50610115600160a060020a03600435166105ec565b3480156101a457600080fd5b506101b06004356105fe565b6040518085815260200184600160a060020a0316600160a060020a0316815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b838110156102115781810151838201526020016101f9565b50505050905090810190601f16801561023e5780820380516001836020036101000a031916815260200191505b509550505050505060405180910390f35b34801561025b57600080fd5b506102646106bc565b60408051600160a060020a039092168252519081900360200190f35b34801561028c57600080fd5b506102646004356024356106cb565b3480156102a757600080fd5b5061015d600160a060020a0360043516610702565b3480156102c857600080fd5b506101156107ad565b61015d6107b3565b3480156102e557600080fd5b506101156107bd565b3480156102fa57600080fd5b506040805160206004803580820135838102808601850190965280855261015d9536959394602494938501929182918501908490808284375050604080516020601f818a01358b0180359182018390048302840183018552818452989b8a359b909a9099940197509195509182019350915081908401838280828437509497506107c39650505050505050565b34801561039357600080fd5b506040805160206004803580820135601f810184900484028501840190955284845261015d9436949293602493928401919081908401838280828437509497506109c59650505050505050565b3480156103ec57600080fd5b50610115610af8565b34801561040157600080fd5b50610115610afe565b34801561041657600080fd5b5061015d600160a060020a0360043516602435610b04565b34801561043a57600080fd5b50610443610b74565b6040805160208082528351818301528351919283929083019185019080838360005b8381101561047d578181015183820152602001610465565b50505050905090810190601f1680156104aa5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b600e5481565b60055481565b3360009081526001602081905260409091205414610552576040805160e560020a62461bcd02815260206004820152602260248201527f4f6e6c792041646d696e732043616e2041646420416e6f746865722041646d6960448201527f6e2e000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a0381163314156105b3576040805160e560020a62461bcd02815260206004820152601a60248201527f41646d696e732043616e6e6f742052656d6f76652053656c662e000000000000604482015290519081900360640190fd5b600160a060020a0316600090815260016020526040812055565b60098054829081106105db57fe5b600091825260209091200154905081565b60016020526000908152604090205481565b600860209081526000918252604091829020805460018083015460028085015460038601805489516101009682161596909602600019011692909204601f81018890048802850188019098528784529396600160a060020a03909216959394939091908301828280156106b25780601f10610687576101008083540402835291602001916106b2565b820191906000526020600020905b81548152906001019060200180831161069557829003601f168201915b5050505050905084565b600f54600160a060020a031681565b600a602052816000526040600020818154811015156106e657fe5b600091825260209091200154600160a060020a03169150829050565b3360009081526001602081905260409091205414610790576040805160e560020a62461bcd02815260206004820152602260248201527f4f6e6c792041646d696e732043616e2041646420416e6f746865722041646d6960448201527f6e2e000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a0316600090815260016020819052604090912055565b60035481565b6002805434019055565b60025481565b600954600019016000805b6000838152600a60205260409020548210156109bd576000838152600a6020526040902080543391908490811061080157fe5b600091825260209091200154600160a060020a03161480610832575033600090815260016020819052604090912054145b801561085257506005546000848152600860205260409020600201540185145b156109b2575060019182016000818152600860205260408120828155909301805473ffffffffffffffffffffffffffffffffffffffff191633179055915b8551811015610907576000838152600a6020526040902086518790839081106108b557fe5b602090810291909101810151825460018082018555600094855292909320909201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a039093169290921790915501610890565b600083815260086020908152604090912060028101879055855161093392600390920191870190610c0b565b506009805460018101825560009182527f6e1540171b6c0c960b71a7020d9f60077f6af931a8bbf590da0223dacf75c7af01849055600e54600280548290039055600f54604051600160a060020a039091169282156108fc02929190818181858888f193505050501580156109ac573d6000803e3d6000fd5b506109bd565b6001909101906107ce565b505050505050565b3360009081526001602081905260409091205414610a53576040805160e560020a62461bcd02815260206004820152602160248201527f4f6e6c792041646d696e732043616e205365742047656e6573697320496e666f60448201527f2e00000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60075460ff1615610ad4576040805160e560020a62461bcd02815260206004820152602260248201527f47656e6573697320496e666f2048617320416c7265616479204265656e20536560448201527f742e000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b8051610ae7906006906020840190610c0b565b50506007805460ff19166001179055565b600d5481565b60045481565b3360009081526001602081905260409091205414610b2157600080fd5b600254811115610b3057600080fd5b600280548290039055604051600160a060020a0383169082156108fc029083906000818181858888f19350505050158015610b6f573d6000803e3d6000fd5b505050565b60068054604080516020601f6002600019610100600188161502019095169490940493840181900481028201810190925282815260609390929091830182828015610c005780601f10610bd557610100808354040283529160200191610c00565b820191906000526020600020905b815481529060010190602001808311610be357829003601f168201915b505050505090505b90565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f10610c4c57805160ff1916838001178555610c79565b82800160010185558215610c79579182015b82811115610c79578251825591602001919060010190610c5e565b50610c85929150610c89565b5090565b610c0891905b80821115610c855760008155600101610c8f5600a165627a7a72305820c27b97e9fe6aea84f9ebf25d6c36765e315cfb94679b407ac9e8597d2911c6620029"

// DeployApplicationChain deploys a new Filestorm contract, binding an instance of ApplicationChain to it.
func DeployApplicationChain(auth *bind.TransactOpts, backend bind.ContractBackend, blockSec *big.Int, flushNumber *big.Int, initial_validators []common.Address, exchangeRate *big.Int) (common.Address, *types.Transaction, *ApplicationChain, error) {
	parsed, err := abi.JSON(strings.NewReader(ApplicationChainABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(ApplicationChainBin), backend, blockSec, flushNumber, initial_validators, exchangeRate)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ApplicationChain{ApplicationChainCaller: ApplicationChainCaller{contract: contract}, ApplicationChainTransactor: ApplicationChainTransactor{contract: contract}, ApplicationChainFilterer: ApplicationChainFilterer{contract: contract}}, nil
}

// ApplicationChain is an auto generated Go binding around an Filestorm contract.
type ApplicationChain struct {
	ApplicationChainCaller     // Read-only binding to the contract
	ApplicationChainTransactor // Write-only binding to the contract
	ApplicationChainFilterer   // Log filterer for contract events
}

// ApplicationChainCaller is an auto generated read-only Go binding around an Filestorm contract.
type ApplicationChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApplicationChainTransactor is an auto generated write-only Go binding around an Filestorm contract.
type ApplicationChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApplicationChainFilterer is an auto generated log filtering Go binding around an Filestorm contract events.
type ApplicationChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ApplicationChainSession is an auto generated Go binding around an Filestorm contract,
// with pre-set call and transact options.
type ApplicationChainSession struct {
	Contract     *ApplicationChain // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ApplicationChainCallerSession is an auto generated read-only Go binding around an Filestorm contract,
// with pre-set call options.
type ApplicationChainCallerSession struct {
	Contract *ApplicationChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// ApplicationChainTransactorSession is an auto generated write-only Go binding around an Filestorm contract,
// with pre-set transact options.
type ApplicationChainTransactorSession struct {
	Contract     *ApplicationChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// ApplicationChainRaw is an auto generated low-level Go binding around an Filestorm contract.
type ApplicationChainRaw struct {
	Contract *ApplicationChain // Generic contract binding to access the raw methods on
}

// ApplicationChainCallerRaw is an auto generated low-level read-only Go binding around an Filestorm contract.
type ApplicationChainCallerRaw struct {
	Contract *ApplicationChainCaller // Generic read-only contract binding to access the raw methods on
}

// ApplicationChainTransactorRaw is an auto generated low-level write-only Go binding around an Filestorm contract.
type ApplicationChainTransactorRaw struct {
	Contract *ApplicationChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewApplicationChain creates a new instance of ApplicationChain, bound to a specific deployed contract.
func NewApplicationChain(address common.Address, backend bind.ContractBackend) (*ApplicationChain, error) {
	contract, err := bindApplicationChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ApplicationChain{ApplicationChainCaller: ApplicationChainCaller{contract: contract}, ApplicationChainTransactor: ApplicationChainTransactor{contract: contract}, ApplicationChainFilterer: ApplicationChainFilterer{contract: contract}}, nil
}

// NewApplicationChainCaller creates a new read-only instance of ApplicationChain, bound to a specific deployed contract.
func NewApplicationChainCaller(address common.Address, caller bind.ContractCaller) (*ApplicationChainCaller, error) {
	contract, err := bindApplicationChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ApplicationChainCaller{contract: contract}, nil
}

// NewApplicationChainTransactor creates a new write-only instance of ApplicationChain, bound to a specific deployed contract.
func NewApplicationChainTransactor(address common.Address, transactor bind.ContractTransactor) (*ApplicationChainTransactor, error) {
	contract, err := bindApplicationChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ApplicationChainTransactor{contract: contract}, nil
}

// NewApplicationChainFilterer creates a new log filterer instance of ApplicationChain, bound to a specific deployed contract.
func NewApplicationChainFilterer(address common.Address, filterer bind.ContractFilterer) (*ApplicationChainFilterer, error) {
	contract, err := bindApplicationChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ApplicationChainFilterer{contract: contract}, nil
}

// bindApplicationChain binds a generic wrapper to an already deployed contract.
func bindApplicationChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ApplicationChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ApplicationChain *ApplicationChainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ApplicationChain.Contract.ApplicationChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ApplicationChain *ApplicationChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ApplicationChain.Contract.ApplicationChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ApplicationChain *ApplicationChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ApplicationChain.Contract.ApplicationChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ApplicationChain *ApplicationChainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ApplicationChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ApplicationChain *ApplicationChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ApplicationChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ApplicationChain *ApplicationChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ApplicationChain.Contract.contract.Transact(opts, method, params...)
}

// FLUSHAMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCaller) FLUSHAMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "FLUSH_AMOUNT")
	return *ret0, err
}

// FLUSHAMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_ApplicationChain *ApplicationChainSession) FLUSHAMOUNT() (*big.Int, error) {
	return _ApplicationChain.Contract.FLUSHAMOUNT(&_ApplicationChain.CallOpts)
}

// FLUSHAMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCallerSession) FLUSHAMOUNT() (*big.Int, error) {
	return _ApplicationChain.Contract.FLUSHAMOUNT(&_ApplicationChain.CallOpts)
}

// FOUNDATIONBLACKHOLEADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_ApplicationChain *ApplicationChainCaller) FOUNDATIONBLACKHOLEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "FOUNDATION_BLACK_HOLE_ADDRESS")
	return *ret0, err
}

// FOUNDATIONBLACKHOLEADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_ApplicationChain *ApplicationChainSession) FOUNDATIONBLACKHOLEADDRESS() (common.Address, error) {
	return _ApplicationChain.Contract.FOUNDATIONBLACKHOLEADDRESS(&_ApplicationChain.CallOpts)
}

// FOUNDATIONBLACKHOLEADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_ApplicationChain *ApplicationChainCallerSession) FOUNDATIONBLACKHOLEADDRESS() (common.Address, error) {
	return _ApplicationChain.Contract.FOUNDATIONBLACKHOLEADDRESS(&_ApplicationChain.CallOpts)
}

// FOUNDATIONMOACREQUIREDAMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCaller) FOUNDATIONMOACREQUIREDAMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "FOUNDATION_MOAC_REQUIRED_AMOUNT")
	return *ret0, err
}

// FOUNDATIONMOACREQUIREDAMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_ApplicationChain *ApplicationChainSession) FOUNDATIONMOACREQUIREDAMOUNT() (*big.Int, error) {
	return _ApplicationChain.Contract.FOUNDATIONMOACREQUIREDAMOUNT(&_ApplicationChain.CallOpts)
}

// FOUNDATIONMOACREQUIREDAMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCallerSession) FOUNDATIONMOACREQUIREDAMOUNT() (*big.Int, error) {
	return _ApplicationChain.Contract.FOUNDATIONMOACREQUIREDAMOUNT(&_ApplicationChain.CallOpts)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) constant returns(uint256)
func (_ApplicationChain *ApplicationChainCaller) Admins(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "admins", arg0)
	return *ret0, err
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) constant returns(uint256)
func (_ApplicationChain *ApplicationChainSession) Admins(arg0 common.Address) (*big.Int, error) {
	return _ApplicationChain.Contract.Admins(&_ApplicationChain.CallOpts, arg0)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) constant returns(uint256)
func (_ApplicationChain *ApplicationChainCallerSession) Admins(arg0 common.Address) (*big.Int, error) {
	return _ApplicationChain.Contract.Admins(&_ApplicationChain.CallOpts, arg0)
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCaller) Balance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "balance")
	return *ret0, err
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_ApplicationChain *ApplicationChainSession) Balance() (*big.Int, error) {
	return _ApplicationChain.Contract.Balance(&_ApplicationChain.CallOpts)
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCallerSession) Balance() (*big.Int, error) {
	return _ApplicationChain.Contract.Balance(&_ApplicationChain.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCaller) ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "chainId")
	return *ret0, err
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(uint256)
func (_ApplicationChain *ApplicationChainSession) ChainId() (*big.Int, error) {
	return _ApplicationChain.Contract.ChainId(&_ApplicationChain.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCallerSession) ChainId() (*big.Int, error) {
	return _ApplicationChain.Contract.ChainId(&_ApplicationChain.CallOpts)
}

// FlushEpoch is a free data retrieval call binding the contract method 0x090aaae5.
//
// Solidity: function flushEpoch() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCaller) FlushEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "flushEpoch")
	return *ret0, err
}

// FlushEpoch is a free data retrieval call binding the contract method 0x090aaae5.
//
// Solidity: function flushEpoch() constant returns(uint256)
func (_ApplicationChain *ApplicationChainSession) FlushEpoch() (*big.Int, error) {
	return _ApplicationChain.Contract.FlushEpoch(&_ApplicationChain.CallOpts)
}

// FlushEpoch is a free data retrieval call binding the contract method 0x090aaae5.
//
// Solidity: function flushEpoch() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCallerSession) FlushEpoch() (*big.Int, error) {
	return _ApplicationChain.Contract.FlushEpoch(&_ApplicationChain.CallOpts)
}

// FlushList is a free data retrieval call binding the contract method 0x18c15211.
//
// Solidity: function flushList(uint256 ) constant returns(uint256)
func (_ApplicationChain *ApplicationChainCaller) FlushList(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "flushList", arg0)
	return *ret0, err
}

// FlushList is a free data retrieval call binding the contract method 0x18c15211.
//
// Solidity: function flushList(uint256 ) constant returns(uint256)
func (_ApplicationChain *ApplicationChainSession) FlushList(arg0 *big.Int) (*big.Int, error) {
	return _ApplicationChain.Contract.FlushList(&_ApplicationChain.CallOpts, arg0)
}

// FlushList is a free data retrieval call binding the contract method 0x18c15211.
//
// Solidity: function flushList(uint256 ) constant returns(uint256)
func (_ApplicationChain *ApplicationChainCallerSession) FlushList(arg0 *big.Int) (*big.Int, error) {
	return _ApplicationChain.Contract.FlushList(&_ApplicationChain.CallOpts, arg0)
}

// FlushMapping is a free data retrieval call binding the contract method 0x4c2f8619.
//
// Solidity: function flushMapping(uint256 ) constant returns(uint256 flushId, address validator, uint256 blockNumber, string blockHash)
func (_ApplicationChain *ApplicationChainCaller) FlushMapping(opts *bind.CallOpts, arg0 *big.Int) (struct {
	FlushId     *big.Int
	Validator   common.Address
	BlockNumber *big.Int
	BlockHash   string
}, error) {
	ret := new(struct {
		FlushId     *big.Int
		Validator   common.Address
		BlockNumber *big.Int
		BlockHash   string
	})
	out := ret
	err := _ApplicationChain.contract.Call(opts, out, "flushMapping", arg0)
	return *ret, err
}

// FlushMapping is a free data retrieval call binding the contract method 0x4c2f8619.
//
// Solidity: function flushMapping(uint256 ) constant returns(uint256 flushId, address validator, uint256 blockNumber, string blockHash)
func (_ApplicationChain *ApplicationChainSession) FlushMapping(arg0 *big.Int) (struct {
	FlushId     *big.Int
	Validator   common.Address
	BlockNumber *big.Int
	BlockHash   string
}, error) {
	return _ApplicationChain.Contract.FlushMapping(&_ApplicationChain.CallOpts, arg0)
}

// FlushMapping is a free data retrieval call binding the contract method 0x4c2f8619.
//
// Solidity: function flushMapping(uint256 ) constant returns(uint256 flushId, address validator, uint256 blockNumber, string blockHash)
func (_ApplicationChain *ApplicationChainCallerSession) FlushMapping(arg0 *big.Int) (struct {
	FlushId     *big.Int
	Validator   common.Address
	BlockNumber *big.Int
	BlockHash   string
}, error) {
	return _ApplicationChain.Contract.FlushMapping(&_ApplicationChain.CallOpts, arg0)
}

// FlushValidatorList is a free data retrieval call binding the contract method 0x635cf62a.
//
// Solidity: function flushValidatorList(uint256 , uint256 ) constant returns(address)
func (_ApplicationChain *ApplicationChainCaller) FlushValidatorList(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "flushValidatorList", arg0, arg1)
	return *ret0, err
}

// FlushValidatorList is a free data retrieval call binding the contract method 0x635cf62a.
//
// Solidity: function flushValidatorList(uint256 , uint256 ) constant returns(address)
func (_ApplicationChain *ApplicationChainSession) FlushValidatorList(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _ApplicationChain.Contract.FlushValidatorList(&_ApplicationChain.CallOpts, arg0, arg1)
}

// FlushValidatorList is a free data retrieval call binding the contract method 0x635cf62a.
//
// Solidity: function flushValidatorList(uint256 , uint256 ) constant returns(address)
func (_ApplicationChain *ApplicationChainCallerSession) FlushValidatorList(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _ApplicationChain.Contract.FlushValidatorList(&_ApplicationChain.CallOpts, arg0, arg1)
}

// GetGenesisInfo is a free data retrieval call binding the contract method 0xfd8735a5.
//
// Solidity: function getGenesisInfo() constant returns(string)
func (_ApplicationChain *ApplicationChainCaller) GetGenesisInfo(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "getGenesisInfo")
	return *ret0, err
}

// GetGenesisInfo is a free data retrieval call binding the contract method 0xfd8735a5.
//
// Solidity: function getGenesisInfo() constant returns(string)
func (_ApplicationChain *ApplicationChainSession) GetGenesisInfo() (string, error) {
	return _ApplicationChain.Contract.GetGenesisInfo(&_ApplicationChain.CallOpts)
}

// GetGenesisInfo is a free data retrieval call binding the contract method 0xfd8735a5.
//
// Solidity: function getGenesisInfo() constant returns(string)
func (_ApplicationChain *ApplicationChainCallerSession) GetGenesisInfo() (string, error) {
	return _ApplicationChain.Contract.GetGenesisInfo(&_ApplicationChain.CallOpts)
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCaller) Period(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ApplicationChain.contract.Call(opts, out, "period")
	return *ret0, err
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_ApplicationChain *ApplicationChainSession) Period() (*big.Int, error) {
	return _ApplicationChain.Contract.Period(&_ApplicationChain.CallOpts)
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_ApplicationChain *ApplicationChainCallerSession) Period() (*big.Int, error) {
	return _ApplicationChain.Contract.Period(&_ApplicationChain.CallOpts)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_ApplicationChain *ApplicationChainTransactor) AddAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ApplicationChain.contract.Transact(opts, "addAdmin", admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_ApplicationChain *ApplicationChainSession) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _ApplicationChain.Contract.AddAdmin(&_ApplicationChain.TransactOpts, admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_ApplicationChain *ApplicationChainTransactorSession) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _ApplicationChain.Contract.AddAdmin(&_ApplicationChain.TransactOpts, admin)
}

// AddFund is a paid mutator transaction binding the contract method 0xa2f09dfa.
//
// Solidity: function addFund() returns()
func (_ApplicationChain *ApplicationChainTransactor) AddFund(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ApplicationChain.contract.Transact(opts, "addFund")
}

// AddFund is a paid mutator transaction binding the contract method 0xa2f09dfa.
//
// Solidity: function addFund() returns()
func (_ApplicationChain *ApplicationChainSession) AddFund() (*types.Transaction, error) {
	return _ApplicationChain.Contract.AddFund(&_ApplicationChain.TransactOpts)
}

// AddFund is a paid mutator transaction binding the contract method 0xa2f09dfa.
//
// Solidity: function addFund() returns()
func (_ApplicationChain *ApplicationChainTransactorSession) AddFund() (*types.Transaction, error) {
	return _ApplicationChain.Contract.AddFund(&_ApplicationChain.TransactOpts)
}

// Flush is a paid mutator transaction binding the contract method 0xd66cfad5.
//
// Solidity: function flush(address[] current_validators, uint256 blockNumber, string blockHash) returns()
func (_ApplicationChain *ApplicationChainTransactor) Flush(opts *bind.TransactOpts, current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _ApplicationChain.contract.Transact(opts, "flush", current_validators, blockNumber, blockHash)
}

// Flush is a paid mutator transaction binding the contract method 0xd66cfad5.
//
// Solidity: function flush(address[] current_validators, uint256 blockNumber, string blockHash) returns()
func (_ApplicationChain *ApplicationChainSession) Flush(current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _ApplicationChain.Contract.Flush(&_ApplicationChain.TransactOpts, current_validators, blockNumber, blockHash)
}

// Flush is a paid mutator transaction binding the contract method 0xd66cfad5.
//
// Solidity: function flush(address[] current_validators, uint256 blockNumber, string blockHash) returns()
func (_ApplicationChain *ApplicationChainTransactorSession) Flush(current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _ApplicationChain.Contract.Flush(&_ApplicationChain.TransactOpts, current_validators, blockNumber, blockHash)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_ApplicationChain *ApplicationChainTransactor) RemoveAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _ApplicationChain.contract.Transact(opts, "removeAdmin", admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_ApplicationChain *ApplicationChainSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _ApplicationChain.Contract.RemoveAdmin(&_ApplicationChain.TransactOpts, admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_ApplicationChain *ApplicationChainTransactorSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _ApplicationChain.Contract.RemoveAdmin(&_ApplicationChain.TransactOpts, admin)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(string genesis) returns()
func (_ApplicationChain *ApplicationChainTransactor) SetGenesisInfo(opts *bind.TransactOpts, genesis string) (*types.Transaction, error) {
	return _ApplicationChain.contract.Transact(opts, "setGenesisInfo", genesis)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(string genesis) returns()
func (_ApplicationChain *ApplicationChainSession) SetGenesisInfo(genesis string) (*types.Transaction, error) {
	return _ApplicationChain.Contract.SetGenesisInfo(&_ApplicationChain.TransactOpts, genesis)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(string genesis) returns()
func (_ApplicationChain *ApplicationChainTransactorSession) SetGenesisInfo(genesis string) (*types.Transaction, error) {
	return _ApplicationChain.Contract.SetGenesisInfo(&_ApplicationChain.TransactOpts, genesis)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(address recv, uint256 amount) returns()
func (_ApplicationChain *ApplicationChainTransactor) WithdrawFund(opts *bind.TransactOpts, recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ApplicationChain.contract.Transact(opts, "withdrawFund", recv, amount)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(address recv, uint256 amount) returns()
func (_ApplicationChain *ApplicationChainSession) WithdrawFund(recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ApplicationChain.Contract.WithdrawFund(&_ApplicationChain.TransactOpts, recv, amount)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(address recv, uint256 amount) returns()
func (_ApplicationChain *ApplicationChainTransactorSession) WithdrawFund(recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _ApplicationChain.Contract.WithdrawFund(&_ApplicationChain.TransactOpts, recv, amount)
}

// SafeMathABI is the input ABI used to generate the binding from.
const SafeMathABI = "[]"

// SafeMathBin is the compiled bytecode used for deploying new contracts.
var SafeMathBin = "0x604c602c600b82828239805160001a60731460008114601c57601e565bfe5b5030600052607381538281f30073000000000000000000000000000000000000000030146080604052600080fd00a165627a7a7230582068ef568f2cd2fe8dcaf8b32b5189f420d233be889e802e955e6b388eca36215c0029"

// DeploySafeMath deploys a new Filestorm contract, binding an instance of SafeMath to it.
func DeploySafeMath(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *SafeMath, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(SafeMathBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// SafeMath is an auto generated Go binding around an Filestorm contract.
type SafeMath struct {
	SafeMathCaller     // Read-only binding to the contract
	SafeMathTransactor // Write-only binding to the contract
	SafeMathFilterer   // Log filterer for contract events
}

// SafeMathCaller is an auto generated read-only Go binding around an Filestorm contract.
type SafeMathCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathTransactor is an auto generated write-only Go binding around an Filestorm contract.
type SafeMathTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathFilterer is an auto generated log filtering Go binding around an Filestorm contract events.
type SafeMathFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SafeMathSession is an auto generated Go binding around an Filestorm contract,
// with pre-set call and transact options.
type SafeMathSession struct {
	Contract     *SafeMath         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SafeMathCallerSession is an auto generated read-only Go binding around an Filestorm contract,
// with pre-set call options.
type SafeMathCallerSession struct {
	Contract *SafeMathCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// SafeMathTransactorSession is an auto generated write-only Go binding around an Filestorm contract,
// with pre-set transact options.
type SafeMathTransactorSession struct {
	Contract     *SafeMathTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// SafeMathRaw is an auto generated low-level Go binding around an Filestorm contract.
type SafeMathRaw struct {
	Contract *SafeMath // Generic contract binding to access the raw methods on
}

// SafeMathCallerRaw is an auto generated low-level read-only Go binding around an Filestorm contract.
type SafeMathCallerRaw struct {
	Contract *SafeMathCaller // Generic read-only contract binding to access the raw methods on
}

// SafeMathTransactorRaw is an auto generated low-level write-only Go binding around an Filestorm contract.
type SafeMathTransactorRaw struct {
	Contract *SafeMathTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSafeMath creates a new instance of SafeMath, bound to a specific deployed contract.
func NewSafeMath(address common.Address, backend bind.ContractBackend) (*SafeMath, error) {
	contract, err := bindSafeMath(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SafeMath{SafeMathCaller: SafeMathCaller{contract: contract}, SafeMathTransactor: SafeMathTransactor{contract: contract}, SafeMathFilterer: SafeMathFilterer{contract: contract}}, nil
}

// NewSafeMathCaller creates a new read-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathCaller(address common.Address, caller bind.ContractCaller) (*SafeMathCaller, error) {
	contract, err := bindSafeMath(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathCaller{contract: contract}, nil
}

// NewSafeMathTransactor creates a new write-only instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathTransactor(address common.Address, transactor bind.ContractTransactor) (*SafeMathTransactor, error) {
	contract, err := bindSafeMath(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SafeMathTransactor{contract: contract}, nil
}

// NewSafeMathFilterer creates a new log filterer instance of SafeMath, bound to a specific deployed contract.
func NewSafeMathFilterer(address common.Address, filterer bind.ContractFilterer) (*SafeMathFilterer, error) {
	contract, err := bindSafeMath(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SafeMathFilterer{contract: contract}, nil
}

// bindSafeMath binds a generic wrapper to an already deployed contract.
func bindSafeMath(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(SafeMathABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.SafeMathCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.SafeMathTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SafeMath *SafeMathCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _SafeMath.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SafeMath *SafeMathTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SafeMath *SafeMathTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SafeMath.Contract.contract.Transact(opts, method, params...)
}
