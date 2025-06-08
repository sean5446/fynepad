package main

import (
	"os"
)

// ReadFileContent returns the contents of a file as a string.
func ReadFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFileContent writes the string data to the specified file path.
func WriteFileContent(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
