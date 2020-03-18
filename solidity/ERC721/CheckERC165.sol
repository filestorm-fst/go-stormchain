pragma solidity >=0.5.16 <0.6.0;

interface ERC165 {
 /// @notice Query if a contract implements an interface
 /// @param interfaceID The interface identifier, as
 ///  specified in ERC-165
 /// @dev Interface identification is specified in
 ///  ERC-165. This function uses less than 30,000 gas.
 /// @return `true` if the contract implements `interfaceID`
 ///  and `interfaceID` is not 0xffffffff, `false` otherwise
 function supportsInterface(bytes4 interfaceID) external view returns (bool);
}

contract CheckERC165 is ERC165
{
    //mapping of all the supported interfaces
    mapping(bytes4 => bool) internal supportedInterfaces;

    constructor() public
    {
        supportedInterfaces[this.supportsInterface.selector] = true;
    }

    function supportsInterface(bytes4 interfaceID) external view returns (bool)
    {
        return supportedInterfaces[interfaceID];
    }
}
