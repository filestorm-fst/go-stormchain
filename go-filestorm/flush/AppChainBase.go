// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package flush

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

// AppChainBaseABI is the input ABI used to generate the binding from.
const AppChainBaseABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"FLUSH_AMOUNT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"flushEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushList\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chainName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"admins\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushMapping\",\"outputs\":[{\"name\":\"flushId\",\"type\":\"uint256\"},{\"name\":\"validator\",\"type\":\"address\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"FOUNDATION_BLACK_HOLE_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushValidatorList\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chainId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"addFund\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastFlushedBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"balance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"current_validators\",\"type\":\"address[]\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"string\"}],\"name\":\"flush\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"genesis\",\"type\":\"string\"}],\"name\":\"setGenesisInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newEpoch\",\"type\":\"uint256\"}],\"name\":\"updateFlushEpoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"FOUNDATION_MOAC_REQUIRED_AMOUNT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"distributeGasFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"period\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"}],\"name\":\"updateChainName\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recv\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGenesisInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"uniqueId\",\"type\":\"uint256\"},{\"name\":\"blockSec\",\"type\":\"uint256\"},{\"name\":\"flushNumber\",\"type\":\"uint256\"},{\"name\":\"initial_validators\",\"type\":\"address[]\"},{\"name\":\"totalSupply\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"

// AppChainBaseBin is the compiled bytecode used for deploying new contracts.
var AppChainBaseBin = "0x6080604081905260006003819055600755678ac7230489e80000600f5566b1a2bc2ec5000060105560118054600160a060020a0319167348328afc8dd45c1c252e7e883fc89bd17ddee7c0179055620014b23881900390819083398101604090815281516020830151918301516060840151608085015160a0860151600f549487019693949293919091019160009081903410156200012557604080517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f4e6f7420456e6f756768204d4f414320746f20437265617465204170706c696360448201527f6174696f6e20436861696e000000000000000000000000000000000000000000606482015290519081900360840190fd5b60008054600160a060020a0319163317905587516200014c9060029060208b0190620002d0565b505050600485905560058490556006839055346003556009805460ff191690556000808052600a6020527f13da86008ba1c6922daee3e07db95305ef49ebced9f5467a0b8613fcc6b343e38190557f13da86008ba1c6922daee3e07db95305ef49ebced9f5467a0b8613fcc6b343e48054600160a060020a03191633179055805b83518110156200023a576000828152600c602052604090208451859083908110620001f457fe5b6020908102919091018101518254600180820185556000948552929093209092018054600160a060020a031916600160a060020a039093169290921790915501620001cd565b6000828152600a6020818152604080842060016002820155815180840192839052858152948790529290915291516200027a9260039092019190620002d0565b5050600b805460018181019092557f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db90191909155600e919091553360009081526020829052604090205550620003759350505050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200031357805160ff191683800117855562000343565b8280016001018555821562000343579182015b828111156200034357825182559160200191906001019062000326565b506200035192915062000355565b5090565b6200037291905b808211156200035157600081556001016200035c565b90565b61112d80620003856000396000f3006080604052600436106101325763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166308925f61811461013c578063090aaae5146101635780631785f53c1461017857806318c152111461019b5780631c93b03a146101b3578063429b62e51461023d5780634c2f86191461025e5780635a35b2fc14610315578063635cf62a1461034657806370480275146103615780639a8a059214610382578063a2f09dfa14610397578063a33f5c4a1461039f578063b69ef8a8146103b4578063d66cfad5146103c9578063d9b11dc814610462578063dd77088e146104bb578063dfe89af2146104d3578063e15f7f93146104e8578063ef78d4fd146104fd578063f095ee6414610512578063f7c8d2211461056b578063fd8735a51461058f575b6003805434019055005b34801561014857600080fd5b506101516105a4565b60408051918252519081900360200190f35b34801561016f57600080fd5b506101516105aa565b34801561018457600080fd5b50610199600160a060020a03600435166105b0565b005b3480156101a757600080fd5b506101516004356106b9565b3480156101bf57600080fd5b506101c86106d8565b6040805160208082528351818301528351919283929083019185019080838360005b838110156102025781810151838201526020016101ea565b50505050905090810190601f16801561022f5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b34801561024957600080fd5b50610151600160a060020a0360043516610763565b34801561026a57600080fd5b50610276600435610775565b6040518085815260200184600160a060020a0316600160a060020a0316815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b838110156102d75781810151838201526020016102bf565b50505050905090810190601f1680156103045780820380516001836020036101000a031916815260200191505b509550505050505060405180910390f35b34801561032157600080fd5b5061032a610833565b60408051600160a060020a039092168252519081900360200190f35b34801561035257600080fd5b5061032a600435602435610842565b34801561036d57600080fd5b50610199600160a060020a0360043516610879565b34801561038e57600080fd5b50610151610924565b61019961092a565b3480156103ab57600080fd5b50610151610934565b3480156103c057600080fd5b5061015161093a565b3480156103d557600080fd5b50604080516020600480358082013583810280860185019096528085526101999536959394602494938501929182918501908490808284375050604080516020601f818a01358b0180359182018390048302840183018552818452989b8a359b909a9099940197509195509182019350915081908401838280828437509497506109409650505050505050565b34801561046e57600080fd5b506040805160206004803580820135601f8101849004840285018401909552848452610199943694929360249392840191908190840183828082843750949750610b4f9650505050505050565b3480156104c757600080fd5b50610199600435610c82565b3480156104df57600080fd5b50610151610d95565b3480156104f457600080fd5b50610199610d9b565b34801561050957600080fd5b50610151610ebc565b34801561051e57600080fd5b506040805160206004803580820135601f8101849004840285018401909552848452610199943694929360249392840191908190840183828082843750949750610ec29650505050505050565b34801561057757600080fd5b50610199600160a060020a0360043516602435610f67565b34801561059b57600080fd5b506101c8610fd2565b60105481565b60065481565b336000908152600160208190526040909120541461063e576040805160e560020a62461bcd02815260206004820152602260248201527f4f6e6c792041646d696e732043616e2041646420416e6f746865722041646d6960448201527f6e2e000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a03811633141561069f576040805160e560020a62461bcd02815260206004820152601a60248201527f41646d696e732043616e6e6f742052656d6f76652053656c662e000000000000604482015290519081900360640190fd5b600160a060020a0316600090815260016020526040812055565b600b8054829081106106c757fe5b600091825260209091200154905081565b6002805460408051602060018416156101000260001901909316849004601f8101849004840282018401909252818152929183018282801561075b5780601f106107305761010080835404028352916020019161075b565b820191906000526020600020905b81548152906001019060200180831161073e57829003601f168201915b505050505081565b60016020526000908152604090205481565b600a60209081526000918252604091829020805460018083015460028085015460038601805489516101009682161596909602600019011692909204601f81018890048802850188019098528784529396600160a060020a03909216959394939091908301828280156108295780601f106107fe57610100808354040283529160200191610829565b820191906000526020600020905b81548152906001019060200180831161080c57829003601f168201915b5050505050905084565b601154600160a060020a031681565b600c6020528160005260406000208181548110151561085d57fe5b600091825260209091200154600160a060020a03169150829050565b3360009081526001602081905260409091205414610907576040805160e560020a62461bcd02815260206004820152602260248201527f4f6e6c792041646d696e732043616e2041646420416e6f746865722041646d6960448201527f6e2e000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b600160a060020a0316600090815260016020819052604090912055565b60045481565b6003805434019055565b60075481565b60035481565b60008060008060035411151561095557600080fd5b600b54600019019250600091505b6000838152600c6020526040902054821015610b47576000838152600c6020526040902080543391908490811061099657fe5b600091825260209091200154600160a060020a031614806109c7575033600090815260016020819052604090912054145b15610b3c575060019182016000818152600a60205260408120828155909301805473ffffffffffffffffffffffffffffffffffffffff191633179055915b8551811015610a7c576000838152600c602052604090208651879083908110610a2a57fe5b602090810291909101810151825460018082018555600094855292909320909201805473ffffffffffffffffffffffffffffffffffffffff1916600160a060020a039093169290921790915501610a05565b6000838152600a60209081526040909120600281018790558551610aa892600390920191870190611069565b50600b805460018101825560009182527f0175b7a638427703f0dbe7bb9bbf987a2551717b34e79f33b5b1008d1fa01db901849055601054600380548290039055601154604051600160a060020a039091169282156108fc02929190818181858888f19350505050158015610b21573d6000803e3d6000fd5b506000838152600a6020526040902060020154600755610b47565b600190910190610963565b505050505050565b3360009081526001602081905260409091205414610bdd576040805160e560020a62461bcd02815260206004820152602160248201527f4f6e6c792041646d696e732043616e205365742047656e6573697320496e666f60448201527f2e00000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b60095460ff1615610c5e576040805160e560020a62461bcd02815260206004820152602260248201527f47656e6573697320496e666f2048617320416c7265616479204265656e20536560448201527f742e000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b8051610c71906008906020840190611069565b50506009805460ff19166001179055565b3360009081526001602081905260409091205414610d10576040805160e560020a62461bcd02815260206004820152602360248201527f4f6e6c792041646d696e732043616e2055706461746520466c7573682045706f60448201527f63682e0000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b610168811015610d90576040805160e560020a62461bcd02815260206004820152603160248201527f466c7573682045706f6368204d75737420626520457175616c20746f206f722060448201527f47726561746572207468616e203336302e000000000000000000000000000000606482015290519081900360840190fd5b600655565b600f5481565b3360009081526001602081905260408220548291829114610e2c576040805160e560020a62461bcd02815260206004820152602360248201527f4f6e6c792041646d696e732043616e204469737472696275746520476173204660448201527f65652e0000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b506000915067016345785d8a00009050815b6000838152600c6020526040902054811015610eb7576000838152600c60205260409020805482908110610e6e57fe5b6000918252602082200154604051600160a060020a039091169184156108fc02918591818181858888f19350505050158015610eae573d6000803e3d6000fd5b50600101610e3e565b505050565b60055481565b3360009081526001602081905260409091205414610f50576040805160e560020a62461bcd02815260206004820152602260248201527f4f6e6c792041646d696e732043616e2055706461746520436861696e204e616d60448201527f652e000000000000000000000000000000000000000000000000000000000000606482015290519081900360840190fd5b8051610f63906002906020840190611069565b5050565b3360009081526001602081905260409091205414610f8457600080fd5b600354811115610f9357600080fd5b600380548290039055604051600160a060020a0383169082156108fc029083906000818181858888f19350505050158015610eb7573d6000803e3d6000fd5b60088054604080516020601f600260001961010060018816150201909516949094049384018190048102820181019092528281526060939092909183018282801561105e5780601f106110335761010080835404028352916020019161105e565b820191906000526020600020905b81548152906001019060200180831161104157829003601f168201915b505050505090505b90565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106110aa57805160ff19168380011785556110d7565b828001600101855582156110d7579182015b828111156110d75782518255916020019190600101906110bc565b506110e39291506110e7565b5090565b61106691905b808211156110e357600081556001016110ed5600a165627a7a72305820584e94e8fa0bf026575ff7ec92754cf4611eaf0a8ef955fcddff936d2635e83c0029"

// DeployAppChainBase deploys a new Filestorm contract, binding an instance of AppChainBase to it.
func DeployAppChainBase(auth *bind.TransactOpts, backend bind.ContractBackend, name string, uniqueId *big.Int, blockSec *big.Int, flushNumber *big.Int, initial_validators []common.Address, totalSupply *big.Int) (common.Address, *types.Transaction, *AppChainBase, error) {
	parsed, err := abi.JSON(strings.NewReader(AppChainBaseABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}

	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AppChainBaseBin), backend, name, uniqueId, blockSec, flushNumber, initial_validators, totalSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AppChainBase{AppChainBaseCaller: AppChainBaseCaller{contract: contract}, AppChainBaseTransactor: AppChainBaseTransactor{contract: contract}, AppChainBaseFilterer: AppChainBaseFilterer{contract: contract}}, nil
}

// AppChainBase is an auto generated Go binding around an Filestorm contract.
type AppChainBase struct {
	AppChainBaseCaller     // Read-only binding to the contract
	AppChainBaseTransactor // Write-only binding to the contract
	AppChainBaseFilterer   // Log filterer for contract events
}

// AppChainBaseCaller is an auto generated read-only Go binding around an Filestorm contract.
type AppChainBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainBaseTransactor is an auto generated write-only Go binding around an Filestorm contract.
type AppChainBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainBaseFilterer is an auto generated log filtering Go binding around an Filestorm contract events.
type AppChainBaseFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainBaseSession is an auto generated Go binding around an Filestorm contract,
// with pre-set call and transact options.
type AppChainBaseSession struct {
	Contract     *AppChainBase     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AppChainBaseCallerSession is an auto generated read-only Go binding around an Filestorm contract,
// with pre-set call options.
type AppChainBaseCallerSession struct {
	Contract *AppChainBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// AppChainBaseTransactorSession is an auto generated write-only Go binding around an Filestorm contract,
// with pre-set transact options.
type AppChainBaseTransactorSession struct {
	Contract     *AppChainBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// AppChainBaseRaw is an auto generated low-level Go binding around an Filestorm contract.
type AppChainBaseRaw struct {
	Contract *AppChainBase // Generic contract binding to access the raw methods on
}

// AppChainBaseCallerRaw is an auto generated low-level read-only Go binding around an Filestorm contract.
type AppChainBaseCallerRaw struct {
	Contract *AppChainBaseCaller // Generic read-only contract binding to access the raw methods on
}

// AppChainBaseTransactorRaw is an auto generated low-level write-only Go binding around an Filestorm contract.
type AppChainBaseTransactorRaw struct {
	Contract *AppChainBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAppChainBase creates a new instance of AppChainBase, bound to a specific deployed contract.
func NewAppChainBase(address common.Address, backend bind.ContractBackend) (*AppChainBase, error) {
	contract, err := bindAppChainBase(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AppChainBase{AppChainBaseCaller: AppChainBaseCaller{contract: contract}, AppChainBaseTransactor: AppChainBaseTransactor{contract: contract}, AppChainBaseFilterer: AppChainBaseFilterer{contract: contract}}, nil
}

// NewAppChainBaseCaller creates a new read-only instance of AppChainBase, bound to a specific deployed contract.
func NewAppChainBaseCaller(address common.Address, caller bind.ContractCaller) (*AppChainBaseCaller, error) {
	contract, err := bindAppChainBase(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AppChainBaseCaller{contract: contract}, nil
}

// NewAppChainBaseTransactor creates a new write-only instance of AppChainBase, bound to a specific deployed contract.
func NewAppChainBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*AppChainBaseTransactor, error) {
	contract, err := bindAppChainBase(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AppChainBaseTransactor{contract: contract}, nil
}

// NewAppChainBaseFilterer creates a new log filterer instance of AppChainBase, bound to a specific deployed contract.
func NewAppChainBaseFilterer(address common.Address, filterer bind.ContractFilterer) (*AppChainBaseFilterer, error) {
	contract, err := bindAppChainBase(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AppChainBaseFilterer{contract: contract}, nil
}

// bindAppChainBase binds a generic wrapper to an already deployed contract.
func bindAppChainBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AppChainBaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AppChainBase *AppChainBaseRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AppChainBase.Contract.AppChainBaseCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AppChainBase *AppChainBaseRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainBase.Contract.AppChainBaseTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AppChainBase *AppChainBaseRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AppChainBase.Contract.AppChainBaseTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AppChainBase *AppChainBaseCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _AppChainBase.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AppChainBase *AppChainBaseTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainBase.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AppChainBase *AppChainBaseTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AppChainBase.Contract.contract.Transact(opts, method, params...)
}

// FLUSHAMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) FLUSHAMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "FLUSH_AMOUNT")
	return *ret0, err
}

// FLUSHAMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) FLUSHAMOUNT() (*big.Int, error) {
	return _AppChainBase.Contract.FLUSHAMOUNT(&_AppChainBase.CallOpts)
}

// FLUSHAMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) FLUSHAMOUNT() (*big.Int, error) {
	return _AppChainBase.Contract.FLUSHAMOUNT(&_AppChainBase.CallOpts)
}

// FOUNDATIONBLACKHOLEADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_AppChainBase *AppChainBaseCaller) FOUNDATIONBLACKHOLEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "FOUNDATION_BLACK_HOLE_ADDRESS")
	return *ret0, err
}

// FOUNDATIONBLACKHOLEADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_AppChainBase *AppChainBaseSession) FOUNDATIONBLACKHOLEADDRESS() (common.Address, error) {
	return _AppChainBase.Contract.FOUNDATIONBLACKHOLEADDRESS(&_AppChainBase.CallOpts)
}

// FOUNDATIONBLACKHOLEADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_AppChainBase *AppChainBaseCallerSession) FOUNDATIONBLACKHOLEADDRESS() (common.Address, error) {
	return _AppChainBase.Contract.FOUNDATIONBLACKHOLEADDRESS(&_AppChainBase.CallOpts)
}

// FOUNDATIONMOACREQUIREDAMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) FOUNDATIONMOACREQUIREDAMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "FOUNDATION_MOAC_REQUIRED_AMOUNT")
	return *ret0, err
}

// FOUNDATIONMOACREQUIREDAMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) FOUNDATIONMOACREQUIREDAMOUNT() (*big.Int, error) {
	return _AppChainBase.Contract.FOUNDATIONMOACREQUIREDAMOUNT(&_AppChainBase.CallOpts)
}

// FOUNDATIONMOACREQUIREDAMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) FOUNDATIONMOACREQUIREDAMOUNT() (*big.Int, error) {
	return _AppChainBase.Contract.FOUNDATIONMOACREQUIREDAMOUNT(&_AppChainBase.CallOpts)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) Admins(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "admins", arg0)
	return *ret0, err
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) Admins(arg0 common.Address) (*big.Int, error) {
	return _AppChainBase.Contract.Admins(&_AppChainBase.CallOpts, arg0)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins(address ) constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) Admins(arg0 common.Address) (*big.Int, error) {
	return _AppChainBase.Contract.Admins(&_AppChainBase.CallOpts, arg0)
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) Balance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "balance")
	return *ret0, err
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) Balance() (*big.Int, error) {
	return _AppChainBase.Contract.Balance(&_AppChainBase.CallOpts)
}

// Balance is a free data retrieval call binding the contract method 0xb69ef8a8.
//
// Solidity: function balance() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) Balance() (*big.Int, error) {
	return _AppChainBase.Contract.Balance(&_AppChainBase.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) ChainId(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "chainId")
	return *ret0, err
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) ChainId() (*big.Int, error) {
	return _AppChainBase.Contract.ChainId(&_AppChainBase.CallOpts)
}

// ChainId is a free data retrieval call binding the contract method 0x9a8a0592.
//
// Solidity: function chainId() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) ChainId() (*big.Int, error) {
	return _AppChainBase.Contract.ChainId(&_AppChainBase.CallOpts)
}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() constant returns(string)
func (_AppChainBase *AppChainBaseCaller) ChainName(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "chainName")
	return *ret0, err
}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() constant returns(string)
func (_AppChainBase *AppChainBaseSession) ChainName() (string, error) {
	return _AppChainBase.Contract.ChainName(&_AppChainBase.CallOpts)
}

// ChainName is a free data retrieval call binding the contract method 0x1c93b03a.
//
// Solidity: function chainName() constant returns(string)
func (_AppChainBase *AppChainBaseCallerSession) ChainName() (string, error) {
	return _AppChainBase.Contract.ChainName(&_AppChainBase.CallOpts)
}

// FlushEpoch is a free data retrieval call binding the contract method 0x090aaae5.
//
// Solidity: function flushEpoch() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) FlushEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "flushEpoch")
	return *ret0, err
}

// FlushEpoch is a free data retrieval call binding the contract method 0x090aaae5.
//
// Solidity: function flushEpoch() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) FlushEpoch() (*big.Int, error) {
	return _AppChainBase.Contract.FlushEpoch(&_AppChainBase.CallOpts)
}

// FlushEpoch is a free data retrieval call binding the contract method 0x090aaae5.
//
// Solidity: function flushEpoch() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) FlushEpoch() (*big.Int, error) {
	return _AppChainBase.Contract.FlushEpoch(&_AppChainBase.CallOpts)
}

// FlushList is a free data retrieval call binding the contract method 0x18c15211.
//
// Solidity: function flushList(uint256 ) constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) FlushList(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "flushList", arg0)
	return *ret0, err
}

// FlushList is a free data retrieval call binding the contract method 0x18c15211.
//
// Solidity: function flushList(uint256 ) constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) FlushList(arg0 *big.Int) (*big.Int, error) {
	return _AppChainBase.Contract.FlushList(&_AppChainBase.CallOpts, arg0)
}

// FlushList is a free data retrieval call binding the contract method 0x18c15211.
//
// Solidity: function flushList(uint256 ) constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) FlushList(arg0 *big.Int) (*big.Int, error) {
	return _AppChainBase.Contract.FlushList(&_AppChainBase.CallOpts, arg0)
}

// FlushMapping is a free data retrieval call binding the contract method 0x4c2f8619.
//
// Solidity: function flushMapping(uint256 ) constant returns(uint256 flushId, address validator, uint256 blockNumber, string blockHash)
func (_AppChainBase *AppChainBaseCaller) FlushMapping(opts *bind.CallOpts, arg0 *big.Int) (struct {
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
	err := _AppChainBase.contract.Call(opts, out, "flushMapping", arg0)
	return *ret, err
}

// FlushMapping is a free data retrieval call binding the contract method 0x4c2f8619.
//
// Solidity: function flushMapping(uint256 ) constant returns(uint256 flushId, address validator, uint256 blockNumber, string blockHash)
func (_AppChainBase *AppChainBaseSession) FlushMapping(arg0 *big.Int) (struct {
	FlushId     *big.Int
	Validator   common.Address
	BlockNumber *big.Int
	BlockHash   string
}, error) {
	return _AppChainBase.Contract.FlushMapping(&_AppChainBase.CallOpts, arg0)
}

// FlushMapping is a free data retrieval call binding the contract method 0x4c2f8619.
//
// Solidity: function flushMapping(uint256 ) constant returns(uint256 flushId, address validator, uint256 blockNumber, string blockHash)
func (_AppChainBase *AppChainBaseCallerSession) FlushMapping(arg0 *big.Int) (struct {
	FlushId     *big.Int
	Validator   common.Address
	BlockNumber *big.Int
	BlockHash   string
}, error) {
	return _AppChainBase.Contract.FlushMapping(&_AppChainBase.CallOpts, arg0)
}

// FlushValidatorList is a free data retrieval call binding the contract method 0x635cf62a.
//
// Solidity: function flushValidatorList(uint256 , uint256 ) constant returns(address)
func (_AppChainBase *AppChainBaseCaller) FlushValidatorList(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "flushValidatorList", arg0, arg1)
	return *ret0, err
}

// FlushValidatorList is a free data retrieval call binding the contract method 0x635cf62a.
//
// Solidity: function flushValidatorList(uint256 , uint256 ) constant returns(address)
func (_AppChainBase *AppChainBaseSession) FlushValidatorList(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _AppChainBase.Contract.FlushValidatorList(&_AppChainBase.CallOpts, arg0, arg1)
}

// FlushValidatorList is a free data retrieval call binding the contract method 0x635cf62a.
//
// Solidity: function flushValidatorList(uint256 , uint256 ) constant returns(address)
func (_AppChainBase *AppChainBaseCallerSession) FlushValidatorList(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _AppChainBase.Contract.FlushValidatorList(&_AppChainBase.CallOpts, arg0, arg1)
}

// GetGenesisInfo is a free data retrieval call binding the contract method 0xfd8735a5.
//
// Solidity: function getGenesisInfo() constant returns(string)
func (_AppChainBase *AppChainBaseCaller) GetGenesisInfo(opts *bind.CallOpts) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "getGenesisInfo")
	return *ret0, err
}

// GetGenesisInfo is a free data retrieval call binding the contract method 0xfd8735a5.
//
// Solidity: function getGenesisInfo() constant returns(string)
func (_AppChainBase *AppChainBaseSession) GetGenesisInfo() (string, error) {
	return _AppChainBase.Contract.GetGenesisInfo(&_AppChainBase.CallOpts)
}

// GetGenesisInfo is a free data retrieval call binding the contract method 0xfd8735a5.
//
// Solidity: function getGenesisInfo() constant returns(string)
func (_AppChainBase *AppChainBaseCallerSession) GetGenesisInfo() (string, error) {
	return _AppChainBase.Contract.GetGenesisInfo(&_AppChainBase.CallOpts)
}

// LastFlushedBlock is a free data retrieval call binding the contract method 0xa33f5c4a.
//
// Solidity: function lastFlushedBlock() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) LastFlushedBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "lastFlushedBlock")
	return *ret0, err
}

// LastFlushedBlock is a free data retrieval call binding the contract method 0xa33f5c4a.
//
// Solidity: function lastFlushedBlock() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) LastFlushedBlock() (*big.Int, error) {
	return _AppChainBase.Contract.LastFlushedBlock(&_AppChainBase.CallOpts)
}

// LastFlushedBlock is a free data retrieval call binding the contract method 0xa33f5c4a.
//
// Solidity: function lastFlushedBlock() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) LastFlushedBlock() (*big.Int, error) {
	return _AppChainBase.Contract.LastFlushedBlock(&_AppChainBase.CallOpts)
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) Period(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "period")
	return *ret0, err
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) Period() (*big.Int, error) {
	return _AppChainBase.Contract.Period(&_AppChainBase.CallOpts)
}

// Period is a free data retrieval call binding the contract method 0xef78d4fd.
//
// Solidity: function period() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) Period() (*big.Int, error) {
	return _AppChainBase.Contract.Period(&_AppChainBase.CallOpts)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_AppChainBase *AppChainBaseTransactor) AddAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "addAdmin", admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_AppChainBase *AppChainBaseSession) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.Contract.AddAdmin(&_AppChainBase.TransactOpts, admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(address admin) returns()
func (_AppChainBase *AppChainBaseTransactorSession) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.Contract.AddAdmin(&_AppChainBase.TransactOpts, admin)
}

// AddFund is a paid mutator transaction binding the contract method 0xa2f09dfa.
//
// Solidity: function addFund() returns()
func (_AppChainBase *AppChainBaseTransactor) AddFund(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "addFund")
}

// AddFund is a paid mutator transaction binding the contract method 0xa2f09dfa.
//
// Solidity: function addFund() returns()
func (_AppChainBase *AppChainBaseSession) AddFund() (*types.Transaction, error) {
	return _AppChainBase.Contract.AddFund(&_AppChainBase.TransactOpts)
}

// AddFund is a paid mutator transaction binding the contract method 0xa2f09dfa.
//
// Solidity: function addFund() returns()
func (_AppChainBase *AppChainBaseTransactorSession) AddFund() (*types.Transaction, error) {
	return _AppChainBase.Contract.AddFund(&_AppChainBase.TransactOpts)
}

// DistributeGasFee is a paid mutator transaction binding the contract method 0xe15f7f93.
//
// Solidity: function distributeGasFee() returns()
func (_AppChainBase *AppChainBaseTransactor) DistributeGasFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "distributeGasFee")
}

// DistributeGasFee is a paid mutator transaction binding the contract method 0xe15f7f93.
//
// Solidity: function distributeGasFee() returns()
func (_AppChainBase *AppChainBaseSession) DistributeGasFee() (*types.Transaction, error) {
	return _AppChainBase.Contract.DistributeGasFee(&_AppChainBase.TransactOpts)
}

// DistributeGasFee is a paid mutator transaction binding the contract method 0xe15f7f93.
//
// Solidity: function distributeGasFee() returns()
func (_AppChainBase *AppChainBaseTransactorSession) DistributeGasFee() (*types.Transaction, error) {
	return _AppChainBase.Contract.DistributeGasFee(&_AppChainBase.TransactOpts)
}

// Flush is a paid mutator transaction binding the contract method 0xd66cfad5.
//
// Solidity: function flush(address[] current_validators, uint256 blockNumber, string blockHash) returns()
func (_AppChainBase *AppChainBaseTransactor) Flush(opts *bind.TransactOpts, current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "flush", current_validators, blockNumber, blockHash)
}

// Flush is a paid mutator transaction binding the contract method 0xd66cfad5.
//
// Solidity: function flush(address[] current_validators, uint256 blockNumber, string blockHash) returns()
func (_AppChainBase *AppChainBaseSession) Flush(current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _AppChainBase.Contract.Flush(&_AppChainBase.TransactOpts, current_validators, blockNumber, blockHash)
}

// Flush is a paid mutator transaction binding the contract method 0xd66cfad5.
//
// Solidity: function flush(address[] current_validators, uint256 blockNumber, string blockHash) returns()
func (_AppChainBase *AppChainBaseTransactorSession) Flush(current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _AppChainBase.Contract.Flush(&_AppChainBase.TransactOpts, current_validators, blockNumber, blockHash)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_AppChainBase *AppChainBaseTransactor) RemoveAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "removeAdmin", admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_AppChainBase *AppChainBaseSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.Contract.RemoveAdmin(&_AppChainBase.TransactOpts, admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(address admin) returns()
func (_AppChainBase *AppChainBaseTransactorSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.Contract.RemoveAdmin(&_AppChainBase.TransactOpts, admin)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(string genesis) returns()
func (_AppChainBase *AppChainBaseTransactor) SetGenesisInfo(opts *bind.TransactOpts, genesis string) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "setGenesisInfo", genesis)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(string genesis) returns()
func (_AppChainBase *AppChainBaseSession) SetGenesisInfo(genesis string) (*types.Transaction, error) {
	return _AppChainBase.Contract.SetGenesisInfo(&_AppChainBase.TransactOpts, genesis)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(string genesis) returns()
func (_AppChainBase *AppChainBaseTransactorSession) SetGenesisInfo(genesis string) (*types.Transaction, error) {
	return _AppChainBase.Contract.SetGenesisInfo(&_AppChainBase.TransactOpts, genesis)
}

// UpdateChainName is a paid mutator transaction binding the contract method 0xf095ee64.
//
// Solidity: function updateChainName(string name) returns()
func (_AppChainBase *AppChainBaseTransactor) UpdateChainName(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "updateChainName", name)
}

// UpdateChainName is a paid mutator transaction binding the contract method 0xf095ee64.
//
// Solidity: function updateChainName(string name) returns()
func (_AppChainBase *AppChainBaseSession) UpdateChainName(name string) (*types.Transaction, error) {
	return _AppChainBase.Contract.UpdateChainName(&_AppChainBase.TransactOpts, name)
}

// UpdateChainName is a paid mutator transaction binding the contract method 0xf095ee64.
//
// Solidity: function updateChainName(string name) returns()
func (_AppChainBase *AppChainBaseTransactorSession) UpdateChainName(name string) (*types.Transaction, error) {
	return _AppChainBase.Contract.UpdateChainName(&_AppChainBase.TransactOpts, name)
}

// UpdateFlushEpoch is a paid mutator transaction binding the contract method 0xdd77088e.
//
// Solidity: function updateFlushEpoch(uint256 newEpoch) returns()
func (_AppChainBase *AppChainBaseTransactor) UpdateFlushEpoch(opts *bind.TransactOpts, newEpoch *big.Int) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "updateFlushEpoch", newEpoch)
}

// UpdateFlushEpoch is a paid mutator transaction binding the contract method 0xdd77088e.
//
// Solidity: function updateFlushEpoch(uint256 newEpoch) returns()
func (_AppChainBase *AppChainBaseSession) UpdateFlushEpoch(newEpoch *big.Int) (*types.Transaction, error) {
	return _AppChainBase.Contract.UpdateFlushEpoch(&_AppChainBase.TransactOpts, newEpoch)
}

// UpdateFlushEpoch is a paid mutator transaction binding the contract method 0xdd77088e.
//
// Solidity: function updateFlushEpoch(uint256 newEpoch) returns()
func (_AppChainBase *AppChainBaseTransactorSession) UpdateFlushEpoch(newEpoch *big.Int) (*types.Transaction, error) {
	return _AppChainBase.Contract.UpdateFlushEpoch(&_AppChainBase.TransactOpts, newEpoch)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(address recv, uint256 amount) returns()
func (_AppChainBase *AppChainBaseTransactor) WithdrawFund(opts *bind.TransactOpts, recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "withdrawFund", recv, amount)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(address recv, uint256 amount) returns()
func (_AppChainBase *AppChainBaseSession) WithdrawFund(recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AppChainBase.Contract.WithdrawFund(&_AppChainBase.TransactOpts, recv, amount)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(address recv, uint256 amount) returns()
func (_AppChainBase *AppChainBaseTransactorSession) WithdrawFund(recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AppChainBase.Contract.WithdrawFund(&_AppChainBase.TransactOpts, recv, amount)
}
