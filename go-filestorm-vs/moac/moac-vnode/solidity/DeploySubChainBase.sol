pragma solidity ^0.4.11;
//David Chen
//Subchain definition for application.


contract SysContract {
    function delayedSend(uint _blk, address _to, uint256 _value, bool bonded) public returns (bool success);
}


contract SubChainProtocolBase {
    enum SCSStatus { notRegistered, performing, withdrawPending, initialPending, withdrawDone, badActor }

    struct SCS {
        address from; //address as id
        uint256 bond;   // value
        uint state; // one of SCSStatus
        uint256 registerBlock;
        uint256 withdrawBlock;
    }

    struct SCSApproval {
        uint bondApproved;
        uint bondedCount;
        address[] subchainAddr;
        uint[] amount;
    }

    mapping(address => SCS) public scsList;
    mapping(address => SCSApproval) public scsApprovalList;

    uint public scsCount;
    string public subChainProtocol;
    uint public bondMin;
    uint public constant PENDING_BLOCK_DELAY = 5; // 8 minutes
    uint public constant WITHDRAW_BLOCK_DELAY = 8640; // one day, given 10s block rate
    SysContract internal constant SYS_CONTRACT = SysContract(0x0000000000000000000000000000000000000065);

    //monitor if subchain is inactive
    //this is used to allow node to exit from zoombie subchain
    mapping(address => uint) public subChainLastActiveBlock;
    mapping(address => uint) public subChainExpireBlock;

    //events
    event Registered(address scs);
    event UnRegistered(address sender);

    address[] public scsArray;
    uint public protocolType;

    //constructor
    function SubChainProtocolBase(string protocol, uint bmin, uint _protocolType) public {
        scsCount = 0;
        subChainProtocol = protocol;
        bondMin = bmin;
        protocolType = _protocolType;
    }

    function() public payable {  // todo: david review
        revert();
    }

    // register for SCS
    // SCS will be notified through 3rd party communication method. SCS will need to register here manually.
    // One protocol base can have several different subchains.
    function register(address scs) public payable returns (bool) {
        //already registered or not enough bond
        require(
            (scsList[scs].state == uint(SCSStatus.notRegistered)
            || scsList[scs].state == uint(SCSStatus.performing))
            && msg.value >= bondMin * 10 ** 18
        );

        addScsId(scs);

        scsList[scs].from = scs;
        if (scsList[scs].state == uint(SCSStatus.notRegistered)) {
            //if not register before, update
            scsList[scs].registerBlock = block.number + PENDING_BLOCK_DELAY;
            scsList[scs].withdrawBlock = 2 ** 256 - 1;
            scsCount++;
            scsList[scs].bond = msg.value;
        } else {
            //add more fund
            scsList[scs].bond += msg.value;            
        }
        scsList[scs].state = uint(SCSStatus.performing);
        return true;
    }

    // withdrawRequest for SCS
    function withdrawRequest() public returns (bool success) {
        //only can withdraw when active
        require(scsList[msg.sender].state == uint(SCSStatus.performing));

        //need to make sure node is not working for any suchain anymore
        require(scsApprovalList[msg.sender].bondedCount == 0 );        

        scsList[msg.sender].withdrawBlock = block.number;
        scsList[msg.sender].state = uint(SCSStatus.withdrawPending);
        scsCount--;

        removeScsId(msg.sender);

        UnRegistered(msg.sender);
        return true;
    }

    function withdraw() public {
        if (
            scsList[msg.sender].state == uint(SCSStatus.withdrawPending)
            && block.number > (scsList[msg.sender].withdrawBlock + WITHDRAW_BLOCK_DELAY)
        ) {
            scsList[msg.sender].state == uint(SCSStatus.withdrawDone);
            scsList[msg.sender].from.transfer(scsList[msg.sender].bond);
        }
    }

    function isPerforming(address _addr) public view returns (bool res) {
        return (scsList[_addr].state == uint(SCSStatus.performing) && scsList[_addr].registerBlock < block.number);
    }

    function getSelectionTarget(uint thousandth, uint minnum) public view returns (uint target) {
        // find a target to choose thousandth/1000 of total scs
        if (minnum < 50) {
            minnum = 50;
        }

        if (scsCount < minnum) {          // or use scsCount* thousandth / 1000 + 1 < minnum
            return 255;
        }

        uint m = thousandth * scsCount / 1000;

        if (m < minnum) {
            m = minnum;
        }

        target = (m * 256 / scsCount + 1) / 2;

        return target;
    }

    function getSelectionTargetByCount(uint targetnum) public view returns (uint target) {
        if (scsCount <= targetnum) {        
            return 255;
        }

        //calculate distance
        target = (targetnum * 256 / scsCount + 1) / 2;

        if (target == 0 ) {
            target = 0;
        }

        return target;
    }


    //display approved scs list
    function approvalAddresses(address addr) public view returns (address[]) {
        address[] memory res = new address[](scsApprovalList[addr].bondedCount);
        for (uint i = 0; i < scsApprovalList[addr].bondedCount; i++) {
            res[i] = (scsApprovalList[addr].subchainAddr[i]);
        }
        return res;
    }

    //display approved amount array
    function approvalAmounts(address addr) public view returns (uint[]) {
        uint[] memory res = new uint[](scsApprovalList[addr].bondedCount);
        for (uint i = 0; i < scsApprovalList[addr].bondedCount; i++) {
            res[i] = (scsApprovalList[addr].amount[i]);
        }
        return res;
    }

    //subchain need to set this before allow nodes to join
    function setSubchainExpireBlock(uint blk) public {
        subChainExpireBlock[msg.sender] = blk;
    }

    //set active block
    function setSubchainActiveBlock() public {
        subChainLastActiveBlock[msg.sender] = block.number;
    }

    //approve the bond to be deduced if act maliciously
    function approveBond(address scs, uint amount, uint8 v, bytes32 r, bytes32 s) public returns (bool) {
        //require subchain is active
        //require( (subChainLastActiveBlock[msg.sender] + subChainExpireBlock[msg.sender])  > block.number);

        //make sure SCS is performing
        if (!isPerforming(scs)) {
            return false;
        }

        //verify signature
        //combine scs and subchain address
        bytes32 hash = sha256(scs, msg.sender);

        //verify signature matches.
        if (ecrecover(hash, v, r, s) != scs) {
            return false;
        }

        //check if bond still available for SCSApproval
        if (scsList[scs].bond < (scsApprovalList[scs].bondApproved + amount)) {
            return false;
        }

        //add subchain info
        scsApprovalList[scs].bondApproved += amount;
        scsApprovalList[scs].subchainAddr.push(msg.sender);
        scsApprovalList[scs].amount.push(amount);
        scsApprovalList[scs].bondedCount++;

        return true;
    }

    //must called from SubChainBase
    function forfeitBond(address scs, uint amount) public payable returns (bool) {
        //require( (subChainLastActiveBlock[msg.sender] + subChainExpireBlock[msg.sender])  > block.number);
        
        //check if subchain is approved
        for (uint i = 0; i < scsApprovalList[scs].bondedCount; i++) {
            if (scsApprovalList[scs].subchainAddr[i] == msg.sender && scsApprovalList[scs].amount[i] == amount) {
                //delete array item by moving the last item in current postion and delete the last one
                scsApprovalList[scs].bondApproved -= amount;
                scsApprovalList[scs].bondedCount--;
                scsApprovalList[scs].subchainAddr[i]
                    = scsApprovalList[scs].subchainAddr[scsApprovalList[scs].bondedCount];
                scsApprovalList[scs].amount[i] = scsApprovalList[scs].amount[scsApprovalList[scs].bondedCount];

                delete scsApprovalList[scs].subchainAddr[scsApprovalList[scs].bondedCount];
                delete scsApprovalList[scs].amount[scsApprovalList[scs].bondedCount];
                scsApprovalList[scs].subchainAddr.length--;
                scsApprovalList[scs].amount.length--;

                //doing the deduction
                scsList[scs].bond -= amount;
                //scsList[scs].state = uint(SCSStatus.badActor);
                msg.sender.transfer(amount);

                return true;
            }
        }

        return false;
    }


    //scs to request to release from a subchain if subchain is not active
    //anyone can request this
    function releaseRequest(address scs, address subchain) public returns (bool) {
        //check subchain info
        for (uint i=0; i < scsApprovalList[scs].bondedCount; i++) {
            if (scsApprovalList[scs].subchainAddr[i] == subchain && 
                (subChainLastActiveBlock[subchain] + subChainExpireBlock[subchain])  < block.number) {
                scsApprovalList[scs].bondApproved -= scsApprovalList[scs].amount[i];
                scsApprovalList[scs].bondedCount--;
                scsApprovalList[scs].subchainAddr[i]
                    = scsApprovalList[scs].subchainAddr[scsApprovalList[scs].bondedCount];
                scsApprovalList[scs].amount[i] = scsApprovalList[scs].amount[scsApprovalList[scs].bondedCount];

                //clear
                delete scsApprovalList[scs].subchainAddr[scsApprovalList[scs].bondedCount];
                delete scsApprovalList[scs].amount[scsApprovalList[scs].bondedCount];
                scsApprovalList[scs].subchainAddr.length--;
                scsApprovalList[scs].amount.length--;

                //DAVID: not send back bond. It only happens in withdraw request. 
                //just make node out of subchain
                return true;
            }
        }
        return false;
    }

    //subchain to request to release a scs from a subchain
    function releaseFromSubchain(address scs, uint amount) public returns (bool) {
        //check subchain info
        for (uint i=0; i < scsApprovalList[scs].bondedCount; i++) {
            if (scsApprovalList[scs].subchainAddr[i] == msg.sender && scsApprovalList[scs].amount[i] == amount) {
                scsApprovalList[scs].bondApproved -= amount;
                scsApprovalList[scs].bondedCount--;
                scsApprovalList[scs].subchainAddr[i]
                    = scsApprovalList[scs].subchainAddr[scsApprovalList[scs].bondedCount];
                scsApprovalList[scs].amount[i] = scsApprovalList[scs].amount[scsApprovalList[scs].bondedCount];

                //clear
                delete scsApprovalList[scs].subchainAddr[scsApprovalList[scs].bondedCount];
                delete scsApprovalList[scs].amount[scsApprovalList[scs].bondedCount];
                scsApprovalList[scs].subchainAddr.length--;
                scsApprovalList[scs].amount.length--;

                //DAVID: not send back bond. It only happens in withdraw request. 
                //just make node out of subchain
                return true;
            }
        }
        return false;
    }

    function addScsId(address scsId) private {
        if (scsList[scsId].state == uint(SCSStatus.notRegistered)) {
            scsArray.push(scsId);
        }
    }

    function removeScsId(address scsId) private {
        uint len = scsArray.length;
        for (uint i=0; i<len; i++) {
            if (scsArray[i] ==  scsId) {
                delete scsArray[i];
            }
        }
    }
}



contract SCSRelay {
    // 0-registeropen, 1-registerclose, 2-createproposal, 3-disputeproposal, 4-approveproposal, 5-registeradd
    function notifySCS(address cnt, uint msgtype) public returns (bool success);
}


contract SubChainBase {
    enum ProposalFlag {noState, pending, disputed, approved, rejected, expired, pendingAccept}
    enum ProposalCheckStatus {undecided, approval, expired}
    enum ConsensusStatus {initStage, workingStage, failure}
    enum SCSRelayStatus {registerOpen, registerClose, createProposal, disputeProposal, approveProposal, registerAdd, regAsMonitor, regAsBackup, updateLastFlushBlk}
    enum SubChainStatus {open, pending, close}

    struct Proposal {
        address proposedBy;
        bytes32 lastApproved;
        bytes32 hash;
        uint start;
        uint end;
        //bytes newState;
        uint[] distributionAmount;
        uint flag; // one of ProposalFlag
        uint startingBlock;
        uint[] voters; //voters index
        uint votecount;
        uint[] badActors;
        address[] viaNodeAddress;
        uint[] viaNodeAmount;
    }

    struct VRS {
        bytes32 r;
        bytes32 s;
        uint8 v;
    }

    struct SyncNode {
        address nodeId;
        string link;
    }

    address public protocol;
    uint public minMember;
    uint public maxMember;
    uint public selTarget;
    uint public consensusFlag; // 0: init stage 1: working stage 2: failure
    uint public flushInRound;
    bytes32 public proposalHashInProgress;
    bytes32 public proposalHashApprovedLast;  //index: 7
    uint internal curFlushIndex;
    uint internal pendingFlushIndex;

    bytes public funcCode;
    bytes internal state;

    uint internal lastFlushBlk;

    address internal owner;

    //nodes list is updated at each successful flush
    uint public nodeCount;
    address[] public nodeList;    //index: 0f

    uint8[2] public randIndex;
    mapping(address => uint ) public nodePerformance;
    mapping(bytes32 => Proposal) public proposals;  //index: 12
    mapping(address => uint) public currentRefundGas;

    uint internal registerFlag;

    uint public proposalExpiration = 12;
    uint public penaltyBond = 10 ** 18; // 1 Moac penalty
    mapping(address=>address) public scsBeneficiary;
    uint public blockReward = 5 * 10 ** 14;    //index: 18
    uint public txReward  = 1 * 10 ** 11;
    uint public viaReward = 1 * 10 ** 13;

    uint public nodeToReleaseCount;
    uint[5] public nodesToRelease;  //nodes wish to withdraw, only allow 5 to release at a time
    mapping(address=>VRS) internal nodesToReleaseVRS;
    uint[] public nodesToDispel;

    address[] public nodesToJoin;  //nodes to be joined
    uint public joinCntMax;
    uint public joinCntNow;
    uint public MONITOR_JOIN_FEE = 1 * 10 ** 16;
    mapping(address=>uint) public nodesWatching;  //nodes watching

    SyncNode[] public syncNodes;
    uint indexAutoRetire;

    uint constant VIANODECNT = 100;
    SCSRelay internal constant SCS_RELAY = SCSRelay(0x000000000000000000000000000000000000000d);
    uint public constant NODE_INIT_PERFORMANCE = 5;
    uint public constant AUTO_RETIRE_COUNT = 2;
    bool public constant AUTO_RETIRE = false;
    address public VnodeProtocolBaseAddr;
    uint public MONITOR_MIN_FEE = 1 * 10 ** 12;
    uint public syncReward = 1 * 10 ** 11;
    uint public MAX_GAS_PRICE = 20 * 10 ** 9;

    uint public DEFLATOR_VALUE = 80; // in 1/millionth: in a year, exp appreciation is 12x
    uint internal subchainstatus;

    //events
    event ReportStatus(string message);
    event TransferAmount(address addr, uint amount);



    //constructor
    function SubChainBase(address proto, address vnodeProtocolBaseAddr, uint min, uint max, uint thousandth, uint flushRound) public {
        VnodeProtocolBaseAddr = vnodeProtocolBaseAddr;
        SubChainProtocolBase protocnt = SubChainProtocolBase(proto);
        selTarget = protocnt.getSelectionTarget(thousandth, min);
        protocnt.setSubchainExpireBlock(flushInRound*5);
        protocnt.setSubchainActiveBlock();

        minMember = min;
        maxMember = max;
        protocol = proto; //address
        consensusFlag = uint(ConsensusStatus.initStage);
        owner = msg.sender;

        flushInRound = flushRound;
        if (flushInRound <= 100) {
            flushInRound = 100;
        }
        lastFlushBlk = 2 ** 256 - 1;

        randIndex[0] = uint8(0);
        randIndex[1] = uint8(1);
        indexAutoRetire = 0;
        subchainstatus = uint(SubChainStatus.open);
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

    function isMemberValid(address addr) public view returns (bool) {
        return nodePerformance[addr] > 0;
    }

    function getSCSRole(address scs) public view returns (uint) {
        uint i = 0;

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
        
        if (matchSelTarget(scs, randIndex[0], randIndex[1])) {
            //ReportStatus("SCS not selected");
            SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
            if (!protocnt.isPerforming(scs)) {
                return 0;
            }
            return 4;
        }
        
        return 0;
    }

    function registerAsMonitor(address monitor) public payable { 
        require(msg.value >= MONITOR_MIN_FEE);
        require(nodesWatching[monitor] == 0); 
        require(monitor != address(0));
        nodesWatching[monitor] = msg.value;
        SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.regAsMonitor));
    }

    //v,r,s are the signature of msg hash(scsaddress+subchainAddr)
    // function registerOwnerSCS(address scs, address beneficiary) public returns (bool) {
    //     require(msg.sender == owner);
    //     require(scs != address(0));
    //     require(beneficiary != address(0));

    //     nodeList.push(scs);
    //     nodeCount++;
    //     nodePerformance[scs] = NODE_INIT_PERFORMANCE;
    //     scsBeneficiary[scs] = beneficiary;
    // }


    //v,r,s are the signature of msg hash(scsaddress+subchainAddr)
    function registerAsSCS(address beneficiary, uint8 v, bytes32 r, bytes32 s) public returns (bool) {
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
        for (uint i=0; i < nodeCount; i++) {
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

        uint i = 0;
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

        SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.regAsBackup));
        return true;
    }

    function BackupUpToDate(uint index) public {
        require( registerFlag == 2 );
        require( nodesToJoin[index] == msg.sender);
        nodePerformance[msg.sender] = NODE_INIT_PERFORMANCE;
    }

    //user can explicitly release
    function requestRelease(uint index) public returns (bool) {
        //only in nodeList can call this function
        require(nodeList[index] == msg.sender);
        //check if full already
        if (nodeToReleaseCount >= 5) {
            return false;
        }

        //check if already requested
        for (uint i = 0; i < nodeToReleaseCount; i++) {
            if (nodesToRelease[i] == index) {
                return false;
            }
        }

        nodesToRelease[nodeToReleaseCount] = index;
        nodeToReleaseCount++;

        return true;
    }

    function registerOpen() public {
        require(msg.sender == owner);
        registerFlag = 1;

        //call precompiled code to invoke action on v-node
        SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.registerOpen));
    }

    function registerClose() public returns (bool) {
        require(msg.sender == owner);
        registerFlag = 0;

        if (nodeCount < minMember) {
            SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
            //release already enrolled scs
            //release already enrolled scs
            for (uint i = nodeCount; i > 0; i--) {
                //release fund
                address cur = nodeList[i - 1];
                protocnt.releaseFromSubchain(
                    cur,
                    penaltyBond
                );

                delete nodeList[i - 1];
            }

            nodeCount = 0;

            return false;
        }

        //now we can start to work now
        lastFlushBlk = block.number;
        curFlushIndex = 0;

        //call precompiled code to invoke action on v-node
        SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.registerClose));
        return true;
    }

    function registerAdd(uint nodeToAdd) public {
        require(msg.sender == owner);
        registerFlag = 2;
        joinCntMax = nodeToAdd;
        joinCntNow = nodesToJoin.length;
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        selTarget = protocnt.getSelectionTargetByCount(nodeToAdd);

        //call precompiled code to invoke action on v-node
        SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.registerAdd)); // todo David
    }

    //|----------|---------|---------|xxx|yyy|zzz|
    function getEstFlushBlock(uint index) public view returns (uint) {
        uint blk = lastFlushBlk + flushInRound;
        //each flusher has [0, 2*expire] to finish
        if (index >= curFlushIndex) {
            blk += (index - curFlushIndex) * 2 * proposalExpiration;
        }
        else {
            blk += (index + nodeCount - curFlushIndex) * 2 * proposalExpiration;
        }

        // if (curblk > (blk+2 * proposalExpiration)) {
        //     uint jump = (curblk-blk)/(2 * proposalExpiration * nodeCount);
        //     if ((curblk-blk) > (2 * proposalExpiration * nodeCount * jump + proposalExpiration)) {
        //         blk = blk + (jump + 1) * (2 * proposalExpiration * nodeCount);
        //     } else {
        //         blk = blk + jump * (2 * proposalExpiration * nodeCount);
        //     }
        // }
        return blk;
    }


    // create proposal
    // bytes32 hash;
    // bytes newState;
    function createProposal(
        uint indexInlist,
        bytes32 lastFlushedHash,
        bytes32 hash,
        uint start, 
        uint end,
        uint[] distAmount,
        uint[] badactors,
        address[] viaNodeAddress,
        uint[] viaNodeAmount
        
    )
        public
        returns (bool)
    {
        uint gasinit = msg.gas; //gasleft();
        require(indexInlist < nodeCount && msg.sender == nodeList[indexInlist]);
        require(block.number >= getEstFlushBlock(indexInlist) && 
                block.number < (getEstFlushBlock(indexInlist)+ 2*proposalExpiration));
        require( viaNodeAddress.length <= VIANODECNT);
        require( viaNodeAddress.length == viaNodeAmount.length);
        require( distAmount.length == nodeCount);
        require( badactors.length < nodeCount/2);
        require( tx.gasprice <= MAX_GAS_PRICE );

        //if already a hash proposal in progress, check if it is set to expire
        if (
            proposals[proposalHashInProgress].flag == uint(ProposalFlag.pending)
        ) {
            //for some reason, lastone is not updated
            //set to expire
            proposals[proposalHashInProgress].flag = uint(ProposalFlag.expired);  //expired.
            //reduce proposer's performance
            if (nodePerformance[proposals[proposalHashInProgress].proposedBy] > 0) {
                nodePerformance[proposals[proposalHashInProgress].proposedBy]--;
            }
        }

        //proposal must based on last approved hash
        if (lastFlushedHash != proposalHashApprovedLast) {
            //ReportStatus("Proposal base bad");

            return false;
        }

        //check if sender is part of SCS list
        if (!isSCSValid(msg.sender)) {
            //ReportStatus("Proposal requester invalid");
            return false;
        }

        //check if proposal is already in
        if (proposals[hash].flag > uint(ProposalFlag.noState)) {
            //ReportStatus("Proposal in progress");
            return false;
        }

        //store it into storage.
        proposals[hash].proposedBy = msg.sender;
        proposals[hash].lastApproved = proposalHashApprovedLast;
        proposals[hash].hash = hash;
        proposals[hash].start = start;
        proposals[hash].end = end;
        //proposals[hash].newState = newState;
        for (uint i=0; i < nodeCount; i++) {
            proposals[hash].distributionAmount.push(distAmount[i]);
        }
        proposals[hash].flag = uint(ProposalFlag.pending);
        proposals[hash].startingBlock = block.number;
        //add into voter list
        proposals[hash].voters.push(indexInlist);
        proposals[hash].votecount++;

        for (i=0; i < badactors.length; i++) {
            proposals[hash].badActors.push(badactors[i]);
        }

        //set via node
        for (i=0; i < viaNodeAddress.length; i++) {
            proposals[hash].viaNodeAddress.push(viaNodeAddress[i]);
            proposals[hash].viaNodeAmount.push(viaNodeAmount[i]);
        }
        

        //notify v-node
        SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.createProposal));

        proposalHashInProgress = hash;
        pendingFlushIndex = indexInlist;
        currentRefundGas[msg.sender] += (gasinit - msg.gas + 21486 ) * tx.gasprice;
        //ReportStatus("Proposal creates ok");
        
        return true;
    }


    //vote on proposal
    function voteOnProposal(uint indexInlist, bytes32 hash) public returns (bool) {
        uint gasinit = msg.gas;
        require(indexInlist < nodeCount && msg.sender == nodeList[indexInlist]);
        require( tx.gasprice <= MAX_GAS_PRICE );
        //check if sender is part of SCS list
        if (!isSCSValid(msg.sender)) {
            //ReportStatus("Voter invalid");
            
            return false;
        }

        //check if proposal is in proper flag state
        if (proposals[hash].flag != uint(ProposalFlag.pending)) {
            //ReportStatus("Voting not ready");
            return false;
        }
        //check if dispute proposal in proper range [0, expire]
        if ((proposals[hash].startingBlock + 2*proposalExpiration) < block.number) {
            //ReportStatus("Proposal expired");
            return false;
        }

        //traverse back to make sure not double vote
        for (uint i=0; i < proposals[hash].votecount; i++) {
            if (proposals[hash].voters[i] == indexInlist) {
                //ReportStatus("Voter already voted");
                return false;
            }
        }

        //add into voter list
        proposals[hash].voters.push(indexInlist);
        proposals[hash].votecount++;

        currentRefundGas[msg.sender] += (gasinit - msg.gas + 21486) * tx.gasprice;
        //ReportStatus("Voter votes ok");
        
        return true;
    }

    function checkProposalStatus(bytes32 hash ) public view returns (uint) {
        if ((proposals[hash].startingBlock + 2*proposalExpiration) < block.number) {
            //expired
            return uint(ProposalCheckStatus.expired);
        }

        //if reaches 50% more agreement
        if ((proposals[hash].votecount * 2) > nodeCount) {
            //more than 50% approval
            return uint(ProposalCheckStatus.approval);
        }

        //undecided
        return uint(ProposalCheckStatus.undecided);
    }

    //request proposal approval
    function requestProposalAction(uint indexInlist, bytes32 hash) public payable returns (bool) {
        uint gasinit = msg.gas;
        require(indexInlist < nodeCount && msg.sender == nodeList[indexInlist]);
        require(proposals[hash].flag == uint(ProposalFlag.pending));
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
        uint chk = checkProposalStatus(hash);
        if (chk == uint(ProposalCheckStatus.undecided)) {
            //ReportStatus("No agreement");
            return false;
        } 
        else if (chk == uint(ProposalCheckStatus.expired)) {
            proposals[hash].flag = uint(ProposalFlag.expired);  //expired.
            //reduce proposer's performance
            address by = proposals[hash].proposedBy;
            if (nodePerformance[by] > 0) {
                nodePerformance[by]--;
            }
            //ReportStatus("Proposal expired");
            
            return false;
        }


        //mark as approved
        proposals[hash].flag = uint(ProposalFlag.approved);
        //reset flag
        proposalHashInProgress = 0x0;
        proposalHashApprovedLast = hash;
        lastFlushBlk = block.number;

        //punish bad actors
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        uint i = 0;
        for (i=0; i<proposals[hash].badActors.length; i++) {
            uint badguy = proposals[hash].badActors[i];
            protocnt.forfeitBond(nodeList[badguy], penaltyBond);
            nodePerformance[nodeList[badguy]] = 0;
            nodesToDispel.push(badguy);
        }

        //for correct voter, increase performance
        for (i = 0; i < proposals[hash].votecount; i++) {
            address vt = nodeList[proposals[hash].voters[i]];
            if (nodePerformance[vt] < NODE_INIT_PERFORMANCE) {
                nodePerformance[vt]++;
            }
        }

        //award to distribution list
        for (i = 0; i < nodeCount; i++) {
            uint targetGas = currentRefundGas[nodeList[i]];
            currentRefundGas[nodeList[i]] = 0;
            nodeList[i].transfer(targetGas);
            TransferAmount(nodeList[i], targetGas);
            targetGas = proposals[hash].distributionAmount[i];
            scsBeneficiary[nodeList[i]].transfer(targetGas);
            TransferAmount(scsBeneficiary[nodeList[i]], targetGas);
            
        }

        //award via nodes
        for (i = 0; i < proposals[hash].viaNodeAddress.length; i++) {
            proposals[hash].viaNodeAddress[i].transfer(proposals[hash].viaNodeAmount[i]);
            TransferAmount(proposals[hash].viaNodeAddress[i], proposals[hash].viaNodeAmount[i]);
        }
        

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

        curFlushIndex = pendingFlushIndex + 1;
        if (curFlushIndex > nodeCount) {
            curFlushIndex = 0;
        }

        //if need toauto retire nodes
        if (AUTO_RETIRE) {
            for (i=0; i<AUTO_RETIRE_COUNT; i++) {
                if (indexAutoRetire >= nodeCount) {
                    indexAutoRetire = 0;
                }
                requestRelease(indexAutoRetire);
                indexAutoRetire ++ ;
            }
        }

        //notify v-node
        SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.approveProposal));

        //make protocol pool to know subchain is active
        protocnt.setSubchainActiveBlock();

        //adjust reward
        adjustReward();
        
        //refund current caller
        msg.sender.transfer((gasinit - msg.gas + 15000) * tx.gasprice);
        //ReportStatus("Request ok");
        
        if (subchainstatus == uint(SubChainStatus.pending)) {
            withdrawal();
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
    function increaseReward(uint percent) private {
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
        uint blk = lastFlushBlk + flushInRound + (nodeCount - 1) * 2 * proposalExpiration;
        
        if (block.number >= blk) {
            lastFlushBlk = block.number;
            SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.updateLastFlushBlk));
        }
    }

    function withdraw(address recv, uint amount) public payable {
        require(owner == msg.sender);

        //withdraw to address
        recv.transfer(amount);
    }
    
    function withdrawal() private {
        subchainstatus = uint(SubChainStatus.close);
        //release fund
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);
        //release already enrolled scs
        for (uint i = nodeCount; i > 0; i--) {
            //release fund
            address cur = nodeList[i-1];
            protocnt.releaseFromSubchain(
                cur,
                penaltyBond
            );

            delete nodeList[i-1];
        }

        nodeCount = 0;

        //refund all to owner
        owner.transfer(this.balance);
        
        //kill self
    }

    function close() public {
        require(owner == msg.sender);

        subchainstatus = uint(SubChainStatus.pending);
        
        if (proposalHashInProgress == 0x0) {
            lastFlushBlk = block.number - flushInRound;
            SCS_RELAY.notifySCS(address(this), uint(SCSRelayStatus.updateLastFlushBlk));
        }
    }

    function addSyncNode(address id, string link) public {
        require(owner == msg.sender);
        syncNodes.push(SyncNode(id, link));
    }

    function removeSyncNode(uint index) public {
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
        uint i = 0;
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
    function applyRemoveNodes(uint nodetype) private {
        SubChainProtocolBase protocnt = SubChainProtocolBase(protocol);

        uint count = nodesToDispel.length;
        if (nodetype == 1) {
            count = nodeToReleaseCount;
        }

        if (count == 0) {
            return;
        }

        // all nodes set 0 at initial, set node to be removed as 1.
        uint[] memory nodeMark = new uint[](nodeCount);
        uint idx = 0;
        uint i = 0;
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
    

    function getindexByte(address a, uint8 randIndex1, uint8 randIndex2) private  pure returns (uint b) {
        uint8 first = uint8(uint(a) / (2 ** (4 * (39 - uint(randIndex1)))) * 2 ** 4);
        uint8 second = uint8(uint8(uint(a) / (2 ** (4 * (39 - uint(randIndex2)))) * 2 ** 4) / 2 ** 4);    // &15
        return uint(byte(first + second));
    }
    
    function matchSelTarget(address addr, uint8 index1, uint8 index2) public view returns (bool) {
        // check if selTargetdist matches.
        uint addr0 = getindexByte(addr, index1, index2);
        uint cont0 = getindexByte(address(this), index1, index2);

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
    
}