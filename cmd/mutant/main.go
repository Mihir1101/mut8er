package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"mutation-tool/internal/models"
	"mutation-tool/internal/processor"
	"mutation-tool/internal/reporter"
	"mutation-tool/internal/utils"
)

func main() {
	start := time.Now()
	if len(os.Args) < 2 {
		fmt.Println("Usage: mutant <path to foundry project>")
		return
	}

	projectPath := os.Args[1]
	contractsPath := filepath.Join(projectPath, "src")
	fmt.Printf("Looking for contracts in: %s\n", contractsPath)

	// Find Solidity files
	files, err := utils.GetSolidityFiles(contractsPath)
	if err != nil {
		fmt.Printf("Error finding Solidity files: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("No Solidity files found in", contractsPath)
		return
	}

	// Initialize processor and reporter
	mutantProcessor := processor.NewMutantProcessor(projectPath)
	reportGenerator := reporter.NewReportGenerator()

	var overallReports []models.ContractMutationReport
	var reportMutex sync.Mutex

	var wg sync.WaitGroup
	for _, file := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			
			// Process mutants for the file
			report, err := mutantProcessor.ProcessMutantsForFile(filePath)
			if err != nil {
				fmt.Printf("Error processing mutants for %s: %v\n", filePath, err)
				return
			}
			
			// Generate Markdown report for this contract
			err = reportGenerator.GenerateMarkdownReport(report)
			if err != nil {
				fmt.Printf("Error generating report for %s: %v\n", filePath, err)
			}

			reportMutex.Lock()
			overallReports = append(overallReports, report)
			reportMutex.Unlock()
		}(file)
	}

	wg.Wait()

	// Generate an overall summary
	err = reportGenerator.GenerateOverallSummaryReport(overallReports)
	if err != nil {
		fmt.Printf("Error generating overall summary: %v\n", err)
	}

	elapsed := time.Since(start)
	fmt.Printf("Mutation testing completed in %s\n", elapsed)
}