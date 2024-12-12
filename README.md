# Mut8er for Solidity Smart Contracts

## Overview

This is a mutation testing tool designed specifically for Solidity smart contracts using the Foundry testing framework. Mutation testing is a powerful technique to evaluate the quality of your test suites by introducing small, systematic changes (mutations) to your source code and analyzing how well your tests detect these changes.

## Key Features

- **Comprehensive Mutation Analysis**
  - Automatically generates mutants by making small modifications to your smart contract code
  - Supports mutations across various operator types to thoroughly test your test suite's effectiveness
  - Provides detailed insights into how robust your tests are against different code variations

- **Processing Capabilities**
  - Recursive discovery of Solidity files within your project directory
  - Parallel processing of contract mutations for improved performance
  - Intelligent mutation rule application

- **Detailed Reporting**
  - Generates comprehensive markdown reports for each contract
  - Produces an overall mutation testing summary
  - Provides granular details about each mutant, including original and modified code, mutation rules, and test outcomes

## Mutation Rules

The tool currently supports mutations for the following operators:

- **Arithmetic Operators**
  - `+` ↔ `-`
  - `*` ↔ `/`

- **Comparison Operators**
  - `>` ↔ `<`
  - `==` ↔ `!=`

- **Logical Operators**
  - `&&` ↔ `||`

## Prerequisites

Before using the mutation testing tool, ensure you have the following installed:

- Go (version 1.18 or later)
- Foundry testing framework
- Access to a Solidity project structured with Foundry

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Mihir1101/mut8er.git
   ```

2. Navigate to the project directory:
   ```bash
   cd mut8er
   ```

3. Download dependencies:
   ```bash
   go mod tidy
   ```

4. Build the project:
   ```bash
   go build ./cmd/mutant
   ```

## Usage

Run the mutation testing tool by providing the path to your Foundry project:

```bash
./mutant /path/to/your/foundry/project
```

## Limitations and Future Improvements

**Current Limitations:**
- The tool currently runs using the `forge test` command, which means it requires a working Foundry project and is incompatible with other testing frameworks
- Running mutations on large projects can be time-consuming
- Limited to predefined mutation rules

**Upcoming Improvements:**
- Option to target specific files or directories
- More granular mutation rule configuration
- Support for more complex mutation strategies
- Performance optimizations for larger projects

## Planned Feature: Targeted Mutation Testing

Working on enhancing the tool to allow more targeted mutation testing:

- Ability to specify a single file for mutation
- Option to pass a directory path to generate mutants for a specific file
- More flexible and efficient mutation generation process

### Planned Usage Example

```bash
# Future version will support targeted mutations
./mutant /path/to/project --file src/MyContract.sol
```

## Reports

After running the mutation testing tool, you'll find detailed reports in the `mutation_reports` directory:

- Individual contract mutation reports (e.g., `MyContract_mutation_report.md`)
- An overall mutation testing summary (`mutation_testing_summary.md`)

These reports provide comprehensive insights into:
- Total number of mutants generated
- Passed and failed mutants
- Detailed breakdown of mutations and their test outcomes

**Note**: This tool is intended to help improve the quality of your test suites. However, catching all mutant does not guarantee complete test coverage or bug-free code.