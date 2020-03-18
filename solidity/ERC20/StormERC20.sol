pragma solidity >=0.4.22;

contract StormToken
{
    string public constant name = "Storm Token";
    string public constant symbol = "STO";
    uint8 public constant decimals = 18;

    struct Account
    {
        uint256 balance;

        //addresses that the account has given an allowance to
        mapping(address => uint256) allowed;
    }

    event Approval
    (
        address indexed tokenOwner,
        address indexed spender,
        uint256 tokens
    );

    event Transfer
    (
        address indexed from,
        address indexed to,
        uint256 tokens
    );

    uint256 _totalSupply;
    mapping(address => Account) accounts;

    constructor(uint256 numOfTokens) public
    {
        _totalSupply = numOfTokens;
        accounts[msg.sender].balance = numOfTokens; //initially give all tokens to constructor sender
    }

    function totalSupply() public view returns (uint256)
    {
        return _totalSupply;
    }

    function balanceOf(address account) public view returns (uint256)
    {
        return accounts[account].balance;
    }

    function transfer(address to, uint256 numOfTokens) public returns (bool)
    {
        //make sure sender has enough tokens to transfer
        require(accounts[msg.sender].balance >= numOfTokens,
                "Sender does not have sufficient funds to transfer");

        //take away tokens from sender
        accounts[msg.sender].balance -= numOfTokens;

        //place tokens in to address
        accounts[to].balance += numOfTokens;

        emit Transfer(msg.sender, to, numOfTokens);

        return true;
    }

    function approve(address delegate, uint numOfTokens) public returns (bool)
    {
        //give delegate an allowance and add to allowed mapping
        accounts[msg.sender].allowed[delegate] = numOfTokens;

        emit Approval(msg.sender, delegate, numOfTokens);

        return true;
    }

    function allowance(address owner, address delegate) public view returns (uint256)
    {
        return accounts[owner].allowed[delegate];
    }

    function transferFrom(address owner, address buyer, uint256 numOfTokens) public returns (bool)
    {
        require(accounts[owner].balance >= numOfTokens,
                "Owner has insufficient tokens");
        require(accounts[owner].allowed[msg.sender] >= numOfTokens,
                "Sender's allowance is not high enough");

        //reduce senders allowance by numOfTokens
        accounts[owner].allowed[msg.sender] -= numOfTokens;

        //reduce owners balance
        accounts[owner].balance -= numOfTokens;

        //transfer numOfTokens to the buyer
        accounts[buyer].balance += numOfTokens;

        //emit a Transfer event
        emit Transfer(owner, buyer, numOfTokens);

        return true;
    }
}
