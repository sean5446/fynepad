package main

import (
	"encoding/json"
	"os"
)

const sessionfile = "session.json"

type TabSessionData struct {
	Text         string `json:"Text"`
	Wrapping     int    `json:"Wrapping"`
	CursorRow    int    `json:"CursorRow"`
	CursorColumn int    `json:"CursorColumn"`
	Title        string `json:"Title"`
	Filepath     string `json:"Filepath"`
}

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

func saveSessionData(tabsData []*TabEntryWithShortcut) {
	var session []TabSessionData
	for _, entry := range tabsData {
		session = append(session, TabSessionData{
			Text:         entry.Text,
			Wrapping:     int(entry.Wrapping),
			CursorRow:    entry.CursorRow,
			CursorColumn: entry.CursorColumn,
			Title:        entry.Title,
			Filepath:     entry.Filepath,
		})
	}

	data, _ := json.MarshalIndent(session, "", "  ")
	os.WriteFile("session.json", data, 0644)
}
