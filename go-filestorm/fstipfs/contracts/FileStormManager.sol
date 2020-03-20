pragma solidity =0.4.24 <0.6.3;
/**
 * @title FileStormManager.sol
 * This is the smart contract to manage filestorm nodes on a blockchain.
 */

contract FileStormManager {

    address internal owner;
    mapping(address => uint) public admins;

    uint256 public BLOCK_SECOND = 5;
    uint256 public DISBURSE_EPOCH = 720 * 24;
    uint256 public staking_amount = 1000 * 10 ** 18;
    uint256 public disburse_amount = 8 * 10 ** 17;

    // create a parameter for withdraw wait block
        // create mapping of nodes, staking amount, total disbursement, ..., time/block, finish block

    constructor() public payable {
        owner = msg.sender;
        admins[msg.sender] = 1;
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

    function updateStakingAmount(uint256 amount) public {
        require(admins[msg.sender] == 1, "Only Admins Can Change Staking Amount.");
        staking_amount = amount;
    }

    function updateDisburseAmount(uint256 amount) public {
        require(admins[msg.sender] == 1, "Only Admins Can Change Disburse Amount.");
        disburse_amount = amount;
    }

    addNode(address nodeAddress, string fileStormId) payable public {
        // token amount should be multiple of 1000
        // addNode to global node mapping
    }

    removeNode(address nodeAddress, string fileStormId) public {
        // check if it's the owner
        // remove node from mapping
        // wait withdraw block to remove and give back tokens.
    }

    disburse() public {
        // 200 million tokens in total mined in 30 years
        // disburse evenly to all participating nodes based on how much the put in stake
        // remove node based on finish block number
    }
}