// Copyright 2019 The go-filestorm Authors
// This file is part of go-filestorm.
//
// go-filestorm is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-filestorm is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-filestorm. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/filestorm/go-filestorm/params"
)

const (
	ipcAPIs  = "admin:1.0 debug:1.0 fst:1.0 fstash:1.0 miner:1.0 net:1.0 personal:1.0 rpc:1.0 shh:1.0 txpool:1.0 web3:1.0"
	httpAPIs = "fst:1.0 net:1.0 rpc:1.0 web3:1.0"
)

// Tests that a node embedded within a console can be started up properly and
// then terminated by closing the input stream.
func TestConsoleWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"

	// Start a storm console, make sure it's cleaned up and terminate the console
	storm := runStorm(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--fsterbase", coinbase, "--shh",
		"console")

	// Gather all the infos the welcome message needs to contain
	storm.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	storm.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	storm.SetTemplateFunc("gover", runtime.Version)
	storm.SetTemplateFunc("gethver", func() string { return params.VersionWithCommit("", "") })
	storm.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	storm.SetTemplateFunc("apis", func() string { return ipcAPIs })

	// Verify the actual welcome message to the required template
	storm.Expect(`
Welcome to the Storm JavaScript console!

instance: Storm/v{{gethver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{.Fsterbase}}
at block: 0 ({{niltime}})
 datadir: {{.Datadir}}
 modules: {{apis}}

> {{.InputLine "exit"}}
`)
	storm.ExpectExit()
}

// Tests that a console can be attached to a running node via various means.
func TestIPCAttachWelcome(t *testing.T) {
	// Configure the instance for IPC attachement
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	var ipc string
	if runtime.GOOS == "windows" {
		ipc = `\\.\pipe\gfst` + strconv.Itoa(trulyRandInt(100000, 999999))
	} else {
		ws := tmpdir(t)
		defer os.RemoveAll(ws)
		ipc = filepath.Join(ws, "storm.ipc")
	}
	// Note: we need --shh because testAttachWelcome checks for default
	// list of ipc modules and shh is included there.
	storm := runStorm(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--fsterbase", coinbase, "--shh", "--ipcpath", ipc)

	waitForEndpoint(t, ipc, 3*time.Second)
	testAttachWelcome(t, storm, "ipc:"+ipc, ipcAPIs)

	storm.Interrupt()
	storm.ExpectExit()
}

func TestHTTPAttachWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P
	storm := runStorm(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--fsterbase", coinbase, "--rpc", "--rpcport", port)

	endpoint := "http://127.0.0.1:" + port
	waitForEndpoint(t, endpoint, 3*time.Second)
	testAttachWelcome(t, storm, endpoint, httpAPIs)

	storm.Interrupt()
	storm.ExpectExit()
}

func TestWSAttachWelcome(t *testing.T) {
	coinbase := "0x8605cdbbdb6d264aa742e77020dcbc58fcdce182"
	port := strconv.Itoa(trulyRandInt(1024, 65536)) // Yeah, sometimes this will fail, sorry :P

	storm := runStorm(t,
		"--port", "0", "--maxpeers", "0", "--nodiscover", "--nat", "none",
		"--fsterbase", coinbase, "--ws", "--wsport", port)

	endpoint := "ws://127.0.0.1:" + port
	waitForEndpoint(t, endpoint, 3*time.Second)
	testAttachWelcome(t, storm, endpoint, httpAPIs)

	storm.Interrupt()
	storm.ExpectExit()
}

func testAttachWelcome(t *testing.T, storm *testgeth, endpoint, apis string) {
	// Attach to a running storm note and terminate immediately
	attach := runStorm(t, "attach", endpoint)
	defer attach.ExpectExit()
	attach.CloseStdin()

	// Gather all the infos the welcome message needs to contain
	attach.SetTemplateFunc("goos", func() string { return runtime.GOOS })
	attach.SetTemplateFunc("goarch", func() string { return runtime.GOARCH })
	attach.SetTemplateFunc("gover", runtime.Version)
	attach.SetTemplateFunc("gethver", func() string { return params.VersionWithCommit("", "") })
	attach.SetTemplateFunc("fsterbase", func() string { return storm.Fsterbase })
	attach.SetTemplateFunc("niltime", func() string { return time.Unix(0, 0).Format(time.RFC1123) })
	attach.SetTemplateFunc("ipc", func() bool { return strings.HasPrefix(endpoint, "ipc") })
	attach.SetTemplateFunc("datadir", func() string { return storm.Datadir })
	attach.SetTemplateFunc("apis", func() string { return apis })

	// Verify the actual welcome message to the required template
	attach.Expect(`
Welcome to the Storm JavaScript console!

instance: Storm/v{{gethver}}/{{goos}}-{{goarch}}/{{gover}}
coinbase: {{fsterbase}}
at block: 0 ({{niltime}}){{if ipc}}
 datadir: {{datadir}}{{end}}
 modules: {{apis}}

> {{.InputLine "exit" }}
`)
	attach.ExpectExit()
}

// trulyRandInt generates a crypto random integer used by the console tests to
// not clash network ports with other tests running cocurrently.
func trulyRandInt(lo, hi int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(hi-lo)))
	return int(num.Int64()) + lo
}
