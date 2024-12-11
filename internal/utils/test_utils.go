package utils

import (
	"os/exec"
	"regexp"
)

// ExtractTestSummary extracts a simple pass/fail summary from the test output
func ExtractTestSummary(output string) string {
	// Match the line containing test result summary
	re := regexp.MustCompile(`Suite result:.*(\d+) passed; (\d+) failed;.*`)
	match := re.FindStringSubmatch(output)
	
	if len(match) > 2 {
		failed := match[2]

		// If no tests failed, return PASS; otherwise, return FAIL
		if failed == "0" {
			return "PASS"
		}
		return "FAIL"
	}
	return "UNKNOWN"
}

// RunForgeTest runs forge tests in the specified project path
func RunForgeTest(projectPath string) (string, error) {
	cmd := exec.Command("forge", "test")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	
	return string(output), err
}