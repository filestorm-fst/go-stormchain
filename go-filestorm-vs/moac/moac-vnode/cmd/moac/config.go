// Copyright 2017 The MOAC-core Authors
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

package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"unicode"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/cmd/utils"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/contracts/release"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/node"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
	"github.com/naoina/toml"
)

var (
	dumpConfigCommand = cli.Command{
		Action:      utils.MigrateFlags(dumpConfig),
		Name:        "dumpconfig",
		Usage:       "Show configuration values",
		ArgsUsage:   "",
		Flags:       append(nodeFlags, rpcFlags...),
		Category:    "MISCELLANEOUS COMMANDS",
		Description: `The dumpconfig command shows configuration values.`,
	}

	configFileFlag = cli.StringFlag{
		Name:  "config",
		Usage: "TOML configuration file",
	}
)

// These settings ensure that TOML keys use the same names as Go struct fields.
var tomlSettings = toml.Config{
	NormFieldName: func(rt reflect.Type, key string) string {
		return key
	},
	FieldToKey: func(rt reflect.Type, field string) string {
		return field
	},
	MissingField: func(rt reflect.Type, field string) error {
		link := ""
		if unicode.IsUpper(rune(rt.Name()[0])) && rt.PkgPath() != "main" {
			link = fmt.Sprintf(", see https://godoc.org/%s#%s for available fields", rt.PkgPath(), rt.Name())
		}
		return fmt.Errorf("field '%s' is not defined in %s%s", field, rt.String(), link)
	},
}

type mcstatsConfig struct {
	URL string `toml:",omitempty"`
}

type moacConfig struct {
	Mc      mc.Config
	Node    node.Config
	Mcstats mcstatsConfig
}

func loadNodeConfig(file string, cfg *moacConfig) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tomlSettings.NewDecoder(bufio.NewReader(f)).Decode(cfg)
	// Add file name to errors that have a line number.
	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(file + ", " + err.Error())
	}
	return err
}

//Remove the shh/scs as default modules
// Added VNODE module to HTTP modules
func defaultNodeConfig() node.Config {
	cfg := node.DefaultConfig
	cfg.Name = clientIdentifier
	cfg.Version = params.VersionWithCommit(gitCommit)
	//Added VNODE to the HTTP module
	cfg.HTTPModules = append(cfg.HTTPModules, "mc", "vnode")
	cfg.WSModules = append(cfg.WSModules, "mc")
	cfg.IPCPath = "moac.ipc"
	return cfg
}

func makeNodeFromConfig(ctx *cli.Context) (*node.Node, moacConfig) {
	// Load defaults.
	cfg := moacConfig{
		Mc:   mc.DefaultConfig,
		Node: defaultNodeConfig(),
	}

	// Load config file.
	if file := ctx.GlobalString(configFileFlag.Name); file != "" {
		log.Debugf("makeNodeFromConfig config %v", file)
		if err := loadNodeConfig(file, &cfg); err != nil {
			utils.Fatalf("%v", err)
		}
	}

	// Apply flags.
	utils.SetNodeConfig(ctx, &cfg.Node)
	//Create the node with the config
	n, err := node.New(&cfg.Node)
	if err != nil {
		utils.Fatalf("Failed to create the protocol stack: %v", err)
	}

	utils.SetMoacConfig(ctx, n, &cfg.Mc)
	if ctx.GlobalIsSet(utils.MoacStatusURLFlag.Name) {
		cfg.Mcstats.URL = ctx.GlobalString(utils.MoacStatusURLFlag.Name)
	}

	return n, cfg
}

func makeFullNode(ctx *cli.Context) *node.Node {
	// generate node.Node object
	n, cfg := makeNodeFromConfig(ctx)
	log.Infof("Connect to networkid = %v", cfg.Mc.NetworkId)
	// register moac service to the node instance
	// which includes all p2p servers and the protocolmanager
	utils.RegisterMoacService(n, &cfg.Mc)

	// Add the MoacNode Stats daemon if requested.
	if cfg.Mcstats.URL != "" {
		utils.RegisterMcStatsService(n, cfg.Mcstats.URL)
	}

	// Add the release oracle service so it boots along with node.
	if err := n.Register(func(ctx *node.ServiceContext) (node.Service, error) {
		config := release.Config{
			Oracle: relOracle,
			Major:  uint32(params.VersionMajor),
			Minor:  uint32(params.VersionMinor),
			Patch:  uint32(params.VersionPatch),
		}
		commit, _ := hex.DecodeString(gitCommit)
		copy(config.Commit[:], commit)
		return release.NewReleaseService(ctx, config)
	}); err != nil {
		utils.Fatalf("Failed to register the Moac release oracle service: %v", err)
	}
	return n
}

// dumpConfig is the dumpconfig command.
func dumpConfig(ctx *cli.Context) error {
	_, cfg := makeNodeFromConfig(ctx)
	comment := ""

	if cfg.Mc.Genesis != nil {
		cfg.Mc.Genesis = nil
		comment += "# Note: this config doesn't contain the genesis block.\n\n"
	}

	out, err := tomlSettings.Marshal(&cfg)
	if err != nil {
		return err
	}
	io.WriteString(os.Stdout, comment)
	os.Stdout.Write(out)
	return nil
}
