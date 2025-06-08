package main

import (
	"os"
)

func readFileContent(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeFileContent(path, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}
