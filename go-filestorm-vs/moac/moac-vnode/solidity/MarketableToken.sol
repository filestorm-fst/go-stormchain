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
contract MarketableToken is StandardToken {
  using SafeMath for uint256;

  //owner holding
  mapping(address => uint256) lockPeriod;
  
  //allocate once
  bool isAllocated;
  address public owner;

  //for minting
  uint256 public curSupplyCap;
  uint256 public curSupply;
  uint256 public curStartBlock;
  uint256 public curEndBlock;
  uint256 public curPrice;
  uint256 public curStage;
  uint256 public totalMint;
  uint256 public curSellLimit;
  bool public minting;

  //price adjustment
  //if sold withing before 1/2 supply term
  //price *= 1.1, supply *= (1+addpercent)
  //if sold before supply term
  //price *= 1.02
  //if sold more than half before supply term ends
  //price *= 0.98
  //if sold less than half before supply term ends
  //price *= 0.91

  //price to MOAC, 10**9 token convert to how many MOAC
  //if 100:1, 100 token = 1 MOAC, priceOneGInMOAC = 10**7
  uint256 public priceOneGInMOAC = 10**7;
   //for how many block to release one batch supply and also change the price
  uint256 public supplyTermInBlock = 60000;  //if 10s a block, about one week 
  //supply amount for each release
  uint256 public supplyAmount = 10000*10**18;
  //suplly increatment
  uint256 public supplyAddPer = 50;
  //public sale lock period
  uint256 public publicSaleLockPeriod = 100;

  function() public payable {
        //only allow protocol send
        // require(protocol == msg.sender);
  }

  // 100000000*10**18; //100Million;
  function MarketableToken(uint256 tokensupply) public {
    //total token supply
    totalSupply_ = tokensupply * 10 ** 18;
    balances[this] = totalSupply_;
    
    owner = msg.sender;
  }

  function updateOwner(address newowner) public {
    require(owner == msg.sender);
    owner = newowner;
  }

  function initToken(address[] addr, uint256[] bals, uint256[] lock) public {
    require(owner == msg.sender && !isAllocated);
    require(addr.length == bals.length);
    require(addr.length == lock.length);

    for( uint i=0; i<addr.length; i++){
      require(balances[this] > bals[i]);
      balances[this] -= bals[i];
      balances[addr[i]] += bals[i];
      lockPeriod[addr[i]] = block.number + lock[i];
    }
    isAllocated = true;
    startMint();
  }

  function startMint() private {
    require(owner == msg.sender);
    
    totalMint = balances[this];
    curSupplyCap = supplyAmount;
    if( curSupplyCap > totalMint ) {
      curSupplyCap = totalMint;
    }
    curSupply = curSupplyCap;
    curStartBlock = block.number;
    curEndBlock = curStartBlock + supplyTermInBlock;
    curPrice = 10**9;
    curSellLimit = curSupply/2;
    curStage = 0;
    //curMintLeft = totalMint;
    minting = true;
  }

  function refresh() public {
    //any one can call this
    if( block.number > curEndBlock ) {

      if( balances[this] < 10**18 * 100 ) {
        //mint is done
        minting = false;
        return;
      }
    
      if( curStage == 3 )
      {
        //increase price, //increase cap
        curPrice = curPrice*110/100;
        curSupplyCap = curSupplyCap*110/100;
      } else if(curStage == 2 ) {
        curPrice = curPrice*102/100;
      } else {
        //sold less than half
        if(curSupplyCap/2 < curSupply ) {
          curPrice = curPrice*91/100;
        } else {
          curPrice = curPrice*98/100;
        }
      }
    
      if(balances[this] < curSupplyCap ){
        curSupplyCap = balances[this];
      }
      curSupply = curSupplyCap;
      curStartBlock = block.number;
      curEndBlock = curStartBlock + supplyTermInBlock;
      curSellLimit = curSupply/2;
      curStage = 0;
      //curMintLeft -= curSupplyCap;
    }
  }

  function buyMintToken(address useraddr, uint256 value) public returns (uint256) {
    require(owner == msg.sender);
    
    require(curSupply > 0 && minting && block.number > curStartBlock && block.number < curEndBlock);
    require(balances[this] > curSupply);
    uint256 token = value * curPrice / priceOneGInMOAC;
    uint256 refund = 0;

    if(token > curSupply){
      //only buy existing and refund back to sender
      uint256 actualMOAC = curSupply * priceOneGInMOAC / curPrice;
      require (value > actualMOAC );
      refund = ( value - actualMOAC);
      balances[this] -= curSupply;
      balances[useraddr] += curSupply;
      curSupply = 0;

      //check sales condition
      if( (curEndBlock-block.number) > (block.number - curStartBlock) ) {
        curStage = 3; // sold quickly
      } else {
        curStage = 2;
      }
    } else {
      balances[this] -= curSupply;
      balances[useraddr] += token;
      curSupply -= token;

      //consider the case that supply is too small not worth a tx, consider it is sold
      if( (curSupply * priceOneGInMOAC / curPrice) < 1){
        if( (curEndBlock-block.number) > (block.number - curStartBlock) ) {
          curStage = 3; // sold quickly
        } else {
          curStage = 2;
        }
      }
      refund = 0;
    }

    lockPeriod[useraddr] = block.number + publicSaleLockPeriod;
    return refund;
  }

  //return amount of MOAC to be refund
  function sellMintTokenPre(address useraddr, uint256 amount) public view returns (uint256) {
    require(owner == msg.sender);
    
    require(minting); //only in minting stage can buy/sell with this method
    require(lockPeriod[useraddr] < block.number);
    require(amount < balances[useraddr]);
    require(amount < curSellLimit);

    uint256 moac = amount* priceOneGInMOAC / curPrice;
    return moac;
  }

  function sellMintToken(address useraddr, uint256 amount) public returns (bool) {
    require(owner == msg.sender);
    balances[useraddr] -= amount;
    balances[this] += amount;
    uint256 moac = amount* priceOneGInMOAC / curPrice;
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

}
