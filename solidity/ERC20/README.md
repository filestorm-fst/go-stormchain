# ERC20 Token

ERC20 is a set of rules defined by the Ethereum community to create a fungible (can be replaced by an equal part: such as a dollar bill can be exchanged with 4 quarters). This set of rules makes it easier for developers to create tokens that can be used with most wallets that support Ether.

An ERC20 token is simply a smart contract that holds information about the token, such as the total supply, the name, etc.

ERC20 states that there are a minimum of 6 functions, and 2 event types that must be implement in order to adhere to the rules of the standard. the 6 required functions are: totalSupply, balanceOf, transfer, approve, allowance, and transferFrom. The 2 event types are: Approval and Transfer. It does not matter how you choose to implement these functions as long as they adhere to the return values specified by the ERC20 standard. You are also free to add any additional functions to your token in order to give it uniqueness.

You can find an example of a basic implementation of an ERC20 token in: StormERC20.sol
