// Copyright 2016 The MOAC-core Authors
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

package types

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/crypto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
)

var emptyAddress = common.StringToAddress("")

func TestEIP155Signing(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	signer := NewPanguSigner(big.NewInt(18))
	tx, err := SignTx(NewTransaction(0, addr, new(big.Int), new(big.Int), new(big.Int), 0, nil, nil), signer, key)
	if err != nil {
		t.Fatal(err)
	}

	from, err := Sender(signer, tx)
	if err != nil {
		t.Fatal(err)
	}
	if from != addr {
		t.Errorf("exected from and address to be equal. Got %x want %x", from, addr)
	}
}

func TestEIP155ChainId(t *testing.T) {
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)

	//Set the signer with chainID = 18
	signer := NewPanguSigner(big.NewInt(18))
	tx, err := SignTx(NewTransaction(0, addr, new(big.Int), new(big.Int), new(big.Int), 0, nil, nil), signer, key)
	if err != nil {
		t.Fatal(err)
	}
	if !tx.Protected() {
		t.Fatal("expected tx to be protected")
	}

	if tx.ChainId().Cmp(signer.chainId) != 0 {
		t.Error("expected chainId to be", signer.chainId, "got", tx.ChainId())
	}

	tx = NewTransaction(0, addr, new(big.Int), new(big.Int), new(big.Int), 0, nil, nil)
	tx, err = SignTx(tx, PanguSigner{}, key)
	if err != nil {
		t.Fatal(err)
	}

	if tx.Protected() {
		t.Error("didn't expect tx to be protected")
	}

	if tx.ChainId().Sign() != 0 {
		t.Error("expected chain id to be 0 got", tx.ChainId())
	}
}

/*
 * Use new Pangu signature algorithm
 * to sign and decode the signature
 * Can be used to verify the Chain3 interface.
 */
func TestEIP155SigningPangu(t *testing.T) {
	// Test vectors come from http://vitalik.ca/files/eip155_testvec.txt
	for i, test := range []struct {
		txRlp, addr string
	}{
		// {"f864808504a817c800825208943535353535353535353535353535353535353535808025a0044852b2a670ade5407e78fb2863c51de9fcb96542a07186fe3aeda6bb8a116da0044852b2a670ade5407e78fb2863c51de9fcb96542a07186fe3aeda6bb8a116d", "0xf0f6f18bca1b28cd68e4357452947e021241e9ce"},
		// {"f864018504a817c80182a410943535353535353535353535353535353535353535018025a0489efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bcaa0489efdaa54c0f20c7adf612882df0950f5a951637e0307cdcb4c672f298b8bc6", "0x23ef145a395ea3fa3deb533b8a9e1b4c6c25d112"},
		// {"f864028504a817c80282f618943535353535353535353535353535353535353535088025a02d7c5bef027816a800da1736444fb58a807ef4c9603b7848673f7e3a68eb14a5a02d7c5bef027816a800da1736444fb58a807ef4c9603b7848673f7e3a68eb14a5", "0x2e485e0c23b4c3c542628a5f672eeab0ad4888be"},
		// {"f865038504a817c803830148209435353535353535353535353535353535353535351b8025a02a80e1ef1d7842f27f2e6be0972bb708b9a135c38860dbe73c27c3486c34f4e0a02a80e1ef1d7842f27f2e6be0972bb708b9a135c38860dbe73c27c3486c34f4de", "0x82a88539669a3fd524d669e858935de5e5410cf0"},
		// {"f865048504a817c80483019a28943535353535353535353535353535353535353535408025a013600b294191fc92924bb3ce4b969c1e7e2bab8f4c93c3fc6d0a51733df3c063a013600b294191fc92924bb3ce4b969c1e7e2bab8f4c93c3fc6d0a51733df3c060", "0xf9358f2538fd5ccfeb848b64a96b743fcc930554"},
		// {"f865058504a817c8058301ec309435353535353535353535353535353535353535357d8025a04eebf77a833b30520287ddd9478ff51abbdffa30aa90a8d655dba0e8a79ce0c1a04eebf77a833b30520287ddd9478ff51abbdffa30aa90a8d655dba0e8a79ce0c1", "0xa8f7aba377317440bc5b26198a363ad22af1f3a4"},
		// {"f866068504a817c80683023e3894353535353535353535353535353535353535353581d88025a06455bf8ea6e7463a1046a0b52804526e119b4bf5136279614e0b1e8e296a4e2fa06455bf8ea6e7463a1046a0b52804526e119b4bf5136279614e0b1e8e296a4e2d", "0xf1f571dc362a0e5b2696b8e775f8491d3e50de35"},
		// {"f867078504a817c807830290409435353535353535353535353535353535353535358201578025a052f1a9b320cab38e5da8a8f97989383aab0a49165fc91c737310e4f7e9821021a052f1a9b320cab38e5da8a8f97989383aab0a49165fc91c737310e4f7e9821021", "0xd37922162ab7cea97c97a87551ed02c9a38b7332"},
		// {"f867088504a817c8088302e2489435353535353535353535353535353535353535358202008025a064b1702d9298fee62dfeccc57d322a463ad55ca201256d01f62b45b2e1c21c12a064b1702d9298fee62dfeccc57d322a463ad55ca201256d01f62b45b2e1c21c10", "0x9bddad43f934d313c2b79ca28a432dd2b7281029"},
		//chainid = 1
		//{"f86f01808509502f900083200b20947312f4b8a4457a36827f185325fd6b66a3f8bb8b880f43fc2c04ee0000008025a0dfa37f170535100ff2c6081c326406fec05d7f9b314dfe256af0bfef6c3ba421a01ff2ce9f8fe5167da27d3de416ddc4844572ca88b9735c39307b430d1e13926a", "0x3c24d7329e92f84f08556ceb6df1cdb0104ca49f"},
		//chainid = 2
		//{"f86f01808509502f900083200b20947312f4b8a4457a36827f185325fd6b66a3f8bb8b8810a741a462780000008028a0b9ea3f64e78835789046ea794c0268bdcf5b5a4aaa2bd34e5a84e7bfbb90121ba0740c141cea882e5a48f5738ea2a24706d54ce25875c74349e8de4bb459961521", "0x3c24d7329e92f84f08556ceb6df1cdb0104ca49f"},
		//chainid =99
		// {"f87002808561c9f368008255f094d814f2ac2c4ca49b33066582e4e97ebae02f2ab9880de0b6b3a764000000808081eba0153c7b95b5fcaf393dbceab35ce03adf98eb67146fd1f9dccaff41ded5378af2a05d8f0be761803d9c5000c8a8cdb3f410abdf75a0f2fc0ac3dfa89af889b98e24", "0x7312F4B8A4457a36827f185325Fd6B66a3f8BB8B"},
		// {"f87049808477359400834c4b4094d814f2ac2c4ca49b33066582e4e97ebae02f2ab988115dd030eb16980000808081eea080ba847c6be25ff568075ad4c4b6c0c5df43ecd2359937edceaa314c5006d3aca07de4d38348db1ed9018a3b10d1e6a520a15ead920d762fb588728293d3d91cb8", "0x01560cd3bac62cc6d7e6380600d9317363400896"},
		//testnet 101,
		{"f8705680840bebc200834c4b4094d814f2ac2c4ca49b33066582e4e97ebae02f2ab988115dd030eb16980000808081eda00fdb06601eb2e2b37a41e60446c0fb2b49d7f10468d3646fc239503a3d10f6e7a072dd7c2ee7cb1c473615ecde5f2a13398ff7d93c37921dd06839482c514972a4", "0xa8863fc8Ce3816411378685223C03DAae9770ebB"},
	} {
		//Setup the signer with chainID
		//99- mainnet, 100 - devnet, 101 - testnet,
		signer := NewPanguSigner(big.NewInt(101))

		var tx *Transaction
		err := rlp.DecodeBytes(common.Hex2Bytes(test.txRlp), &tx)
		fmt.Println(tx.String())
		log.Info("[core/tx_pool.go->TxPool.add] tx=%v", tx.String())
		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		from, err := Sender(signer, tx)

		if err != nil {
			t.Errorf("%d: %v", i, err)
			continue
		}

		addr := common.HexToAddress(test.addr)
		if from != addr {
			t.Errorf("%d: expected %x got %x", i, addr, from)
		}
		fmt.Println("=======================")
	}
}

func TestChainId(t *testing.T) {
	key, _ := defaultTestKey()

	tx := NewTransaction(0, common.Address{}, new(big.Int), new(big.Int), new(big.Int), 0, nil, nil)

	var err error
	tx, err = SignTx(tx, NewPanguSigner(big.NewInt(101)), key)
	if err != nil {
		t.Fatal(err)
	}

	_, err = Sender(NewPanguSigner(big.NewInt(100)), tx)
	if err != ErrInvalidChainId {
		t.Error("expected error:", ErrInvalidChainId)
	}

	_, err = Sender(NewPanguSigner(big.NewInt(101)), tx)
	if err != nil {
		t.Error("expected no error")
	}
}
