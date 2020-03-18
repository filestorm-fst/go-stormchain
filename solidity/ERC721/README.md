# ERC721 Token

ERC721 is a set of rules that expands on the base ERC20 standards to give it more function and security. One of the key differences between ERC20 and ERC721 is that where ERC20 tokens are fungible, ERC721 are non-fungible. The best way to think of a non-fungible token is to think of it as a trading card, where each trading card has a unique identifier making each one, one of a kind. This opens up the door to many implementation possibilities, one fun example being CryptoKitties.

The ERC721 standards require you to implement a minimum of 10 functions and 3 event types. The 10 required functions are: totalSupply, balanceOf, ownerOf, isApprovedForAll, setApprovalForAll, getApproved, approve, transferFrom, and 2 different versions of safeTransferFrom. The 3 required events are Approval, Transfer, and ApprovalForAll you are of course free to add any additional functions that will make your tokens even more unique.

You can find a basic implementation of an ERC721 token in: StormToken721
