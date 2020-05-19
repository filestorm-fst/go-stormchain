pragma solidity ^0.4.11;
pragma experimental ABIEncoderV2;
//David Chen
//Subchain protocol definition for consenus.
//It also contains the logic to allow user to join this
//pool with bond
import "./SubChainProtocolBase.sol";


contract SCSRelay {
    // 0-registeropen, 1-registerclose, 2-createproposal, 3-disputeproposal, 4-approveproposal, 5-registeradd
    function notifySCS(address cnt, uint256 msgtype) public returns (bool success);
}

//ATO new
contract Token {
    function allowance(address _owner, address _spender) public view returns (uint256);
    function transferFrom(address _from, address _to, uint256 _value) public returns (bool);
    function transfer(address _to, uint256 _value) public returns (bool);
    function balanceOf(address _owner) public view returns (uint256 balance);
    function totalSupply() public view returns (uint256);
    function decimals() public view returns (uint256);
}


contract SubChainBase {
    enum ProposalFlag {noState, pending, disputed, approved, rejected, expired, pendingAccept}
    enum ProposalCheckStatus {undecided, approval, expired}
    enum ConsensusStatus {initStage, workingStage, failure}
    enum SCSRelayStatus {registerOpen, registerClose, createProposal, disputeProposal, approveProposal, registerAdd, regAsMonitor, regAsBackup, updateLastFlushBlk, distributeProposal, reset, uploadRedeemData, requestEnterAndRedeem, requestRelease}
    enum SubChainStatus {open, pending, close}

    struct Proposal {
        address proposedBy;
        bytes32 lastApproved;
        bytes32 hash;
        uint256 start;
        uint256 end;
        //bytes newState;
        uint256[] distributionAmount;
        uint256 flag; // one of ProposalFlag
        uint256 startingBlock;
        uint256[] voters; //voters index
        uint256 votecount;
        uint256[] badActors;
        address viaNodeAddress;
        uint256 preRedeemNum;
        address[] redeemAddr;
        uint256[] redeemAmt;
        address[] minerAddr;
        uint256 distributeFlag;
        address[] redeemAgreeList;
    }
    
    struct SyncNode {
        address nodeId;
        string link;
    }
    
    struct holdings {
        address[] userAddr;
        uint256[] amount;
        uint256[] time;
    }

    struct RedeemRecords {
        uint256[] redeemAmount;
        uint256[] redeemtime;
    }
    
    struct VnodeInfo {
        address protocol;
        uint256[] members;
        uint256[] rewards;    //0:blockReward; 1:txReward; 2:viaReward
        uint256 proposalExpiration;
        address VnodeProtocolBaseAddr;
        uint256 penaltyBond;
        uint256 subchainstatus;
        address owner;
        uint256 BALANCE;
        uint256[] redeems;

        address[] nodeList; 
        address[] nodesToJoin;
    }

     struct MonitorInfo {
        address owner;
        address from; // address as id
        uint256 bond; // value
        string link;  // ip:prort
    } 

    address public protocol;
    uint256 public minMember;
    uint256 public maxMember;
    uint256 public selTarget;
    uint256 public consensusFlag; // 0: init stage 1: working stage 2: failure
    uint256 public flushInRound;
    uint256 public initialFlushInRound;

    bytes32 public proposalHashInProgress;
    bytes32 public proposalHashApprovedLast;  //index: 7
    uint256 internal curFlushIndex;
    uint256 internal pendingFlushIndex;
    
    bytes public unUsedPara;

    uint256 internal lastFlushBlk;

    address internal owner;

    //nodes list is updated at each successful flush
    uint256 public nodeCount;
    address[] public nodeList;    //index: 0f

    uint8[2] public randIndex;
    mapping(address => uint256 ) public nodePerformance;
    mapping(bytes32 => Proposal) public proposals;  //index: 12
    mapping(address => uint256) public currentRefundGas;

    uint256 internal registerFlag;

    uint256 public constant proposalExpiration = 24;
    uint256 public constant penaltyBond = 10 ** 18; // 1 Moac penalty
    mapping(address=>address) public scsBeneficiary;
    uint256 public blockReward = 5 * 10 ** 14;    //index: 18
    uint256 public txReward  = 1 * 10 ** 11;
    uint256 public viaReward = 1 * 10 ** 16;

    uint256 public nodeToReleaseCount;
    uint256[5] public nodesToRelease;  //nodes wish to withdraw, only allow 5 to release at a time
    uint256[] public nodesToDispel;

    address[] public nodesToJoin;  //nodes to be joined
    uint256 public joinCntMax;
    uint256 public joinCntNow;
    mapping(address=>uint256) public nodesWatching;  //nodes watching

    SyncNode[] public syncNodes;
    uint256 indexAutoRetire;

    SCSRelay internal constant SCS_RELAY = SCSRelay(0x000000000000000000000000000000000000000d);
    uint256 public constant NODE_INIT_PERFORMANCE = 5;
    uint256 public constant AUTO_RETIRE_COUNT = 2;
    bool public constant AUTO_RETIRE = false;
    address public VnodeProtocolBaseAddr;
    uint256 public constant MONITOR_MIN_FEE = 1 * 10 ** 12;
    uint256 public syncReward = 1 * 10 ** 11;
    uint256 public constant MAX_GAS_PRICE = 20 * 10 ** 9;

    uint256 public constant DEFLATOR_VALUE = 80; // in 1/millionth: in a year, exp appreciation is 12x

    uint256 internal subchainstatus;
    uint256 public BALANCE = 0;    //index: 30
    address public ERCAddr;
    uint256 public ERCRate = 1;
    uint256 public ERCDecimals = 0;
    //temp holdingplace whenentering microchain
    holdings internal holdingPool;
    uint256 public per_recharge_num = 250;
    uint256 public recharge_cycle = 6;
    
    // inidicator of fund needed
    uint256 public contractNeedFund;

    mapping(address=>RedeemRecords) internal records;
    MonitorInfo[] public monitors;
    
    uint256 public constant MAX_DELETE_NUM = 5;
    uint256 public dappRedeemPos = 0;
    uint256 public per_upload_redeemdata_num = 160;//167
    uint256 public per_redeemdata_num = 110;//140
    uint256 public constant max_redeemdata_num = 500;

    uint256 public constant maxFlushInRound = 500;
    uint256 public constant txNumInFlush = 100;
    
    uint256 public totalOperation;
    uint256 public totalBond;

    address public foundation = 0x0000000000000000000000000000000000000001;
    uint256 public constant settlement_block_delay = 50;

    //events
    event ReportStatus(string message);
    event TransferAmount(address addr, uint256 amount);



    //constructor
    function SubChainBase(address proto, address vnodeProtocolBaseAddr, address ercAddr,  uint256 ercRate, uint256 min, uint256 max, uint256 thousandth, uint256 flushRound) public {
        require(min == 1 || min == 3 || min == 5 || min == 7);
        require(max == 11 || max == 21 || max == 31 || max == 51 || max == 99);
        require(flushRound >= 40  && flushRound <= 500);
         
        flushInRound = flushRound;
        initialFlushInRound = flushRound;
        VnodeProtocolBaseAddr = vnodeProtocolBaseAddr;
        SubChainProtocolBase protocnt = SubChainProtocolBase(proto);
        selTarget = protocnt.getSelectionTarget(thousandth, min);
        protocnt.setSubchainActiveBlock();

        ERCAddr = ercAddr;
        ERCRate = ercRate;
        Token erc20 = Token(ERCAddr);
        ERCDecimals = erc20.decimals();
         if (ERCDecimals > 18) {
            revert();
        }
        ERCDecimals = 18 - ERCDecimals;
        uint256 total = erc20.totalSupply();
        BALANCE = total * ERCRate * (10 ** (ERCDecimals));
        if (total != BALANCE / ERCRate / (10 ** (ERCDecimals))) {
            revert();
        }

        minMember = min;
        maxMember = max;
        protocol = proto; //address
        consensusFlag = uint256(ConsensusStatus.initStage);
        owner = msg.sender;
        foundation = msg.sender;

        protocnt.setSubchainExpireBlock(flushInRound*5);
        lastFlushBlk = 2 ** 256 - 1;
        
        randIndex[0] = uint8(0);
        randIndex[1] = uint8(1);
        indexAutoRetire = 0;
        subchainstatus = uint256(SubChainStatus.open);
    }

    function() public payable {
        //only allow protocol send
        require(protocol == msg.sender);
    }

    function setOwner() public {
        // todo david, how can owner be 0
        if (owner == address(0)) {
            owner = msg.sender;
        }
    }

    function getFlushStatus() public view returns (bool) {
        uint256 blk = lastFlushBlk + flushInRound + nodeCount * 2 * proposalExpiration;
        return (block.number <= blk);
    }
    
    function updatePerRechargeNum(uint256 num) public {
        require(owner == msg.sender);

        per_recharge_num = num;
    }
    
    function updateRechargeCycle(uint256 num) public {
        require(owner == msg.sender);

        recharge_cycle = num;
    }
    
    function updatePerUploadRedeemNum(uint256 num) public {
        require(owner == msg.sender);

        per_upload_redeemdata_num = num;
    }

    function updatePerRedeemNum(uint256 num) public {
        require(owner == msg.sender);

        per_redeemdata_num = num;
    }

    function isMemberValid(address addr) public view returns (bool) {
        return nodePerformance[addr] > 0;
    }
    
    function getVnodeInfo() public view returns (VnodeInfo) {
        VnodeInfo vnodeinfo;
        
        vnodeinfo.protocol = protocol;
        uint256[] memory members = new uint256[](2);
        members[0] = minMember;
        members[1] = maxMember;
        vnodeinfo.members = members;
        uint256[] memory rewards = new uint256[](3);
        rewards[0] = blockReward;    //index: 18
        rewards[1] = txReward;
        rewards[2] = viaReward;
        vnodeinfo.rewards = rewards;
        vnodeinfo.proposalExpiration = proposalExpiration;
        vnodeinfo.VnodeProtocolBaseAddr = VnodeProtocolBaseAddr;
        vnodeinfo.penaltyBond = penaltyBond;
        vnodeinfo.subchainstatus = subchainstatus;
        vnodeinfo.owner = owner;
        vnodeinfo.BALANCE = BALANCE;
        uint256[] memory redeems = new uint256[](4);
        redeems[0] = dappRedeemPos;
        redeems[1] = per_upload_redeemdata_num;
        redeems[2] = per_redeemdata_num;
        redeems[3] = max_redeemdata_num;
        vnodeinfo.redeems = redeems;

        vnodeinfo.nodeList = nodeList;
        vnodeinfo.nodesToJoin = nodesToJoin;

        return vnodeinfo;
    }
    
    function getProposal(uint256 types) public view returns (Proposal) {
        if (types == 1) {
            return proposals[proposalHashInProgress];
        } else if (types == 2) {
            return proposals[proposalHashApprovedLast];
        }
    }

    function getSCSRole(address scs) public view returns (uint256) {
        uint256 i = 0;

        for (i = 0; i < nodeList.length; i++) {
            if (nodeList[i] == scs) {
                return 1;
            }
        }
        
        if (nodesWatching[scs] >= 10**9) {
            return 2;
        }
        
        for (i = 0; i < nodesToJoin.length; i++) {
            if (nodesToJoin[i] == scs) {
                return 3;
            }
        }
        
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        if (protocnt.isPerforming(scs)) {
            if (matchSelTarget(scs, randIndex[0], randIndex[1])) {
                return 4;
            }
        }
        
        return 0;
    }

    function registerAsMonitor(address monitor, string link) public payable {
        require(subchainstatus == uint256(SubChainStatus.open));
        require(msg.value >= MONITOR_MIN_FEE);
        require(nodesWatching[monitor] == 0); 
        require(monitor != address(0));
        require(getSCSRole(monitor) == 4 || getSCSRole(monitor) == 0);
        nodesWatching[monitor] = msg.value;
        totalBond += msg.value;

        // Add MonitorInfo        
        monitors.push(MonitorInfo(msg.sender, monitor, msg.value, link));
               
        SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.regAsMonitor));
    }

    function getMonitorInfo() public view returns (address[], string[]) {
        uint256 cnt = monitors.length;
        address[] memory addrlist = new address[](cnt);
        string[] memory strlist = new string[](cnt);
        uint256 i = 0;
        for (i = 0; i < cnt; i++) {
            addrlist[i] = monitors[i].from;
            strlist[i] = monitors[i].link;
        }

        return (addrlist, strlist);
    }

    function removeMonitorInfo(address monitor) public {
        uint256 i = 0;
        uint256 cnt = monitors.length;
        for (i = cnt; i > 0; i--) {
            if (monitors[i-1].from == monitor) {
                if (msg.sender == owner || msg.sender == monitors[i-1].owner)
                { 
                    // withdraw                
                    monitor.transfer(monitors[i-1].bond);
                    totalBond -= monitors[i-1].bond;

                    // delete
                    monitors[i-1] = monitors[cnt-1];
                    delete monitors[cnt-1];
                    monitors.length--;
                    nodesWatching[monitor] = 0;
                    delete nodesWatching[monitor];
                }
                return;
            }
        }
    }

    //v,r,s are the signature of msg hash(scsaddress+subchainAddr)
    function registerAsSCS(address beneficiary, uint8 v, bytes32 r, bytes32 s) public returns (bool) {
        require(subchainstatus == uint256(SubChainStatus.open));
        require(getSCSRole(msg.sender) == 4);
        if (registerFlag != 1) {
            //ReportStatus("Register not open");
            return false;
        }
        //check if valid registered in protocol pool
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        if (!protocnt.isPerforming(msg.sender)) {
            //ReportStatus("SCS not performing");
            return false;
        }

        if (!matchSelTarget(msg.sender, randIndex[0], randIndex[1])) {
            //ReportStatus("SCS not selected");            
            return false;
        }

        // if reach max, reject
        if (nodeCount > maxMember) {
            //ReportStatus("max nodes reached");            
            return false;
        }

        //check if node already registered
        for (uint256 i=0; i < nodeCount; i++) {
            if (nodeList[i] == msg.sender) {
                //ReportStatus("Node already registered");
                return false;
            }
        }

        //make sure msg.sender approve bond deduction
        if (!protocnt.approveBond(msg.sender, penaltyBond, v, r, s)) {
            //ReportStatus("Bond approval failed.");            
            return false;
        }

        nodeList.push(msg.sender);
        nodeCount++;
        nodePerformance[msg.sender] = NODE_INIT_PERFORMANCE;

        //todo: refund gas
        //msg.sender.send(gasleft() * tx.gasprice);

        if (beneficiary == address(0)) {
            scsBeneficiary[msg.sender] = msg.sender;
        }
        else {
            scsBeneficiary[msg.sender] = beneficiary;
        }

        //ReportStatus("Reg successful");

        return true;
    }

    //v,r,s are the signature of msg hash(scsaddress+subchainAddr)
    function registerAsBackup(address beneficiary, uint8 v, bytes32 r, bytes32 s) public returns (bool) {
        require(subchainstatus == uint256(SubChainStatus.open));
        require(getSCSRole(msg.sender) == 4);
        if (registerFlag != 2) {
            return false;
        }

        //check if valid registered in protocol pool
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        if (!protocnt.isPerforming(msg.sender)) {
            return false;
        }

        if (!matchSelTarget(msg.sender, randIndex[0], randIndex[1])) {
            return false;
        }

        //if reach max, reject
        if (joinCntNow >= joinCntMax) {
            return false;
        }

        uint256 i = 0;
        //check if node already registered
        for (i = 0; i < nodeCount; i++) {
            if (nodeList[i] == msg.sender) {
                return false;
            }
        }

        for (i = 0; i < nodesToJoin.length; i++) {
            if (nodesToJoin[i] == msg.sender) {
                return false;
            }
        }

        //make sure msg.sender approve bond deduction
        if (!protocnt.approveBond(msg.sender, penaltyBond, v, r, s)) {
            return false;
        }

        nodesToJoin.push(msg.sender);
        joinCntNow++;
        //set to performance to 0 since backup node has no block synced yet. 
        nodePerformance[msg.sender] = 0;//NODE_INIT_PERFORMANCE;

        //todo: refund gas
        //msg.sender.send(gasleft() * tx.gasprice);

        if (beneficiary == address(0)) {
            scsBeneficiary[msg.sender] = msg.sender;
        }
        else {
            scsBeneficiary[msg.sender] = beneficiary;
        }

        SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.regAsBackup));
        return true;
    }

    function BackupUpToDate(uint256 index) public {
        require( registerFlag == 2 );
        require( nodesToJoin[index] == msg.sender);
        nodePerformance[msg.sender] = NODE_INIT_PERFORMANCE;
    }

    //user can explicitly release
    function requestRelease(uint256 senderType, uint256 index) public returns (bool) {
        //only in nodeList and scsBeneficiary can call this function
        if (senderType == 1) {
            if (nodeList[index] != msg.sender) {
                return false;
            }
        } else if (senderType == 2) {
            if (scsBeneficiary[nodeList[index]] != msg.sender) {
                return false;
            }
        } else {
            return false;
        }
        //check if full already
        if (nodeToReleaseCount >= 5) {
            return false;
        }

        //check if already requested
        for (uint256 i = 0; i < nodeToReleaseCount; i++) {
            if (nodesToRelease[i] == index) {
                return false;
            }
        }

        nodesToRelease[nodeToReleaseCount] = index;
        nodeToReleaseCount++;

        return true;
    }
    
    //user can explicitly release
    function requestReleaseImmediate(uint256 senderType, uint256 index) public returns (bool) {
        //only in nodeList and scsBeneficiary can call this function
        if (senderType == 1) {
            if (nodeList[index] != msg.sender) {
                return false;
            }
        } else if (senderType == 2) {
            if (scsBeneficiary[nodeList[index]] != msg.sender) {
                return false;
            }
        } else {
            return false;
        }

        if (block.number <= lastFlushBlk + flushInRound + nodeCount * 2 * proposalExpiration) {
            return false;
        }
        

        address cur = nodeList[index];
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        if (!protocnt.releaseFromSubchain(
            cur,
            penaltyBond
        ) ) {
            return false;
        }
        
        nodeCount--;
        nodeList[index] = nodeList[nodeCount];
        delete nodeList[nodeCount];
        nodeList.length--;
        SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.requestRelease));

        return true;
    }

    function registerOpen() public {
        require(subchainstatus == uint256(SubChainStatus.open));
        require(msg.sender == owner);
        registerFlag = 1;

        //call precompiled code to invoke action on v-node
        SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.registerOpen));
    }

    function registerClose() public returns (bool) {
        require(subchainstatus == uint256(SubChainStatus.open));
        require(msg.sender == owner);
        registerFlag = 0;

        if (nodeCount < minMember) {
            SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
            //release already enrolled scs
            //release already enrolled scs
            for (uint256 i = nodeCount; i > 0; i--) {
                //release fund
                address cur = nodeList[i - 1];
                if (protocnt.releaseFromSubchain(
                    cur,
                    penaltyBond
                )) {
                    delete nodeList[i - 1];
                    nodeCount--;
                }
            }
            return false;
        }

        //now we can start to work now
        lastFlushBlk = block.number;
        curFlushIndex = 0;

        //call precompiled code to invoke action on v-node
        SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.registerClose));
        return true;
    }

    function registerAdd(uint256 nodeToAdd) public {
        require(subchainstatus == uint256(SubChainStatus.open));
        require(msg.sender == owner);
        require(joinCntNow + nodeCount < maxMember);
        registerFlag = 2;
        joinCntMax = maxMember - joinCntNow - nodeCount;
        joinCntNow = nodesToJoin.length;
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        selTarget = protocnt.getSelectionTargetByCount(nodeToAdd);

        //call precompiled code to invoke action on v-node
        SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.registerAdd)); // todo David
    }
    
    function getFlushInfo() public view returns (uint256) {
        
        for (uint256 i=1; i <= nodeCount; i++) {
            uint256 blk = lastFlushBlk + flushInRound + i * 2 * proposalExpiration;
            
            if (blk > block.number) {
                return blk - block.number;
            }
        }
        
        return 0;
    }

    function getholdingPool() public constant returns (address[]) {
        
        return holdingPool.userAddr;
    }

    function getEnteringAmount(address userAddr, uint256 holdingPoolPos) public constant returns (address[] enteringAddr, uint256[] enteringAmt, uint256[] enteringtime, uint256[] rechargeParam) {
        uint256 i;
        uint256 j = 0;
        
        if (userAddr != address(0)) {
            for (i = holdingPoolPos; i < holdingPool.userAddr.length; i++) {
                if (holdingPool.userAddr[i] == userAddr) {
                    j++;
                }
            }
        } else {
            j = holdingPool.userAddr.length - holdingPoolPos;
        }
        
        address[] memory addrs = new address[](j);
        uint256[] memory amounts = new uint256[](j);
        uint256[] memory times = new uint256[](j);
        uint256[] memory params = new uint256[](2);
        j = 0;
        for (i = holdingPoolPos; i < holdingPool.userAddr.length; i++) {
            if (userAddr != address(0)) {
                if (holdingPool.userAddr[i] == userAddr) {
                    amounts[j] = holdingPool.amount[i];
                    times[j] = holdingPool.time[i];
                    j++;
                }
            } else {
                addrs[j] = holdingPool.userAddr[i];
                amounts[j] = holdingPool.amount[i];
                times[j] = holdingPool.time[i];
                j++;
            }
        }
        params[0] = per_recharge_num;
        params[1] = recharge_cycle;
        return (addrs, amounts, times, params);
    }
    
    function getRedeemRecords(address userAddr) public view returns (RedeemRecords) {
        
        return records[userAddr];
    }

    //|----------|---------|---------|xxx|yyy|zzz|
    function getEstFlushBlock(uint256 index) public view returns (uint256) {
        uint256 blk = lastFlushBlk + flushInRound;
        //each flusher has [0, 2*expire] to finish
        if (index >= curFlushIndex) {
            blk += (index - curFlushIndex) * 2 * proposalExpiration;
        }
        else {
            blk += (index + nodeCount - curFlushIndex) * 2 * proposalExpiration;
        }
        return blk;
    }


    // create proposal
    // bytes32 hash;
    // bytes newState;
    function createProposal(
        uint256 indexInlist,
        bytes32[] hashlist,
        uint256[] blocknum,
        uint256[] distAmount,
        uint256[] badactors,
        address viaNodeAddress,
        uint256 preRedeemNum
    )
        public
        returns (bool)
    {
        uint256 gasinit = msg.gas; //gasleft();
        require(indexInlist < nodeCount && msg.sender == nodeList[indexInlist]);
        require(block.number >= getEstFlushBlock(indexInlist) && 
                block.number < (getEstFlushBlock(indexInlist)+ 2*proposalExpiration));
        require( distAmount.length == nodeCount);
        require( badactors.length < nodeCount/2);
        require( tx.gasprice <= MAX_GAS_PRICE );
        require( contractNeedFund < totalOperation );

        //if already a hash proposal in progress, check if it is set to expire
        if (
            proposals[proposalHashInProgress].flag == uint256(ProposalFlag.pending)
        ) {
            //for some reason, lastone is not updated
            //set to expire
            proposals[proposalHashInProgress].flag = uint256(ProposalFlag.expired);  //expired.
            //reduce proposer's performance
            if (nodePerformance[proposals[proposalHashInProgress].proposedBy] > 0) {
                nodePerformance[proposals[proposalHashInProgress].proposedBy]--;
            }
        }

        //proposal must based on last approved hash
        if (hashlist[0] != proposalHashApprovedLast) {
            //ReportStatus("Proposal base bad");

            return false;
        }

        //check if sender is part of SCS list
        if (!isSCSValid(msg.sender)) {
            //ReportStatus("Proposal requester invalid");
            return false;
        }

        bytes32 curhash = hashlist[1];
        //check if proposal is already in
        if (proposals[curhash].flag > uint256(ProposalFlag.noState)) {
            //ReportStatus("Proposal in progress");
            return false;
        }

        //store it into storage.
        proposals[curhash].proposedBy = msg.sender;
        proposals[curhash].lastApproved = proposalHashApprovedLast;
        proposals[curhash].hash = curhash;
        proposals[curhash].start = blocknum[0];
        proposals[curhash].end = blocknum[1];
        //proposals[hash].newState = newState;
        uint256 i=0;
        for (i=0; i < nodeCount; i++) {
            proposals[curhash].distributionAmount.push(distAmount[i]);
            proposals[curhash].minerAddr.push(nodeList[i]);
        }
        proposals[curhash].flag = uint256(ProposalFlag.pending);
        proposals[curhash].startingBlock = getEstFlushBlock(indexInlist);
        //add into voter list
        proposals[curhash].voters.push(indexInlist);
        proposals[curhash].votecount++;

        proposals[curhash].redeemAgreeList.push(msg.sender);

        for (i=0; i < badactors.length; i++) {
            proposals[curhash].badActors.push(badactors[i]);
        }

        //set via nodeelse
        proposals[curhash].viaNodeAddress = viaNodeAddress;
        
        // ErcMapping ss;
        proposals[curhash].preRedeemNum = preRedeemNum;
        proposals[curhash].distributeFlag = 0;
        
        //notify v-node
        if (preRedeemNum == 0) {
            SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.createProposal));
        } else {
            SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.uploadRedeemData));
        }

        proposalHashInProgress = curhash;
        pendingFlushIndex = indexInlist;
        currentRefundGas[msg.sender] += (gasinit - msg.gas + 21486 ) * tx.gasprice;
        //ReportStatus("Proposal creates ok");
        
        return true;
    }

    function UploadRedeemData(
        address[] redeemAddr,
        uint256[] redeemAmt
    )
        public
        returns (bool)
    {
        Proposal storage prop = proposals[proposalHashInProgress];
        uint256 gasinit = msg.gas; //gasleft();
        require(msg.sender == prop.proposedBy);
        require( redeemAddr.length + prop.redeemAddr.length <= prop.preRedeemNum);
        require( tx.gasprice <= MAX_GAS_PRICE );

        
        for (uint256 i=0; i < redeemAddr.length; i++) {
            prop.redeemAddr.push(redeemAddr[i]);
            prop.redeemAmt.push(redeemAmt[i]);
        }
        
        //notify v-node
        if (prop.redeemAddr.length == prop.preRedeemNum) {
            SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.createProposal));
        } else {
            SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.uploadRedeemData));
        }

        currentRefundGas[msg.sender] += (gasinit - msg.gas + 21486 ) * tx.gasprice;
        //ReportStatus("Proposal creates ok");
        
        return true;
    }

    //vote on proposal
    function voteOnProposal(uint256 indexInlist, bytes32 hash, bool redeem) public returns (bool) {
        uint256 gasinit = msg.gas;
        Proposal storage prop = proposals[hash];
        
        require(indexInlist < nodeCount && msg.sender == nodeList[indexInlist]);
        require( tx.gasprice <= MAX_GAS_PRICE );
        //check if sender is part of SCS list
        if (!isSCSValid(msg.sender)) {
            //ReportStatus("Voter invalid");
            
            return false;
        }

        //check if proposal is in proper flag state
        if (prop.flag != uint256(ProposalFlag.pending)) {
            //ReportStatus("Voting not ready");
            return false;
        }
        //check if dispute proposal in proper range [0, expire]
        if ((prop.startingBlock + 2*proposalExpiration) < block.number) {
            //ReportStatus("Proposal expired");
            return false;
        }

        //traverse back to make sure not double vote
        for (uint256 i=0; i < prop.votecount; i++) {
            if (prop.voters[i] == indexInlist) {
                //ReportStatus("Voter already voted");
                return false;
            }
        }

        if (redeem) {
            prop.redeemAgreeList.push(msg.sender);
        }

        //add into voter list
        prop.voters.push(indexInlist);
        prop.votecount++;
        
        currentRefundGas[msg.sender] += (gasinit - msg.gas + 21486) * tx.gasprice;
        //ReportStatus("Voter votes ok");
        
        return true;
    }

    function checkProposalStatus(bytes32 hash ) public view returns (uint256) {
        if ((proposals[hash].startingBlock + 2*proposalExpiration) < block.number) {
            //expired
            return uint256(ProposalCheckStatus.expired);
        }

        //if reaches 50% more agreement
        if ((proposals[hash].votecount * 2) > nodeCount) {
            //more than 50% approval
            return uint256(ProposalCheckStatus.approval);
        }

        //undecided
        return uint256(ProposalCheckStatus.undecided);
    }
    
    function revenueDistribution(bytes32 hash ) private returns (bool) {
        Proposal storage prop = proposals[hash];
        address cur;
         //check if contract has enough fund
        uint256 totalamount = 0;
        
        totalamount += viaReward;
        
        for (uint256 i = 0; i < prop.minerAddr.length; i++) {
            cur = prop.minerAddr[i];
            totalamount += currentRefundGas[cur];
            totalamount += prop.distributionAmount[i];
        }
         //if not enough amount, halt proposal
        if( totalamount > totalOperation ) {
            //set global flag
            contractNeedFund = totalamount;
            return false;
        }
        
         //doing actual distribution
        uint256 size;
        address viaAddr = prop.viaNodeAddress;
        assembly {
            size := extcodesize(viaAddr)
        }
        if (size == 0) {
            prop.viaNodeAddress.transfer(viaReward);
            totalOperation -= viaReward;
            TransferAmount(prop.viaNodeAddress, viaReward);
        }
        
        uint256 amts;
        for ( i = 0; i < prop.minerAddr.length; i++) {
            cur = prop.minerAddr[i];
            uint256 targetGas = currentRefundGas[cur];
            currentRefundGas[cur] = 0;
            cur.transfer(targetGas);
            totalOperation -= targetGas;
            TransferAmount(cur, targetGas);
            targetGas = prop.distributionAmount[i];
            scsBeneficiary[cur].transfer(targetGas);
            amts += targetGas;
            totalOperation -= targetGas;
            TransferAmount(scsBeneficiary[cur], targetGas);
        }

        uint256 txNum = (amts - blockReward * (prop.end - prop.start + 1)) / txReward;
        if (txNum <= txNumInFlush) {
            flushInRound += 40;
            if (flushInRound > maxFlushInRound) {
                flushInRound = maxFlushInRound;
            }
        } else {
            flushInRound = flushInRound / 2;
            if (flushInRound < initialFlushInRound) {
                flushInRound = initialFlushInRound;
            }
        }
        return true;
    }
    
    //request proposal approval
    function requestProposalAction(uint256 indexInlist, bytes32 hash) public returns (bool) {
        uint256 gasinit = msg.gas;
        Proposal storage prop = proposals[hash];

        require(indexInlist < nodeCount && msg.sender == nodeList[indexInlist]);
        require(prop.flag == uint256(ProposalFlag.pending));
        require( tx.gasprice <= MAX_GAS_PRICE );

        //check if sender is part of SCS list
        if (!isSCSValid(msg.sender)) {
            //ReportStatus("Requester not permitted");
            return false;
        }

        //make sure the proposal to be approved is the correct proposal in progress
        if (proposalHashInProgress != hash) {
            //ReportStatus("Request incorrect.");
             return false;
        }

        //check if ready to accept
        uint256 chk = checkProposalStatus(hash);
        if (chk == uint256(ProposalCheckStatus.undecided)) {
            //ReportStatus("No agreement");
            return false;
        } 
        else if (chk == uint256(ProposalCheckStatus.expired)) {
            prop.flag = uint256(ProposalFlag.expired);  //expired.
            //reduce proposer's performance
            address by = prop.proposedBy;
            if (nodePerformance[by] > 0) {
                nodePerformance[by]--;
            }
            //ReportStatus("Proposal expired");
            
            return false;
        }

        //punish bad actors
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        uint256 i = 0;
        for (i=0; i<prop.badActors.length; i++) {
            uint256 badguy = prop.badActors[i];
            nodePerformance[nodeList[badguy]] = 0;
        }
        
        //punish nodePerformance is 0
        if (nodesToDispel.length < MAX_DELETE_NUM) {
            uint256 num = MAX_DELETE_NUM - nodesToDispel.length;
            for (i=0; i<nodeCount; i++) {
                if (num == 0) {
                    break;
                }
                if (nodePerformance[nodeList[i]] == 0) {
                    nodesToDispel.push(i);
                    protocnt.forfeitBond(nodeList[i], penaltyBond);
                    totalOperation += penaltyBond;
                    num--;
                }
            }
        }

        //for correct voter, increase performance
        for (i = 0; i < prop.votecount; i++) {
            address vt = nodeList[prop.voters[i]];
            if (nodePerformance[vt] < NODE_INIT_PERFORMANCE) {
                nodePerformance[vt]++;
            }
        }

        //award to distribution list
        //in following request action

        //award via nodes
        //in following request action
         //token redeem
        //token redeem is done in following request action
        

        //remove bad nodes
        applyRemoveNodes(0);

        //remove node to release
        applyRemoveNodes(1);

        //update randIndex
        bytes32 randseed = sha256(hash, block.number);
        randIndex[0] = uint8(randseed[0]) / 8;
        randIndex[1] = uint8(randseed[1]) / 8;

        //if some nodes want to join in
        if (registerFlag == 2) {
            applyJoinNodes();
        }

        //if need toauto retire nodes
        if (AUTO_RETIRE) {
            for (i=0; i<AUTO_RETIRE_COUNT; i++) {
                if (indexAutoRetire >= nodeCount) {
                    indexAutoRetire = 0;
                }
                requestRelease(1, indexAutoRetire);
                indexAutoRetire ++ ;
            }
        }
        
        //notify v-node
        if (prop.redeemAddr.length == 0) {
            if (revenueDistribution(hash)) {
                flushEnd(hash);
            }
        } else {
            SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.approveProposal));
        }
        
        //make protocol pool to know subchain is active
        protocnt.setSubchainActiveBlock();

        //adjust reward
        adjustReward();
        
        //update flag
        prop.distributeFlag = 1;
        //refund current caller
        msg.sender.transfer((gasinit - msg.gas + 15000) * tx.gasprice);
        totalOperation = this.balance - totalBond;
        
        //refund all to owner
        if (subchainstatus == uint256(SubChainStatus.close)) {
            owner.transfer(this.balance - totalBond);
            totalOperation = 0;
        }
        return true;
    }
    
    function flushEnd(bytes32 hash ) private {
        Proposal storage prop = proposals[hash];
        
        //setflag
        prop.distributeFlag = 2;
        //mark as approved
        prop.flag = uint256(ProposalFlag.approved);
        //reset flag
        proposalHashInProgress = 0x0;
        proposalHashApprovedLast = hash;
        lastFlushBlk = block.number;

        curFlushIndex = pendingFlushIndex + 1;
        if (curFlushIndex > nodeCount) {
            curFlushIndex = 0;
        }

        if (subchainstatus == uint256(SubChainStatus.pending)) {
            withdrawal();
        }
        
        SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.distributeProposal));
    }
    
    function requestEnterAndRedeemAction(bytes32 hash) public returns (bool) {
        uint256 gasinit = msg.gas;
        //any one can request
        Proposal storage prop = proposals[hash];
        require(BALANCE != 0);
        require(prop.distributeFlag == 1);
        require(prop.flag == uint256(ProposalFlag.pending));

        uint256 chk = checkProposalStatus(hash);
        if (chk == uint256(ProposalCheckStatus.undecided)) {
            //ReportStatus("No agreement");
            return false;
        }
        else if (chk == uint256(ProposalCheckStatus.expired)) {
            prop.flag = uint256(ProposalFlag.expired);  //expired.
            //reduce proposer's performance
            address by = prop.proposedBy;
            if (nodePerformance[by] > 0) {
                nodePerformance[by]--;
            }
            //ReportStatus("Proposal expired");
            
            return false;
        }
        
        //redeem tokens
        uint256 i;
        bool res = true;
        if (prop.redeemAgreeList.length > nodeCount/2 && prop.preRedeemNum != 0 && ERCAddr != address(0)) {
            uint256 len = prop.preRedeemNum;
            if (len > per_redeemdata_num) {
                len = per_redeemdata_num;
            }
            
            uint256 pos = prop.redeemAddr.length - prop.preRedeemNum;
            address addr;
            uint256 amt;
            Token erc20 = Token(ERCAddr);
            uint256 balance = erc20.balanceOf(address(this));
            for (i = pos; i < pos + len; i++) {
                addr =  prop.redeemAddr[i];
                amt = prop.redeemAmt[i]/ERCRate/(10**ERCDecimals);
                balance = balance - amt;
                if( balance >= 0 ) {
                    if (!erc20.transfer(addr, amt)) {
                        return false;
                    }
                } else {
                    return false;
                }
                records[addr].redeemAmount.push(prop.redeemAmt[i]);
                records[addr].redeemtime.push(now);
                prop.preRedeemNum --;
                dappRedeemPos ++;
            }
            res = false;
        }
        if (res) {
            if (revenueDistribution(hash)) {
                flushEnd(hash);
            }
        } else {
            SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.requestEnterAndRedeem));
        }
        
        msg.sender.transfer((gasinit - msg.gas + 15000) * tx.gasprice);
        totalOperation = this.balance - totalBond;
        
        if (subchainstatus == uint256(SubChainStatus.close)) {
            owner.transfer(this.balance - totalBond);
            totalOperation = 0;
        }
        return true;
    }
    
    function adjustReward() private {
        blockReward = blockReward - blockReward * DEFLATOR_VALUE / 10 ** 6;    
        txReward = txReward - txReward * DEFLATOR_VALUE / 10 ** 6;    
        viaReward = viaReward - viaReward * DEFLATOR_VALUE / 10 ** 6;    
        syncReward = syncReward - syncReward * DEFLATOR_VALUE / 10 ** 6;    
    }

    //to increase reward if deflator is too much
    function increaseReward(uint256 percent) private {
        require(owner == msg.sender);
        blockReward = blockReward + blockReward * percent / 100;    
        txReward = txReward - txReward * percent / 100;    
        viaReward = viaReward - viaReward * percent / 100;    
        syncReward = syncReward - syncReward * percent / 100;    
    }

    function addFund() public payable {
        // do nothing
        //ReportStatus("fund added" );
        require(owner == msg.sender);
        totalOperation += msg.value;
        if( totalOperation > contractNeedFund ) {
            contractNeedFund = 0;
            uint256 blk = lastFlushBlk + flushInRound + nodeCount * 2 * proposalExpiration;
            
            if (block.number >= blk) {
                lastFlushBlk = block.number;
                SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.updateLastFlushBlk));
            }
        }
    }

    function withdraw(address recv, uint256 amount) public {
        require(owner == msg.sender);
        require(totalOperation >= amount);
        //withdraw to address
        recv.transfer(amount);
        totalOperation -= amount;
    }

    function withdrawal() private {
        subchainstatus = uint256(SubChainStatus.close);
        registerFlag = 0;
        //release fund
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        //release already enrolled scs
        for (uint256 i = nodeCount; i > 0; i--) {
            //release fund
            address cur = nodeList[i-1];
            if (protocnt.releaseFromSubchain(
                cur,
                penaltyBond
            )) {
                delete nodeList[i - 1];
                nodeCount--;
            }
        }

        for (i = joinCntNow; i > 0; i--) {
            cur = nodesToJoin[i-1];
            if (protocnt.releaseFromSubchain(
                cur,
                penaltyBond
            )) {
                delete nodesToJoin[i-1];
                nodesToJoin.length --;
                joinCntNow--;
            }
        }
        joinCntMax = 0;

        //kill self
    }

    function close() public {
        require(owner == msg.sender);

        subchainstatus = uint256(SubChainStatus.pending);
        
        if (proposalHashInProgress == 0x0) {
            lastFlushBlk = block.number - flushInRound;
            SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.updateLastFlushBlk));
        }
    }

    function reset() public {
        require(owner == msg.sender);
        //Comment out the logic that relies on the flush state when function reset executes
        // uint256 blk = lastFlushBlk + flushInRound + nodeCount * 2 * proposalExpiration;
        // if (block.number >= blk) {
            lastFlushBlk = block.number;
            flushInRound = 60;
            SCS_RELAY.notifySCS(address(this), uint256(SCSRelayStatus.reset));
        // }
    }

    function addSyncNode(address id, string link) public {
        require(owner == msg.sender);
        syncNodes.push(SyncNode(id, link));
    }

    function removeSyncNode(uint256 index) public {
        require(owner == msg.sender && syncNodes.length > index);
        syncNodes[index] = syncNodes[syncNodes.length - 1];
        delete syncNodes[syncNodes.length - 1];
        syncNodes.length--;
    }

    function isSCSValid(address addr) private view returns (bool) {
        if (!isMemberValid(addr)) {
            return false;
        }

        //check if valid registered in protocol pool
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        if (!protocnt.isPerforming(addr)) {
            return false;
        }
        return true;
    }

    function applyJoinNodes() private {
        uint256 i = 0;
        for (i = joinCntNow; i > 0; i--) {
            if( nodePerformance[nodesToJoin[i-1]] == NODE_INIT_PERFORMANCE) {
                nodeList.push(nodesToJoin[i-1]);
                nodeCount++;

                //delete node
                nodesToJoin[i-1] = nodesToJoin[nodesToJoin.length-1];
                delete nodesToJoin[nodesToJoin.length-1];
                nodesToJoin.length --;
            }
        }

        joinCntNow = nodesToJoin.length;
        if( joinCntNow == 0 ) {
            joinCntMax = 0;
            registerFlag = 0;
        }
    }

    // reuse this code for remove bad node or other volunteerly leaving node
    // nodetype 0: bad node, 1: volunteer leaving node
    function applyRemoveNodes(uint256 nodetype) private {
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);

        uint256 count = nodesToDispel.length;
        if (nodetype == 1) {
            count = nodeToReleaseCount;
        }

        if (count == 0) {
            return;
        }

        // all nodes set 0 at initial, set node to be removed as 1.
        uint256[] memory nodeMark = new uint256[](nodeCount);
        uint256 idx = 0;
        uint256 i = 0;
        for (i = 0; i < count; i++) {
            if (nodetype == 0) {
                //bad ones
                nodeMark[nodesToDispel[i]] = 1;
            }
            else {
                idx = nodesToRelease[i];
                //volunteer leaving, only were not marked as bad ones
                if (nodeMark[idx] == 0) {
                    nodeMark[idx] = 1;
                    //release fund
                    address cur = nodeList[idx];
                    protocnt.releaseFromSubchain(
                        cur,
                        penaltyBond
                    );
                }
            }
        }

        //adjust to update nodeList
        for (i = nodeCount; i > 0; i--) {
            if (nodeMark[i-1] == 1) {
                //swap with last element
                // remove node from list
                nodeCount--;
                nodeList[i-1] = nodeList[nodeCount];
                delete nodeList[nodeCount];
                nodeList.length--;
                //nodesToDispel.length--;
            }

            // if (i == 0) {
            //     break;
            // }
            // else {
            //     i--;
            // }
        }

        //clear nodesToDispel and nodesToRelease array
        if (nodetype == 0) {
            //clear bad ones
            nodesToDispel.length = 0 ;
        } else {
            //clear release count
            nodeToReleaseCount = 0;
        }
    }
    
    // ATO new
    function rebuildFromLastFlushPoint() public {
        require(msg.sender == owner);
        //notifyscs
        //set flushindex
        curFlushIndex = 0;
    }


    function getindexByte(address a, uint8 randIndex1, uint8 randIndex2) private  pure returns (uint256 b) {
        uint8 first = uint8(uint256(a) / (2 ** (4 * (39 - uint256(randIndex1)))) * 2 ** 4);
        uint8 second = uint8(uint8(uint256(a) / (2 ** (4 * (39 - uint256(randIndex2)))) * 2 ** 4) / 2 ** 4);    // &15
        return uint256(byte(first + second));
    }
    
    function matchSelTarget(address addr, uint8 index1, uint8 index2) public view returns (bool) {
        // check if selTargetdist matches.
        uint256 addr0 = getindexByte(addr, index1, index2);
        uint256 cont0 = getindexByte(address(this), index1, index2);

        if (selTarget == 255) {
            return true;
        }

        if (addr0 >= cont0) {
            if ((addr0 - cont0) <= selTarget) {
                return true;
            }
            else {
                if (cont0 - selTarget < 0) {
                    if ((addr0 - cont0) >= 256 - selTarget) {
                        //lower half round to top,  addr0 -256 >= cont0 -selTarget
                        return true;
                    }
                    return false;
                }
                return false;
            }
        }
        else {
            //addr0 < cont0
            if ((cont0 - addr0) <= selTarget) {
                return true;
            }
            else {
                if (cont0 + selTarget >= 256) {
                    if ((cont0 - addr0) >= 256 - selTarget) {
                        //top half round to bottom,   addr0 +256  <= (selTarget+cont0)
                        return true;
                    }
                    return false;
                }
                return false;
            }
        }

        return true;
    }
    
    function buyMintToken(uint256 _value) public returns (bool) {
        if (ERCAddr != address(0)) {
            Token erc20 = Token(ERCAddr);
            uint256 allowedamt = erc20.allowance(msg.sender, address(this));
            if (_value <= allowedamt) {
                if (erc20.transferFrom(msg.sender, address(this), _value)) {
                    return requestEnterMicrochain(msg.sender, _value*ERCRate*(10**ERCDecimals));
                }
            }
        }
        return false;
    }
    
    function buyMintTokenTransfer(address _to, uint256 _value) public returns (bool) {
        if (ERCAddr != address(0)) {
            Token erc20 = Token(ERCAddr);
            uint256 allowedamt = erc20.allowance(msg.sender, address(this));
            if (_value <= allowedamt) {
                if (erc20.transferFrom(msg.sender, address(this), _value)) {
                    return requestEnterMicrochain(_to, _value*ERCRate*(10**ERCDecimals));
                }
            }
        }
        return false;
    }
    
    function requestEnterMicrochain(address _to, uint256 amount) private returns (bool){
        if( ERCAddr != address(0)) {
            holdingPool.userAddr.push(_to);
            holdingPool.amount.push(amount);
            holdingPool.time.push(now);

            return true;
        }               
        return false;
    }

    function settlement() public {
        require(foundation == msg.sender || owner == msg.sender);
        uint256 blk = lastFlushBlk + flushInRound + nodeCount * 2 * proposalExpiration + settlement_block_delay;
        if (block.number > blk) {
            Token erc20 = Token(ERCAddr);
            uint256 balance = erc20.balanceOf(address(this));
            erc20.transfer(foundation, balance);
        }
        
    }

    function updateFoundation(address _addr) public {
        require(foundation == msg.sender || owner == msg.sender);
        foundation = _addr;
    }
}