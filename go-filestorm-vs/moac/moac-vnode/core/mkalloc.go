// Copyright 2017 The MOAC-core Authors
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

// +build none

/*

   The mkalloc tool creates the genesis allocation constants in genesis_alloc.go
   It outputs a const declaration that contains an RLP-encoded list of (address, balance) tuples.

       go run mkalloc.go genesis.json

*/
package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"

	"github.com/filestorm/go-filestorm/moac/moac-vnode/core"
	"github.com/filestorm/go-filestorm/moac/moac-lib/rlp"
)

type allocItem struct {
	Addr, Balance *big.Int
	Code          []byte
}

type allocList []allocItem

func (a allocList) Len() int           { return len(a) }
func (a allocList) Less(i, j int) bool { return a[i].Addr.Cmp(a[j].Addr) < 0 }
func (a allocList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

/*
 * Added the code section
 */
func makelist(g *core.Genesis) allocList {
	a := make(allocList, 0, len(g.Alloc))
	for addr, account := range g.Alloc {
		fmt.Println("account storage:", account.Storage)
		fmt.Println("account code:", account.Code)
		fmt.Println("account Nonce:", account.Nonce)
		// if len(account.Storage) > 0 || len(account.Code) > 0 || account.Nonce != 0 {
		if len(account.Storage) > 0 || account.Nonce != 0 {
			panic(fmt.Sprintf("can't encode account %x", addr))
		}
		a = append(a, allocItem{addr.Big(), account.Balance, account.Code})
	}
	sort.Sort(a)
	return a
}

func makealloc(g *core.Genesis) string {
	a := makelist(g)
	data, err := rlp.EncodeToBytes(a)
	if err != nil {
		panic(err)
	}
	return strconv.QuoteToASCII(string(data))
}

/*
 * Added file outputs
 file, err := os.OpenFile("output.txt", os.O_WRONLY|os.O_CREATE, 0666)
    if err != nil {
        fmt.Println("File does not exists or cannot be created")
        os.Exit(1)
    }
    defer file.Close()

    w := bufio.NewWriter(file)
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    i := r.Perm(5)
    fmt.Fprintf(w, "%v\n", i)

    w.Flush()
*/
func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: mkalloc genesis.json")
		os.Exit(1)
	}

	g := new(core.Genesis)
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	if err := json.NewDecoder(file).Decode(g); err != nil {
		panic(err)
	}
	fmt.Println("const devAllocData =", makealloc(g))
	return
	//Prepare an output file and send the RLP encoded string
	//This file need to be put into existing genesis_alloc.go
	//to be used
	outfile, err := os.OpenFile("genesis_tmpout.txt", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("File does not exists or cannot be created")
		os.Exit(1)
	}
	defer outfile.Close()
	w := bufio.NewWriter(outfile)
	fmt.Fprintf(w, "const devAllocData =%v\n", makealloc(g))

	w.Flush()
	// fmt.FPrintln("const devAllocData =", makealloc(g))
}
