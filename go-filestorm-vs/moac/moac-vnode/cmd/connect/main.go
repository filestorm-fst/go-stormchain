package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts/keystore"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/core/types"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mcclient"
)

var (
	// accJSONFlag = flag.String("account.json", "", "Key json file to fund user requests with")
	// accPassFlag = flag.String("account.pass", "", "Decryption password to access faucet funds")

	nodeIpFlag   = flag.String("nodeIp", "gateway.moac.io/mainnet", "Connecting node")
	contractFlag = flag.String("contract", "", "Registration contract")
)

func main() {

	strFunc := "0x3a5ee7e5"
	strNum := "0000000000000000000000000000000000000000000000000000000000000168"
	strHash := "00000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000042307839316331396530373966333936336634626138653138393239646536393762656536336431336437396165333066666337346131393465376233306366613065000000000000000000000000000000000000000000000000000000000000"

	data := []byte(strFunc + strNum + strHash)

	flag.Parse()

	nodeIp := "http://" + (*nodeIpFlag)

	if *contractFlag == "" {
		fmt.Println("Contract address was not provided.")
		return
	}

	contractAddress := common.HexToAddress(*contractFlag)

	fmt.Println("start")

	// fmt.Printf("Please enter password: ")
	// bytePassword, _ := terminal.ReadPassword(0)
	// pass := string(bytePassword)
	// fmt.Println()

	pass := "11223344"

	keystoreFile := "/Users/xiannongfu/StormChain/Moac/data/keystore/UTC--2020-03-04T00-21-07.280235000Z--a71848bf3847b493e1dc9ad334c48a49ff3b0e67"
	tmpPath := "/Users/xiannongfu/tmp"

	ks := keystore.NewKeyStore(tmpPath, keystore.StandardScryptN, keystore.StandardScryptP)
	blob, err := ioutil.ReadFile(keystoreFile)
	if err != nil {
		//		log.Crit("Failed to read account key contents", "file", *accJSONFlag, "err", err)
		fmt.Println("Failed to read keystore contents", "file", "err", err)
	}

	acc, err := ks.Import(blob, pass, pass)
	if err != nil {
		//		log.Crit("Failed to import faucet signer account", "err", err)
		fmt.Println("Failed to import keystore account", "err", err)
	}
	ks.Unlock(acc, pass)

	client, err := mcclient.Dial(nodeIp)

	if err != nil {
		fmt.Println("error: ", err)
	}

	balance, _ := client.BalanceAt(context.Background(), ks.Accounts()[0].Address, nil)
	nonce, _ := client.NonceAt(context.Background(), ks.Accounts()[0].Address, nil)
	price, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		fmt.Println("Suggested Gas Price Err: ", err)
		price = big.NewInt(20000000000)
	}

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		fmt.Println("chain Id err ", err)
	}
	fmt.Println(chainID)

	fmt.Println("balance|nonce|price: ", balance, nonce, price)

	tobalance, _ := client.BalanceAt(context.Background(), contractAddress, nil)
	fmt.Println("to account balance:", tobalance)

	//price = big.NewInt(35000000)
	//amount := big.NewInt(1).Mul(big.NewInt(2000), big.NewInt(100000000000000000))
	amount := big.NewInt(0)

	// Before EIP 161, Nonce starts from 0. MOAC didn't update EIP 161
	tx := types.NewTransaction(nonce, contractAddress, amount, big.NewInt(35000), price, 0, nil, data)
	signed, err := ks.SignTx(ks.Accounts()[0], tx, chainID)
	if err != nil {
		fmt.Println("Sign error: ", err)
	}

	if err := client.SendTransaction(context.Background(), signed); err != nil {
		fmt.Println("Sent Transaction error: ", err)
	}

	err1 := os.Remove(keystoreFile)
	if err1 != nil {
		fmt.Println("Failed to remove keystore.")
	}

	fmt.Println("end")

}
