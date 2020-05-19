package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/filestorm/go-filestorm/moac/chain3go"
	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/common/hexutil"
	"github.com/filestorm/go-filestorm/moac/moac-lib/crypto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/crypto/sha3"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts/abi"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts/abi/bind"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts/keystore"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mcclient"
	"log"
	"math"
	"math/big"
	"strings"
	"testing"
)

var clientIp = "http://47.115.27.232:8502"
var moacClientIp = "http://47.112.217.66:8547"
var getWay  = "http://gateway.moac.io/mainnet"
var testGetWay  = "http://gateway.moac.io/testnet"

func TestNewAccount(t *testing.T)  {
	privateKeyStruct, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(privateKeyStruct)
	privateKey := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Println(privateKey) //7c0f8e7d6a872bbd0ae138787ef4a35a449351d85055b9357d4a536f75ed229c

	publicKey := privateKeyStruct.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Println(hexutil.Encode(publicKeyBytes)[4:])  //535d0b9e22ee1be115349138c57e368d1493c184cfc3c0d6e372d2efbe88ea1a196aaeed247a0a4c01d3d08753e08679a8e35b1e224cfccdc0165ed699a5713d

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Println(address)  //0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de

	hash := sha3.NewKeccak256()
	hash.Write(publicKeyBytes[1:])
	fmt.Println(hexutil.Encode(hash.Sum(nil)[12:])) //0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de
}

func TestClientBalance(t *testing.T)  {
	client, err := mcclient.Dial(moacClientIp)
	if err != nil {
		log.Fatal(err)
		return
	}

	//query account balance
	account := common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	fstValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(fstValue)
}

func TestClientDeployContract(t *testing.T) {

	client, err := mcclient.Dial(clientIp)
	if err != nil {
		log.Fatal(err)
		return
	}

	account := common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de")

	//1.set gasPrice and gasLimit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(4000000) // in units
	//2.set value (fst)
	value := new(big.Int)
	value.SetString("20000000000000000000", 10)
	//value := big.NewInt(0)
	//3 get nonce
	nonce, err := client.PendingNonceAt(context.Background(), account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("nonce: %d" , nonce)
	fmt.Println()

	//load privateKey
	privateKey, err := crypto.HexToECDSA("7c0f8e7d6a872bbd0ae138787ef4a35a449351d85055b9357d4a536f75ed229c")
	if err != nil {
		log.Fatal(err)
	}
	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = big.NewInt(int64(gasLimit))
	auth.GasPrice = gasPrice

	//deploy contract
	var initialValidators []common.Address
	initialValidators = append(initialValidators, common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de"))
	initialValidators = append(initialValidators, common.HexToAddress("0x9d817e62c998d274c7f95083205e0b76c48e5ae6"))
	initialValidators = append(initialValidators, common.HexToAddress("0x336fc8e106a86e1f8cde7c5e56a7c8a0e23039f3"))

	//address, tx, instance, err := DeployApplicationChain(auth,client,big.NewInt(15),big.NewInt(200),initialValidators,big.NewInt(1))
	if err != nil {
		log.Fatal(err)
	}
	//input := "1.0"
	//address, tx, instance, err := chain3go.DeployStore(auth,client,input)
	//if err != nil {
	//	log.Fatal(err)
	//}


	//fmt.Println(address.Hex())   // 0xF2D99BCA7D993dB11f33fdF2235D7Af7c9Ab4824  0xad4cE0F1Ad5D6f44595CFA5250Aad642A27132e8
	//fmt.Println(tx.Hash().Hex()) // 0x3fc8300a54d81a00afc2a830a29edd81a3259a49b9ff3e4733f7946703992ae8  0x9fe06fcec0fda9a9fa56f05f13a3294a4ae07ed716cf3c5552ed3dd9ac0643f2

	//_ = instance
}


func TestClientCallContract(t *testing.T) {
	client, err := mcclient.Dial(clientIp)
	if err != nil {
		log.Fatal(err)
		return
	}

	address := common.HexToAddress("0xe37a718B024Bbb76eA5942C6A6B50EA909019cea")
	instance, err := NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}
	key := [32]byte{}
	copy(key[:], []byte("foo"))
	genesis, err := instance.Items(nil,key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(genesis)

}


func TestClientCheckTrans(t *testing.T) {
	client, err := mcclient.Dial(clientIp)
	if err != nil {
		log.Fatal(err)
		return
	}
	txHash := common.HexToHash("0x1ed625cda1eca19974e5a54dbebc4e02653192e27fe5b70ad141fcb4538e38fb")
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err != nil {
		log.Fatal(err)
	}

	marshal, err := json.Marshal(receipt)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s",marshal)
}

func TestKeystoreToPrivateKey(t *testing.T)  {
	storeKey := `{"address":"7d8958276355603e1d6b08cf9ddfa68f245eacad","crypto":{"cipher":"aes-128-ctr","ciphertext":"c3d7c0d073418428af5418ad1ba57e947671fffc4869241885536752a4848e6a","cipherparams":{"iv":"f0f08f515b3def87436e0e9bf371c5e6"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"0d81a32aef9f1914afb7be221b08dc3ade1614a53ef2ca32c944fc0e4f1ffe9a"},"mac":"9b338ab7ad4c39b6676b0ed1b83c1e07b32a886333ad6c8687d6f6993f1c416b"},"id":"94e8f65f-6335-4abc-b363-448361d274d5","version":3}`
	password := "shuqian2020!"
	key, err := keystore.DecryptKey([]byte(storeKey), password)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	privateKey := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Printf("%s",privateKey)
}

func TestMOACClientDeployContract(t *testing.T) {
	client, err := mcclient.Dial(moacClientIp)
	if err != nil {
		log.Fatal(err)
		return
	}
	account := common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de")

	//1.set gasPrice and gasLimit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(400000) // in units
	//2.set value (fst)
	//value := new(big.Int)
	//value.SetString("20000000000000000000", 10)
	value := big.NewInt(0)
	//3 get nonce
	//nonce, err := rpcClient.PendingNonceAt(context.Background(), account)
	nonce, err := client.PendingNonceAt(context.Background(),account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("nonce: %d" , nonce)
	fmt.Println()

	//load privateKey
	privateKey, err := crypto.HexToECDSA("7c0f8e7d6a872bbd0ae138787ef4a35a449351d85055b9357d4a536f75ed229c")
	if err != nil {
		log.Fatal(err)
	}
	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = big.NewInt(int64(gasLimit))
	auth.GasPrice = gasPrice

	//deploy contract
	//var initialValidators []common.Address
	//initialValidators = append(initialValidators, common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de"))
	//initialValidators = append(initialValidators, common.HexToAddress("0x9d817e62c998d274c7f95083205e0b76c48e5ae6"))
	//initialValidators = append(initialValidators, common.HexToAddress("0x336fc8e106a86e1f8cde7c5e56a7c8a0e23039f3"))
	//
	//address, tx, instance, err := chain3go.DeployApplicationChain(auth,client,big.NewInt(15),big.NewInt(200),initialValidators,big.NewInt(1))
	//if err != nil {
	//	log.Fatal(err)
	//}
	input := "1.0"
	address, tx, instance, err := DeployStore(auth,client,input)
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println(address.Hex())   // 0xF2D99BCA7D993dB11f33fdF2235D7Af7c9Ab4824  0xad4cE0F1Ad5D6f44595CFA5250Aad642A27132e8
	fmt.Println(tx.Hash().Hex()) // 0x3fc8300a54d81a00afc2a830a29edd81a3259a49b9ff3e4733f7946703992ae8  0x9fe06fcec0fda9a9fa56f05f13a3294a4ae07ed716cf3c5552ed3dd9ac0643f2

	_ = instance
}

func TestDeployStore(t *testing.T) {
	client, err := mcclient.Dial(getWay)
	if err != nil {
		log.Fatal(err)
		return
	}
	account := common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de")

	//1.set gasPrice and gasLimit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(400000) // in units
	//2.set value (fst)
	//value := new(big.Int)
	//value.SetString("20000000000000000000", 10)
	value := big.NewInt(0)
	//3 get nonce
	//nonce, err := rpcClient.PendingNonceAt(context.Background(), account)
	nonce, err := client.PendingNonceAt(context.Background(),account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("nonce: %d" , nonce)
	fmt.Println()

	//load privateKey
	privateKey, err := crypto.HexToECDSA("7c0f8e7d6a872bbd0ae138787ef4a35a449351d85055b9357d4a536f75ed229c")
	if err != nil {
		log.Fatal(err)
	}
	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = big.NewInt(int64(gasLimit))
	auth.GasPrice = gasPrice

	parsed, err := abi.JSON(strings.NewReader(StoreABI))
	if err != nil {
		panic(err)
	}

	input := "1.0"
	address, tx, contract, err := bind.DeployContract(auth, parsed, common.FromHex(StoreBin), client, input)
	if err != nil {
		panic(err)
	}

	fmt.Println(address.Hex())    //0x64eF1F64332b1a70869aeE541941e8399903EFD1
	fmt.Println(tx.Hash().Hex())  //0x2f6d092987f2937e98b3ffac077885cf4cc57571fb0caad4b1711de1c1271b0c

	_ = contract
}

func TestClientCallContract2(t *testing.T) {
	client, err := mcclient.Dial(moacClientIp)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA("7c0f8e7d6a872bbd0ae138787ef4a35a449351d85055b9357d4a536f75ed229c")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = big.NewInt(300000) // in units
	auth.GasPrice = gasPrice

	address := common.HexToAddress("0x64eF1F64332b1a70869aeE541941e8399903EFD1")
	instance, err := NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	tx, err := instance.SetItem(auth, key, value)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", tx.Hash().Hex()) // tx sent: 0x8d490e535678e9a24360e955d75b27ad307bdfb97a1dca51d0f3035dcee3e870

	result, err := instance.Items(nil, key)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(result[:])) // "bar"

}


func TestClientDeployContract2(t *testing.T) {
	client, err := mcclient.Dial(testGetWay)
	if err != nil {
		log.Fatal(err)
		return
	}
	ketstoreStr := `{"address":"15b5cf1860998d78c3f0808082f6f3ce5209e7de","crypto":{"cipher":"aes-128-ctr","ciphertext":"37b5ec3d0d6d5a49bbdb9d46b3d80574b985fa6e7dcb7f76f6ea21a4dade84b7","cipherparams":{"iv":"0ec47b6fdbc4a7611e5e4c407d9cb0d7"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"7f4cf4808c7f5245d87254018b49931762415dc318b5054ba5c16b96885b0797"},"mac":"3aa1b7e2878aff517aa3439b1cf80701ea5dc3e9c21000dfa857865374e9f1a3"},"id":"e79de0c7-050d-4fe5-8fec-714f60faaad5","version":3}`
	password := "123"
	var initialValidators []common.Address
	initialValidators = append(initialValidators, common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de"))
	initialValidators = append(initialValidators, common.HexToAddress("0x9d817e62c998d274c7f95083205e0b76c48e5ae6"))
	initialValidators = append(initialValidators, common.HexToAddress("0x336fc8e106a86e1f8cde7c5e56a7c8a0e23039f3"))
	contract, tx, err := chain3go.ClientDeployContract(client,ketstoreStr,password,"",big.NewInt(123456),big.NewInt(15),big.NewInt(200),initialValidators,big.NewInt(1000))
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("contract is :" ,contract)   //0xe3036C23c2e63A04B85EB67507407C23AFbBd976
	fmt.Println("tx hash is :" , tx)    //0x914931d977e657a1ee711add4174d0d6fdedeaa3129a01e218bd19eaf3267c9f
}


func TestClientCallContract4(t *testing.T) {
	client, err := mcclient.Dial(testGetWay)
	if err != nil {
		log.Fatal(err)
		return
	}
	ketstoreStr := `{"address":"15b5cf1860998d78c3f0808082f6f3ce5209e7de","crypto":{"cipher":"aes-128-ctr","ciphertext":"37b5ec3d0d6d5a49bbdb9d46b3d80574b985fa6e7dcb7f76f6ea21a4dade84b7","cipherparams":{"iv":"0ec47b6fdbc4a7611e5e4c407d9cb0d7"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"7f4cf4808c7f5245d87254018b49931762415dc318b5054ba5c16b96885b0797"},"mac":"3aa1b7e2878aff517aa3439b1cf80701ea5dc3e9c21000dfa857865374e9f1a3"},"id":"e79de0c7-050d-4fe5-8fec-714f60faaad5","version":3}`
	password := "123"
	var initialValidators []common.Address
	initialValidators = append(initialValidators, common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de"))
	initialValidators = append(initialValidators, common.HexToAddress("0x9d817e62c998d274c7f95083205e0b76c48e5ae6"))
	initialValidators = append(initialValidators, common.HexToAddress("0x336fc8e106a86e1f8cde7c5e56a7c8a0e23039f3"))
	contract := "0xe3036C23c2e63A04B85EB67507407C23AFbBd976"
	txHash, err := chain3go.ClientFlush(client, ketstoreStr, password, contract, initialValidators, big.NewInt(4417457), "0x6a5c557306ad3e220763c3aaf74c65c2faff9834a7f478149beb65b3a39ff924")
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("contract is :" ,contract)   //0xe3036C23c2e63A04B85EB67507407C23AFbBd976
	fmt.Println("tx hash is :" , txHash)    //0x51dd61753a4020aa6b4e41ecfa5732707e1fbdd3eadfbc178c01cc38496f867d
}

func TestMoacTrans(t *testing.T) {
	client, err := mcclient.Dial(getWay)
	if err != nil {
		log.Fatal(err)
		return
	}
	privateKey, err := crypto.HexToECDSA("980f9d3cd1b8cfc3bd33c153b39f5ada9f3c8e2d805cffa96e7bb33d6529ed0f")
	if err != nil {
		log.Fatal(err)
	}

	via := common.HexToAddress("7d8958276355603e1d6b08cf9ddfa68f245eacad")

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	value := new(big.Int)
	value.SetString("50000000000000000000", 10)
	gasLimit := uint64(3000000)                // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//toAddress := common.HexToAddress("0xfc6f4451cee24139da955794061c60beb305f925")
	//toAddress := common.HexToAddress("0x82b071e55366d31c9a3ab2c3274d2db934516a2a")
	toAddress := common.HexToAddress("0xd02443b8d564fed4aD332Cd52508b69b511dF5B8")
	var data []byte

	tx := types.NewTransaction(nonce, toAddress, value, big.NewInt(int64(gasLimit)), gasPrice, 0, &via , data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	signedTx, err := types.SignTx(tx, types.NewPanguSigner(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}

func TestAddFund(t *testing.T)  {
	client, err := mcclient.Dial(moacClientIp)
	privateKey,err :=chain3go.GetPrivateKey(`{"address":"7d8958276355603e1d6b08cf9ddfa68f245eacad","crypto":{"cipher":"aes-128-ctr","ciphertext":"c3d7c0d073418428af5418ad1ba57e947671fffc4869241885536752a4848e6a","cipherparams":{"iv":"f0f08f515b3def87436e0e9bf371c5e6"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"0d81a32aef9f1914afb7be221b08dc3ade1614a53ef2ca32c944fc0e4f1ffe9a"},"mac":"9b338ab7ad4c39b6676b0ed1b83c1e07b32a886333ad6c8687d6f6993f1c416b"},"id":"94e8f65f-6335-4abc-b363-448361d274d5","version":3}`,"shuqian2020!")
	if err != nil {
		log.Fatal(err)
	}
	//load privateKey
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal(err)
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(key)
	auth.Nonce = big.NewInt(int64(nonce))
	value := new(big.Int)
	value.SetString("500000000000000000000", 10)
	auth.Value = value
	auth.GasLimit = big.NewInt(int64(3000000))
	auth.GasPrice = gasPrice
	contractAddress := common.HexToAddress("0x4d1CeAA9D8FB6daf3b5A6c74739E73662Db860Cf")
	instance, err := chain3go.NewAppChainBase(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	tx, err := instance.AddFund(auth)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(tx.Hash().Hex())
}