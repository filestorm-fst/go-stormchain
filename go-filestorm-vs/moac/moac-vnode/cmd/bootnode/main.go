// Copyright 2015 The MOAC-core Authors
// This file is part of MOAC-core.
//
// MOAC-core is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// MOAC-core is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with MOAC-core. If not, see <http://www.gnu.org/licenses/>.

// bootnode runs a bootstrap node for the MoacNode Discovery Protocol.
package main

import (
	"crypto/ecdsa"
	"crypto/sha512"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/crypto"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/cmd/utils"
	discover "github.com/filestorm/go-filestorm/moac/moac-vnode/p2p/discover"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/p2p/discv5"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/p2p/nat"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/p2p/netutil"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
)

func main() {
	var (
		listenAddr      = flag.String("addr", ":30301", "listen address")
		genKey          = flag.String("genkey", "", "generate a node key")
		genKeyForTarget = flag.String("genkeyfortarget", "", "generate a node key which is closest to the target")
		vmodule         = flag.String("vmodule", "", "log verbosity pattern")
		nodeKeyFile     = flag.String("nodekey", "", "private key filename")
		nodeKeyHex      = flag.String("nodekeyhex", "", "private key as hex (for testing)")
		natdesc         = flag.String("nat", "none", "port mapping mechanism (any|none|upnp|pmp|extip:<IP>)")
		netrestrict     = flag.String("netrestrict", "", "restrict network communication to the given IP networks (CIDR masks)")
		runv5           = flag.Bool("v5", false, "run a v5 topic discovery bootnode")
		writeAddr       = flag.Bool("writeaddress", false, "write out the node's pubkey hash and quit")
		verbosity       = flag.Int("verbosity", int(log.LvlInfo), "log verbosity (0-9)")
		networkid       = flag.Int("networkid", params.MainNetworkId, "network id for this bootnode, default to 99")

		nodeKey *ecdsa.PrivateKey
		err     error
	)
	flag.Parse()

	glogger := log.NewGlogHandler(log.StreamHandler(os.Stderr, log.TerminalFormat(false)))
	glogger.Verbosity(log.Lvl(*verbosity))
	glogger.Vmodule(*vmodule)
	log.Root().SetHandler(glogger)
	NetworkID := params.MainNetworkId

	natm, err := nat.Parse(*natdesc)
	if err != nil {
		utils.Fatalf("-nat: %v", err)
	}
	switch {
	case *networkid != params.MainNetworkId:
		NetworkID = *networkid
	case *genKey != "":
		nodeKey, err = crypto.GenerateKey()
		if err != nil {
			utils.Fatalf("could not generate key: %v", err)
		}
		if err = crypto.SaveECDSA(*genKey, nodeKey); err != nil {
			utils.Fatalf("%v", err)
		}
		return
	case *genKeyForTarget != "":
		nodeIDToPrikey := make(map[discover.NodeID]*ecdsa.PrivateKey)
		// generate target based on subnet id
		var subnetID discover.NodeID
		subnetHash512 := sha512.New().Sum([]byte(*genKeyForTarget))
		copy(subnetID[:], subnetHash512[:])
		target := crypto.Keccak256Hash(subnetID[:])
		fmt.Printf("target input: %s\n", *genKeyForTarget)
		fmt.Printf("target: %v\n", common.Bytes2Hex(target[:]))
		top := 10
		n := 200000
		randomKeys := []*discover.Node{}
		close := &discover.NodesByDistance{Target: target}
		for prikey := range utils.GeneratePrikeys(n) {
			// convert prikey to node
			nodeid := discover.PubkeyID(&prikey.PublicKey)
			node := &discover.Node{
				IP:  nil,
				UDP: 0,
				TCP: 0,
				ID:  nodeid,
			}
			node.SetSha()
			close.Push(node, top)
			nodeIDToPrikey[node.ID] = prikey

			// also generate n random keys
			if len(randomKeys) < top {
				randomKeys = append(randomKeys, node)
			}
		}

		fmt.Printf("Closest node(private) keys:\n")
		for i, n := range close.GetEntries() {
			n_sha := crypto.Keccak256Hash(n.ID[:])
			fmt.Printf(
				"%d, private key: %v, verify: %v\n",
				i+1,
				hex.EncodeToString(crypto.FromECDSA(nodeIDToPrikey[n.ID])),
				hex.EncodeToString(n_sha[:]),
			)
		}
		fmt.Print("\nRandom node(private) keys:\n")
		for i, n := range randomKeys {
			fmt.Printf(
				"%d, private key: %v, address: %x\n",
				i+1,
				hex.EncodeToString(crypto.FromECDSA(nodeIDToPrikey[n.ID])),
				crypto.PubkeyToAddress(nodeIDToPrikey[n.ID].PublicKey),
			)
		}

		return
	case *nodeKeyFile == "" && *nodeKeyHex == "":
		utils.Fatalf("Use -nodekey or -nodekeyhex to specify a private key")
	case *nodeKeyFile != "" && *nodeKeyHex != "":
		utils.Fatalf("Options -nodekey and -nodekeyhex are mutually exclusive")
	case *nodeKeyFile != "":
		if nodeKey, err = crypto.LoadECDSA(*nodeKeyFile); err != nil {
			utils.Fatalf("-nodekey: %v", err)
		}
	case *nodeKeyHex != "":
		if nodeKey, err = crypto.HexToECDSA(*nodeKeyHex); err != nil {
			utils.Fatalf("-nodekeyhex: %v", err)
		}
	}

	if *writeAddr {
		fmt.Printf("%v\n", discover.PubkeyID(&nodeKey.PublicKey))
		os.Exit(0)
	}

	var restrictList *netutil.Netlist
	if *netrestrict != "" {
		restrictList, err = netutil.ParseNetlist(*netrestrict)
		if err != nil {
			utils.Fatalf("-netrestrict: %v", err)
		}
	}

	if *runv5 {
		if _, err := discv5.ListenUDP(nodeKey, *listenAddr, natm, "", restrictList); err != nil {
			utils.Fatalf("%v", err)
		}
	} else {
		if _, err := discover.ListenUDP(nodeKey, *listenAddr, natm, "", restrictList, uint64(NetworkID), false); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	select {}
}
