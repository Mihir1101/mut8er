package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"regexp"
)

// MutationRule defines a struct for mutation logic
type MutationRule struct {
	Original string
	Mutant   string
}

type MutantDetails struct {
	OriginalLine string
	MutatedLine  string
	TestOutcome  string
	RuleApplied  MutationRule
}

type ContractMutationReport struct {
	FileName         string
	TotalMutants     int
	PassedMutants    int
	FailedMutants    int
	MutantDetails    []MutantDetails
}

var (
	mutationRules = []MutationRule{
		{"+", "-"},
		{"-", "+"},
		{">", "<"},
		{"<", ">"},
		{"*", "/"},
		{"/", "*"},
	}
	mutantsDir = "mutants"
	reportDir  = "mutation_reports"
	fileMutex  sync.Mutex
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mutant <path to foundry project>")
		return
	}

	// Ensure report directory exists
	os.MkdirAll(reportDir, os.ModePerm)

	projectPath := os.Args[1]
	contractsPath := filepath.Join(projectPath, "src")
	fmt.Printf("Looking for contracts in: %s\n", contractsPath)

	files := getSolidityFiles(contractsPath)
	if len(files) == 0 {
		fmt.Println("No Solidity files found in", contractsPath)
		return
	}

	var overallReports []ContractMutationReport
	var reportMutex sync.Mutex

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			report := processMutantsForFile(filePath, projectPath)
			
			reportMutex.Lock()
			overallReports = append(overallReports, report)
			reportMutex.Unlock()

			// Generate Markdown report for this contract
			generateMarkdownReport(report)
		}(file)
	}

	wg.Wait()

	// Generate an overall summary
	generateOverallSummaryReport(overallReports)
}

func processMutantsForFile(filePath, projectPath string) ContractMutationReport {
	originalCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Printf("Failed to read file %s: %v\n", filePath, err)
		return ContractMutationReport{}
	}

	code := string(originalCode)
	lines := strings.Split(code, "\n")
	
	report := ContractMutationReport{
		FileName: filepath.Base(filePath),
	}

	var mutantDetailsList []MutantDetails
	var mutantsMutex sync.Mutex
	var wg sync.WaitGroup

	for lineIndex, line := range lines {
		// Skip irrelevant lines
		if strings.HasPrefix(strings.TrimSpace(line), "//") || 
			strings.HasPrefix(strings.TrimSpace(line), "pragma") || 
			strings.Contains(line, "SPDX-License-Identifier") ||
			strings.Contains(line, "++") || 
			strings.Contains(line, "--") {
			continue
		}

		for _, rule := range mutationRules {
			if strings.Contains(line, rule.Original) {
				wg.Add(1)
				go func(lineIndex int, line string, rule MutationRule) {
					defer wg.Done()

					// Create mutant
					mutatedLine := strings.Replace(line, rule.Original, rule.Mutant, 1)
					mutantCode := replaceLine(code, lineIndex, mutatedLine)

					// Test the mutant
					fileMutex.Lock()
					err := ioutil.WriteFile(filePath, []byte(mutantCode), 0644)
					if err != nil {
						fileMutex.Unlock()
						return
					}

					cmd := exec.Command("forge", "test")
					cmd.Dir = projectPath
					output, _ := cmd.CombinedOutput()
					
					testSummary := extractTestSummary(string(output))
					
					// Restore original file
					ioutil.WriteFile(filePath, originalCode, 0644)
					fileMutex.Unlock()

					// Create mutant details
					mutantDetails := MutantDetails{
						OriginalLine: line,
						MutatedLine:  mutatedLine,
						TestOutcome:  testSummary,
						RuleApplied:  rule,
					}

					mutantsMutex.Lock()
					mutantDetailsList = append(mutantDetailsList, mutantDetails)
					mutantsMutex.Unlock()
				}(lineIndex, line, rule)
			}
		}
	}

	wg.Wait()

	// Populate report
	report.TotalMutants = len(mutantDetailsList)
	report.MutantDetails = mutantDetailsList
	report.PassedMutants = countPassedMutants(mutantDetailsList)
	report.FailedMutants = report.TotalMutants - report.PassedMutants

	return report
}

func generateMarkdownReport(report ContractMutationReport) {
	if report.TotalMutants == 0 {
		return
	}

	// Create markdown content
	markdownContent := fmt.Sprintf("# Mutation Testing Report for %s\n\n", report.FileName)
	markdownContent += fmt.Sprintf("## Summary\n")
	markdownContent += fmt.Sprintf("- **Total Mutants**: %d\n", report.TotalMutants)
	markdownContent += fmt.Sprintf("- **Passed Mutants**: %d\n", report.PassedMutants)
	markdownContent += fmt.Sprintf("- **Failed Mutants**: %d\n\n", report.FailedMutants)

	markdownContent += "## Mutant Details\n\n"

	for i, mutant := range report.MutantDetails {
		markdownContent += fmt.Sprintf("### Mutant %d\n", i+1)
		markdownContent += "#### Original Line\n"
		markdownContent += fmt.Sprintf("```solidity\n%s\n```\n", mutant.OriginalLine)
		markdownContent += "#### Mutated Line\n"
		markdownContent += fmt.Sprintf("```solidity\n%s\n```\n", mutant.MutatedLine)
		markdownContent += fmt.Sprintf("#### Mutation Rule\n")
		markdownContent += fmt.Sprintf("- Original: `%s`\n", mutant.RuleApplied.Original)
		markdownContent += fmt.Sprintf("- Mutant: `%s`\n", mutant.RuleApplied.Mutant)
		markdownContent += fmt.Sprintf("#### Test Outcome: **%s**\n\n", mutant.TestOutcome)
	}

	// Write to file
	reportPath := filepath.Join(reportDir, fmt.Sprintf("%s_mutation_report.md", 
		strings.TrimSuffix(report.FileName, ".sol")))
	
	err := ioutil.WriteFile(reportPath, []byte(markdownContent), 0644)
	if err != nil {
		fmt.Printf("Error writing report for %s: %v\n", report.FileName, err)
	}
}

func generateOverallSummaryReport(reports []ContractMutationReport) {
	markdownContent := "# Overall Mutation Testing Summary\n\n"
	markdownContent += "## Contract Mutation Statistics\n\n"

	totalContracts := 0
	totalMutants := 0
	totalPassedMutants := 0
	totalFailedMutants := 0

	for _, report := range reports {
		if report.TotalMutants > 0 {
			totalContracts++
			totalMutants += report.TotalMutants
			totalPassedMutants += report.PassedMutants
			totalFailedMutants += report.FailedMutants

			markdownContent += fmt.Sprintf("### %s\n", report.FileName)
			markdownContent += fmt.Sprintf("- Total Mutants: %d\n", report.TotalMutants)
			markdownContent += fmt.Sprintf("- Passed Mutants: %d\n", report.PassedMutants)
			markdownContent += fmt.Sprintf("- Failed Mutants: %d\n\n", report.FailedMutants)
		}
	}

	markdownContent += "## Overall Summary\n"
	markdownContent += fmt.Sprintf("- **Total Contracts Analyzed**: %d\n", totalContracts)
	markdownContent += fmt.Sprintf("- **Total Mutants**: %d\n", totalMutants)
	markdownContent += fmt.Sprintf("- **Total Passed Mutants**: %d\n", totalPassedMutants)
	markdownContent += fmt.Sprintf("- **Total Failed Mutants**: %d\n", totalFailedMutants)

	// Write overall summary
	summaryPath := filepath.Join(reportDir, "mutation_testing_summary.md")
	err := ioutil.WriteFile(summaryPath, []byte(markdownContent), 0644)
	if err != nil {
		fmt.Printf("Error writing overall summary: %v\n", err)
	}
}

func countPassedMutants(mutants []MutantDetails) int {
	passed := 0
	for _, mutant := range mutants {
		if mutant.TestOutcome == "PASS" {
			passed++
		}
	}
	return passed
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

// extractTestSummary extracts a simple pass/fail summary from the test output
func extractTestSummary(output string) string {
	// Match the line containing test result summary
	re := regexp.MustCompile(`Suite result:.*(\d+) passed; (\d+) failed;.*`)
	match := re.FindStringSubmatch(output)
	if len(match) > 2 {
		failed := match[2]

		// If no tests failed, return PASS; otherwise, return FAIL
		if failed == "0" {
			return "MUTANT SURVIVED"
		}
		return "MUTANT GOT CAUGHT"
	}
	return "UNKNOWN"
}


// replaceLine replaces a specific line in the code with a mutated version
func replaceLine(code string, lineIndex int, newLine string) string {
	lines := strings.Split(code, "\n")
	lines[lineIndex] = newLine
	return strings.Join(lines, "\n")
}