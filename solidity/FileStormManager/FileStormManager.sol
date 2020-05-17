pragma solidity ^0.5.2;
/**
 * @title FileStormManager.sol
 * This is the smart contract to manage filestorm nodes on a blockchain.
 */

contract FileStormManager {

    address internal owner;
    mapping(address => uint) public admins;

    uint256 public BLOCK_SECOND = 5;
    uint256 public DISBURSE_EPOCH = 720 * 24; //once a day
    uint256 public staking_amount = 1000 * 10 ** 18;
    uint256 public disburse_amount = 8 * 10 ** 17;

    // create a parameter for withdraw wait block
    // create mapping of nodes, staking amount, total disbursement, ..., time/block, finish block

    //total disbursement
    uint256 constant public totalDisbursement = 2 * 10 ** 8; //200 million

    //mapping of nodes in network
    mapping(address => bool) nodes;

    //mapping of staking amount for each node in network
    mapping(address => uint256) stakingAmount;

    //mapping for each address to their filestorm ID
    mapping(string => address) stormIDOwner;

    //mapping for finish block for each removed node
    mapping(address => uint256) finishBlock;

    //array to keep track of all nodes to be removed
    address[] nodeRemovals;

    //array of all nodes to be disbursed
    address[] disburseNodes;

    constructor() public payable {
        owner = msg.sender;
        admins[msg.sender] = 1;
    }

    function addAdmin(address admin) public {
        require(admins[msg.sender] == 1, "Only Admins Can Add Another Admin.");
        admins[admin] = 1;
    }

    function removeAdmin(address admin) public {
        require(admins[msg.sender] == 1, "Only Admins Can Remove Another Admin.");
        require(admin != msg.sender, "Admins Cannot Remove Self.");
        admins[admin] = 0;
    }

    function updateStakingAmount(uint256 amount) public {
        require(admins[msg.sender] == 1, "Only Admins Can Change Staking Amount.");
        staking_amount = amount;
    }

    function updateDisburseAmount(uint256 amount) public {
        require(admins[msg.sender] == 1, "Only Admins Can Change Disburse Amount.");
        disburse_amount = amount;
    }

    function addNode(address nodeAddress, string memory fileStormId) payable public {
        //make sure staking amount is greater than the minimum
        require(msg.value >= staking_amount, "Staking amount not high enough");

        // token amount should be multiple of 1000
        require(msg.value % 1000 == 0, "Token Amount Should Be A Multiple Of 1000.");

        // addNode to global node mapping
        nodes[nodeAddress] = true;
        stakingAmount[nodeAddress] = msg.value;
        stormIDOwner[fileStormId] = nodeAddress;
        disburseNodes.push(nodeAddress);

    }

    function removeNode(address nodeAddress, string memory fileStormId) public {
        // check if it's the owner
        require(nodeAddress == msg.sender, "Only owner of node can remove node");
        require(nodeAddress == stormIDOwner[fileStormId], "FileStormId does not match the nodeAddress.");

        // remove node from mapping
        delete nodes[nodeAddress]; //remove from nodes in network
        delete stormIDOwner[fileStormId]; //remove stormID from network

        // wait withdraw block to remove and give back tokens.
        nodeRemovals.push(nodeAddress); //add to array and disburse funcion will take care of releasing stake ammount
        finishBlock[nodeAddress] = block.number + DISBURSE_EPOCH;
    }

    function disburse() public {
        // 200 million tokens in total mined in 30 years
        // disburse evenly to all participating nodes based on how much the put in stake
        uint256 nodeAmount = disburseNodes.length;
        for(uint256 i = 0; i < nodeAmount; i++)
        {
          address payable temp = address(uint160(disburseNodes[i]));
          temp.transfer(disburse_amount);
        }

        // remove node based on finish block number
        uint256 currentBlock = block.number;
        uint256 length = nodeRemovals.length;
        for(uint256 i = 0; i < length; i++)
        {
          address removeAddress = nodeRemovals[i];
          if(finishBlock[removeAddress] >= currentBlock)
          {
            //give back stake
            address payable temp = address(uint160(removeAddress));
            temp.transfer(stakingAmount[removeAddress]);
            //remove from array
            delete nodeRemovals[i];
            delete stakingAmount[removeAddress];
            delete finishBlock[removeAddress];
          }
        }
    }

    //added functions for testing purposes
    function isNode(address node) public view returns (bool)
    {
      return nodes[node];
    }

    function nodeStakeAmount(address node) public view returns (uint256)
    {
      return stakingAmount[node];
    }

    function isAdmin(address addr) public view returns (bool)
    {
      if(admins[addr] == 1) return true;
      return false;
    }

    function updateDisburseEpoch(uint256 time) public
    {
      DISBURSE_EPOCH = time;
    }
}
