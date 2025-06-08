package main

import (
	"encoding/json"
	"fmt"
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
func saveSession(entries []*TabData) error {
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
func loadSession() ([]SessionEntry, error) {
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

	// If Filepath is set, try to read real content
	for i, tab := range session {
		if tab.Filepath != "" {
			content, err := readFileContent(tab.Filepath)
			if err == nil {
				session[i].Text = content
			} else {
				fmt.Println("Warning: Could not open file from session:", err)
			}
		}
	}

	return session, nil
}
