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

	//var amt = leftPad(web3.toHex(web3.toWei(amount)).slice(2).toString(16),64,0);
	//var strData = '';
		
	web3.mc.sendTransaction(
		{
			from: src,
			value:web3.toSha(amount,'mc'),
			to: tgtaddr,
			gas: "2000000",
			gasPrice: web3.mc.gasPrice,
			data: strData
		});
		
	console.log('sending from:' + 	src + ' to:' + tgtaddr  + ' with data:' + strData);

}

function sendRaw() {
	var privateKey = new Buffer('e331b6d69882b4cb4ea581d88e0b604039a3de5967688d3dcffdd2270c0fd109', 'hex')

	var rawTx = {
	  nonce: '0x00',
	  gasPrice: '0x09184e72a000', 
	  gasLimit: '0x2710',
	  to: '0x0000000000000000000000000000000000000000', 
	  value: '0x00', 
	  data: '0x7f7465737432000000000000000000000000000000000000000000000000000000600057'
	}

	var tx = new Tx(rawTx);
	tx.sign(privateKey);

	var serializedTx = tx.serialize();

	//console.log(serializedTx.toString('hex'));
	//f889808609184e72a00082271094000000000000000000000000000000000000000080a47f74657374320000000000000000000000000000000000000000000000000000006000571ca08a8bbf888cfa37bbf0bb965423625641fc956967b81d12e23709cead01446075a01ce999b56a8a88504be365442ea61239198e23d1fce7d00fcfc5cd3b44b7215f

	chain3.mc.sendRawTransaction('0x' + serializedTx.toString('hex'), function(err, hash) {
	  if (!err)
	    console.log(hash); // "0x7f9fade1c0d57a7af66ab4ead79fade1c0d57a7af66ab4ead7c2c2eb7b11a91385"
	});

}

function send0()
{
	web3.personal.unlockAccount('0x7b39c87b38e73f2b80f3cec5e1efa50c5c4d7936','test123',0);
	var src = "0x7b39c87b38e73f2b80f3cec5e1efa50c5c4d7936";
	var tgtaddr = "0xd5a122da73e752b5958eec74b0904be0206395f9";
	sendtx(src, tgtaddr, '10' ,'');
}
	
function send101()
{
	web3.personal.unlockAccount('0x80e206cbbbeaa9d3be73da35e9edab68664deafe','',0);
	var src = "0x80e206cbbbeaa9d3be73da35e9edab68664deafe";
	var tgtaddr = "0x0000000000000000000000000000000000000065";
	sendtx(src, tgtaddr, '100','' );
	
}

function send100()
{
	web3.personal.unlockAccount('0x80e206cbbbeaa9d3be73da35e9edab68664deafe','',0);
	var src = "0x80e206cbbbeaa9d3be73da35e9edab68664deafe";
	var tgtaddr = "0x0000000000000000000000000000000000000064";
	sendtx(src, tgtaddr, '100','' );
	
}
function send1()
{
	web3.personal.unlockAccount('0x80e206cbbbeaa9d3be73da35e9edab68664deafe','',0);
	var src = "0x80e206cbbbeaa9d3be73da35e9edab68664deafe";
	var tgtaddr = "0x2cd2248437277c75be4c4b36b8033f75e10e414f";
	sendtx(src, tgtaddr, '18','' );
	
}

function delayedsend()
{
	web3.personal.unlockAccount('0x2cd2248437277c75be4c4b36b8033f75e10e414f','',0);
	var src = "0x2cd2248437277c75be4c4b36b8033f75e10e414f";
	var tgtaddr = "0x0000000000000000000000000000000000000065";
	sendtx(src, tgtaddr, '0','0xdef0412f00000000000000000000000000000000000000000000000000000000000000c8000000000000000000000000c88918182ce7cca8942da3b1dbffb03516538c930000000000000000000000000000000000000000000000000000000000000100' );
	
}

function delayedsend2()
{
	web3.personal.unlockAccount('0x2cd2248437277c75be4c4b36b8033f75e10e414f','',0);
	var src = "0x2cd2248437277c75be4c4b36b8033f75e10e414f";
	var tgtaddr = "0x0000000000000000000000000000000000000065";
	sendtx(src, tgtaddr, '0','0xdef0412f000000000000000000000000000000000000000000000000000000000000012c00000000000000000000000045dcae837cb3740a33790786c53a0f51f01c14be00000000000000000000000000000000000000000000000010A741A462780000' );
	
}

function query()
{
	web3.personal.unlockAccount('0x80e206cbbbeaa9d3be73da35e9edab68664deafe','',0);
	var src = "0x80e206cbbbeaa9d3be73da35e9edab68664deafe";
	var tgtaddr = "0x0000000000000000000000000000000000000065";
	sendtx(src, tgtaddr, '0','0xaea722e8' );
	
}
function sendtask()
{
	web3.personal.unlockAccount('0x80e206cbbbeaa9d3be73da35e9edab68664deafe','',0);
	var src = "0x80e206cbbbeaa9d3be73da35e9edab68664deafe";
	var tgtaddr = "0x0000000000000000000000000000000000000065";
	sendtx(src, tgtaddr, '0','0x650aed4f' );
	
}

//blk	
//00000000000000000000000000000000000000000000000000000000000002ee
//_to
//000000000000000000000000c88918182ce7cca8942da3b1dbffb03516538c93	
//value
//0000000000000000000000000000000000000000000010000000000000000033	
//bonded
//0000000000000000000000000000000000000000000000000000000000000000	
//together
//0000000000000000000000000000000000000000000000000000000000000065000000000000000000000000c88918182ce7cca8942da3b1dbffb03516538c9300000000000000000000000000000000000000000000100000000000000000330000000000000000000000000000000000000000000000000000000000000000	
//0xc8d0d29a00000000000000000000000000000000000000000000000000000000000003ee000000000000000000000000c88918182ce7cca8942da3b1dbffb03516538c9300000000000000000000000000000000000000000000100000000000000000330000000000000000000000000000000000000000000000000000000000000000
//0xdef0412f00000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000100001000000000000000000000000000000000000000000000000000000000000000a	
	
	
	
	
	
	