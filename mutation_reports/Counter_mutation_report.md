# Mutation Testing Report for Counter.sol

## Summary
- **Total Mutants**: 3
- **Passed Mutants**: 0
- **Failed Mutants**: 3

## Mutant Details

### Mutant 1
#### Original Line
```solidity
        return(a*b);
```
#### Mutated Line
```solidity
        return(a/b);
```
#### Mutation Rule
- Original: `*`
- Mutant: `/`
#### Test Outcome: **MUTANT GOT CAUGHT**

### Mutant 2
#### Original Line
```solidity
        return(a<b);
```
#### Mutated Line
```solidity
        return(a>b);
```
#### Mutation Rule
- Original: `<`
- Mutant: `>`
#### Test Outcome: **MUTANT GOT CAUGHT**

### Mutant 3
#### Original Line
```solidity
        return(a+b);
```
#### Mutated Line
```solidity
        return(a-b);
```
#### Mutation Rule
- Original: `+`
- Mutant: `-`
#### Test Outcome: **MUTANT GOT CAUGHT**

