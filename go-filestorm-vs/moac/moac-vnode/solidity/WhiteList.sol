pragma solidity ^0.4.11;
//David Chen, Xiong Liu

contract WhiteListOld {
    function getWhiteList(uint pageNum, uint pageSize) public view returns (bytes32[]);
    function getWhiteFunc(uint pageNum, uint pageSize) public view returns (bytes32[]);
}

contract WhiteList {

    address public owner;
    mapping(address => uint) public admins;
    mapping(bytes32 => uint) public whiteList;
    bytes32[] public whiteListArray;
    mapping(bytes32 => uint) public whiteFunc;
    bytes32[] public whiteFuncArray;


    //constructor
    function WhiteList(address whiteListOldAddr) public {
        owner = msg.sender;
        
        WhiteListOld whiteListOld = WhiteListOld(whiteListOldAddr);
        whiteListArray = whiteListOld.getWhiteList(0,1000);
        uint i=0;
        for (i=0; i<whiteListArray.length; i++) {
            whiteList[whiteListArray[i]] = i+1;
        }

        whiteFuncArray = whiteListOld.getWhiteFunc(0,1000);
        for (i=0; i<whiteFuncArray.length; i++) {
            whiteFunc[whiteFuncArray[i]] = i+1;
        }
    }

    function() public payable { 
        revert();
    }

    function addAdmin(address admin) public {
        require(msg.sender == owner);
        admins[admin] = 1;
    }

    function removeAdmin(address admin) public {
        require(msg.sender == owner);
        admins[admin] = 0;
    }

    // whiteList
    function register(bytes32 hash) public {
        require(msg.sender == owner || admins[msg.sender] == 1);
        require(whiteList[hash]==0);
        
        whiteListArray.push(hash);
        whiteList[hash] = whiteListArray.length;
    }

    function remove(bytes32 hash) public {
        require(msg.sender == owner || admins[msg.sender] == 1);
        uint count = whiteListArray.length;
        uint index = whiteList[hash];
        if (index > 0) {
            if (index == count) {
                delete whiteListArray[index-1];
                whiteListArray.length--;
            } else if (index < count) {
                whiteListArray[index-1] = whiteListArray[count-1];
                whiteList[whiteListArray[count-1]] = index;
                delete whiteListArray[count-1];
                whiteListArray.length--;
            }
            delete whiteList[hash];
        }        
    }

    function getWhiteList(uint pageNum, uint pageSize) public view returns (bytes32[]) {
        require(pageSize != 0);
		uint start = pageNum*pageSize;
		uint end = (pageNum+1)*pageSize;
		uint count = whiteListArray.length;
        bytes32[] memory memList;
        uint i = 0;
		if (start<count) {
			if (end>count) {
			    end = count;
		    }

            memList = new bytes32[](end-start);		
            for (i=start; i<end; i++) {
                memList[i-start] = whiteListArray[i];
            }
        }

		return memList;
    }

    function isValid(bytes32 hashFirst, bytes32 hashSecond) public view returns (bool) {
        if (whiteList[hashFirst] > 0) {
            return true;
        }
       
        if (whiteList[hashSecond] > 0) {
            return true;
        }

        return false;
    }

    // whiteFunc
    function registerFunc(bytes32 hash) public {
        require(msg.sender == owner || admins[msg.sender] == 1);
        require(whiteFunc[hash]==0);
        
        whiteFuncArray.push(hash);
        whiteFunc[hash] = whiteFuncArray.length;
    }

    function removeFunc(bytes32 hash) public {
        require(msg.sender == owner || admins[msg.sender] == 1);
        uint count = whiteFuncArray.length;
        uint index = whiteFunc[hash];
        if (index > 0) {
            if (index == count) {
                delete whiteFuncArray[index-1];
                whiteFuncArray.length--;
            } else if (index < count) {
                whiteFuncArray[index-1] = whiteFuncArray[count-1];
                whiteFunc[whiteFuncArray[count-1]] = index;
                delete whiteFuncArray[count-1];
                whiteFuncArray.length--;
            }
            delete whiteFunc[hash];
        }        
    }

    function getWhiteFunc(uint pageNum, uint pageSize) public view returns (bytes32[]) {
        require(pageSize != 0);
		uint start = pageNum*pageSize;
		uint end = (pageNum+1)*pageSize;
		uint count = whiteFuncArray.length;
        bytes32[] memory memFun; 
		if (start<count) {
            if (end>count) {
                end = count;
            }

            memFun = new bytes32[](end-start);		
            for (uint i=start; i<end; i++) {
                memFun[i-start] = whiteFuncArray[i];
            }
        }

		return memFun;
    }

    function isValidFunc(bytes32 hash) public view returns (bool) {
        if (whiteFunc[hash] > 0) {
            return true;
        }

        return false;
    }

    function getWhiteInfo(uint pageNum, uint pageSize) public view returns (bytes32[],bytes32[]) {
        require(pageSize != 0);
		uint start = pageNum*pageSize;
		uint end = (pageNum+1)*pageSize;
		uint count = whiteListArray.length;
        bytes32[] memory memList;
        uint i = 0;
		if (start<count && count>0) {
			if (end>count) {
			    end = count;
		    }

            memList = new bytes32[](end-start);		
            for (i=start; i<end; i++) {
                memList[i-start] = whiteListArray[i];
            }
        }

		start = pageNum*pageSize;
		end = (pageNum+1)*pageSize;
		count = whiteFuncArray.length;
        bytes32[] memory memFun; 
		if (start<count && count>0) {
            if (end>count) {
                end = count;
            }

            memFun = new bytes32[](end-start);		
            for (i=start; i<end; i++) {
                memFun[i-start] = whiteFuncArray[i];
            }
        }

		return (memList, memFun);
    }
}
