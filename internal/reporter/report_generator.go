package reporter

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"mutation-tool/internal/models"
)

const (
	reportDir = "mutation_reports"
)

type ReportGenerator struct{}

func NewReportGenerator() *ReportGenerator {
	// Ensure report directory exists
	os.MkdirAll(reportDir, os.ModePerm)
	return &ReportGenerator{}
}

func (rg *ReportGenerator) GenerateMarkdownReport(report models.ContractMutationReport) error {
	if report.TotalMutants == 0 {
		return nil
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
	
	return os.WriteFile(reportPath, []byte(markdownContent), 0644)
}

func (rg *ReportGenerator) GenerateOverallSummaryReport(reports []models.ContractMutationReport) error {
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
	return os.WriteFile(summaryPath, []byte(markdownContent), 0644)
}