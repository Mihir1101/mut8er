package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetSolidityFiles recursively finds all Solidity files in the directory
func GetSolidityFiles(rootPath string) []string {
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

// ReadFile reads the content of a file
func ReadFile(filePath string) ([]byte, error) {
	return ioutil.ReadFile(filePath)
}

// WriteFile writes content to a file
func WriteFile(filePath string, data []byte) error {
	return ioutil.WriteFile(filePath, data, 0644)
}

// ReplaceLine replaces a specific line in the code with a mutated version
func ReplaceLine(code string, lineIndex int, newLine string) string {
	lines := strings.Split(code, "\n")
	lines[lineIndex] = newLine
	return strings.Join(lines, "\n")
}
