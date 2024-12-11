package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// GetSolidityFiles recursively finds all Solidity files in the directory
func GetSolidityFiles(rootPath string) ([]string, error) {
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
		return nil, fmt.Errorf("error walking through files: %v", err)
	}
	
	return files, nil
}

// ReadFileContent reads the entire content of a file
func ReadFileContent(filePath string) (string, error) {
	originalCode, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	return string(originalCode), nil
}

// WriteFileContent writes content to a file
func WriteFileContent(filePath string, content []byte) error {
	return ioutil.WriteFile(filePath, content, 0644)
}