set binFolder to "~/go/src/github.com/innowells/moac-vnode/build/bin"
set moacFolder to "~/go/src/github.com/innowells/moac-vnode/cmd/moac"
set scsFolder to "~/go/src/github.com/innowells/moac-scs/scsserver"
set scsBuildFolder to "~/go/src/github.com/innowells/moac-scs/build"
set jsFolder to "~/go/src/github.com/innowells/moac-vnode/solidity/"
set testSolFilename to "SubChainCombinedTest.sol"
set contractProtocolBaseName to "SubChainProtocolBase.sol"
set contractBaseName to "SubChainBase.sol"
set contractName to "SubChainP1.sol"

set runTill to 10

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
				write text "/usr/local/Cellar/go/1.10/libexec/bin/go build -i -buildmode=exe -o " & scsBuildFolder & "/scs_server ."
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
				write text "rm ./scs1/server/_logs/scs.log"
				write text "rm ./scs2/server/_logs/scs.log"
				
			end tell
		end if
		if runTill ³ 2 then
			create tab with default profile
			set compileTab to current tab
			tell current session of compileTab
				write text "echo -ne \"\\033]0;\"compile\"\\007\""
				-- 1) compile solidity contract
				write text "cd " & jsFolder
				write text "node compileScript.js ./" & testSolFilename
				-- write text "node compileScript.js ./" & contractBaseName
				-- write text "node compileScript.js ./" & contractName
			end tell
			delay 20
		end if
		if runTill ³ 3 then
			create tab with default profile
			set moac1Tab to current tab
			tell current session of moac1Tab
				write text "echo -ne \"\\033]0;\"moac 1\"\\007\""
				write text "cd " & binFolder
				write text "cd moac"
				write text "./moac --datadir \"" & binFolder & "/moac/data\" --verbosity 4 --port 30301 --rpcport 8101 --jspath " & jsFolder
			end tell
			delay 5
		end if
		if runTill ³ 4 then
			create tab with default profile
			set moac1aTab to current tab
			tell current session of moac1aTab
				write text "echo -ne \"\\033]0;\"moac 1a\"\\007\""
				write text "cd " & binFolder
				write text "cd moac"
				write text "./moac attach ipc:/Users/biajee/go/src/github.com/innowells/moac-vnode/build/bin/moac/data/moac.ipc"
				delay 1
				write text "personal.newAccount(\"test123\")"
				write text "miner.start(1)"
			end tell
			delay 1
		end if
		if runTill ³ 5 then
			create tab with default profile
			set moac2Tab to current tab
			tell current session of moac2Tab
				write text "echo -ne \"\\033]0;\"moac 2\"\\007\""
				write text "cd " & binFolder
				write text "cd moac2"
				write text "./moac --datadir \"" & binFolder & "/moac2/data\" --verbosity 4 --port 30302 --rpcport 8102 --jspath " & jsFolder
			end tell
			delay 5
		end if
		if runTill ³ 6 then
			create tab with default profile
			set moac2aTab to current tab
			tell current session of moac2aTab
				write text "echo -ne \"\\033]0;\"moac 2a\"\\007\""
				write text "cd " & binFolder
				write text "cd moac2"
				write text "./moac attach ipc:/Users/biajee/go/src/github.com/innowells/moac-vnode/build/bin/moac2/data/moac.ipc"
				delay 1
				write text "admin.addPeer(\"enode://17ad78ad3b84091719f8b1464b02177587071b8eca482a926434fe5f736deaa4de4312851573c61d996be1386fa2c14807e79c72fb985394902c472ffe48a201@[::]:30301\")"
				delay 10
				write text "personal.newAccount(\"test123\")"
				write text "miner.start(1)"
			end tell
			delay 1
		end if
		if runTill ³ 7 then
			create tab with default profile
			set moac3Tab to current tab
			tell current session of moac3Tab
				write text "echo -ne \"\\033]0;\"moac 3\"\\007\""
				write text "cd " & binFolder
				write text "cd moac3"
				write text "./moac --datadir \"" & binFolder & "/moac3/data\" --verbosity 4 --port 30303 --rpcport 8103 --jspath " & jsFolder
			end tell
			delay 5
		end if
		if runTill ³ 8 then
			--			create window with default profile
			create tab with default profile
			set scs1Tab to current tab
			tell current session of scs1Tab
				write text "echo -ne \"\\033]0;\"scs 1\"\\007\""
				write text "cd " & binFolder
				write text "cd scs1/server"
				write text "./scs_server"
			end tell
			delay 1
		end if
		if runTill ³ 9 then
			create tab with default profile
			set scs2Tab to current tab
			tell current session of scs2Tab
				write text "echo -ne \"\\033]0;\"scs 2\"\\007\""
				write text "cd " & binFolder
				write text "cd scs2/server"
				write text "./scs_server"
			end tell
			delay 5
		end if
		if runTill ³ 10 then
			create tab with default profile
			set moac3aTab to current tab
			tell current session of moac3aTab
				write text "echo -ne \"\\033]0;\"moac 3a\"\\007\""
				write text "cd " & binFolder
				write text "cd moac3"
				write text "./moac attach ipc:/Users/biajee/go/src/github.com/innowells/moac-vnode/build/bin/moac3/data/moac.ipc --jspath " & jsFolder
				delay 1
				write text "admin.addPeer(\"enode://17ad78ad3b84091719f8b1464b02177587071b8eca482a926434fe5f736deaa4de4312851573c61d996be1386fa2c14807e79c72fb985394902c472ffe48a201@[::]:30301\")"
				write text "admin.addPeer(\"enode://e0f0a386c8885568a5af96cb8979d24b23e789df79acaef9d227f46cd30069fc02434e9de22cfb39b8db78a76dbe9b34a5399002663b8b6c72e61a4b9ecd384d@[::]:30302\")"
				delay 10
				write text "personal.newAccount(\"test123\")"
				write text "miner.start(1)"
				delay 10
				write text "loadScript('./subchainTest.js')"
				-- write text "simpleStorageContractTest()"
				-- write text "fullSubchainContractTest()"
				write text "directCallTest()"
			end tell
		end if
	end tell
end tell