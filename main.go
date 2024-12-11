package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

// MutationRule defines a struct for mutation logic
type MutationRule struct {
	Original string
	Mutant   string
}

// Predefined mutation rules
var mutationRules = []MutationRule{
	{"+", "-"},
	{"-", "+"},
	{">", "<"},
	{"<", ">"},
	{"*", "/"},
	{"/", "*"},
}

var mutantsDir = "mutants"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mutant <path to foundry project>")
		return
	}

	projectPath := os.Args[1]
	contractsPath := filepath.Join(projectPath, "src")
	fmt.Printf("Looking for contracts in: %s\n", contractsPath)

	files := getSolidityFiles(contractsPath)
	if len(files) == 0 {
		fmt.Println("No Solidity files found in", contractsPath)
		return
	}

	var wg sync.WaitGroup

	for _, file := range files {
		fmt.Printf("Processing file: %s\n", file)
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			if err := processFile(filePath, projectPath); err != nil {
				fmt.Printf("Error processing file %s: %v\n", filePath, err)
			}
		}(file)
	}

	wg.Wait()
}

// getSolidityFiles recursively finds all Solidity files in the directory
func getSolidityFiles(rootPath string) []string {
	var files []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".sol") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking through files: %v\n", err)
	}
	return files
}

// processFile handles creating and testing mutants for a given Solidity file
func processFile(filePath, projectPath string) error {
	originalCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	code := string(originalCode)
	mutants := generateMutants(code, mutationRules)
	if len(mutants) == 0 {
		fmt.Printf("No mutants created for %s\n", filePath)
		return nil
	}

	var wg sync.WaitGroup

	for i, mutant := range mutants {
		wg.Add(1)
		go func(mutantCode string, index int) {
			defer wg.Done()
			fmt.Printf("Testing mutant %d for file: %s\n", index+1, filePath)
			if err := testMutant(filePath, mutantCode, index+1, projectPath, originalCode); err != nil {
				fmt.Printf("Error testing mutant %d: %v\n", index+1, err)
			}
		}(mutant, i)
	}

	wg.Wait()
	return nil
}

// testMutant replaces the file with a mutant, runs tests, and restores the original
func testMutant(filePath string, mutant string, mutantNumber int,projectPath string, originalCode []byte) error {
    // Backup the original file
    backupPath := filePath + ".mutant"+ fmt.Sprint(mutantNumber) + ".backup"
    if _, err := os.Stat(backupPath); os.IsNotExist(err) {
        err = ioutil.WriteFile(backupPath, originalCode, 0644)
        if err != nil {
            return fmt.Errorf("failed to backup original file: %v", err)
        }
    }
    fmt.Printf("Backup created: %s\n", backupPath)

    // Replace the file with the mutant
    err := ioutil.WriteFile(filePath, []byte(mutant), 0644)
    if err != nil {
        return fmt.Errorf("failed to write mutant to file: %v", err)
    }

    // Run tests
    cmd := exec.Command("forge", "test")
    cmd.Dir = projectPath
    output, err := cmd.CombinedOutput()
    fmt.Printf("Test output:\n%s\n", string(output))

    // Restore the original file
    restoreErr := os.Rename(backupPath, filePath)
    if restoreErr != nil {
        return fmt.Errorf("failed to restore original file: %v", restoreErr)
    }

    return err
}


// generateMutants applies mutation rules to create all possible mutants of the code
func generateMutants(code string, rules []MutationRule) []string {
	lines := strings.Split(code, "\n")
	var mutants []string

	for lineIndex, line := range lines {
		// Skip irrelevant lines
		if strings.HasPrefix(strings.TrimSpace(line), "//") || strings.HasPrefix(strings.TrimSpace(line), "pragma") || strings.Contains(line, "SPDX-License-Identifier") {
			continue
		}

		// Skip lines containing "++" or "--"
		if strings.Contains(line, "++") || strings.Contains(line, "--") {
			// Do not mutate lines with "++" or "--"
			continue
		}

		// Apply each mutation rule
		for _, rule := range rules {
			if strings.Contains(line, rule.Original) {
				mutatedLine := strings.Replace(line, rule.Original, rule.Mutant, 1)
				mutants = append(mutants, replaceLine(code, lineIndex, mutatedLine))
			}
		}
	}

	return mutants
}

// replaceLine replaces a specific line in the code with a mutated version
func replaceLine(code string, lineIndex int, newLine string) string {
	lines := strings.Split(code, "\n")
	lines[lineIndex] = newLine
	return strings.Join(lines, "\n")
}
