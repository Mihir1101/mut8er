package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"mutation-tool/mutations"
	"mutation-tool/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mutant <path to foundry project>")
		return
	}

	projectPath := os.Args[1]
	contractsPath := filepath.Join(projectPath, "src")
	fmt.Printf("Looking for contracts in: %s\n", contractsPath)

	files := utils.GetSolidityFiles(contractsPath)
	if len(files) == 0 {
		fmt.Println("No Solidity files found in", contractsPath)
		return
	}

	// Create a WaitGroup for concurrent execution of test runs
	var wg sync.WaitGroup

	// Channel to collect the results of tests
	results := make(chan string)

	// Create a separate file to log mutants that pass all tests
	passFile, err := os.Create("mutants_passed.txt")
	if err != nil {
		fmt.Printf("Error creating pass file: %v\n", err)
		return
	}
	defer passFile.Close()

	// Process each file and test the mutants in parallel
	for _, file := range files {
		fmt.Printf("Processing file: %s\n", file)
		err := processFile(file, projectPath, results, &wg)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", file, err)
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the results channel after all testing is done
	close(results)

	// Write the passed mutants to the file
	for result := range results {
		_, err := passFile.WriteString(result + "\n")
		if err != nil {
			fmt.Printf("Error writing to pass file: %v\n", err)
		}
	}
}

// processFile handles creating mutants and running tests concurrently
func processFile(filePath, projectPath string, results chan<- string, wg *sync.WaitGroup) error {
	originalCode, err := utils.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	code := string(originalCode)
	fmt.Printf("Read file: %s\n", filePath)

	mutants := mutations.GenerateMutants(code, mutations.MutationRules)
	if len(mutants) == 0 {
		fmt.Printf("No mutants created for %s\n", filePath)
		return nil
	}

	// Launch a goroutine for each mutant to run the Foundry tests concurrently
	for i, mutant := range mutants {
		mutantFilePath := fmt.Sprintf("mutants/%s_mutant_%d.sol", strings.TrimSuffix(filepath.Base(filePath), ".sol"), i+1)
		err := utils.WriteFile(mutantFilePath, []byte(mutant))
		if err != nil {
			fmt.Printf("Failed to write mutant file: %v\n", err)
			continue
		}
		fmt.Printf("Mutant created: %s\n", mutantFilePath)

		// Increment the WaitGroup counter
		wg.Add(1)

		// Run the tests in a separate goroutine
		go runTests(filePath, mutantFilePath, projectPath, results, wg)
	}

	return nil
}

// runTests replaces the contract with the mutant and runs Foundry tests
func runTests(originalContractPath, mutantFilePath, projectPath string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Backup the original contract before replacing it
	backupPath := originalContractPath + ".bak"
	err := os.Rename(originalContractPath, backupPath)
	if err != nil {
		fmt.Printf("Error backing up contract: %v\n", err)
		return
	}
	defer os.Rename(backupPath, originalContractPath) // Restore the original contract after test

	// Replace the original contract with the mutant
	err = os.Rename(mutantFilePath, originalContractPath)
	if err != nil {
		fmt.Printf("Error replacing contract with mutant: %v\n", err)
		return
	}

	// Change the working directory to the Foundry project directory
	originalDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}
	defer os.Chdir(originalDir) // Ensure we restore the original directory after tests

	err = os.Chdir(projectPath)
	if err != nil {
		fmt.Printf("Error changing directory to project path: %v\n", err)
		return
	}

	// Execute Foundry's test command (run from the foundry project directory)
	cmd := exec.Command("forge", "test")
	cmd.Dir = projectPath // Ensure we're running `forge test` in the Foundry project directory
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Error running tests for mutant %s: %v\n", mutantFilePath, err)
		return
	}

	// Check if tests passed (can be modified based on output)
	if strings.Contains(string(output), "All tests passed!") {
		results <- mutantFilePath
	} else {
		fmt.Printf("Tests failed for mutant %s\n", mutantFilePath)
	}
}
