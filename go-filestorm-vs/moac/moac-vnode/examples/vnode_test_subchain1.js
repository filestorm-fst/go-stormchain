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

function sendtx(src, tgtaddr, amount, strData) {

	mc.sendTransaction(
		{
			from: src,
			value:chain3.toSha(amount,'mc'),
			to: tgtaddr,
			gas: "2000000",
			gasPrice: chain3.mc.gasPrice,
			data: strData
		});
		
	console.log('sending from:' + 	src + ' to:' + tgtaddr  + ' amount:' + amount + ' with data:' + strData);

}

/*
 * Send a direct through the scs nodes
 *
*/
function senddirecttx(src, tgtaddr, amount, strData) {

	mc.sendTransaction(
		{
			from: src,
			value:chain3.toSha(amount,'mc'),
			to: tgtaddr,
			gas: "2000000",
			gasPrice: chain3.mc.gasPrice,
			shardingflag: 1
			data: strData
		});
		
	console.log('sending from:' + 	src + ' to:' + tgtaddr  + ' amount:' + amount + ' with data:' + strData);

}

function sendsystx(src, tgtaddr, amount, strData) {

	mc.sendTransaction(
		{
			from: src,
			value:chain3.toSha(amount,'mc'),
			to: tgtaddr,
			gas: "2000000",
			gasPrice: chain3.mc.gasPrice,
			SystemContract: 1, 
			data: strData
		});
		
	console.log('sending from:' + 	src + ' to:' + tgtaddr  + ' amount:' + amount + ' with data:' + strData);

}

function testnull(address) {

	var testnull2 = subchainbaseContract.new(
	   proto,
	   min,
	   max,
	   thousandth,
	   flushRound,
	   {
		 from: mc.accounts[0], 
		 data: '0x', 
		
		 gas: '5800000'
	   }, function (e, contract){
		console.log(e, contract);
		if (typeof contract.address !== 'undefined') {
			 console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
		}
	 })
}

function testnull2(){
		personal.unlockAccount(mc.accounts[0], '', 0);
	sendtx(mc.accounts[0], '0xb2d1fcd62f89c2798e37481d7f883611100eac4a', 0, '0x08bf29a8');

}

function testdifficulty(){
	for( i=48000; i<mc.blockNumber; i+=1 )
	{
		console.log(mc.getBlock(i).difficulty);
	}
		
}

function testblocktime(){
	for( i=(mc.blockNumber-50); i<mc.blockNumber; i+=1 )
	{
		console.log(mc.getBlock(i).timestamp);
	}
		
}


function testnull3(){
		personal.unlockAccount(mc.accounts[1], '', 0);
	

var distributionContract = chain3.mc.contract([{"constant":true,"inputs":[],"name":"MaxBalanceCanGetDrop","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"DropAmount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"GetMoac","outputs":[{"name":"","type":"bool"}],"payable":true,"stateMutability":"payable","type":"function"},{"inputs":[],"payable":true,"stateMutability":"payable","type":"constructor"},{"payable":true,"stateMutability":"payable","type":"fallback"}]);

var distribution = distributionContract.new(
   {
     from: chain3.mc.accounts[1], 
     data: '0x', 
	 value:chain3.toSha(10,'mc'),
     gas: '4700000'
   }, function (e, contract){
    console.log(e, contract);
    if (typeof contract.address !== 'undefined') {
         console.log('Contract mined! address: ' + contract.address + ' transactionHash: ' + contract.transactionHash);
    }
 })	
}

function transferbal(address) {
		personal.unlockAccount(mc.accounts[0], '', 0);
		sendtx(mc.accounts[0], address, 2);
	
}

function registertopool(address) {
		//personal.unlockAccount(address, '', 0);
		var registerdata = "0x4420e486000000000000000000000000"+address.substring(2);
		sendtx(mc.accounts[0], subchainprotocolbase.address, 12, registerdata);
}


function registeropen() {
		//personal.unlockAccount(mc.accounts[0], '', 0);
		//var benefit = "0x4420e486000000000000000000000000"+mc.accounts[i].substring(2);
		sendtx(mc.accounts[0], subchainbase.address, 0, "0x5defc56c");
}

function addfundtosubchain() {
		//personal.unlockAccount(mc.accounts[0], '', 0);
		//var benefit = "0x4420e486000000000000000000000000"+mc.accounts[i].substring(2);
		sendtx(mc.accounts[0], subchainbase.address, 1, "0xa2f09dfa");
}


function registerownerSCS(addr) {
		//personal.unlockAccount(mc.accounts[0], '', 0);
		var params = "0x7f732395000000000000000000000000"+addr.substring(2)+"000000000000000000000000"+addr.substring(2);
		sendtx(mc.accounts[0], subchainbase.address, 0, params);
}

function scsregister(address) {
		//personal.unlockAccount(address, '', 0);
		var registerdata = "0x99d53d43000000000000000000000000"+address.substring(2)+"0000000000000000000000000000000000000000000000000000000000000000"+"0000000000000000000000000000000000000000000000000000000000000000"+"0000000000000000000000000000000000000000000000000000000000000000";
		sendtx(address, subchainbase.address, 0, registerdata);
}

function scsregister12() {
	for( var i=0; i<12; i++ )
	{
		scsregister(mc.accounts[i+1]);
	}
}

function registerclose() {
		//personal.unlockAccount(mc.accounts[0], '', 0);
		//var benefit = "0x4420e486000000000000000000000000"+mc.accounts[i].substring(2);
		sendtx(mc.accounts[0], subchainbase.address, 0, "0x69f3576f");
}

function propCreate() {
		personal.unlockAccount(mc.accounts[2], '', 0);
		var propdata = "0x2253d465"+
		"0000000000000000000000000000000000000000000000000000000000000000"+"0100000000000000000000000000000000000000000000000000000000000000"+"0000000000000000000000000000000000000000000000000000000000000060"+"000000000000000000000000000000000000000000000000000000000000000c"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000"+"00000000000000000000000000000000000000000000000000038D7EA4C68000";
		sendtx(mc.accounts[2], subchainbase.address, 0, propdata);
}



