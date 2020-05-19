var cache = [
  '',
  ' ',
  '  ',
  '   ',
  '    ',
  '     ',
  '      ',
  '       ',
  '        ',
  '         '
];

function leftPad (str, len, ch) {
	// convert `str` to `string`
	str = str + '';
	// `len` is the `pad`'s length now
	len = len - str.length;
	// doesn't need to pad
	if (len <= 0) return str;
	// `ch` defaults to `' '`
	if (!ch && ch !== 0) ch = ' ';
	// convert `ch` to `string`
	ch = ch + '';
	// cache common use cases
	if (ch === ' ' && len < 10) return cache[len] + str;
	// `pad` starts with an empty string
	var pad = '';
	// loop
	while (true) {
		// add `ch` to `pad` if `len` is odd
		if (len & 1) pad += ch;
		// divide `len` by 2, ditch the remainder
		len >>= 1;
		// "double" the `ch` so this operation count grows logarithmically on `len`
		// each time `ch` is "doubled", the `len` would need to be "doubled" too
		// similar to finding a value in binary search tree, hence O(log(n))
		if (len) ch += ch;
		// `len` is 0, exit the loop
		else break;
	}
	// pad `str`!
	return pad + str;
}

function sendtx(src, tgtaddr, amount, strData, callback) {

	//var amt = leftPad(chain3.toHex(chain3.toWei(amount)).slice(2).toString(16),64,0);
	//var strData = '';
		
	chain3.mc.sendTransaction(
		{
			from: src,
			value:chain3.toWei(amount,'mcer'),
			to: tgtaddr,
			gas: "100000",
			gasPrice: chain3.mc.gasPrice,
			data: strData
		},
		callback);
		
	console.log('sending from:' + 	src + ' to:' + tgtaddr  + ' with data:' + strData);

}


//SubChainBase functions
function registerSCS1(SubChainP1Address, SubChainP1Abi, callback)
{
	console.log("starting registerSCS1");
	var scsAccount = '';
	if (mc.accounts.length <= 1) {
		console.log("here10");
		scsAccount = personal.newAccount("test123");
	} else {
		console.log("here11");
		scsAccount = mc.accounts[1];
	} 

	console.log("here12");
	chain3.personal.unlockAccount(scsAccount,'test123',0);

	console.log("here13");
	MyContract1 = chain3.mc.contract(SubChainP1Abi);
	console.log("here14", SubChainP1Address);
	var contractInstance = MyContract1.at(SubChainP1Address);

	console.log("here15");
	var value = contractInstance.RegisterAsSCS.call().then(function(result) {
		if (callback) {
			callback(result);
		}
	});

	console.log("registerSCS1 returned", value);

	return value;
}

function registerSCS2(SubChainP1Address, SubChainP1Abi)
{
	console.log("starting registerSCS2");
	var scsAccount = '';
	if (mc.accounts.length <= 2) {
		console.log("here20");
		scsAccount = personal.newAccount("test123");
	} else {
		console.log("here21");
		scsAccount = mc.accounts[2];
	} 

	console.log("here22");
	chain3.personal.unlockAccount(scsAccount,'test123',0);

	console.log("here23");
	MyContract2 = chain3.mc.contract(SubChainP1Abi);
	console.log("here24", SubChainP1Address);
	var contractInstance = MyContract2.at(SubChainP1Address);

	console.log("here25");
	var value = contractInstance.RegisterAsSCS.call();

	console.log("registerSCS2 returned", value);

	return value;
}

function callFuncCode(SubChainP1Address, SubChainP1Abi, input)
{
	console.log("starting callFuncCode");
	var mainAccount = mc.coinbase;

	chain3.personal.unlockAccount(mainAccount,'test123',0);

	var myContract = chain3.mc.contract(SubChainP1Abi);
	var contractInstance = myContract.at(SubChainP1Address);

		var contractInstance = MyContract.new(
		{from: mainAccount, data: input, gas: 4000000, queryFlag: 0, shardingFlag: 1},
		function(e, contract) {
			console.log("callFuncCode responded");
			console.log(e, contract);
			if (typeof contract.address !== 'undefined') {
				var event = contract.logging();

			    // monitor
			    event.watch(function(error, result){
				    console.log("Event are as following:-------");
				    
				     console.log(JSON.stringify(result));
				     console.log(result.args.msg);

				    console.log("Event ending-------");
				});


		        // contract.test.call();

			}
		}
	);


	console.log("callFuncCode returned", value);

	return value;
}


function retrieveFuncCode(SubChainP1Address, SubChainP1Abi)
{
	console.log("starting retrieveFuncCode");
	var mainAccount = mc.coinbase;

	chain3.personal.unlockAccount(mainAccount,'test123',0);

	var myContract = chain3.mc.contract(SubChainP1Abi);
	var contractInstance = myContract.at(SubChainP1Address);

	var value = contractInstance.get(FuncToHex("funccode")).call();

	console.log("retrieveFuncCode returned", value);

	return value;
}


function registerClose(SubChainP1Address, SubChainP1Abi)
{
	console.log("starting registerClose");
	var mainAccount = mc.coinbase;

	chain3.personal.unlockAccount(mainAccount,'test123',0);

	var myContract = chain3.mc.contract(SubChainP1Abi);
	var contractInstance = myContract.at(SubChainP1Address);

	var value = contractInstance.registerClose.call();

	console.log("registerOpen registerClose", value);

	return value;
}


function registerOpen(SubChainP1Address, SubChainP1Abi)
{
	console.log("starting registerOpen");
	var mainAccount = mc.coinbase;

	chain3.personal.unlockAccount(mainAccount,'test123',0);

	var myContract = chain3.mc.contract(SubChainP1Abi);
	var contractInstance = myContract.at(SubChainP1Address);

	var value = contractInstance.RegisterOpen.call(function(e,c) {
		console.log("registerOpen returned", e, c);
	});

	console.log("registerOpen called", value);

	return value;
}

function protocolregister(scsaddr,scspasswd,protocoladdr,callback)
{
	console.log("protocolregister");
	chain3.personal.unlockAccount(scsaddr,scspasswd,0);
	// sendtx(scsaddr, protocoladdr, '11','0x4420e4860000000000000000000000007b5f7eb099545f63b5a496723f1f94770b4ac959');

	// 1) loadfile
	var contractName = 'SubChainP1';
	loadScript(contractName + '.abi');
	loadScript(contractName + '.bin');
	var abi = JSON.stringify(contractAbi);
	var bytecode = contractBytecode;

	var abiJson = JSON.parse(abi);
	MyContract = chain3.mc.contract(abiJson);
	var contractInstance = MyContract.at(protocoladdr);

	contractInstance.register("0x78e754EAe2f46A816c2adB5c98fb47f2Da9B00B6",
		function (e, c) {
			cosnole.log("protocolregister Returned", e, c);
		}
	);


	// sendtx(scsaddr, protocoladdr, '11','0x4420e486000000000000000000000000' + scsaddr.substr(2, 100), callback);
}

function loadSubChainP1(SubChainProtocolBaseAddress)
{
	var mainAccount = mc.coinbase;

	// 1) loadfile
	var contractName = 'SubChainP1';
	loadScript(contractName + '.abi');
	loadScript(contractName + '.bin');
	var abi = JSON.stringify(contractAbi);
	var bytecode = contractBytecode;
	console.log("output", abi);
	console.log("output", bytecode);

	// 2) deploy subchain contract
	var mainAccount = mc.coinbase;
	chain3.personal.unlockAccount(mainAccount,'test123',0);

	var abiJson = JSON.parse(abi);
	MyContract = chain3.mc.contract(abiJson);

	var min = 100;
	var max = 2000;
	var thousandth = 200;
	var flushRound = 10;
	var contractInstance = MyContract.new(
		SubChainProtocolBaseAddress,
		min,
		max,
		thousandth,
		flushRound,
		{from: mainAccount, data: bytecode, gas: 8000000, queryFlag: 0, shardingFlag: 0},
		function(e, contract) {
			console.log("loadSubChainP1 responded");
			console.log(e, contract);
			if (typeof contract.address !== 'undefined') {
				console.log('SubChainP1 Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
				// protocolregister(mc.coinbase, 'test123', contract.address, 
				// 	function (err, c) {
				// 		console.log("protocolregister returned", err, c);

				// 		registerOpen(contract.address, abiJson);
				// 	});

				registerOpen(contract.address, abiJson);
				// registerSCS1(contract.address, abiJson);
				// registerSCS2(contract.address, abiJson);
				var event = contract.logging();

			    // monitor
			    event.watch(function(error, result){
				    console.log("Event are as following:-------");
				    
				     console.log(JSON.stringify(result));
				     console.log(result.args.msg);

				    console.log("Event ending-------");
				});


		        // contract.test.call();

			}
		}
	);
	return contractInstance;

}

function directCallTest() {
	var mainAccount = mc.coinbase;
	chain3.personal.unlockAccount(mainAccount,'test123',0);
	var amount = 100;

	while (mc.accounts.length < 3) {
		personal.newAccount("test123");
	}

	console.log("mc.accounts", mc.accounts);
	var account1 = mc.accounts[1];

	chain3.mc.sendTransaction(
		{
			from: mainAccount,
			value: chain3.toSha(amount,'mc'),
			to: account1,
			gas: "100000",
			gasPrice: chain3.mc.gasPrice,
			data: '',
			shardingFlag: 1
		},
		function(e,c) {
			console.log("e", e, "c", c);
		}
	);


}


function sendSomeMC() {
	var mainAccount = mc.coinbase;
	chain3.personal.unlockAccount(mainAccount,'test123',0);
	var amount = 100;

	while (mc.accounts.length < 3) {
		personal.newAccount("test123");
	}

	console.log("mc.accounts", mc.accounts);
	var account1 = mc.accounts[1];
	var account2 = mc.accounts[2];

	chain3.mc.sendTransaction(
	{
		from: mainAccount,
		value: chain3.toSha(amount,'mc'),
		to: account1,
		gas: "100000",
		gasPrice: chain3.mc.gasPrice,
		data: ''
	});

	chain3.mc.sendTransaction(
	{
		from: mainAccount,
		value:chain3.toSha(amount,'mc'),
		to: account2,
		gas: "100000",
		gasPrice: chain3.mc.gasPrice,
		data: ''
	});

}

function activateFnNofityScs() {
	var mainAccount = mc.coinbase;
	chain3.personal.unlockAccount(mainAccount,'test123',0);

	chain3.mc.sendTransaction(
	{
		from: mainAccount,
		value: '1',
		to: '0x000000000000000000000000000000000000000d',
		gas: "100000",
		gasPrice: chain3.mc.gasPrice,
		data: ''
	});
}

function simpleStorageContractTest() {
	// 1) loadfile
	var contractName = 'SimpleStorage';
	loadScript(contractName + '.abi');
	loadScript(contractName + '.bin');
	var abi = JSON.stringify(contractAbi);
	var bytecode = contractBytecode;
	console.log("output", abi);
	console.log("output", bytecode);

	// 2) deploy subchain contract
	var mainAccount = mc.coinbase;
	chain3.personal.unlockAccount(mainAccount,'test123',0);

	var abiJson = JSON.parse(abi);
	MyContract = chain3.mc.contract(abiJson);

	var contractInstance = MyContract.new(
		{from: mainAccount, data: bytecode, gas: 8000000, queryFlag: 0, shardingFlag: 0},
		function(e, contract) {
			console.log("simpleStorageContractTest responded");
			console.log(e, contract);
			if (typeof contract.address !== 'undefined') {
				console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
				//loadSubChainP1(contract.address);
				//activateFnNofityScs();
				var event = contract.logging();


			    // monitor
			    event.watch(function(error, result){
				    console.log("Event are as following:-------");
				    
				     console.log(JSON.stringify(result));
				     console.log(result.args.msg);

				    console.log("Event ending-------");
				});


		        // contract.test.call();

			}
		}
	);
	return contractInstance;

}

function fullSubchainContractTest() {
	// 1) loadfile
	var contractName = 'SubChainProtocolBase';
	loadScript(contractName + '.abi');
	loadScript(contractName + '.bin');
	var abi = JSON.stringify(contractAbi);
	var bytecode = contractBytecode;
	console.log("output", abi);
	console.log("output", bytecode);

	// 2) deploy subchain contract
	var mainAccount = mc.coinbase;
	chain3.personal.unlockAccount(mainAccount,'test123',0);

	var abiJson = JSON.parse(abi);
	MyContract = chain3.mc.contract(abiJson);

	var min = 100;
	var max = 2000;
	var thousandth = 200;
	var protocol = "test";
	var contractInstance = MyContract.new(
		protocol,
		min,
		// max,
		// thousandth,
		{from: mainAccount, data: bytecode, gas: 8000000, queryFlag: 0, shardingFlag: 0},
		function(e, contract) {
			console.log("fullSubchainContractTest SubChainProtocolBase responded");
			console.log(e, contract);
			if (typeof contract.address !== 'undefined') {
				console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
				loadSubChainP1(contract.address);
				activateFnNofityScs();
				var event = contract.logging();


			    // monitor
			    event.watch(function(error, result){
				    console.log("Event are as following:-------");
				    
				     console.log(JSON.stringify(result));
				     console.log(result.args.msg);

				    console.log("Event ending-------");
				});


		        // contract.test.call();

			}
		}
	);
	return contractInstance;
}

function initMC() {
	if (mc.accounts.length < 1) {
		console.log('Cannot find coinbase');
		personal.newAccount('test123')
		.then(function(){
			setTimeout(
			miner.start(1),
			10000
			);
		});
	} else {
		console.log('Found coinbase', mc.coinbase);
		miner.start(1);
	}
}