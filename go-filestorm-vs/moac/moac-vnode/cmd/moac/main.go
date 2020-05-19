// Copyright 2014 The MOAC-core Authors
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

// moac is the official command-line client for MoacNode.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	_ "net/http/pprof"

	"github.com/filestorm/go-filestorm/moac/moac-lib/common"
	"github.com/filestorm/go-filestorm/moac/moac-lib/log"
	"github.com/filestorm/go-filestorm/moac/moac-lib/metrics"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/accounts/keystore"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/cmd/utils"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/console"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/internal/debug"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mc"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/mcclient"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/node"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/params"
	"github.com/filestorm/go-filestorm/moac/moac-vnode/vnode"
	"gopkg.in/urfave/cli.v1"
)

const (
	clientIdentifier = "moac" // Client identifier to advertise over the network
)

var (
	// Git SHA1 commit hash of the release (set via linker flags)
	gitCommit = ""
	// MoacNode address of the Moac release oracle.
	relOracle = common.HexToAddress("0xfa7b9770ca4cb04296cac84f37736d4041251cdf")
	// The app that holds all commands and flags.
	app = utils.NewApp(gitCommit, "the MOAC-vnode command line interface")
	// flags that configure the node
	nodeFlags = []cli.Flag{
		utils.IdentityFlag,
		utils.UnlockedAccountFlag,
		utils.PasswordFileFlag,
		utils.BootnodesFlag,
		utils.BootnodesV4Flag,
		utils.BootnodesV5Flag,
		utils.SubnetBootnodesFlag,
		utils.DataDirFlag,
		utils.KeyStoreDirFlag,
		utils.NoUSBFlag,
		utils.EthashCacheDirFlag,
		utils.EthashCachesInMemoryFlag,
		utils.EthashCachesOnDiskFlag,
		utils.EthashDatasetDirFlag,
		utils.EthashDatasetsInMemoryFlag,
		utils.EthashDatasetsOnDiskFlag,
		utils.TxPoolNoLocalsFlag,
		utils.TxPoolJournalFlag,
		utils.TxPoolRejournalFlag,
		utils.TxPoolPriceLimitFlag,
		utils.TxPoolPriceBumpFlag,
		utils.TxPoolAccountSlotsFlag,
		utils.TxPoolGlobalSlotsFlag,
		utils.TxPoolAccountQueueFlag,
		utils.TxPoolGlobalQueueFlag,
		utils.TxPoolLifetimeFlag,
		utils.FastSyncFlag,
		utils.LightModeFlag,
		utils.SyncModeFlag,
		utils.LightServFlag,
		utils.LightPeersFlag,
		utils.LightKDFFlag,
		utils.CacheFlag,
		utils.TrieCacheGenFlag,
		utils.ListenPortFlag,
		utils.MaxPeersFlag,
		utils.MaxPendingPeersFlag,
		utils.MoacbaseFlag,
		utils.GasPriceFlag,
		utils.MinerThreadsFlag,
		utils.MiningEnabledFlag,
		utils.TargetGasLimitFlag,
		utils.NATFlag,
		utils.NoDiscoverFlag,
		utils.DiscoveryV5Flag,
		utils.NetrestrictFlag,
		utils.NodeKeyFileFlag,
		utils.NodeKeyHexFlag,
		utils.DevModeFlag,
		utils.TestnetFlag,
		utils.VMEnableDebugFlag,
		utils.NetworkIdFlag,
		utils.RPCCORSDomainFlag,
		utils.MoacStatusURLFlag,
		utils.MetricsEnabledFlag,
		utils.FakePoWFlag,
		utils.NoCompactionFlag,
		utils.GpoBlocksFlag,
		utils.GpoPercentileFlag,
		utils.ExtraDataFlag,
		configFileFlag,
	}

	rpcFlags = []cli.Flag{
		utils.RPCEnabledFlag,
		utils.RPCListenAddrFlag,
		utils.RPCPortFlag,
		utils.RPCApiFlag,
		utils.WSEnabledFlag,
		utils.WSListenAddrFlag,
		utils.WSPortFlag,
		utils.WSApiFlag,
		utils.WSAllowedOriginsFlag,
		utils.IPCDisabledFlag,
		utils.IPCPathFlag,
	}

	vnodeConfigFlags = []cli.Flag{
		utils.VnodeConfig,
	}
	// whisperFlags = []cli.Flag{
	// 	utils.WhisperEnabledFlag,
	// 	utils.WhisperMaxMessageSizeFlag,
	// 	utils.WhisperMinPOWFlag,
	// }
)

func init() {
	// Initialize the CLI app and start Moac
	app.Action = moac
	app.HideVersion = true // we have a command to print the version
	app.Copyright = "Copyright 2017-2019 MOAC Tech Inc."
	app.Commands = []cli.Command{
		// See chaincmd.go:
		initCommand,
		importCommand,
		exportCommand,
		removedbCommand,
		dumpCommand,
		// See monitorcmd.go:
		monitorCommand,
		// See accountcmd.go:
		accountCommand,
		walletCommand,
		// See consolecmd.go:
		consoleCommand,
		attachCommand,
		javascriptCommand,
		// See misccmd.go:
		makecacheCommand,
		makedagCommand,
		versionCommand,
		bugCommand,
		licenseCommand,
		// See config.go
		dumpConfigCommand,
	}
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Flags = append(app.Flags, nodeFlags...)
	app.Flags = append(app.Flags, rpcFlags...)
	app.Flags = append(app.Flags, consoleFlags...)
	app.Flags = append(app.Flags, debug.Flags...)
	//Add custom flags
	app.Flags = append(app.Flags, vnodeConfigFlags...)
	// app.Flags = append(app.Flags, whisperFlags...)

	app.Before = func(ctx *cli.Context) error {
		runtime.GOMAXPROCS(runtime.NumCPU())
		if err := debug.Setup(ctx); err != nil {
			return err
		}
		// Start system runtime metrics collection
		go metrics.CollectProcessMetrics(3 * time.Second)

		utils.SetupNetwork(ctx)
		return nil
	}

	app.After = func(ctx *cli.Context) error {
		debug.Exit()
		console.Stdin.Close() // Resets terminal mode.
		return nil
	}
}

//Start the main node with correct version number
// 2019/07/16, removed the pprof thread
func main() {

	// Remove these lines to avoid 6060 port usage
	// go func() {
	// 	origlog.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	fmt.Printf("Start MOAC %s ...\n", params.VersionWithName)
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

// moac is the main entry point into the system if no special subcommand is ran.
// It creates a default node based on the command line arguments and runs it in
// blocking mode, waiting for it to be shut down.
func moac(ctx *cli.Context) error {
	configFilePath := ctx.GlobalString(utils.VnodeConfig.Name)
	log.Infof("vnode file %v\n", configFilePath)
	// Start the VNODE service with SCS
	go vnode.VnodeServiceStart(configFilePath)

	// Load the VNODE config
	go initLog(ctx)
	log.Debugf("[cmd/moac/main.go->moac] about to call makeFullNode ctx")
	// Create the full node with input config
	node := makeFullNode(ctx)

	// Main node started
	startNode(ctx, node)
	node.Wait()

	return nil
}

/*
 * Add build directory for the log file
 * if '_logs' directory not exists.
 */
func initLog(ctx *cli.Context) error {
	absLogPath, _ := filepath.Abs("./_logs")
	absLogPath = absLogPath + "/moac"
	if err := os.MkdirAll(filepath.Dir(absLogPath), os.ModePerm); err != nil {
		log.Errorf("Error %v in create %s\n", err, absLogPath)
		return err
	}

	lvl := log.Lvl(ctx.GlobalInt(utils.VerbosityFlag.Name))
	log.LogRotate(absLogPath, lvl)

	return nil
}

// startNode boots up the system node and all registered protocols, after which
// it unlocks any requested accounts, and starts the RPC/IPC interfaces and the
// miner.
func startNode(ctx *cli.Context, node *node.Node) {
	// Start up the node itself
	utils.StartNode(node)

	// Unlock any account specifically requested
	ks := node.AccountManager().Backends(keystore.KeyStoreType)[0].(*keystore.KeyStore)

	passwords := utils.MakePasswordList(ctx)
	unlocks := strings.Split(ctx.GlobalString(utils.UnlockedAccountFlag.Name), ",")
	for i, account := range unlocks {
		if trimmed := strings.TrimSpace(account); trimmed != "" {
			unlockAccount(ctx, ks, trimmed, i, passwords)
		}
	}
	// Register wallet event handlers to open and auto-derive wallets
	events := make(chan accounts.WalletEvent, 16)
	node.AccountManager().Subscribe(events)

	go func() {
		// Create an chain state reader for self-derivation
		rpcClient, err := node.Attach()
		if err != nil {
			utils.Fatalf("Failed to attach to self: %v", err)
		}
		stateReader := mcclient.NewClient(rpcClient)

		// Open any wallets already attached
		for _, wallet := range node.AccountManager().Wallets() {
			if err := wallet.Open(""); err != nil {
				log.Warn("Failed to open wallet", "url", wallet.URL(), "err", err)
			}
		}
		// Listen for wallet event till termination
		for event := range events {
			switch event.Kind {
			case accounts.WalletArrived:
				if err := event.Wallet.Open(""); err != nil {
					log.Warn("New wallet appeared, failed to open", "url", event.Wallet.URL(), "err", err)
				}
			case accounts.WalletOpened:
				status, _ := event.Wallet.Status()
				log.Infof("New wallet appeared url=%v status=%v", event.Wallet.URL(), status)

				if event.Wallet.URL().Scheme == "ledger" {
					event.Wallet.SelfDerive(accounts.DefaultLedgerBaseDerivationPath, stateReader)
				} else {
					event.Wallet.SelfDerive(accounts.DefaultBaseDerivationPath, stateReader)
				}

			case accounts.WalletDropped:
				log.Infof("Old wallet dropped url=%v", event.Wallet.URL())
				event.Wallet.Close()
			}
		}
	}()
	// Start auxiliary services if enabled
	if ctx.GlobalBool(utils.MiningEnabledFlag.Name) {
		// Mining only makes sense if a full MoacService node is running
		// get the running instance of moacservice inside node and assign it to moac
		var moacServ *mc.MoacService
		if err := node.Service(&moacServ); err != nil {
			utils.Fatalf("moac service not running: %v", err)
		}
		// Use a reduced number of threads if requested
		if threads := ctx.GlobalInt(utils.MinerThreadsFlag.Name); threads > 0 {
			type threaded interface {
				SetThreads(threads int)
			}
			if th, ok := moacServ.Engine().(threaded); ok {
				th.SetThreads(threads)
			}
		}
		// Set the gas price to the limits from the CLI and start mining
		moacServ.TxPool().SetGasPrice(utils.GlobalBig(ctx, utils.GasPriceFlag.Name))
		if err := moacServ.StartMining(true); err != nil {
			utils.Fatalf("Failed to start mining: %v", err)
		}
	}
}
