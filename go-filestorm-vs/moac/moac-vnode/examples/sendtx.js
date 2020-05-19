/*
 * Example programs to test under MOAC console
 *  
 */
var src=mc.accounts[0];
var des=mc.accounts[1];

function initB() {
    var t =["","","",""];
t[0]="0xb037ff99dc90ebe31e456998751a942fe21ad0ac";
t[1]="0xb4692258ceb200ea02dfdcfc37161ea4fa2b9de3";
t[2]="0x0c7b8639e63ff3d38ea46f4cefe7d7957f83f157";
t[3]="0x0e6d92aa7a3fe4aeacc27b14c493577bbc079223";
var totalBal = 0;
for (var acctNum in t) {
    var acct = t[acctNum];
    var acctBal = chain3.fromSha(mc.getBalance(acct), "mc");
    totalBal += parseFloat(acctBal);
    console.log("Init account: \t" + acct + " \tbalance: " + acctBal + " mc");
}
console.log("  Init balance: " + totalBal + " mc ");
/*
"0xb037ff99dc90ebe31e456998751a942fe21ad0ac": { "balance": "45000000000000000000000000"},
"0xb4692258ceb200ea02dfdcfc37161ea4fa2b9de3": { "balance": "15000000000000000000000000"},
"0x0c7b8639e63ff3d38ea46f4cefe7d7957f83f157": { "balance": "15000000000000000000000000"},
"0x0e6d92aa7a3fe4aeacc27b14c493577bbc079223": { "balance": "75000000000000000000000000"}
*/
}

/*
 * Display some numbers
*/
function checkB() {
    var totalBal = 0;
    // for (var acctNum in mc.accounts) {
    for (var acctNum =0; acctNum < 10; acctNum ++) {
        var acct = mc.accounts[acctNum];
        var acctBal = chain3.fromSha(mc.getBalance(acct), "mc");
        totalBal += parseFloat(acctBal);
        console.log("  mc.accounts[" + acctNum + "]: \t" + acct + " \tbalance: " + acctBal + " mc");
    }



    /*var blocknumber = mc.blockNumber;
    console.log("  Total balance: " + totalBal + " mc in " +blocknumber);
    console.log("Ave mc per block", totalBal/blocknumber)
*/
};
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

    //var amt = leftPad(chain3.toHex(chain3.toSha(amount)).slice(2).toString(16),64,0);
    //var strData = '';
        
    chain3.mc.sendTransaction(
        {
            from: src,
            value:chain3.toSha(amount,'mc'),
            to: tgtaddr,
            gas: "2000000",
            gasPrice: chain3.mc.gasPrice,
            data: strData//,
            // shardingFlag:0,
            // via : '0x0000000000000000000000000000000000000000'
        }, function (e, transactionHash){
            if (!e) {
                 console.log('Transaction hash: ' + transactionHash);
            }else{
                console.log('Error:'+e);
            }
         });
        
    console.log('sending from:' +   src + ' to:' + tgtaddr  + ' amount:' + amount + ' with data:' + strData);

}

function Send(src, passwd, target, value, indata)
{
    chain3.personal.unlockAccount(src, passwd, 0);
    sendtx(src, target, value, indata );
    
}

function FutureSend(src, passwd, target, value, block)
{
    //address must start with 0x
    // if( !(target.substring(0,2) == "0x" || target.substring(0,2) == "0X" ))
    // {   
    //     console.log("error target address format, expect 0x");
    //     return;
    // }

    var str = "000000000000000000000000AAAA";   
    var strtgt = str.replace("AAAA", target.substring(2));
        
    var amt = leftPad(chain3.toHex(chain3.toSha(value, 'mc')).slice(2).toString(16),64,0);

    var blkstr = leftPad(chain3.toHex(block).slice(2).toString(16),64,0);

    var strData = "0xdef0412f";
    strData = strData + blkstr + strtgt + amt;

    chain3.personal.unlockAccount(src, passwd, 0);
    var src = src;
    var cntaddr = "0x0000000000000000000000000000000000000065";
    sendtx(src, cntaddr, '0', strData );
    
}


scs1="0x1b65cE1A393FFd5960D2ce11E7fd6fDB9e991945";
scs2="0xecd1e094ee13d0b47b72f5c940c17bd0c7630326";
scs3="0x50C15fafb95968132d1a6ee3617E99cCa1FCF059";

var registerdata = "0x4420e486000000000000000000000000ECd1e094Ee13d0B47b72F5c940C17bD0c7630326";
function registeropen(subAddress) {
    sendtx(chain3.mc.coinbase, subAddress, 0, "0x5defc56c");
}

function addfundtosubchain(src, subAddress, value) {
    sendtx(src, subAddress, value, "0xa2f09dfa");
}

function registerclose() {
    sendtx(chain3.mc.coinbase, subAddress, 0, "0x69f3576f");
}

function registertopool(contractadd, scsaddress) {
    var registerdata = "0x4420e486000000000000000000000000"+scsaddress.substring(2);
    sendtx(chain3.mc.coinbase, contractadd, 12, registerdata);
}

