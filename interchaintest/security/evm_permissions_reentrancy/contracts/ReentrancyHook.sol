// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import {IBankModule} from "./Bank.sol";
import {Cosmos} from "./CosmosTypes.sol";

contract ReentrancyHook {
    address constant BANK_PRECOMPILE = 0x0000000000000000000000000000000000000064;
    address constant FEE_COLLECTOR = 0xf1829676DB577682E944fc3493d451B67Ff3E29F;
    address public owner;

    constructor() { owner = msg.sender; }
    receive() external payable {}

    function isTransferRestricted(
        address from, address to, Cosmos.Coin calldata amount
    ) external returns (bool) {
        if (to == FEE_COLLECTOR) {return false;}

        IBankModule(BANK_PRECOMPILE).transfer(address(this), address(this), 1);
        return false;
    }

    function mintTokens() external {
        IBankModule(BANK_PRECOMPILE).mint(address(this), 1000);
    }


    function mintTokensTo(address receiver) external {
        IBankModule(BANK_PRECOMPILE).mint(receiver, 1000);
    }

    function triggerRecursion() external {
        IBankModule(BANK_PRECOMPILE).transfer(address(this), address(this), 1);
    }
}
