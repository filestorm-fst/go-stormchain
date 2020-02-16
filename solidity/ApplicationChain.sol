pragma solidity =0.4.24 <0.6.3;
/**
 * @title ApplicationChain.sol
 * Stormchain can be registered on other blockchains as a
 * deligated application chain to support an application
 * using the tokens from the other chain.
 * This is the first step of building an application chain.
 */

contract ApplicationChain {

    using SafeMath for uint256;

    address internal owner;
    mapping(address => uint) public admins;
    
    uint256 public balance = 0;
    uint256 public chainId;
    uint256 public period;
    uint256 public flushEpoch;

    mapping(uint256=>flushRound) public flushMapping;
    uint256[] public flushList;
    
    mapping(uint256=>address[]) public flushValidatorList;
    
    struct flushRound{
        uint256 flushId;
        address validator;
        uint256 blockNumber;
        string blockHash;
    }

    string extraData;
    uint256 tokenTotal;

    uint256 public FOUNDATION_MOAC_REQUIRED_AMOUNT = 20 * 10 ** 18;
    uint256 public FLUSH_AMOUNT = 1 * 10 ** 17;
    address public FOUNDATION_BLACK_HOLE_ADDRESS = 0x48328afc8dd45c1c252e7e883fc89bd17ddee7c0;

    // construct a chain with no token
    constructor(uint256 blockSec, uint256 flushNumber, address[] initial_validators, uint256 exchangeRate) public payable {
        
        require(
            msg.value >= FOUNDATION_MOAC_REQUIRED_AMOUNT,
            "Not Enough MOAC to Create Application Chain"
        );

        owner = msg.sender;
        chainId = block.number;
        period = blockSec;
        flushEpoch = flushNumber;
        balance = msg.value;
        
        uint256 flushId = 0;
        flushMapping[flushId].flushId = flushId;
        flushMapping[flushId].validator = msg.sender;
        for (uint i=0; i<initial_validators.length; i++){
            flushValidatorList[flushId].push(initial_validators[i]);
        }
        flushMapping[flushId].blockNumber = 1;
        flushMapping[flushId].blockHash = "";      
        flushList.push(flushId);

        tokenTotal = msg.value * exchangeRate; 

        admins[msg.sender] = 1;
        
        // OPTIONAL: Give some gas fee to initial validators.
        // uint256 GAS_FEE = 2 * 10 ** 14;
        // uint256 gasNeededEach = balance / FLUSH_AMOUNT * GAS_FEE / initial_validators.length;
        // uint256 gasNeededTotal = gasNeededEach * initial_validators.length;
        // for (i=0; i<initial_validators.length; i++){
        //     initial_validators[i].transfer(gasNeededEach);
        // }
        // End of Option
    }

    function addAdmin(address admin) public {
        require(admins[msg.sender] == 1, "Only Admins Can Add Another Admin.");
        admins[admin] = 1;
    }

    function removeAdmin(address admin) public {
        require(admins[msg.sender] == 1, "Only Admins Can Add Another Admin.");
        require(admin != msg.sender, "Admins Cannot Remove Self.");
        admins[admin] = 0;
    }

    function addFund() public payable {
        balance += msg.value;
    }
    
    function withdrawFund(address recv, uint amount) public {
        require(admins[msg.sender] == 1);
        //require(admins[recv] == 1);
        require(amount <= balance);
        
        balance -= amount;
        recv.transfer(amount);
    }

    function flush(address[] current_validators, uint256 blockNumber, string blockHash) public {
        uint256 flushId = flushList.length-1;
        for (uint i=0; i<flushValidatorList[flushId].length; i++){
            if ((flushValidatorList[flushId][i]==msg.sender ||
                admins[msg.sender] == 1 ) &&
                flushMapping[flushId].blockNumber + flushEpoch == blockNumber){

                flushId++;
                flushMapping[flushId].flushId = flushId;
                flushMapping[flushId].validator = msg.sender;
                for (uint j=0; j<current_validators.length; j++){
                        flushValidatorList[flushId].push(current_validators[j]);
                    }
                flushMapping[flushId].blockNumber = blockNumber;
                flushMapping[flushId].blockHash = blockHash;  
                flushList.push(flushId);
                
                // sent MOAC to FOUNDATION_BLACK_HOLE_ADDRESS;
                balance -= FLUSH_AMOUNT;
                FOUNDATION_BLACK_HOLE_ADDRESS.transfer(FLUSH_AMOUNT);
                return;
            }
        }
    }
    
    function getGenesisInfo() public view returns (string) {
        
        string memory validString = "";
        string memory allocString = "";
        for (uint i=0; i<flushValidatorList[1].length; i++){
            validString = string(abi.encodePacked(validString, addr2str(flushValidatorList[1][i])));
        }

        uint256 averageAmt = tokenTotal / flushValidatorList[1].length;
        uint256 remainingAmt = tokenTotal - (averageAmt * (flushValidatorList[1].length - 1));
        string memory comma = ",";
        for (i=0; i<flushValidatorList[1].length; i++){
            if (i==flushValidatorList[1].length-1){
                averageAmt = remainingAmt;
                comma = "";
            }
            allocString = string(abi.encodePacked(allocString, 
            '    "',
            addr2str(flushValidatorList[1][i]),
            '": {',
            '      "balance": "',
            uint2str(averageAmt),
            '"',
            '    }',
            comma
            ));
        }

        return string(abi.encodePacked('{',
        '{',
        ' "config": {',
        ' "chainId": ',
        uint2str(chainId),
        ',',
        '  "pbft": {',
        '   "period": ',
        uint2str(period),
        ',',
        '   "epoch": "36000",',
        '   "epochFlush": "',
        uint2str(flushEpoch),
        '"',
        '  }',
        ' },',
        '  "nonce": "0x0",',
        '  "timestamp": "0x5de22b51",',
        '  "extraData": "0x0000000000000000000000000000000000000000000000000000000000000000',
        validString,
        '0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",',
        '  "gasLimit": "0x47b760",',
        '  "difficulty": "0x1",',
        '  "mixHash": "0x0000000000000000000000000000000000000000000000000000000000000000",',
        '  "coinbase": "0x0000000000000000000000000000000000000000",',
        '  "alloc": {',
        allocString,
        '},',
        '  "number": "0x0",',
        '  "gasUsed": "0x0",',
        '  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000"',
        '}'));
    }
    
    function uint2str(uint i) internal pure returns (string){
        if (i == 0) return "0";
        uint j = i;
        uint length;
        while (j != 0){
            length++;
            j /= 10;
        }
        bytes memory bstr = new bytes(length);
        uint k = length - 1;
        while (i != 0){
            bstr[k--] = byte(48 + i % 10);
            i /= 10;
        }
        return string(bstr);
    }

    function addr2str(address _addr) public pure returns(string) {
        bytes32 value = bytes32(uint256(_addr));
        bytes memory alphabet = "0123456789abcdef";
    
        bytes memory str = new bytes(51);

        for (uint i = 0; i < 20; i++) {
            str[i*2] = alphabet[uint(uint8(value[i + 12] >> 4))];
            str[i*2+1] = alphabet[uint(uint8(value[i + 12] & 0x0f))];
        }
        return string(str);
    }
}


/**
 * @title SafeMath
 * @dev Math operations with safety checks that revert on error
 */
library SafeMath {

  /**
  * @dev Multiplies two numbers, reverts on overflow.
  */
  function mul(uint256 _a, uint256 _b) internal pure returns (uint256) {
    // Gas optimization: this is cheaper than requiring 'a' not being zero, but the
    // benefit is lost if 'b' is also tested.
    // See: https://github.com/OpenZeppelin/openzeppelin-solidity/pull/522
    if (_a == 0) {
      return 0;
    }

    uint256 c = _a * _b;
    require(c / _a == _b);

    return c;
  }

  /**
  * @dev Integer division of two numbers truncating the quotient, reverts on division by zero.
  */
  function div(uint256 _a, uint256 _b) internal pure returns (uint256) {
    require(_b > 0); // Solidity only automatically asserts when dividing by 0
    uint256 c = _a / _b;
    // assert(_a == _b * c + _a % _b); // There is no case in which this doesn't hold

    return c;
  }

  /**
  * @dev Subtracts two numbers, reverts on overflow (i.e. if subtrahend is greater than minuend).
  */
  function sub(uint256 _a, uint256 _b) internal pure returns (uint256) {
    require(_b <= _a);
    uint256 c = _a - _b;

    return c;
  }

  /**
  * @dev Adds two numbers, reverts on overflow.
  */
  function add(uint256 _a, uint256 _b) internal pure returns (uint256) {
    uint256 c = _a + _b;
    require(c >= _a);

    return c;
  }

  /**
  * @dev Divides two numbers and returns the remainder (unsigned integer modulo),
  * reverts when dividing by zero.
  */
  function mod(uint256 a, uint256 b) internal pure returns (uint256) {
    require(b != 0);
    return a % b;
  }
}
