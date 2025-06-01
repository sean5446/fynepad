package main

import (
	"encoding/json"
	"os"
)

const sessionfile = "session.json"

type TabData struct {
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

func loadSessionData() ([]TabData, error) {
	file, err := os.Open(sessionfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var tabs []TabData
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tabs)
	if err != nil {
		return nil, err
	}

	return tabs, nil
}

func saveSessionData(tabs []TabData) error {
	file, err := os.Create(sessionfile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Pretty print
	return encoder.Encode(tabs)
}
