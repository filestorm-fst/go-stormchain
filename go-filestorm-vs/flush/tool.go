// chain3go project tool.go
package flush

import (
	"context"
	"crypto/ecdsa"
	"github.com/filestorm/go-filestorm/accounts/abi/bind"
	"github.com/filestorm/go-filestorm/accounts/keystore"
	"github.com/filestorm/go-filestorm/common"
	"github.com/filestorm/go-filestorm/common/hexutil"
	"github.com/filestorm/go-filestorm/crypto"
	"github.com/filestorm/go-filestorm/fstclient"
	"github.com/pborman/uuid"
	"github.com/filestorm/go-filestorm/log"
	"math/big"

)

//根据keystore字符串获取私钥
func GetPrivateKey(storeKey string, password string) (privateKey string,err error) {
	key, err := keystore.DecryptKey([]byte(storeKey), password)
	if err != nil {
		return "",err
	}
	privateKeyBytes := crypto.FromECDSA(key.PrivateKey)
	privateKey = hexutil.Encode(privateKeyBytes)[2:]
	return privateKey,nil
}

//根据私钥获取keystore字符串
func GetKeystoreStr(privateKey, password string) (string, string, error) {

	var err error
	var testKey keystore.Key
	testKey.PrivateKey, err = crypto.HexToECDSA(privateKey)
	if err != nil {
		return "", "", err
	}
	testKey.Address = crypto.PubkeyToAddress(testKey.PrivateKey.PublicKey)
	testKey.Id = uuid.NewRandom()

	scryptN := keystore.StandardScryptN
	scryptP := keystore.StandardScryptP

	var keyJson []byte
	keyJson, err = keystore.EncryptKey(&testKey, password, scryptN, scryptP)

	return string(keyJson), testKey.Address.Hex(), err
}



func ClientDeployContract(client *fstclient.Client ,keystore string,password string, name string,uniqueId *big.Int, blockSec *big.Int, flushNumber *big.Int, initialValidators []common.Address, totalSupply *big.Int) (contract string,tx string,err error){
	privateKey,err :=GetPrivateKey(keystore,password)
	if err != nil {
		log.Error(err.Error())
		return "","",err
	}
	//load privateKey
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Error(err.Error())
		return "","",err
	}
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("error casting public key to ECDSA")
		return "","",err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	//1.set gasPrice and gasLimit
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error(err.Error())
		return "","",err
	}
	//2.set value (fst)
	value := new(big.Int)
	value.SetString("10000000000000000000", 10)
	//3 get nonce
	nonce, err := client.PendingNonceAt(context.Background(), common.Address(common.HexToAddress(address)))
	if err != nil {
		log.Error(err.Error())
		return "","",err
	}

	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(key)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = value
	auth.GasLimit = uint64(9000000)
	auth.GasPrice = gasPrice

	contractAddress, transaction, instance, err := DeployAppChainBase(auth, client, name,uniqueId,blockSec,flushNumber,initialValidators,totalSupply)
	if err != nil {
		log.Error(err.Error())
		return "","",err
	}

	_ = instance
	return contractAddress.Hex(),transaction.Hash().Hex(),nil
}


func ClientFlush(client *fstclient.Client ,keystore string,password string, contract string ,validators []common.Address,blockNumber *big.Int, blockHash string) (txHash string,err error){
	privateKey,err :=GetPrivateKey(keystore,password)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}
	//load privateKey
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("error casting public key to ECDSA")
		return "",err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(key)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice
	contractAddress := common.HexToAddress(contract)
	instance, err := NewAppChainBase(contractAddress, client)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	tx, err := instance.Flush(auth,validators,blockNumber,blockHash)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	return tx.Hash().Hex(),nil
}


func DistributeGasFee(client *fstclient.Client ,keystore string,password string, contract string) (txHash string,err error){
	privateKey,err :=GetPrivateKey(keystore,password)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}
	//load privateKey
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("error casting public key to ECDSA")
		return "",err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(key)
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice
	contractAddress := common.HexToAddress(contract)
	instance, err := NewAppChainBase(contractAddress, client)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	tx, err := instance.DistributeGasFee(auth)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	return tx.Hash().Hex(),nil
}




func AddFund(client *fstclient.Client ,keystore string,password string, contract string) (txHash string,err error){
	privateKey,err :=GetPrivateKey(keystore,password)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}
	//load privateKey
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}
	publicKey := key.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Error("error casting public key to ECDSA")
		return "",err
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(address))
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	//auth, err := bind.NewTransactor(strings.NewReader(keystore), password)
	auth := bind.NewKeyedTransactor(key)
	auth.Nonce = big.NewInt(int64(nonce))
	value := new(big.Int)
	value.SetString("10000000000000000000", 10)
	auth.Value = value
	auth.GasLimit = uint64(3000000)
	auth.GasPrice = gasPrice
	contractAddress := common.HexToAddress(contract)
	instance, err := NewAppChainBase(contractAddress, client)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	tx, err := instance.AddFund(auth)
	if err != nil {
		log.Error(err.Error())
		return "",err
	}

	return tx.Hash().Hex(),nil
}


