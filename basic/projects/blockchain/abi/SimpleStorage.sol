// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract SimpleStorage {
    uint256 private value;

    event ValueChanged(uint256 oldValue, uint256 newValue);

    function store(uint256 newValue) public {
        uint256 old = value;
        value = newValue;
        emit ValueChanged(old, newValue);
    }

    function retrieve() public view returns (uint256) {
        return value;
    }
}