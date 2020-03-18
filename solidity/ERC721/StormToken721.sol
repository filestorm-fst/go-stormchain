pragma solidity >=0.5.16 <0.6.0;

import "./ERC721.sol";
import "./CheckERC165.sol";
import "./ERC721TokenReceiver.sol";

contract StormToken721 is ERC721, CheckERC165
{

    //tokens have uint ID's from 1 to maxSupply

    event Approval
    (
        address indexed tokenOwner,
        address indexed spender,
        uint256 tokenID
    );

    event Transfer
    (
        address indexed from,
        address indexed to,
        uint256 tokenID
    );

    event ApprovalForAll
    (
        address indexed owner,
        address indexed operator,
        bool approved
    );

    uint256 public maxSupply;

    //creator of the contract
    address creator;

    //mapping of the unique tokens each address has
    mapping(address => mapping(uint256 => uint256)) private ownerTokens;

    //mapping of token count for each address
    mapping(address => uint256) private ownedTokens;

    //mapping of ERC721 tokens to owner
    mapping(uint256 => address) private tokenOwner;

    //mapping from owner to approved addresses
    mapping(address => mapping(address=>bool)) private approvedAddresses;

    //mapping of token data, so you don't need to store data on each token (expensive)
    mapping(uint256 => string) tokenData;

    //keep track of who can transfer the token
    mapping (uint256 => address) internal allowance;

    constructor(uint256 _maxSupply) public
    {
        creator = msg.sender;

        //all tokens start with creator of contract
        ownedTokens[creator] = _maxSupply;

        maxSupply = _maxSupply;

        //add to interface check
        //bytes4() is used for overloaded functions in interface
        supportedInterfaces[
                            this.balanceOf.selector ^
                            this.ownerOf.selector ^
                            bytes4(keccak256("safeTransferFrom(address,address,uint256"))^
                            bytes4(keccak256("safeTransferFrom(address,address,uint256,bytes"))^
                            this.transferFrom.selector ^
                            this.approve.selector ^
                            this.setApprovalForAll.selector ^
                            this.getApproved.selector ^
                            this.isApprovedForAll.selector
                            ] = true;
    }

    function name() public pure returns (string memory)
    {
        return "StormToken";
    }

    function symbol() public pure returns (string memory)
    {
        return "STO";
    }

    function decimals() public pure returns (uint8)
    {
        return 18;
    }

    function totalSupply() public view returns (uint256)
    {
        return maxSupply;
    }

    function balanceOf(address account) public view returns (uint256)
    {
        return ownedTokens[account];
    }

    function issueTokens(uint256 _extraTokens) public
    {
        //make sure only creator can issue tokens
        require(msg.sender == creator,
                "Sender is not the creator");

        ownedTokens[creator] += _extraTokens;

        //emit a transfer event for all tokens isued
        for(uint256 i = maxSupply; i < maxSupply + _extraTokens; i++)
        {
            tokenOwner[i] = creator;
            emit Transfer(address(0x0), creator, i);
        }

        maxSupply += _extraTokens;
    }

    function ownerOf(uint256 _tokenID) public view returns (address)
    {
        //make sure the token exists
        require(tokenExists(_tokenID) == true,
                "Token does not exist");

        //if nobody owns it, return the creator of the contract
        if(tokenOwner[_tokenID] == address(0))
        {
            return creator;
        }
        return tokenOwner[_tokenID];
    }

    function isApprovedForAll(address owner, address operator) external view returns (bool)
    {
        return approvedAddresses[owner][operator];
    }

    function setApprovalForAll(address operator, bool approved) public
    {
        approvedAddresses[msg.sender][operator] = approved;
        emit ApprovalForAll(msg.sender, operator, approved);
    }

    function getApproved(uint256 _tokenID) public view returns (address)
    {
        require(tokenExists(_tokenID) == true,
                "Token does not exist");
        return allowance[_tokenID];
    }

    function approve(address _to, uint256 _tokenID) external payable
    {
        address owner = ownerOf(_tokenID);
        require(owner == msg.sender || approvedAddresses[owner][msg.sender] == true,
                "Sender does not have the authority");
        allowance[_tokenID] = _to;
        emit Approval(owner, _to, _tokenID);
    }

    function transferFrom(address _from, address _to, uint256 _tokenID) external payable
    {
        address owner = ownerOf(_tokenID);
        require (owner == msg.sender || allowance[_tokenID] == msg.sender || approvedAddresses[owner][msg.sender],
                 "Sender is not authorized to transfer");

        require(owner == _from,
                "From address must be the owner of the token");

        require(_to != address(0x0),
                "To address is 0");

        tokenOwner[_tokenID] == _to;
        ownedTokens[_to]++;
        ownedTokens[_from]--;

        if(allowance[_tokenID] != address(0x0))
        {
          delete allowance[_tokenID];
        }
        emit Transfer(_from, _to, _tokenID);
    }

    function safeTransferFrom(address _from, address _to, uint256 _tokenID, bytes calldata data) external payable
    {
      address owner = ownerOf(_tokenID);
      require (owner == msg.sender || allowance[_tokenID] == msg.sender || approvedAddresses[owner][msg.sender],
               "Sender is not authorized to transfer");

      require(owner == _from,
              "From address must be the owner of the token");

      require(_to != address(0x0),
              "To address is 0");

      tokenOwner[_tokenID] == _to;
      ownedTokens[_to]++;
      ownedTokens[_from]--;

      if(allowance[_tokenID] != address(0x0))
      {
        delete allowance[_tokenID];
      }
      emit Transfer(_from, _to, _tokenID);

        uint32 size;
        //check size of _to and store into size
        assembly
        {
            size :=extcodesize(_to)
        }

        if(size > 0)
        {
            //_to address belongs to a contract
            ERC721TokenReceiver receiver = ERC721TokenReceiver(_to);

            require(receiver.onERC721Received(msg.sender,_from,_tokenID,data) == bytes4(keccak256("onERC721Received(address,address,uint256,bytes)")));
        }

    }

    function safeTransferFrom(address _from, address _to, uint256 _tokenID) external payable
    {
      address owner = ownerOf(_tokenID);
      require (owner == msg.sender || allowance[_tokenID] == msg.sender || approvedAddresses[owner][msg.sender],
               "Sender is not authorized to transfer");

      require(owner == _from,
              "From address must be the owner of the token");

      require(_to != address(0x0),
              "To address is 0");

      tokenOwner[_tokenID] == _to;
      ownedTokens[_to]++;
      ownedTokens[_from]--;

      if(allowance[_tokenID] != address(0x0))
      {
        delete allowance[_tokenID];
      }
      emit Transfer(_from, _to, _tokenID);

        uint32 size;
        //check size of _to and store into size
        assembly
        {
            size :=extcodesize(_to)
        }

        if(size > 0)
        {
            //_to address belongs to a contract
            ERC721TokenReceiver receiver = ERC721TokenReceiver(_to);
            bytes memory data = "";
            require(receiver.onERC721Received(msg.sender,_from,_tokenID,data) == bytes4(keccak256("onERC721Received(address,address,uint256,bytes)")));
        }
    }

    function tokenExists(uint256 _tokenID) private view returns (bool)
    {
        return _tokenID != 0 && _tokenID <= maxSupply;
    }
}
