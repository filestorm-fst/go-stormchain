pragma solidity ^0.4.11;
//David Chen
//Sharding System Contract
//

import "./SubChainProtocolBase.sol";

contract SubChainProtocolPBFT is SubChainProtocolBase  {

   
    uint public BlockTime;
    
    //constructor
    function SubChainProtocolPBFT(string protocol, uint bmin, uint blktime) SubChainProtocolBase(protocol, bmin) public {
		BlockTime = blktime;
    }


    function() public {
        revert();
    }

}
