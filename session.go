package main

import (
	"encoding/json"
	"os"
)

const sessionfile = "session.json"


func loadSessionData() ([]*TabEntryWithShortcut, error) {
	file, err := os.Open(sessionfile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tabsData)
	if err != nil {
		return nil, err
	}

	println("Loaded session data:", tabsData)
	result := make([]*TabEntryWithShortcut, len(tabsData))
	copy(result, tabsData)
	return result, nil
}

func saveSessionData(tabsData []*TabEntryWithShortcut) error {
	file, err := os.Create(sessionfile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	// update the text content of each tab
	for _, tab := range tabsData {
		println("Saving tab:", tab.Title, "with text:", tab.Entry.Text)
		tab.Text = tab.Entry.Text
	}

	println("Saving session:", tabsData)
	return encoder.Encode(tabsData)
}
