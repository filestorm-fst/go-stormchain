const StormToken = artifacts.require("./StormToken721.sol");

contract("StormToken721", function(accounts) {

  var stormToken;
  var creator = accounts[0];

  it("Balance of creator should equal maxSupply", function(){
    return StormToken.deployed().then(function(instance){
      stormToken = instance;
      return stormToken.balanceOf(creator);
    }).then(function(res){
      assert.equal(100000, res.toNumber(), "Creator's balance should be 100,000");
    });
  });

  it("Creator can issue new tokens", function(){
    return StormToken.deployed().then(function(instance){
      stormToken = instance;
    }).then(function(res){
      stormToken.issueTokens(100);
      return stormToken.balanceOf(creator);
    }).then(function(res){
      assert.equal(100100, res.toNumber(), "Creator's balance should be 100,100");
      return stormToken.totalSupply();
    }).then(function(res){
      assert.equal(res.toNumber(), 100100, "Total supply should be 100,100");
      return stormToken.ownerOf(100001);
    }).then(function(res){
      assert.equal(res, creator, "Creator should own token #100001");
    });
  });

  it("Can Transfer from one address to another using approve", function(){
    return StormToken.deployed().then(function(instance){
      stormToken = instance;
    }).then(function(res){
      stormToken.approve(accounts[1], 1);
      return stormToken.getApproved(1);
    }).then(function(res){
      assert.equal(res, accounts[1]);
      stormToken.transferFrom(creator, accounts[1], 1);
      return stormToken.balanceOf(creator);
    }).then(function(res){
      assert.equal(100099, res.toNumber(), "Creator's balance should be 100099");
      return stormToken.balanceOf(accounts[1]);
    }).then(function(res){
      assert.equal(1, res.toNumber(), "Account 1's balance should be 1");
    });
  });

  it("Can Transfer from one address to another using setApprovalForAll", function(){
    return StormToken.deployed().then(function(instance){
      stormToken = instance;
    }).then(function(res){
      stormToken.setApprovalForAll(accounts[1], true);
      return stormToken.isApprovedForAll(creator, accounts[1]);
    }).then(function(res){
      assert.equal(res, true);
      stormToken.transferFrom(creator, accounts[1], 100);
      stormToken.transferFrom(creator, accounts[1], 120);
      stormToken.transferFrom(creator, accounts[1], 130);
      stormToken.transferFrom(creator, accounts[1], 140);
      return stormToken.balanceOf(creator);
    }).then(function(res){
      assert.equal(100095, res.toNumber(), "Creator's balance should be 100095");
      return stormToken.balanceOf(accounts[1]);
    }).then(function(res){
      assert.equal(5, res.toNumber(), "Account 1's balance should be 5");
    });
  });
});
