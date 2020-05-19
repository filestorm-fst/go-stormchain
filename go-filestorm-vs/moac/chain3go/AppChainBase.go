// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package chain3go

import (
	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts/abi"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts/abi/bind"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"math/big"
	"strings"
)

// AppChainBaseABI is the input ABI used to generate the binding from.
const AppChainBaseABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"FLUSH_AMOUNT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"flushEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"removeAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushList\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chainName\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"admins\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushMapping\",\"outputs\":[{\"name\":\"flushId\",\"type\":\"uint256\"},{\"name\":\"validator\",\"type\":\"address\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"FOUNDATION_BLACK_HOLE_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"flushValidatorList\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"admin\",\"type\":\"address\"}],\"name\":\"addAdmin\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chainId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"addFund\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastFlushedBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"balance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"current_validators\",\"type\":\"address[]\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"blockHash\",\"type\":\"string\"}],\"name\":\"flush\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"genesis\",\"type\":\"string\"}],\"name\":\"setGenesisInfo\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newEpoch\",\"type\":\"uint256\"}],\"name\":\"updateFlushEpoch\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"FOUNDATION_MOAC_REQUIRED_AMOUNT\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"distributeGasFee\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"period\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"name\",\"type\":\"string\"}],\"name\":\"updateChainName\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recv\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"withdrawFund\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getGenesisInfo\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"uniqueId\",\"type\":\"uint256\"},{\"name\":\"blockSec\",\"type\":\"uint256\"},{\"name\":\"flushNumber\",\"type\":\"uint256\"},{\"name\":\"initial_validators\",\"type\":\"address[]\"},{\"name\":\"totalSupply\",\"type\":\"uint256\"}],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"}]"

// AppChainBaseBin is the compiled bytecode used for deploying new contracts.
const AppChainBaseBin = `608060405260006003556000600755678ac7230489e80000600f5566b1a2bc2ec500006010557348328afc8dd45c1c252e7e883fc89bd17ddee7c0601160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060405162001e1938038062001e19833981018060405281019080805182019291906020018051906020019092919080519060200190929190805190602001909291908051820192919060200180519060200190929190505050600080600f54341015151562000177576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252602b8152602001807f4e6f7420456e6f756768204d4f414320746f20437265617465204170706c696381526020017f6174696f6e20436861696e00000000000000000000000000000000000000000081525060400191505060405180910390fd5b336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508760029080519060200190620001cf92919062000405565b50866004819055508560058190555084600681905550346003819055506000600960006101000a81548160ff0219169083151502179055506000915081600a60008481526020019081526020016000206000018190555033600a600084815260200190815260200160002060010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600090505b83518110156200032757600c60008381526020019081526020016000208482815181101515620002ac57fe5b9060200190602002015190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555050808060010191505062000280565b6001600a6000848152602001908152602001600020600201819055506020604051908101604052806000815250600a600084815260200190815260200160002060030190805190602001906200037f92919062000405565b50600b82908060018154018082558091505090600182039060005260206000200160009091929091909150555082600e8190555060018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020819055505050505050505050620004b4565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106200044857805160ff191683800117855562000479565b8280016001018555821562000479579182015b82811115620004785782518255916020019190600101906200045b565b5b5090506200048891906200048c565b5090565b620004b191905b80821115620004ad57600081600090555060010162000493565b5090565b90565b61195580620004c46000396000f300608060405260043610610133576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806308925f6114610145578063090aaae5146101705780631785f53c1461019b57806318c15211146101de5780631c93b03a1461021f578063429b62e5146102af5780634c2f8619146103065780635a35b2fc146103ed578063635cf62a1461044457806370480275146104bb5780639a8a0592146104fe578063a2f09dfa14610529578063a33f5c4a14610533578063b69ef8a81461055e578063d66cfad514610589578063d9b11dc81461063f578063dd77088e146106a8578063dfe89af2146106d5578063e15f7f9314610700578063ef78d4fd14610717578063f095ee6414610742578063f7c8d221146107ab578063fd8735a5146107f8575b34600360008282540192505081905550005b34801561015157600080fd5b5061015a610888565b6040518082815260200191505060405180910390f35b34801561017c57600080fd5b5061018561088e565b6040518082815260200191505060405180910390f35b3480156101a757600080fd5b506101dc600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610894565b005b3480156101ea57600080fd5b5061020960048036038101908080359060200190929190505050610a5c565b6040518082815260200191505060405180910390f35b34801561022b57600080fd5b50610234610a7f565b6040518080602001828103825283818151815260200191508051906020019080838360005b83811015610274578082015181840152602081019050610259565b50505050905090810190601f1680156102a15780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b3480156102bb57600080fd5b506102f0600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610b1d565b6040518082815260200191505060405180910390f35b34801561031257600080fd5b5061033160048036038101908080359060200190929190505050610b35565b604051808581526020018473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200183815260200180602001828103825283818151815260200191508051906020019080838360005b838110156103af578082015181840152602081019050610394565b50505050905090810190601f1680156103dc5780820380516001836020036101000a031916815260200191505b509550505050505060405180910390f35b3480156103f957600080fd5b50610402610c1d565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b34801561045057600080fd5b506104796004803603810190808035906020019092919080359060200190929190505050610c43565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b3480156104c757600080fd5b506104fc600480360381019080803573ffffffffffffffffffffffffffffffffffffffff169060200190929190505050610c90565b005b34801561050a57600080fd5b50610513610db3565b6040518082815260200191505060405180910390f35b610531610db9565b005b34801561053f57600080fd5b50610548610dcb565b6040518082815260200191505060405180910390f35b34801561056a57600080fd5b50610573610dd1565b6040518082815260200191505060405180910390f35b34801561059557600080fd5b5061063d6004803603810190808035906020019082018035906020019080806020026020016040519081016040528093929190818152602001838360200280828437820191505050505050919291929080359060200190929190803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050610dd7565b005b34801561064b57600080fd5b506106a6600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050611129565b005b3480156106b457600080fd5b506106d3600480360381019080803590602001909291905050506112eb565b005b3480156106e157600080fd5b506106ea611471565b6040518082815260200191505060405180910390f35b34801561070c57600080fd5b50610715611477565b005b34801561072357600080fd5b5061072c61162d565b6040518082815260200191505060405180910390f35b34801561074e57600080fd5b506107a9600480360381019080803590602001908201803590602001908080601f0160208091040260200160405190810160405280939291908181526020018383808284378201915050505050509192919290505050611633565b005b3480156107b757600080fd5b506107f6600480360381019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919080359060200190929190505050611729565b005b34801561080457600080fd5b5061080d6117e2565b6040518080602001828103825283818151815260200191508051906020019080838360005b8381101561084d578082015181840152602081019050610832565b50505050905090810190601f16801561087a5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b60105481565b60065481565b60018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054141515610970576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f4f6e6c792041646d696e732043616e2041646420416e6f746865722041646d6981526020017f6e2e00000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b3373ffffffffffffffffffffffffffffffffffffffff168173ffffffffffffffffffffffffffffffffffffffff1614151515610a14576040517f08c379a000000000000000000000000000000000000000000000000000000000815260040180806020018281038252601a8152602001807f41646d696e732043616e6e6f742052656d6f76652053656c662e00000000000081525060200191505060405180910390fd5b6000600160008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555050565b600b81815481101515610a6b57fe5b906000526020600020016000915090505481565b60028054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610b155780601f10610aea57610100808354040283529160200191610b15565b820191906000526020600020905b815481529060010190602001808311610af857829003601f168201915b505050505081565b60016020528060005260406000206000915090505481565b600a6020528060005260406000206000915090508060000154908060010160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1690806002015490806003018054600181600116156101000203166002900480601f016020809104026020016040519081016040528092919081815260200182805460018160011615610100020316600290048015610c135780601f10610be857610100808354040283529160200191610c13565b820191906000526020600020905b815481529060010190602001808311610bf657829003601f168201915b5050505050905084565b601160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b600c60205281600052604060002081815481101515610c5e57fe5b906000526020600020016000915091509054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b60018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054141515610d6c576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f4f6e6c792041646d696e732043616e2041646420416e6f746865722041646d6981526020017f6e2e00000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b60018060008373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000208190555050565b60045481565b34600360008282540192505081905550565b60075481565b60035481565b600080600080600354111515610dec57600080fd5b6001600b80549050039250600091505b600c600084815260200190815260200160002080549050821015611120573373ffffffffffffffffffffffffffffffffffffffff16600c600085815260200190815260200160002083815481101515610e5157fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff161480610edc575060018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054145b1561111357828060010193505082600a60008581526020019081526020016000206000018190555033600a600085815260200190815260200160002060010160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550600090505b855181101561100257600c60008481526020019081526020016000208682815181101515610f8857fe5b9060200190602002015190806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550508080600101915050610f5e565b84600a60008581526020019081526020016000206002018190555083600a60008581526020019081526020016000206003019080519060200190611047929190611884565b50600b839080600181540180825580915050906001820390600052602060002001600090919290919091505550601054600360008282540392505081905550601160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc6010549081150290604051600060405180830381858888f193505050501580156110f0573d6000803e3d6000fd5b50600a600084815260200190815260200160002060020154600781905550611121565b8180600101925050610dfc565b5b505050505050565b60018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054141515611205576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260218152602001807f4f6e6c792041646d696e732043616e205365742047656e6573697320496e666f81526020017f2e0000000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b60001515600960009054906101000a900460ff1615151415156112b6576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f47656e6573697320496e666f2048617320416c7265616479204265656e20536581526020017f742e00000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b80600890805190602001906112cc929190611884565b506001600960006101000a81548160ff02191690831515021790555050565b60018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020541415156113c7576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001807f4f6e6c792041646d696e732043616e2055706461746520466c7573682045706f81526020017f63682e000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b6101688110151515611467576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260318152602001807f466c7573682045706f6368204d75737420626520457175616c20746f206f722081526020017f47726561746572207468616e203336302e00000000000000000000000000000081525060400191505060405180910390fd5b8060068190555050565b600f5481565b600080600060018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200190815260200160002054141515611558576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260238152602001807f4f6e6c792041646d696e732043616e204469737472696275746520476173204681526020017f65652e000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b6000925067016345785d8a00009150600090505b600c60008481526020019081526020016000208054905081101561162857600c6000848152602001908152602001600020818154811015156115aa57fe5b9060005260206000200160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff166108fc839081150290604051600060405180830381858888f1935050505015801561161a573d6000803e3d6000fd5b50808060010191505061156c565b505050565b60055481565b60018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414151561170f576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004018080602001828103825260228152602001807f4f6e6c792041646d696e732043616e2055706461746520436861696e204e616d81526020017f652e00000000000000000000000000000000000000000000000000000000000081525060400191505060405180910390fd5b8060029080519060200190611725929190611884565b5050565b60018060003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000205414151561177657600080fd5b600354811115151561178757600080fd5b806003600082825403925050819055508173ffffffffffffffffffffffffffffffffffffffff166108fc829081150290604051600060405180830381858888f193505050501580156117dd573d6000803e3d6000fd5b505050565b606060088054600181600116156101000203166002900480601f01602080910402602001604051908101604052809291908181526020018280546001816001161561010002031660029004801561187a5780601f1061184f5761010080835404028352916020019161187a565b820191906000526020600020905b81548152906001019060200180831161185d57829003601f168201915b5050505050905090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106118c557805160ff19168380011785556118f3565b828001600101855582156118f3579182015b828111156118f25782518255916020019190600101906118d7565b5b5090506119009190611904565b5090565b61192691905b8082111561192257600081600090555060010161190a565b5090565b905600a165627a7a72305820b774ee93ab81df2cb13319612382aae2c771f2ea67cac01fff4cc261d46d55060029`

// DeployAppChainBase deploys a new MoacNode contract, binding an instance of AppChainBase to it.
func DeployAppChainBase(auth *bind.TransactOpts, backend bind.ContractBackend, name string, uniqueId *big.Int, blockSec *big.Int, flushNumber *big.Int, initial_validators []common.Address, totalSupply *big.Int) (common.Address, *types.Transaction, *AppChainBase, error) {
	parsed, err := abi.JSON(strings.NewReader(AppChainBaseABI))
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(AppChainBaseBin), backend, name, uniqueId, blockSec, flushNumber, initial_validators, totalSupply)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &AppChainBase{AppChainBaseCaller: AppChainBaseCaller{contract: contract}, AppChainBaseTransactor: AppChainBaseTransactor{contract: contract}}, nil
}

// AppChainBase is an auto generated Go binding around an MoacNode contract.
type AppChainBase struct {
	AppChainBaseCaller     // Read-only binding to the contract
	AppChainBaseTransactor // Write-only binding to the contract
}

// AppChainBaseCaller is an auto generated read-only Go binding around an MoacNode contract.
type AppChainBaseCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainBaseTransactor is an auto generated write-only Go binding around an MoacNode contract.
type AppChainBaseTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AppChainBaseSession is an auto generated Go binding around an MoacNode contract,
// with pre-set call and transact options.
type AppChainBaseSession struct {
	Contract     *AppChainBase     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AppChainBaseCallerSession is an auto generated read-only Go binding around an MoacNode contract,
// with pre-set call options.
type AppChainBaseCallerSession struct {
	Contract *AppChainBaseCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// AppChainBaseTransactorSession is an auto generated write-only Go binding around an MoacNode contract,
// with pre-set transact options.
type AppChainBaseTransactorSession struct {
	Contract     *AppChainBaseTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// AppChainBaseRaw is an auto generated low-level Go binding around an MoacNode contract.
type AppChainBaseRaw struct {
	Contract *AppChainBase // Generic contract binding to access the raw methods on
}

// AppChainBaseCallerRaw is an auto generated low-level read-only Go binding around an MoacNode contract.
type AppChainBaseCallerRaw struct {
	Contract *AppChainBaseCaller // Generic read-only contract binding to access the raw methods on
}

// AppChainBaseTransactorRaw is an auto generated low-level write-only Go binding around an MoacNode contract.
type AppChainBaseTransactorRaw struct {
	Contract *AppChainBaseTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAppChainBase creates a new instance of AppChainBase, bound to a specific deployed contract.
func NewAppChainBase(address common.Address, backend bind.ContractBackend) (*AppChainBase, error) {
	contract, err := bindAppChainBase(address, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AppChainBase{AppChainBaseCaller: AppChainBaseCaller{contract: contract}, AppChainBaseTransactor: AppChainBaseTransactor{contract: contract}}, nil
}

// NewAppChainBaseCaller creates a new read-only instance of AppChainBase, bound to a specific deployed contract.
func NewAppChainBaseCaller(address common.Address, caller bind.ContractCaller) (*AppChainBaseCaller, error) {
	contract, err := bindAppChainBase(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &AppChainBaseCaller{contract: contract}, nil
}

// NewAppChainBaseTransactor creates a new write-only instance of AppChainBase, bound to a specific deployed contract.
func NewAppChainBaseTransactor(address common.Address, transactor bind.ContractTransactor) (*AppChainBaseTransactor, error) {
	contract, err := bindAppChainBase(address, nil, transactor)
	if err != nil {
		return nil, err
	}
	return &AppChainBaseTransactor{contract: contract}, nil
}

// bindAppChainBase binds a generic wrapper to an already deployed contract.
func bindAppChainBase(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AppChainBaseABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor), nil
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

// FLUSH_AMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) FLUSH_AMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "FLUSH_AMOUNT")
	return *ret0, err
}

// FLUSH_AMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) FLUSH_AMOUNT() (*big.Int, error) {
	return _AppChainBase.Contract.FLUSH_AMOUNT(&_AppChainBase.CallOpts)
}

// FLUSH_AMOUNT is a free data retrieval call binding the contract method 0x08925f61.
//
// Solidity: function FLUSH_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) FLUSH_AMOUNT() (*big.Int, error) {
	return _AppChainBase.Contract.FLUSH_AMOUNT(&_AppChainBase.CallOpts)
}

// FOUNDATION_BLACK_HOLE_ADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_AppChainBase *AppChainBaseCaller) FOUNDATION_BLACK_HOLE_ADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "FOUNDATION_BLACK_HOLE_ADDRESS")
	return *ret0, err
}

// FOUNDATION_BLACK_HOLE_ADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_AppChainBase *AppChainBaseSession) FOUNDATION_BLACK_HOLE_ADDRESS() (common.Address, error) {
	return _AppChainBase.Contract.FOUNDATION_BLACK_HOLE_ADDRESS(&_AppChainBase.CallOpts)
}

// FOUNDATION_BLACK_HOLE_ADDRESS is a free data retrieval call binding the contract method 0x5a35b2fc.
//
// Solidity: function FOUNDATION_BLACK_HOLE_ADDRESS() constant returns(address)
func (_AppChainBase *AppChainBaseCallerSession) FOUNDATION_BLACK_HOLE_ADDRESS() (common.Address, error) {
	return _AppChainBase.Contract.FOUNDATION_BLACK_HOLE_ADDRESS(&_AppChainBase.CallOpts)
}

// FOUNDATION_MOAC_REQUIRED_AMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseCaller) FOUNDATION_MOAC_REQUIRED_AMOUNT(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _AppChainBase.contract.Call(opts, out, "FOUNDATION_MOAC_REQUIRED_AMOUNT")
	return *ret0, err
}

// FOUNDATION_MOAC_REQUIRED_AMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) FOUNDATION_MOAC_REQUIRED_AMOUNT() (*big.Int, error) {
	return _AppChainBase.Contract.FOUNDATION_MOAC_REQUIRED_AMOUNT(&_AppChainBase.CallOpts)
}

// FOUNDATION_MOAC_REQUIRED_AMOUNT is a free data retrieval call binding the contract method 0xdfe89af2.
//
// Solidity: function FOUNDATION_MOAC_REQUIRED_AMOUNT() constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) FOUNDATION_MOAC_REQUIRED_AMOUNT() (*big.Int, error) {
	return _AppChainBase.Contract.FOUNDATION_MOAC_REQUIRED_AMOUNT(&_AppChainBase.CallOpts)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins( address) constant returns(uint256)
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
// Solidity: function admins( address) constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) Admins(arg0 common.Address) (*big.Int, error) {
	return _AppChainBase.Contract.Admins(&_AppChainBase.CallOpts, arg0)
}

// Admins is a free data retrieval call binding the contract method 0x429b62e5.
//
// Solidity: function admins( address) constant returns(uint256)
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
// Solidity: function flushList( uint256) constant returns(uint256)
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
// Solidity: function flushList( uint256) constant returns(uint256)
func (_AppChainBase *AppChainBaseSession) FlushList(arg0 *big.Int) (*big.Int, error) {
	return _AppChainBase.Contract.FlushList(&_AppChainBase.CallOpts, arg0)
}

// FlushList is a free data retrieval call binding the contract method 0x18c15211.
//
// Solidity: function flushList( uint256) constant returns(uint256)
func (_AppChainBase *AppChainBaseCallerSession) FlushList(arg0 *big.Int) (*big.Int, error) {
	return _AppChainBase.Contract.FlushList(&_AppChainBase.CallOpts, arg0)
}

// FlushMapping is a free data retrieval call binding the contract method 0x4c2f8619.
//
// Solidity: function flushMapping( uint256) constant returns(flushId uint256, validator address, blockNumber uint256, blockHash string)
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
// Solidity: function flushMapping( uint256) constant returns(flushId uint256, validator address, blockNumber uint256, blockHash string)
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
// Solidity: function flushMapping( uint256) constant returns(flushId uint256, validator address, blockNumber uint256, blockHash string)
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
// Solidity: function flushValidatorList( uint256,  uint256) constant returns(address)
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
// Solidity: function flushValidatorList( uint256,  uint256) constant returns(address)
func (_AppChainBase *AppChainBaseSession) FlushValidatorList(arg0 *big.Int, arg1 *big.Int) (common.Address, error) {
	return _AppChainBase.Contract.FlushValidatorList(&_AppChainBase.CallOpts, arg0, arg1)
}

// FlushValidatorList is a free data retrieval call binding the contract method 0x635cf62a.
//
// Solidity: function flushValidatorList( uint256,  uint256) constant returns(address)
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
// Solidity: function addAdmin(admin address) returns()
func (_AppChainBase *AppChainBaseTransactor) AddAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "addAdmin", admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(admin address) returns()
func (_AppChainBase *AppChainBaseSession) AddAdmin(admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.Contract.AddAdmin(&_AppChainBase.TransactOpts, admin)
}

// AddAdmin is a paid mutator transaction binding the contract method 0x70480275.
//
// Solidity: function addAdmin(admin address) returns()
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
// Solidity: function flush(current_validators address[], blockNumber uint256, blockHash string) returns()
func (_AppChainBase *AppChainBaseTransactor) Flush(opts *bind.TransactOpts, current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "flush", current_validators, blockNumber, blockHash)
}

// Flush is a paid mutator transaction binding the contract method 0xd66cfad5.
//
// Solidity: function flush(current_validators address[], blockNumber uint256, blockHash string) returns()
func (_AppChainBase *AppChainBaseSession) Flush(current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _AppChainBase.Contract.Flush(&_AppChainBase.TransactOpts, current_validators, blockNumber, blockHash)
}

// Flush is a paid mutator transaction binding the contract method 0xd66cfad5.
//
// Solidity: function flush(current_validators address[], blockNumber uint256, blockHash string) returns()
func (_AppChainBase *AppChainBaseTransactorSession) Flush(current_validators []common.Address, blockNumber *big.Int, blockHash string) (*types.Transaction, error) {
	return _AppChainBase.Contract.Flush(&_AppChainBase.TransactOpts, current_validators, blockNumber, blockHash)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(admin address) returns()
func (_AppChainBase *AppChainBaseTransactor) RemoveAdmin(opts *bind.TransactOpts, admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "removeAdmin", admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(admin address) returns()
func (_AppChainBase *AppChainBaseSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.Contract.RemoveAdmin(&_AppChainBase.TransactOpts, admin)
}

// RemoveAdmin is a paid mutator transaction binding the contract method 0x1785f53c.
//
// Solidity: function removeAdmin(admin address) returns()
func (_AppChainBase *AppChainBaseTransactorSession) RemoveAdmin(admin common.Address) (*types.Transaction, error) {
	return _AppChainBase.Contract.RemoveAdmin(&_AppChainBase.TransactOpts, admin)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(genesis string) returns()
func (_AppChainBase *AppChainBaseTransactor) SetGenesisInfo(opts *bind.TransactOpts, genesis string) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "setGenesisInfo", genesis)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(genesis string) returns()
func (_AppChainBase *AppChainBaseSession) SetGenesisInfo(genesis string) (*types.Transaction, error) {
	return _AppChainBase.Contract.SetGenesisInfo(&_AppChainBase.TransactOpts, genesis)
}

// SetGenesisInfo is a paid mutator transaction binding the contract method 0xd9b11dc8.
//
// Solidity: function setGenesisInfo(genesis string) returns()
func (_AppChainBase *AppChainBaseTransactorSession) SetGenesisInfo(genesis string) (*types.Transaction, error) {
	return _AppChainBase.Contract.SetGenesisInfo(&_AppChainBase.TransactOpts, genesis)
}

// UpdateChainName is a paid mutator transaction binding the contract method 0xf095ee64.
//
// Solidity: function updateChainName(name string) returns()
func (_AppChainBase *AppChainBaseTransactor) UpdateChainName(opts *bind.TransactOpts, name string) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "updateChainName", name)
}

// UpdateChainName is a paid mutator transaction binding the contract method 0xf095ee64.
//
// Solidity: function updateChainName(name string) returns()
func (_AppChainBase *AppChainBaseSession) UpdateChainName(name string) (*types.Transaction, error) {
	return _AppChainBase.Contract.UpdateChainName(&_AppChainBase.TransactOpts, name)
}

// UpdateChainName is a paid mutator transaction binding the contract method 0xf095ee64.
//
// Solidity: function updateChainName(name string) returns()
func (_AppChainBase *AppChainBaseTransactorSession) UpdateChainName(name string) (*types.Transaction, error) {
	return _AppChainBase.Contract.UpdateChainName(&_AppChainBase.TransactOpts, name)
}

// UpdateFlushEpoch is a paid mutator transaction binding the contract method 0xdd77088e.
//
// Solidity: function updateFlushEpoch(newEpoch uint256) returns()
func (_AppChainBase *AppChainBaseTransactor) UpdateFlushEpoch(opts *bind.TransactOpts, newEpoch *big.Int) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "updateFlushEpoch", newEpoch)
}

// UpdateFlushEpoch is a paid mutator transaction binding the contract method 0xdd77088e.
//
// Solidity: function updateFlushEpoch(newEpoch uint256) returns()
func (_AppChainBase *AppChainBaseSession) UpdateFlushEpoch(newEpoch *big.Int) (*types.Transaction, error) {
	return _AppChainBase.Contract.UpdateFlushEpoch(&_AppChainBase.TransactOpts, newEpoch)
}

// UpdateFlushEpoch is a paid mutator transaction binding the contract method 0xdd77088e.
//
// Solidity: function updateFlushEpoch(newEpoch uint256) returns()
func (_AppChainBase *AppChainBaseTransactorSession) UpdateFlushEpoch(newEpoch *big.Int) (*types.Transaction, error) {
	return _AppChainBase.Contract.UpdateFlushEpoch(&_AppChainBase.TransactOpts, newEpoch)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(recv address, amount uint256) returns()
func (_AppChainBase *AppChainBaseTransactor) WithdrawFund(opts *bind.TransactOpts, recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AppChainBase.contract.Transact(opts, "withdrawFund", recv, amount)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(recv address, amount uint256) returns()
func (_AppChainBase *AppChainBaseSession) WithdrawFund(recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AppChainBase.Contract.WithdrawFund(&_AppChainBase.TransactOpts, recv, amount)
}

// WithdrawFund is a paid mutator transaction binding the contract method 0xf7c8d221.
//
// Solidity: function withdrawFund(recv address, amount uint256) returns()
func (_AppChainBase *AppChainBaseTransactorSession) WithdrawFund(recv common.Address, amount *big.Int) (*types.Transaction, error) {
	return _AppChainBase.Contract.WithdrawFund(&_AppChainBase.TransactOpts, recv, amount)
}
