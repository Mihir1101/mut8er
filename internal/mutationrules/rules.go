package mutationrules

import (
	"strings"
	"mutation-tool/internal/models"
)

// GetDefaultMutationRules returns a list of predefined mutation rules
func GetDefaultMutationRules() []models.MutationRule {
	return []models.MutationRule{
		{"+", "-"},
		{"-", "+"},
		{">", "<"},
		{"<", ">"},
		{"*", "/"},
		{"/", "*"},
		{"==", "!="},
		{"!=", "=="},
		{"&&", "||"},
		{"||", "&&"},
	}
}

// IsMutationCandidate checks if a line is suitable for mutation
func IsMutationCandidate(line string) bool {
	// Skip irrelevant lines
	skipPatterns := []string{
		"//", "pragma", "import", "SPDX-License-Identifier", 
		"++", "--", "+=", "-=", "*=", "/=",
	}

	for _, pattern := range skipPatterns {
		if strings.Contains(line, pattern) {
			return false
		}
	}

	return true
}