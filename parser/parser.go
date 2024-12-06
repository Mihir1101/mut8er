package parser

import (
	"strings"
)

// GetLeadingSpaces returns the leading spaces in a given line.
func GetLeadingSpaces(line string) string {
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

// IsRelevantLine checks if the line should be considered for mutation
func IsRelevantLine(line string) bool {
	return !strings.HasPrefix(strings.TrimSpace(line), "//") &&
		!strings.HasPrefix(strings.TrimSpace(line), "pragma") &&
		!strings.Contains(line, "SPDX-License-Identifier") &&
		!strings.Contains(line, "++") &&
		!strings.Contains(line, "--")
}
