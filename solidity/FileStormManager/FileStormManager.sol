// SPDX-License-Identifier: MIT

pragma solidity ^0.6.0;
/**
 * @title FileStormManager.sol
 * This is the smart contract to manage filestorm nodes on a blockchain.
 */

import "ReentrancyGuard.sol";

contract FileStormManager is ReentrancyGuard {

    address internal owner;
    mapping(address => uint) public admins;

    uint256 public BLOCK_SECOND = 5;
    uint256 public DISBURSE_EPOCH = 720 * 24; //once a day
    uint256 public DISBURSE_TOTAL = 400; // disburse only 400 times. 
    uint256 public staking_amount = 500 * 10 ** 18; // staking amount per terabyte.
    uint256 public staking_limit = 96; // one machine can only store 96TB at most. 

    // add limit for each register type
    enum RegisterType {consensus, store, retrieval, instrastructure}

    mapping(uint => uint256) public disburseMapping;

    enum NodeStatus {registered, working, error, ended}

    struct Node {
      string  nodeId;
      address nodeAddress;
      address beneficiary;
      uint256 stakingAmount;
      uint256 disburseAmount;
      uint256 registerBlock;
      uint registerType;
      uint nodeStatus; 
      uint256 nextRewardBlock;
      uint256 disburseCount;
      uint256 disbursedTotal;
    }

    mapping(address => Node) public nodeMapping;

    constructor() public payable {
      owner = msg.sender;
      admins[msg.sender] = 1;
      
      disburseMapping[uint(RegisterType.consensus)] = 25 * 10 ** 17;
      disburseMapping[uint(RegisterType.retrieval)] = 35 * 10 ** 17;
      disburseMapping[uint(RegisterType.store)] = 375 * 10 ** 16;
      disburseMapping[uint(RegisterType.instrastructure)] = 4 * 10 ** 18;
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

    function updateStakingLimit(uint256 limit) public {
      require(admins[msg.sender] == 1, "Only Admins Can Change Staking Limit.");
      staking_limit = limit;
    }

    function updateDisburseAmount(uint registerType, uint256 amount) public {
      require(admins[msg.sender] == 1, "Only Admins Can Change Disburse Amount.");
      disburseMapping[registerType] = amount;
    }

    function addNode(string memory nodeId, address nodeAddress, address beneficiary, uint256 stakingCount, uint registerType) public payable {

      require(admins[msg.sender] == 1 || msg.sender == nodeAddress || msg.sender == beneficiary
      , "Only Admins or Node Owners Can Register a Node.");

      require(nodeMapping[nodeAddress].nodeStatus == 0 || nodeMapping[nodeAddress].nodeStatus != uint(NodeStatus.ended)
      , "Node is already registered.");
      
      require(disburseMapping[registerType]>0, "Wrong registration type.");
      require(stakingCount <= staking_limit, "staking total reached limit.");

      nodeMapping[nodeAddress].nodeId = nodeId;
      nodeMapping[nodeAddress].nodeAddress = nodeAddress;
      nodeMapping[nodeAddress].beneficiary = beneficiary;
      nodeMapping[nodeAddress].stakingAmount = staking_amount * stakingCount;
      nodeMapping[nodeAddress].disburseAmount = disburseMapping[registerType] * stakingCount;
      nodeMapping[nodeAddress].registerType = registerType;
      nodeMapping[nodeAddress].registerBlock = block.number;
      nodeMapping[nodeAddress].nodeStatus = uint(NodeStatus.registered);
      nodeMapping[nodeAddress].nextRewardBlock = block.number + DISBURSE_EPOCH;
      nodeMapping[nodeAddress].disburseCount = 0;
      nodeMapping[nodeAddress].disbursedTotal = 0;
      
    }

    function updateNodeStatus(address nodeAddress, uint nodeStatus) public {
      require(admins[msg.sender] == 1, "Only Admins Can Update Node Status.");
      nodeMapping[nodeAddress].nodeStatus = nodeStatus;
    }


    function disburse(address nodeAddress) nonReentrant public {
      require(admins[msg.sender] == 1 || msg.sender == nodeAddress || msg.sender == nodeMapping[nodeAddress].beneficiary
      , "Only Admins or Node Owners Can Request Disbursement.");

      require(nodeMapping[nodeAddress].nextRewardBlock <= block.number
      , "Next Disbursement block not reached.");

      if (nodeMapping[nodeAddress].nodeStatus == uint(NodeStatus.working)) {
        nodeMapping[nodeAddress].beneficiary.call{value: nodeMapping[nodeAddress].disburseAmount};
      }
      nodeMapping[nodeAddress].nextRewardBlock = block.number + DISBURSE_EPOCH;
      nodeMapping[nodeAddress].disburseCount ++;
      nodeMapping[nodeAddress].disbursedTotal += nodeMapping[nodeAddress].disburseAmount;

      if (nodeMapping[nodeAddress].disburseCount == DISBURSE_TOTAL) {
        nodeMapping[nodeAddress].nodeStatus = uint(NodeStatus.ended);
      }
    }
    
}
