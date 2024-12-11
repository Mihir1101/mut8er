# Mutation Testing Tool for Solidity Smart Contracts

## Overview

This is a mutation testing tool designed specifically for Solidity smart contracts using the Foundry testing framework. The tool helps identify potential weaknesses in your test suites by automatically generating mutants (modified versions) of your smart contract code and running tests against them.

## Features

- Recursive Solidity file discovery
- Multiple mutation rules (arithmetic, comparison, logical operators)
- Parallel processing of contract mutations
- Detailed markdown reporting
- Overall mutation testing summary

## Prerequisites

- Go (1.18 or later)
- Foundry
- Access to a Solidity project

## Installation

1. Clone the repository
2. Run `go mod tidy` to download dependencies
3. Build the project: `go build ./cmd/mutant`

## Usage

```bash
./mutant /path/to/your/foundry/project
```

## Mutation Rules

The tool currently supports mutations for:
- Arithmetic operators (`+` ↔ `-`, `*` ↔ `/`)
- Comparison operators (`>` ↔ `<`, `==` ↔ `!=`)
- Logical operators (`&&` ↔ `||`)

## Reports

After running, mutation reports will be generated in the `mutation_reports` directory:
- Individual contract mutation reports
- An overall mutation testing summary
