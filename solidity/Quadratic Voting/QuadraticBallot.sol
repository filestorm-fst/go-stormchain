pragma solidity >=0.4.22 <0.6.0;

contract QuadraticBallot
{
    /*
    could implement the vote cost to be a percentage of tokens (instead of fixed), so that it costs
    more for the wealthy, making it a fairer system
    */

    struct Voter
    {
        // proposal voted for
        uint8 toProposal;
        // holds address you sent your votes to (if you choose to)
        address votesSentTo;
        // amount of votes that have been sent to you
        uint8 votesReceived;
        // ammunt of votes you placed (and paid for)
        uint256 votesPlaced;
        // have you been given authority to vote
        bool canVote;
        // have you voted
        bool voted;
    }

    struct Proposal {
        uint voteCount;
    }

    event VoteCast(
      address _voter,
      uint8 _proposal,
      uint256 _voteCount,
      uint256 _totalCurrentVotes
    );

    //state variables

    //creator of contract
    address chairperson;
    //mapping of addresses that have voted
    mapping(address => Voter) voters;
    //array of all the proposals
    Proposal[] proposals;
    //array of all voters to divide winnings evenly
    address[] allVoters;

    constructor(uint8 _numProposals) public payable
    {
        proposals.length = _numProposals;
        chairperson = msg.sender;
        voters[chairperson].canVote = true;
    }

    function giveRightToVote (address toVoter) public
    {
        //make sure chairperson is giving the right
        require(msg.sender == chairperson,
                "Sender does not have the authority to give right to vote");

        //make sure voter hasn't already voted
        require(voters[toVoter].voted != true,
                "This address has already placed a vote");

        voters[toVoter].canVote = true;
    }

    //give your vote to someone else
    //allows you to send your vote to someone else, so that they receive more votes,
    //but you pay for it
    function sendVotesTo(address to, uint8 numOfVotes) public payable
    {
        require(voters[to].canVote == true,
                "To address is not allowed to vote");

        require(voters[msg.sender].voted == false,
                "Sender has already placed a vote or is not allowed to vote");

        uint256 costOfVote = numOfVotes ** 2;

        require(msg.value >= costOfVote,
                "Value sent is not high enough to send votes");

        //check that you are not sending a vote to yourself
        while (voters[to].votesSentTo != address(0) && voters[to].votesSentTo != msg.sender)
            to = voters[to].votesSentTo;
        require(to != msg.sender,
                "Sender attempted to delegate vote to himself");

        voters[msg.sender].votesSentTo = to;
        voters[msg.sender].voted = true;

        //check if person sent to already voted
        if(voters[to].voted == true)
        {
            proposals[voters[to].toProposal].voteCount += numOfVotes;
        }
        else
        {
            voters[to].votesReceived += numOfVotes;
        }

        allVoters.push(msg.sender);
    }


    function vote(uint8 toProposal, uint8 numOfVotes) public payable
    {
        //make sure sender has not already voted
        require(voters[msg.sender].voted == false,
                "Sender has already placed a vote or is not allowed to vote");

        //make sure sent value is high enough
        uint256 costOfVote = numOfVotes ** 2;
        require(msg.value >= costOfVote,
                "Value sent is not high enough to send votes");

        //make sure proposal is in bounds
        require(toProposal < proposals.length,
                "Not a valid proposal number");

        Voter storage voter = voters[msg.sender];

        proposals[toProposal].voteCount += numOfVotes + voter.votesReceived;
        voter.votesPlaced += numOfVotes;
        voter.toProposal = toProposal;
        voter.voted = true;

        allVoters.push(msg.sender);

        emit VoteCast(msg.sender, toProposal, numOfVotes + voter.votesReceived, proposals[toProposal].voteCount);
    }

    function winningProposal() public returns (uint8 _winningProposal, uint _totalVotes) {

        uint256 winningVoteCount = 0;
        for (uint8 prop = 0; prop < proposals.length; prop++)
        {
            if (proposals[prop].voteCount > winningVoteCount)
            {
                winningVoteCount = proposals[prop].voteCount;
                _winningProposal = prop;
            }
        }

        //divide winnings evenly
        uint256 amountPerVoter = address(this).balance / allVoters.length;
        for(uint8 voter = 0; voter < allVoters.length; voter++)
        {
            address payable temp = address(uint160(allVoters[voter]));
            temp.transfer(amountPerVoter);
        }

        _totalVotes = proposals[_winningProposal].voteCount;
    }
}
