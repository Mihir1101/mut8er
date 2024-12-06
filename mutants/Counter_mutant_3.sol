// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.13;

contract Counter {
    uint256 public number;

    function setNumber(uint256 newNumber) public {
        number = newNumber;
    }

    function increment() public {
        number++;
    }

    function decrement() public {
        number--;
    }

    function add(uint256 a, uint256 b) pure public returns (uint256){
        // @mutant return(a+b);
        // @mutant return(a-b);
        return(a+b);
    }

    function multiply(uint256 a, uint256 b) pure public returns (uint256){
        return(a*b);
    }

    function compareLess(uint256 a, uint256 b) pure public returns (bool){
        // @mutant return(a<b);
        return(a>b);
    }
}

