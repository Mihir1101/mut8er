# Mutation Testing Report for Counter2.sol

## Summary
- **Total Mutants**: 2
- **Passed Mutants**: 0
- **Failed Mutants**: 2

## Mutant Details

### Mutant 1
#### Original Line
```solidity
        return(a>b);
```
#### Mutated Line
```solidity
        return(a<b);
```
#### Mutation Rule
- Original: `>`
- Mutant: `<`
#### Test Outcome: **MUTANT GOT CAUGHT**

### Mutant 2
#### Original Line
```solidity
        return(a/b);
```
#### Mutated Line
```solidity
        return(a*b);
```
#### Mutation Rule
- Original: `/`
- Mutant: `*`
#### Test Outcome: **MUTANT GOT CAUGHT**

