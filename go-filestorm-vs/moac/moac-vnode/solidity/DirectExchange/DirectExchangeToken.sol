pragma solidity ^0.4.11;

import "./StandardToken.sol";
import "./SafeMath.sol";


/**
 * @title Standard ERC20 token
 *
 * @dev Implementation of the basic standard token.
 * https://github.com/ethereum/EIPs/issues/20
 * Based on code by FirstBlood: https://github.com/Firstbloodio/token/blob/master/smart_contract/FirstBloodToken.sol
 */
contract DirectExchangeToken is StandardToken {
  using SafeMath for uint256;

  address public owner;

  uint256 public priceOneGInMOAC;
  
  function() public payable {
        //only allow protocol send
        // require(protocol == msg.sender);
    }
  
  function DirectExchangeToken(uint256 tokensupply, uint256 exchangerate) public {
    //total token supply
    totalSupply_ = tokensupply * 10 ** 18;
    balances[this] = totalSupply_;
    priceOneGInMOAC = exchangerate;
    
    owner = msg.sender;
  }

  function updateOwner(address newowner) public {
    require(owner == msg.sender);
    owner = newowner;
  }

  function buyMintToken(address useraddr, uint256 value) public returns (uint256, uint256) {
    require(owner == msg.sender);

    uint256 token = value * priceOneGInMOAC;
    uint256 refund = 0;
    if(token > balances[this]){
      refund = ( token - balances[this]) / priceOneGInMOAC;
      balances[useraddr] += balances[this];
      token = balances[this];
      balances[this] = 0;
    } else {
      balances[this] -= token;
      balances[useraddr] += token;
      refund = 0;
    }
    return (refund, token);
  }

  //return amount of MOAC to be refund
  function sellMintTokenPre(address useraddr, uint256 amount) public view returns (uint256) {
    require(owner == msg.sender);
    require(amount <= balances[useraddr]);

    uint256 moac = amount / priceOneGInMOAC;
    return moac;
  }

  function sellMintToken(address useraddr, uint256 amount) public returns (bool) {
    require(owner == msg.sender);
    
    balances[useraddr] -= amount;
    balances[this] += amount;
    uint256 moac = amount / priceOneGInMOAC;
    useraddr.transfer( moac );
    return true;
  }

  function requestEnterMicrochain(address useraddr, uint256 amount) public returns (bool) {
    require(owner == msg.sender);
    require(balances[useraddr] >= amount);

    balances[useraddr] -= amount;
    // balances[this] += amount;
    // holdingPool[block.number].userAddr.push(useraddr);
    // holdingPool[block.number].amount.push(amount);
    // holdingPoolList.push(useraddr);
    //event
    return true;
  }

  function redeemFromMicroChain(address[] addr, uint256[] bals) public returns (bool){
    require(owner == msg.sender);
    require(addr.length == bals.length);
    for( uint i=0; i<addr.length; i++ ) {
      // require(balances[this] >= bals[i]);
      // balances[this] -= bals[i];
      balances[addr[i]] += bals[i];
      //event
    }
    return true;
  }

  function withdrawTokenMoac(address useraddr, uint256 moac)  public {
    require(owner == msg.sender);

    useraddr.transfer( moac );
  }
}
