pragma solidity ^0.4.13;
//pragma experimental ABIEncoderV2;


//MicroChain definition for application.
contract DeChatManagement {
    enum BoardStatus {noState, working, hiding}
    
	struct BoardInfo {
        address subchainAddr;
        address deployLwSolAdmin;
        address marketableTokenAddr;
        bytes32 rpcIp;
        bytes32 boardName;
        bytes32 picPath;
        uint boardStatus;
        uint exchangeRate;
    }
    
    BoardInfo[] public boardList;
    address internal owner;
    
	function DeChatManagement() public payable {
		owner = msg.sender;
	}
	
	function creatBoard(address subchainAddr, address deployLwSolAdmin, address marketableTokenAddr, bytes32 rpcIp, bytes32 boardName, bytes32 picPath, uint exchangeRate) public {
	    require(owner == msg.sender);
        boardList.push(BoardInfo(subchainAddr, deployLwSolAdmin, marketableTokenAddr, rpcIp, boardName, picPath, uint(BoardStatus.working), exchangeRate));
    }
    
    function getBoardlist(uint status) public constant returns (address[], bytes32[], uint[]) {
        uint i;
        uint j = 0;
        for (i = 0; i < boardList.length; i++) {
            if(boardList[i].boardStatus == status) {
                j++;
            }
        }
        
        address[] memory addrlist = new address[](j*3);
        bytes32[] memory byteslist = new bytes32[](j*3);
        uint[] memory uintlist = new uint[](j*2);
        j = 0;
        for (i = 0; i < boardList.length; i++) {
            if(boardList[i].boardStatus == status) {
                addrlist[j*3] = boardList[i].subchainAddr;
                addrlist[j*3+1] = boardList[i].deployLwSolAdmin;
                addrlist[j*3+2] = boardList[i].marketableTokenAddr;
                byteslist[j*3] = boardList[i].rpcIp;
                byteslist[j*3+1] = boardList[i].boardName;
                byteslist[j*3+2] = boardList[i].picPath;
                uintlist[j*2] = boardList[i].boardStatus;
                uintlist[j*2+1] = boardList[i].exchangeRate;
                j++;
            }
        }
        return (addrlist, byteslist, uintlist);
    }
    
    function updateBoardStatus(uint status, address subchainAddr) public {
        require(owner == msg.sender);
        
        uint i;
        for (i = 0; i < boardList.length; i++) {
            if(boardList[i].subchainAddr == subchainAddr) {
                boardList[i].boardStatus = status;
                break;
            }
        }
    }

    function updateBoardName(bytes32 boardName, address subchainAddr) public {
        require(owner == msg.sender);
        
        uint i;
        for (i = 0; i < boardList.length; i++) {
            if(boardList[i].subchainAddr == subchainAddr) {
                boardList[i].boardName = boardName;
                break;
            }
        }
    }
}