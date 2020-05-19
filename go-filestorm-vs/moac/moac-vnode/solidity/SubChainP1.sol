pragma solidity ^0.4.11;
//David Chen
//Sharding System Contract
//
import './SubChainBase.sol';

contract SubChainP1 is SubChainBase {

	
    //bytes funcode;
    
    
    //constructor
    function SubChainP1(address protoaddress, uint min, uint max, uint thousandth, uint flushRound) 
		SubChainBase(protoaddress, min, max, thousandth, flushRound) public {

    }

    function DeployCode(bytes code) public {
      require(msg.sender == owner);
		  funcode = code;
    }

}


//"e102caa8": "ConsensusFlag()",
//"c69c650a": "IsMemberValid(address)",
//"4af3ae97": "MaxMember()",
//"84eaef89": "MinMember()",
//"6a19f6ad": "NodeCount()",
//"98fe20af": "Protocol()",
//"7c390fc8": "SelTarget()",
//"9a34d888": "UpdateConsensus(uint256,uint256,uint256,address[])"
