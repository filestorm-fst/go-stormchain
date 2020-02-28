/* TPS Testing Script 
 * Name: tps.js
 * Open storm console, and do loadScript("tps.js")
 * unlock fst.coinbase account, and make sure it has enough token
 * run testTransactions(transCount)
 * We can run this script on multiple storm nodes.
 */

// blockSec - how many seconds are in a block.
// transCount - test how many transactions
function testTransactions(blockSec, transCount) {

    fromAddr = fst.coinbase;
    toAddr = "0x0000000000000000000000000000000000000000";
    strData = "0x";

    startBlock = fst.blockNumber;
	for (i=0; i<transCount; i++) {
	    fst.sendTransaction(
		{
			from:   fromAddr,
            		to:     toAddr,
			value:  transCount,
			gas:    21000,
			gasPrice: 1,
			data:   strData
		}, function (e, transactionHash){
            if (!e) {
                 console.log('Transaction hash: ' + transactionHash);
            }
         });
	    console.log('sending from:' + fromAddr + ' to:' + toAddr + ' amount:' + transCount + ' with data:' + strData);
	}
    finishBlock = fst.blockNumber
    console.log('start block:' + startBlock);
    console.log('finish block:' + finishBlock);

    total = 0; b = 0;
    for (j=startBlock-1; j<=finishBlock+10-blockSec; j++) {
        if (fst.getBlockTransactionCount(j)!=null) {
            c = fst.getBlockTransactionCount(j)
            console.log('Block: ' + j + '   Transaction Count:' + c);
            total += c;
            if (c > 0)
                b+=blockSec;
        }
        else
            break;
    }
    console.log('TPS: ' + total/b); 
}
