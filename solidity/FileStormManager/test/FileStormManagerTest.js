const FileStormManager = artifacts.require("./FileStormManager.sol");

contract("FileStormManager", function(accounts) {

  var filestorm;
  var creator = accounts[0];
  var account1 = accounts[1];
  var account2 = accounts[2];
  var prevAmount;
  var currAmount;

  it("Can add creator as a node", function(){
    return FileStormManager.deployed().then(function(instance){
      filestorm = instance;
      filestorm.addNode(creator, 'creator', {from: creator, value: web3.utils.toWei('5000', 'ether')});
    }).then(function(){
      return filestorm.isNode(creator);
    }).then(function(res){
      assert.equal(res, true, "Creator should be a node.");
      return filestorm.nodeStakeAmount(creator);
    }).then(function(res){
      assert.equal(web3.utils.fromWei(res, 'ether'), 5000, "Staking amount should be 5000 ether.");
    });
  });

  it("Can remove node", function(){
    return FileStormManager.deployed().then(function(instance){
      filestorm = instance;
      filestorm.updateDisburseEpoch(0);
      filestorm.removeNode(creator, "creator");
      return filestorm.isNode(creator);
    }).then(function(res){
      assert.equal(res, false, "Creator is still a node");
    });
  });

  it("Disburse function adds ether to balance", function(){
    return FileStormManager.deployed().then(function(instance){
      filestorm = instance;
      filestorm.addNode(creator, "creator", {from: creator,  value: web3.utils.toWei('1000', 'ether')});
      return filestorm.isNode(creator);
    }).then(function(res){
      assert.equal(res, true, "Creator is not a node");
      return web3.eth.getBalance(creator);
    }).then(function(res){
      prevAmount = res;
      return filestorm.disburse();
    }).then(function(res){
      return web3.eth.getBalance(creator);
    }).then(function(res){
      currAmount = res;
      assert.equal(prevAmount < currAmount, true, "Did not gain ether")
    });
  });

  it("Removing node gets stake amount back", function(){
    FileStormManager.deployed().then(function(instance){
      filestorm = instance;
      filestorm.removeNode(creator, "creator");
      return filestorm.isNode(creator);
    }).then(function(res){
      assert.equal(res, false, "Creator is still a node after being removed");
      return web3.eth.getBalance(creator);
    }).then(function(res){
      prevAmount = res;
      return filestorm.disburse();
    }).then(function(res){
      return web3.eth.getBalance(creator);
    }).then(function(res){
      currAmount = res;
      //make sure currAmount - prevAmount > original staked value
      assert.equal((currAmount - prevAmount) > web3.utils.toWei('1000', 'ether'), true,
                    "Creator did not get staked amount back");
    });
  });
});
