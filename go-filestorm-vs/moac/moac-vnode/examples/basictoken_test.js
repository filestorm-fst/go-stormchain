var tokenabi='[{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_value","type":"uint256"}],"name":"burn","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_value","type":"uint256"}],"name":"burnFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"},{"name":"_extraData","type":"bytes"}],"name":"approveAndCall","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"initialSupply","type":"uint256"},{"name":"tokenName","type":"string"},{"name":"tokenSymbol","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn","type":"event"}]';
//Contract mined! address: 0x416d3656637cb9673c574e0b11a23e9b9934b914 transactionHash: 0xc4a32e947f1a38050302cceb143d2abffaad20bfbe31b1cd23a1d6891d221154
//var tokenaddress="0x0a674edac2ccd47ae2a3197ea3163aa81087fbd1";//testnet 101
var tokenaddress="0x416d3656637cb9673c574e0b11a23e9b9934b914";//devnet 100
var tokenContract=mc.contract(JSON.parse(tokenabi));
var t1=tokenContract.at(tokenaddress);
var add0=mc.accounts[0]; //'0xa8863fc8ce3816411378685223c03daae9770ebb'
var add1=mc.accounts[1]; //'0xd814f2ac2c4ca49b33066582e4e97ebae02f2ab9'
function sendtoken(src, passwd, tgtaddr, amount, strData) {

    //var amt = leftPad(chain3.toHex(chain3.toSha(amount)).slice(2).toString(16),64,0);
    //var strData = '';
    chain3.personal.unlockAccount(src, passwd, 0);
    t1.transfer(tgtaddr, amount)
    console.log('sending from:' +   src + ' to:' + tgtaddr  + ' amount:' + amount + ' '+t1.name());

}
function tokenB(){
var totalBal = 0;
for (var acctNum in mc.accounts) {
    var acct = mc.accounts[acctNum];
    var acctBal = t1.balanceOf(acct);
    totalBal += parseFloat(acctBal);
    console.log("  mc.accounts[" + acctNum + "]: \t" + acct + " \tbalance: " + acctBal);
}
console.log("  Total balance: " + totalBal);
}
