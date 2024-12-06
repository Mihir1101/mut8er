package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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

// Main entry point
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

	for _, file := range files {
		fmt.Printf("Processing file: %s\n", file)
		err := processFile(file)
		if err != nil {
			fmt.Printf("Error processing file %s: %v\n", file, err)
		}
	}

	// Create mutants directory if it doesn't exist
	if err := os.MkdirAll(mutantsDir, 0755); err != nil {
		fmt.Printf("Error creating mutants directory: %v\n", err)
		return
	}
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

// processFile handles creating mutants for a given Solidity file
func processFile(filePath string) error {
	originalCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	code := string(originalCode)
	fmt.Printf("Read file: %s\n", filePath)

	mutants := generateMutants(code, mutationRules)
	if len(mutants) == 0 {
		fmt.Printf("No mutants created for %s\n", filePath)
		return nil
	}

	for i, mutant := range mutants {
		mutantFilePath := fmt.Sprintf("%s/%s_mutant_%d.sol", mutantsDir, strings.TrimSuffix(filepath.Base(filePath), ".sol"), i+1)
		err := ioutil.WriteFile(mutantFilePath, []byte(mutant), 0644)
		if err != nil {
			fmt.Printf("Failed to write mutant file: %v\n", err)
			continue
		}
		fmt.Printf("Mutant created: %s\n", mutantFilePath)
	}

	return nil
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
				mutants = append(mutants, generateMutantTag(code, lineIndex, line, mutatedLine))
			}
		}
	}

	return mutants
}

// generateMutantTag creates the tag for a mutated line to keep track of original and mutated versions
func generateMutantTag(code string, lineIndex int, originalLine, mutatedLine string) string {
	// Get the leading spaces from the original line
	leadingSpaces := getLeadingSpaces(originalLine)

	// Tagging the mutation with the original line and the mutated version
	tag := fmt.Sprintf("%s// @mutant %s\n%s", leadingSpaces, originalLine[len(leadingSpaces):], mutatedLine)
	return replaceLine(code, lineIndex, tag)
}

// getLeadingSpaces extracts the leading spaces from a line
func getLeadingSpaces(line string) string {
	// Count the number of leading spaces
	leadingSpaces := ""
	for _, char := range line {
		if char == ' ' {
			leadingSpaces += " "
		} else {
			break
		}
	}
	return leadingSpaces
}


// replaceLine replaces a specific line in the code with a mutated version
func replaceLine(code string, lineIndex int, newLine string) string {
	lines := strings.Split(code, "\n")
	lines[lineIndex] = newLine
	return strings.Join(lines, "\n")
}