package main

import (
	"encoding/json"

	"fyne.io/fyne/v2"
)

type SessionData struct {
	WindowState WindowState `json:"WindowState"`
	TabDetail   []TabDetail `json:"TabDetail"`
	RecentFiles []string    `json:"RecentFiles"`
}

type WindowState struct {
	Width       float32 `json:"Width"`
	Height      float32 `json:"Height"`
	TabSelected int     `json:"TabSelected"`
	// X int `json:"X"`
	// Y int `json:"Y"`
}

var defaultWindowState = WindowState{
	Width:       1000,
	Height:      800,
	TabSelected: 0,
}

type TabDetail struct {
	Text         string `json:"Text"`
	Wrapping     int    `json:"Wrapping"`
	CursorRow    int    `json:"CursorRow"`
	CursorColumn int    `json:"CursorColumn"`
	Title        string `json:"Title"`
	FilePath     string `json:"FilePath"`
}

const sessionFile = "session.json"


func saveSession(tabManager *TabManager, menuManager *MenuManager, w fyne.Window) error {
	var tabDetail []TabDetail
	for _, data := range tabManager.tabsData {
		entry := data.entry
		textValue := entry.Text
		if entry.filePath != "" {
			textValue = ""
		}
		tabDetail = append(tabDetail, TabDetail{
			Text:         textValue,
			Wrapping:     int(entry.Wrapping),
			CursorRow:    entry.CursorRow,
			CursorColumn: entry.CursorColumn,
			Title:        entry.title,
			FilePath:     entry.filePath,
		})
	}

	sessionData := SessionData{
		WindowState: WindowState{
			Width:       w.Canvas().Size().Width,
			Height:      w.Canvas().Size().Height,
			TabSelected: tabManager.tabs.SelectedIndex(),
		},
		TabDetail:   tabDetail,
		RecentFiles: menuManager.recentFiles,
	}

	bytes, err := json.MarshalIndent(sessionData, "", "  ")
	if err != nil {
		return err
	}

	return writeFileContent(sessionFile, string(bytes))
}

func loadSession() ([]TabDetail, WindowState, []string, error) {
	content, err := readFileContent(sessionFile)
	// error reading or empty content
	if err != nil || content == "" {
		return []TabDetail{}, WindowState{}, []string{}, err
	}

	// there was content, but its not valid json
	var sessionData SessionData
	if err := json.Unmarshal([]byte(content), &sessionData); err != nil {
		return []TabDetail{}, WindowState{}, []string{}, err
	}

	// refill Text content if file paths exist
	tabDetails := make([]TabDetail, 0, len(sessionData.TabDetail))
	for _, tab := range sessionData.TabDetail {
		if tab.FilePath != "" {
			content, err := readFileContent(tab.FilePath)
			if err == nil {
				tab.Text = content
				tabDetails = append(tabDetails, tab)
			}
		} else {
			tabDetails = append(tabDetails, tab)
		}
	}

	return tabDetails, sessionData.WindowState, sessionData.RecentFiles, nil
}
