package processor

import (
	"path/filepath"
	"strings"
	"sync"

	"mutation-tool/internal/models"
	"mutation-tool/internal/mutationrules"
	"mutation-tool/internal/utils"
)

type MutantProcessor struct {
	ProjectPath string
	mutantsMutex sync.Mutex
	fileMutex   sync.Mutex
}

func NewMutantProcessor(projectPath string) *MutantProcessor {
	return &MutantProcessor{
		ProjectPath: projectPath,
	}
}

func (mp *MutantProcessor) ProcessMutantsForFile(filePath string) (models.ContractMutationReport, error) {
	// Read original file content
	originalCode, err := utils.ReadFileContent(filePath)
	if err != nil {
		return models.ContractMutationReport{}, err
	}

	lines := strings.Split(originalCode, "\n")
	
	report := models.ContractMutationReport{
		FileName: filepath.Base(filePath),
	}

	var mutantDetailsList []models.MutantDetails
	var wg sync.WaitGroup

	// Get mutation rules
	mutationRules := mutationrules.GetDefaultMutationRules()

	for lineIndex, line := range lines {
		// Skip lines not suitable for mutation
		if !mutationrules.IsMutationCandidate(line) {
			continue
		}

		for _, rule := range mutationRules {
			if strings.Contains(line, rule.Original) {
				wg.Add(1)
				go func(lineIndex int, line string, rule models.MutationRule) {
					defer wg.Done()

					// Create mutant
					mutatedLine := strings.Replace(line, rule.Original, rule.Mutant, 1)
					mutantCode := mp.replaceLine(originalCode, lineIndex, mutatedLine)

					// Test the mutant
					mp.fileMutex.Lock()
					err := utils.WriteFileContent(filePath, []byte(mutantCode))
					if err != nil {
						mp.fileMutex.Unlock()
						return
					}

					output, _ := utils.RunForgeTest(mp.ProjectPath)
					testSummary := utils.ExtractTestSummary(string(output))
					
					// Restore original file
					utils.WriteFileContent(filePath, []byte(originalCode))
					mp.fileMutex.Unlock()

					// Create mutant details
					mutantDetails := models.MutantDetails{
						OriginalLine: line,
						MutatedLine:  mutatedLine,
						TestOutcome:  testSummary,
						RuleApplied:  rule,
					}

					mp.mutantsMutex.Lock()
					mutantDetailsList = append(mutantDetailsList, mutantDetails)
					mp.mutantsMutex.Unlock()
				}(lineIndex, line, rule)
			}
		}
	}

	wg.Wait()

	// Populate report
	report.TotalMutants = len(mutantDetailsList)
	report.MutantDetails = mutantDetailsList
	report.PassedMutants = mp.countPassedMutants(mutantDetailsList)
	report.FailedMutants = report.TotalMutants - report.PassedMutants

	return report, nil
}

func (mp *MutantProcessor) replaceLine(code string, lineIndex int, newLine string) string {
	lines := strings.Split(code, "\n")
	lines[lineIndex] = newLine
	return strings.Join(lines, "\n")
}

func (mp *MutantProcessor) countPassedMutants(mutants []models.MutantDetails) int {
	passed := 0
	for _, mutant := range mutants {
		if mutant.TestOutcome == "PASS" {
			passed++
		}
	}
	return passed
}