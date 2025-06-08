package main

import (
	"encoding/json"
	"os"

	"fyne.io/fyne/v2"
)

type SessionEntry struct {
	Text         string        `json:"Text"`
	Wrapping     fyne.TextWrap `json:"Wrapping"`
	CursorRow    int           `json:"CursorRow"`
	CursorColumn int           `json:"CursorColumn"`
	Title        string        `json:"Title"`
	Filepath     string        `json:"Filepath"`
}

const sessionFile = "session.json"

// SaveSession stores the session data to disk.
func SaveSession(entries []*TabData) error {
	var session []SessionEntry
	for _, data := range entries {
		entry := data.Entry
		session = append(session, SessionEntry{
			Text:         entry.Text,
			Wrapping:     entry.Wrapping,
			CursorRow:    entry.CursorRow,
			CursorColumn: entry.CursorColumn,
			Title:        entry.Title,
			Filepath:     entry.Filepath,
		})
	}

	bytes, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(sessionFile, bytes, 0644)
}

// LoadSession reads the saved tabs from disk.
func LoadSession() ([]SessionEntry, error) {
	if _, err := os.Stat(sessionFile); os.IsNotExist(err) {
		return nil, nil // No session file = clean start
	}

	bytes, err := os.ReadFile(sessionFile)
	if err != nil {
		return nil, err
	}

	var session []SessionEntry
	if err := json.Unmarshal(bytes, &session); err != nil {
		return nil, err
	}

	return session, nil
}
