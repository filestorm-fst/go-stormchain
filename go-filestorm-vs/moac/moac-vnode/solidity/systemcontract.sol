pragma solidity ^0.4.11;
//David Chen
//Moac System Contract
//

contract GasSpendCnt {
  function spendGas (uint256) public;
}

contract Precompiled12 {
  function delegateSend (address, address, uint256) public returns (uint);
}

contract MoacSysContract {

    struct DelayedSendTask {
        uint block; // block# task executed on
		address from; //from address
        address to; // target address
        uint256 value;   // value
		bool revertable; // if revertable
    }
    
    struct DelayedSendTaskBonded {
        uint block; // block# task executed on
		address from; //from address
        address to; // target address
        uint256 value;   // value
		bool revertable; // if revertable
		bool locked;
		bytes32 hashUnLock; // 
		address unlockaddr;
    }
  
    
	uint constant QUEUESIZE  = 50;
	Precompiled12 constant PREC12 = Precompiled12(0x0c); 
	GasSpendCnt constant GASSPENT = GasSpendCnt(0x0e); 
	//system call address
	address constant SYSCALLADDR = address(100);
	address constant SYSCNTADDR = address(101);

	
    mapping(uint => DelayedSendTask[]) queueToSend;
    mapping(uint => DelayedSendTaskBonded[]) queueToSendBonded;
  //  mapping(uint => address[]) queryQueue;
    //mapping(uint => uint) sendQueueSize;

    //version
    uint curVersion;
   

    //events
    event DelayedSend(uint _blk, address _to,  uint256 _value, bool revertable);
    event DelayedSendBonded( uint _blk, address _to, uint256 _value, bool revertable, bool _locked, bytes32 _hashlock, address _unlockaddr );
	event DelayArrayFull(uint _blk);
	event DelayArrayRegistered(uint _blk, address _from);
    
 
    function queryDelaySend( uint _blk, address _from ) public view returns ( address, uint, bool) {
		
		for (uint i = 0; i < queueToSend[_blk].length; i++ ) {
			if ( queueToSend[_blk][i].from == _from ) {
				return (queueToSend[_blk][i].to, queueToSend[_blk][i].value, queueToSend[_blk][i].revertable);
			}
		}

		return (0, 0, false);
	}

	function queryDelaySendBonded( uint _blk, address _from ) public view returns (address, uint, bool, bool, bytes32, address) {
		
		for (uint i = 0; i < queueToSendBonded[_blk].length; i++ ) {
			if ( queueToSendBonded[_blk][i].from == _from ) {
				return (queueToSendBonded[_blk][i].to, queueToSendBonded[_blk][i].value, queueToSendBonded[_blk][i].revertable, queueToSendBonded[_blk][i].locked, queueToSendBonded[_blk][i].hashUnLock, queueToSendBonded[_blk][i].unlockaddr);
			}
		}

		return (0, 0, false, false, 0, 0);
	}

    // used for queue user request and send it later
    function  delayedSend( uint _blk, address _to,  uint256 _value, bool revertable) private returns (bool success) {
		// if pass block.number
		if ( _blk <= block.number )
			return false;

        //if already scheduled at that block
		for (uint i = 0; i < queueToSend[_blk].length; i++ ) {
			if ( queueToSend[_blk][i].from == msg.sender ) {
				//only one sender per block
				DelayArrayRegistered(_blk, msg.sender);
				return false;
			}
		}
				
		//make it more expense to add to task if more tasks in that block
		//so that system is not overloaded. 
		//user to choose less crowded block for delaysend if want to use less gas
		//this should be set exponentially increase 
		GASSPENT.spendGas(queueToSend[_blk].length);

		//update
		DelayedSendTask memory task = DelayedSendTask(_blk, msg.sender,_to, _value, revertable);
		queueToSend[_blk].push(task);
        DelayedSend(_blk, _to,  _value, revertable);
		
		return true;
    }
    
    // used for queue user request and send it later
    function  cancelDelayedSend(uint _blk) public returns (bool success)   {
		// if pass block.number
		if ( _blk <= block.number )
			return false;

        //find scheduled at that block
		for (uint i = 0; i < queueToSend[_blk].length; i++ ) {
			if ( queueToSend[_blk][i].from == msg.sender && queueToSend[_blk][i].revertable) {
				//just clear
				queueToSend[_blk][i].from = address(0);
				queueToSend[_blk][i].to = address(0);
				queueToSend[_blk][i].value = 0;
				return true;
			}
		}
				
        //find scheduled at that block
		for ( i = 0; i < queueToSendBonded[_blk].length; i++ ) {
			if ( queueToSendBonded[_blk][i].from == msg.sender && queueToSendBonded[_blk][i].revertable) {
				//just clear
				queueToSendBonded[_blk][i].from = address(0);
				queueToSendBonded[_blk][i].to = address(0);
				queueToSendBonded[_blk][i].value = 0;
				return true;
			}
		}
					
		return false;
    }
    
    // used for queue user request and send it later
    function  delayedSendBonded( uint _blk, address _to, uint256 _value, bool revertable, bool locked, bytes32 hashlock, address unlockaddr ) public payable returns (bool success)   {
		// if pass block.number
		if ( _blk <= block.number )
			return false;

		if (_value != msg.value )
			return false;

        //if already scheduled at that block
		for (uint i = 0; i < queueToSendBonded[_blk].length; i++ ) {
			if ( queueToSendBonded[_blk][i].from == msg.sender && queueToSendBonded[_blk][i].hashUnLock == hashlock ) {
				//only one sender per block
				DelayArrayRegistered(_blk, msg.sender);
				return false;
			}
		}
				
		//make it more expense to add to task if more tasks in that block
		//so that system is not overloaded. 
		//user to choose less crowded block for delaysend if want to use less gas
		//this should be set exponentially increase 
		GASSPENT.spendGas(queueToSendBonded[_blk].length);

		//update
		DelayedSendTaskBonded memory task = DelayedSendTaskBonded(_blk, msg.sender,_to, _value, revertable, locked, hashlock, unlockaddr);
		queueToSendBonded[_blk].push(task);
        DelayedSendBonded(_blk, _to, _value, revertable, locked, hashlock, unlockaddr);
		
		return true;
    }

	// send task
    function sendTask( ) public {

		//this can only be called by system call addr
		if ( msg.sender != SYSCALLADDR )
			return;
			
		//check each task and send out
		for ( uint i = 0; i < queueToSend[block.number].length; i++ ) {
			if ( queueToSend[block.number][i].from != address(0) ) {
				PREC12.delegateSend(queueToSend[block.number][i].from, queueToSend[block.number][i].to, queueToSend[block.number][i].value);
			}
		}

		//check each task and send out
		for ( i = 0; i < queueToSendBonded[block.number].length; i++ ) {
			if ( queueToSendBonded[block.number][i].from != address(0) ) {
				if (queueToSendBonded[block.number][i].locked ) {
					//refund to owner
					PREC12.delegateSend(SYSCNTADDR, queueToSendBonded[block.number][i].from, queueToSendBonded[block.number][i].value);
				} else {
					//send totarget
					PREC12.delegateSend(SYSCNTADDR, queueToSendBonded[block.number][i].to, queueToSendBonded[block.number][i].value);
				}

				return;
			}
		}    
	}
        
 	// send task
    function unlockTask( uint blk, address sender, bytes32 hash ) public returns (bool success) {
		//check task
		for ( uint i = 0; i < queueToSendBonded[blk].length; i++ ) {
			if ( queueToSendBonded[blk][i].from == sender && queueToSendBonded[blk][i].hashUnLock == hash ) {
				if (!queueToSendBonded[blk][i].locked || queueToSendBonded[blk][i].unlockaddr == msg.sender) {
					queueToSendBonded[blk][i].locked = false;
					return true;
				}

				return false;
			}
		}

		return false;
    }
       
    // send task
    function exec( ) public {
		if ( msg.sender != SYSCALLADDR )
			return;
		
		sendTask();
	}
	
    function() public {
        revert();
    }

}




