set binFolder to "~/go/src/github.com/innowells/moac-vnode/build/bin"
set moacFolder to "~/go/src/github.com/innowells/moac-vnode/cmd/moac"
set scsFolder to "~/go/src/github.com/innowells/moac-scs/scsserver"
set scsBuildFolder to "~/go/src/github.com/innowells/moac-scs/build"
set jsFolder to "~/go/src/github.com/innowells/moac-vnode/solidity/"
set walletFolder to "~/Documents/Code/MyMoacWallet/app"
set testSolFilename to "SubChainCombinedTest.sol"
set contractProtocolBaseName to "SubChainProtocolBase.sol"
set contractBaseName to "SubChainBase.sol"
set contractName to "SubChainP1.sol"

set runTill to 4

tell application "iTerm"
	activate
	tell current window
		if runTill ³ 1 then
			--1) copy all necessary files 
			create tab with default profile
			set prepTab to current tab
			tell current session of prepTab
				write text "echo -ne \"\\033]0;\"prep\"\\007\""
				-- 1.1) compile all go projects
				write text "cd " & scsFolder
				write text "/usr/local/Cellar/go/1.10/libexec/bin/go build -i -o " & scsBuildFolder & "/scs_server ."
				write text "cd " & moacFolder
				write text "/usr/local/Cellar/go/1.10/libexec/bin/go build -i -buildmode=exe -o " & moacFolder & "/moac ."
				
				-- 1.2) copy files to running folders
				write text "mkdir " & binFolder
				write text "cd " & binFolder
				write text "mkdir moac"
				write text "mkdir moac/_logs"
				write text "mkdir moac2"
				write text "mkdir moac2/_logs"
				write text "mkdir moac3"
				write text "mkdir moac3/_logs"
				
				write text "mkdir scs1"
				write text "mkdir scs1/config"
				write text "mkdir scs1/server"
				write text "mkdir scs1/_logs"
				write text "mkdir scs2"
				write text "mkdir scs2/config"
				write text "mkdir scs2/server"
				write text "mkdir scs2/_logs"
				
				write text "cp ../../cmd/moac/moac ./moac/"
				write text "cp ../../genesis_testnet.json ./moac/genesis.json"
				write text "cp ../../cmd/moac/moac ./moac2/"
				write text "cp ../../genesis_testnet.json ./moac2/genesis.json"
				write text "cp ../../cmd/moac/moac ./moac3/"
				write text "cp ../../genesis_testnet.json ./moac3/genesis.json"
				write text "cp " & scsBuildFolder & "/scs_server ./scs1/server/"
				write text "cp " & scsBuildFolder & "/scs_server ./scs2/server/"
				
				write text "rm ./moac/_logs/moac.log"
				write text "rm ./moac2/_logs/moac.log"
				write text "rm ./moac3/_logs/moac.log"
				delay 1
				write text "rm -rf ./moac/data/"
				delay 1
				write text "rm -rf ./moac2/data/"
				delay 1
				write text "rm -rf ./moac3/data/"
				delay 1
				write text "rm ./scs1/server/_logs/scs.log"
				write text "rm ./scs2/server/_logs/scs.log"
				
			end tell
		end if
		if runTill ³ 2 then
			create tab with default profile
			set moac1Tab to current tab
			tell current session of moac1Tab
				write text "echo -ne \"\\033]0;\"moac 1\"\\007\""
				write text "cd " & binFolder
				write text "cd moac"
				write text "./moac --datadir \"" & binFolder & "/moac/data\" init ./genesis.json"
				delay 5
				write text "./moac --datadir \"" & binFolder & "/moac/data\" --port 30301 --exec \"chain3.admin.nodeInfo.enode\" | tail -n 1 > enode.txt"
			end tell
		end if
		delay 5
		if runTill ³ 3 then
			create tab with default profile
			set moac1Tab to current tab
			tell current session of moac1Tab
				write text "echo -ne \"\\033]0;\"moac 2\"\\007\""
				write text "cd " & binFolder
				write text "cd moac2"
				write text "./moac --datadir \"" & binFolder & "/moac2/data\" init ./genesis.json"
				delay 5
				write text "./moac --datadir \"" & binFolder & "/moac2/data\" --port 30302 --exec \"chain3.admin.nodeInfo.enode\" | tail -n 1 > enode.txt"
			end tell
		end if
		delay 5
		if runTill ³ 4 then
			create tab with default profile
			set moac1Tab to current tab
			tell current session of moac1Tab
				write text "echo -ne \"\\033]0;\"moac 3\"\\007\""
				write text "cd " & binFolder
				write text "cd moac3"
				write text "./moac --datadir \"" & binFolder & "/moac3/data\" init ./genesis.json"
				delay 5
				write text "./moac --datadir \"" & binFolder & "/moac3/data\" --port 30303 -exec \"chain3.admin.nodeInfo.enode\" | tail -n 1 > enode.txt"
			end tell
		end if
	end tell
end tell