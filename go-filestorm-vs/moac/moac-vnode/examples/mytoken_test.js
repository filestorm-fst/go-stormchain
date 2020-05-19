var tokenabi='[{"constant":false,"inputs":[{"name":"newSellPrice","type":"uint256"},{"name":"newBuyPrice","type":"uint256"}],"name":"setPrices","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"name","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"}],"name":"approve","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"totalSupply","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transferFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"decimals","outputs":[{"name":"","type":"uint8"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_value","type":"uint256"}],"name":"burn","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"sellPrice","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"balanceOf","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"target","type":"address"},{"name":"mintedAmount","type":"uint256"}],"name":"mintToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_value","type":"uint256"}],"name":"burnFrom","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"buyPrice","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"symbol","outputs":[{"name":"","type":"string"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[],"name":"buy","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"transfer","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"frozenAccount","outputs":[{"name":"","type":"bool"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"_spender","type":"address"},{"name":"_value","type":"uint256"},{"name":"_extraData","type":"bytes"}],"name":"approveAndCall","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"},{"name":"","type":"address"}],"name":"allowance","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"amount","type":"uint256"}],"name":"sell","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"target","type":"address"},{"name":"freeze","type":"bool"}],"name":"freezeAccount","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"inputs":[{"name":"initialSupply","type":"uint256"},{"name":"tokenName","type":"string"},{"name":"tokenSymbol","type":"string"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"target","type":"address"},{"indexed":false,"name":"frozen","type":"bool"}],"name":"FrozenFunds","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":true,"name":"to","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Transfer","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"from","type":"address"},{"indexed":false,"name":"value","type":"uint256"}],"name":"Burn","type":"event"}]';
//testnet v0.3
//> null [object Object]
//Contract mined! address: 0x0a674edac2ccd47ae2a3197ea3163aa81087fbd1 
//transactionHash: 0x0e6424970fff9c2b17fd2984496d72ab1503a9275a5bdd052b84aa7d3db2d2a0
//PanGu 0.8
//Contract mined! address: 0x2258355e41be8435408c57fbb4976cab8a0b16ca 
//transactionHash: 0x5cc35efdfa4bb1c86919eb4c1179b3e1b2a54ee7541ecaa778492b97f501b5de
// var tokenaddress='0x2258355e41be8435408c57fbb4976cab8a0b16ca';
// var tokenaddress='0x9b3B50c86B9756C7e269Bce572664fFD788Be113';//testnetwork id = 101
// var tokenaddress='0x1d12aec505caa2b8513cdad9a13e9d4806a1b704';//mc02
//mc04
//Contract mined! address: 0x967c7f17478cd8e09ce8b5a0d13ec4878695f24a 
//transactionHash: 0x914d18438c47b88330e04234068fe4b574bbc4a3f06e7a44c6f9273712ad4231
var tokenaddress='0x1d12aec505caa2b8513cdad9a13e9d4806a1b704';//mc02
//2018/05/15, deployed new erc20 for testing.
//Contract mined! address: 0xf2f4eec6c2adfcf780aae828de0b25f86506ffae transactionHash: 0x7bb694c3462764cb113e9b742faaf06adc728e70b607f8b7aa95207ee32b1c5e
tokenaddress='0xf2f4eec6c2adfcf780aae828de0b25f86506ffae';
//var tokenaddress='0x28c937495e79e4731e65ad635ea0375750626f80';//devnet, networkid=100
var tokenContract=mc.contract(JSON.parse(tokenabi));
var tcalls=tokenContract.at(tokenaddress);
var add0=mc.accounts[0]; //'0xa8863fc8ce3816411378685223c03daae9770ebb'
var add1=mc.accounts[1]; //'0xd814f2ac2c4ca49b33066582e4e97ebae02f2ab9'
function sendtoken(src, passwd, tgtaddr, amount, strData) {

    // var amount = leftPad(chain3.toHex(chain3.toSha(amount)).slice(2).toString(16),64,0);
    // var strData = '';
    chain3.personal.unlockAccount(src, passwd, 0);
    tcalls.transfer(tgtaddr, amount, {from:src}, function(err, hash){
      console.log(err, hash);
    });
    console.log('sending from:' +   src + ' to:' + tgtaddr  + ' amount:' + amount + ' '+tcalls.name());

}
//Display the tokenInfo
function tokenInfo(){
  console.log("Token Info\nfull name:", tcalls.name());
  console.log("   symbol:", tcalls.symbol());
  console.log("   supply:",tcalls.totalSupply());
  console.log("   owners:", tcalls.owner());
}
function checkToken(inacct){
  console.log(inacct, " has ",tcalls.balanceOf(inacct), tcalls.symbol());
}
function tokenB(){
var totalBal = 0;
for (var acctNum in mc.accounts) {
    var acct = mc.accounts[acctNum];
    var acctBal = tcalls.balanceOf(acct);
    totalBal += parseFloat(acctBal);
    console.log("  mc.accounts[" + acctNum + "]: \t" + acct + " \tbalance: " + acctBal);
}
console.log("  Total balance: " + totalBal);
}
