package mutations

import (
	"fmt"
	"strings"

	"mutation-tool/parser"
	"mutation-tool/utils"
)

// GenerateMutants applies mutation rules to create all possible mutants of the code
func GenerateMutants(code string, rules []MutationRule) []string {
	lines := strings.Split(code, "\n")
	var mutants []string

	for lineIndex, line := range lines {
		// Skip irrelevant lines
		if !parser.IsRelevantLine(line) {
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
	leadingSpaces := parser.GetLeadingSpaces(originalLine)

	// Tagging the mutation with the original line and the mutated version
	tag := fmt.Sprintf("%s// @mutant %s\n%s", leadingSpaces, originalLine[len(leadingSpaces):], mutatedLine)
	return utils.ReplaceLine(code, lineIndex, tag)
}
