package tests

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/filestorm/go-filestorm/accounts/abi/bind"
	"github.com/filestorm/go-filestorm/accounts/keystore"
	"github.com/filestorm/go-filestorm/common"
	"github.com/filestorm/go-filestorm/common/hexutil"
	"github.com/filestorm/go-filestorm/crypto"
	"github.com/filestorm/go-filestorm/fstclient"
	"github.com/filestorm/go-filestorm/moac/chain3go"
	"github.com/tonnerre/golang-go.crypto/sha3"
	"log"
	"math"
	"math/big"
	"testing"
)

var clientIp = "http://47.115.27.232:8502"
var moacClientIp = "http://47.112.217.66:8547"

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
	client, err := fstclient.Dial(clientIp)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("gasPrice: %d",gasPrice.Int64())

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

	client, err := fstclient.Dial(clientIp)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

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
	auth.GasLimit = gasLimit
	auth.GasPrice = gasPrice

	//deploy contract
	var initialValidators []common.Address
	initialValidators = append(initialValidators, common.HexToAddress("0x15B5cF1860998d78c3f0808082F6F3Ce5209e7de"))
	initialValidators = append(initialValidators, common.HexToAddress("0x9d817e62c998d274c7f95083205e0b76c48e5ae6"))
	initialValidators = append(initialValidators, common.HexToAddress("0x336fc8e106a86e1f8cde7c5e56a7c8a0e23039f3"))

	address, tx, instance, err := chain3go.DeployApplicationChain(auth,client,big.NewInt(15),big.NewInt(200),initialValidators,big.NewInt(1))
	if err != nil {
		log.Fatal(err)
	}
	//input := "1.0"
	//address, tx, instance, err := chain3go.DeployStore(auth,client,input)
	//if err != nil {
	//	log.Fatal(err)
	//}


	fmt.Println(address.Hex())   // 0xF2D99BCA7D993dB11f33fdF2235D7Af7c9Ab4824  0xad4cE0F1Ad5D6f44595CFA5250Aad642A27132e8
	fmt.Println(tx.Hash().Hex()) // 0x3fc8300a54d81a00afc2a830a29edd81a3259a49b9ff3e4733f7946703992ae8  0x9fe06fcec0fda9a9fa56f05f13a3294a4ae07ed716cf3c5552ed3dd9ac0643f2

	_ = instance
}

func TestClientCallContract(t *testing.T) {
	client, err := fstclient.Dial(clientIp)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

	address := common.HexToAddress("0x56433d39e397bAcD6c8351Ed5E3a6864147999Df")
	instance, err := chain3go.NewStore(address, client)
	if err != nil {
		log.Fatal(err)
	}

	genesis, err := instance.Version(nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(genesis)
}


func TestClientCheckTrans(t *testing.T) {
	client, err := fstclient.Dial(clientIp)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer client.Close()

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
	storeKey := `{"address":"d02443b8d564fed4ad332cd52508b69b511df5b8","crypto":{"cipher":"aes-128-ctr","ciphertext":"0ce1a5520297c7c28a4a4956ca80b274d7345d32fc07e4682c4c7f88da715455","cipherparams":{"iv":"d168b4c58ac0c942d40155409c993ed4"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"dfb19c95aadbe7b2b94caf3eac9f94f08f6b1e6c6e16f6fd93f3585d8bc6dbe2"},"mac":"15828209a3c40a1b0ac01ea0d67a9c5b8b9834f487869b8e5092a59dadbd0124"},"id":"7297a260-ae92-4cd6-b768-bc9a8e2a11f4","version":3}`
	password := "GZC15527185733"
	key, err := keystore.DecryptKey([]byte(storeKey), password)
	if err != nil {
		log.Fatal(err)
	}
	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	privateKey := hexutil.Encode(privateKeyBytes)[2:]
	fmt.Printf("%s",privateKey)
}

func TestMOACClientDeployContract(t *testing.T) {
	rpcClient := chain3go.NewRpcClient(moacClientIp, 101)
	client, err := fstclient.Dial(moacClientIp)
	//account := common.HexToAddress("0xd02443b8d564fed4ad332cd52508b69b511df5b8")

	//1.set gasPrice and gasLimit
	gasPrice, err := rpcClient.MC_gasPrice()
	if err != nil {
		log.Fatal(err)
	}
	gasLimit := uint64(300000) // in units
	//2.set value (fst)
	//value := new(big.Int)
	//value.SetString("20000000000000000000", 10)
	value := big.NewInt(0)
	//3 get nonce
	//nonce, err := rpcClient.PendingNonceAt(context.Background(), account)
	nonce, err := rpcClient.MC_getTransactionCount("0xd02443b8d564fed4ad332cd52508b69b511df5b8", "pending")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("nonce: %d" , nonce)
	fmt.Println()

	//load privateKey
	privateKey, err := crypto.HexToECDSA("f0765fdcab6dba9ca014b2a2a97f6bef8d679ad69f3695f9baed7d9abf25bc94")
	if err != nil {
		log.Fatal(err)
	}
	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(privateKey)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = gasLimit
	auth.GasPrice = big.NewInt(gasPrice)

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
	address, tx, instance, err := chain3go.DeployStore(auth,client,input)
	if err != nil {
		log.Fatal(err)
	}


	fmt.Println(address.Hex())   // 0xF2D99BCA7D993dB11f33fdF2235D7Af7c9Ab4824  0xad4cE0F1Ad5D6f44595CFA5250Aad642A27132e8
	fmt.Println(tx.Hash().Hex()) // 0x3fc8300a54d81a00afc2a830a29edd81a3259a49b9ff3e4733f7946703992ae8  0x9fe06fcec0fda9a9fa56f05f13a3294a4ae07ed716cf3c5552ed3dd9ac0643f2

	_ = instance
}
