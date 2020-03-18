const QuadraticBallot = artifacts.require("./QuadraticBallot.sol");

contract("QuadraticBallot", function(accounts) {

  var quadraticBallot;
  var chairperson = accounts[0];
  var person1 = '0x2DD05e6512f26a816cbd1A7551A968d6c1B3BC6f'
  var person2 = '0x222160DEC7F5717f43dE8b2bFA272021bA37FB35'
  var person3 = accounts[3];
  it("Should return proposal one as winner", function(){
    return QuadraticBallot.deployed().then(function(instance){
      quadraticBallot = instance;
      quadraticBallot.vote(1, 1, {from: chairperson, value: web3.toWei(1, "ether")});
      //console.log("Balance after placing vote: " + web3.fromWei(web3.eth.getBalance(chairperson), 'ether'));
      return quadraticBallot.winningProposal.call();
    }).then(function(res){
      //console.log("Balance after vote is finsihed: " + web3.fromWei(web3.eth.getBalance(chairperson), 'ether'));
      assert.equal(1, res[0], "Winning proposal should have been 1");
    });
  });

  it("Person 1 should be able to vote", function(){
    return QuadraticBallot.deployed().then(function(instance){
      quadraticBallot = instance;
      return quadraticBallot.giveRightToVote(person1, {from: chairperson});
    }).then(function(){
      quadraticBallot.vote(2, 3, {from:person1, value: web3.toWei(9, 'ether')});
      return quadraticBallot.winningProposal.call();
    }).then(function(res){
      assert.equal(2, res[0].toNumber(), "Winning proposal should have been 2");
    });
  });

  it('Person 2 should not be able to vote', function(){
    return QuadraticBallot.deployed().then(function(instance){
      quadraticBallot = instance;
      quadraticBallot.vote(1, 2, {from: person2, value: web3.toWei(4, 'ether')});
    });
      quadraticBallot.vote(9, 1, {from: chairperson, value: web3.toWei(1, 'ether')});
      return quadraticBallot.winningProposal.call().then(function(res){
        assert.equal(res[0].toNumber(), 9, "Winning Proposal should've been 9");
      });
    });

  it('Person 1 cannot give permission to Person 2', function(){
    return QuadraticBallot.deployed().then(function(instance){
      quadraticBallot = instance;
      quadraticBallot.giveRightToVote(person1, {from: chairperson});
      quadraticBallot.giveRightToVote(person2, {from: person1});
      quadraticBallot.vote(2, 1, {from: person2, value: web3.toWei(4, 'ether')});
    });
      quadraticBallot.vote(1, 1, {from: person1, value: web3.toWei(1, 'ether')});
      return quadraticBallot.winningProposal.call().then(function(res){
        assert.equal(res[0], 2, "Winning proposal should've been 2");
      });
  });
});
