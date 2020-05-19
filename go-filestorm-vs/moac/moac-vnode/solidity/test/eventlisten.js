	var abi = JSON.parse('[{"constant":false,"inputs":[{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"joinSCS","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_from","type":"address"},{"name":"_value","type":"uint256"}],"name":"enroll","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"sendTask","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"queryTask","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[],"name":"checkUpgradable","outputs":[{"name":"success","type":"bool"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"_blk","type":"uint256"},{"name":"_to","type":"address"},{"name":"_value","type":"uint256"}],"name":"delayedSend","outputs":[{"name":"success","type":"bool"}],"payable":true,"stateMutability":"payable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"payable":false,"stateMutability":"nonpayable","type":"fallback"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_blk","type":"uint256"},{"indexed":false,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"DelayedSend","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"hash","type":"bytes32"},{"indexed":false,"name":"expireblk","type":"uint256"}],"name":"UpgradeProposal","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_from","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"Enroll","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_to","type":"address"},{"indexed":false,"name":"_value","type":"uint256"}],"name":"JoinSCS","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_blk","type":"uint256"}],"name":"DelayArrayFull","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"_blk","type":"uint256"},{"indexed":false,"name":"from","type":"address"}],"name":"DelayArrayRegistered","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"info","type":"string"}],"name":"LogStep","type":"event"}]');

    var contract = web3.eth.contract(abi).at("0x0000000000000000000000000000000000000065");
    var myEvent = contract.LogStep();//contract.Evt({},{fromBlock: 0, toBlock: 'latest'});
    myEvent.watch(function(error, result){
        
		if( !error )
		{

			console.log(result.blockNumber + ":" + result.event+":");
			for( key in result.args){
				console.log("\t" + result.args[key]);
			}			
		}
		else{
			console.log("err:" + error);
		}
		

    });
