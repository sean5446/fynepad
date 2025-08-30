package main

import (
	"encoding/json"

	"fyne.io/fyne/v2"
)

type SessionData struct {
	WindowState    WindowState    `json:"WindowState"`
	SessionEntries []SessionEntry `json:"SessionEntries"`
	RecentFiles    []string       `json:"RecentFiles"`
}

type WindowState struct {
	Width  float32 `json:"Width"`
	Height float32 `json:"Height"`
	// X int `json:"X"`
	// Y int `json:"Y"`
}

var defaultWindowState = WindowState{
	Width:  1000,
	Height: 800,
}

type SessionEntry struct {
	Text         string        `json:"Text"`
	Wrapping     fyne.TextWrap `json:"Wrapping"`
	CursorRow    int           `json:"CursorRow"`
	CursorColumn int           `json:"CursorColumn"`
	Title        string        `json:"Title"`
	Filepath     string        `json:"Filepath"`
}

const sessionFile = "session.json"

func saveSession(entries []*TabData, recentFiles []string, windowState WindowState) error {
	var tabSessions []SessionEntry
	for _, data := range entries {
		entry := data.Entry
		tabSessions = append(tabSessions, SessionEntry{
			Text:         entry.Text,
			Wrapping:     entry.Wrapping,
			CursorRow:    entry.CursorRow,
			CursorColumn: entry.CursorColumn,
			Title:        entry.Title,
			Filepath:     entry.Filepath,
		})
	}

	sessionData := SessionData{
		WindowState: WindowState{
			Width:  windowState.Width,
			Height: windowState.Height,
		},
		SessionEntries: tabSessions,
		RecentFiles:    recentFiles,
	}

	bytes, err := json.MarshalIndent(sessionData, "", "  ")
	if err != nil {
		return err
	}

	return writeFileContent(sessionFile, string(bytes))
}

func loadSession() ([]SessionEntry, WindowState, []string, error) {
	content, err := readFileContent(sessionFile)
	// error reading or empty content
	if err != nil || content == "" {
		return []SessionEntry{}, WindowState{}, []string{}, err
	}

	// there was content, but its not valid json
	var sessionData SessionData
	if err := json.Unmarshal([]byte(content), &sessionData); err != nil {
		return []SessionEntry{}, WindowState{}, []string{}, err
	}

	// Refill Text content if file paths exist
	for i, tab := range sessionData.SessionEntries {
		if tab.Filepath != "" {
			content, err := readFileContent(tab.Filepath)
			if err == nil {
				sessionData.SessionEntries[i].Text = content
			}
		}
	}

	return sessionData.SessionEntries, sessionData.WindowState, sessionData.RecentFiles, nil
}
