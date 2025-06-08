package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

const recentFilePath = "recent.json"

type menuManager struct {
	window      fyne.Window
	tabManager  *TabManager
	recentFiles []string
}

func newMenuManager(w fyne.Window, tm *TabManager) *menuManager {
	m := &menuManager{
		window:     w,
		tabManager: tm,
	}
	m.loadRecentFiles()
	m.buildMenu()
	return m
}

func (m *menuManager) buildMenu() {
	openItem := fyne.NewMenuItem("Open...", func() {
		m.tabManager.showOpenFileDialogWithCallback(m.openedFileCallback)
	})

	quitItem := fyne.NewMenuItem("Quit", func() {
		_ = saveSession(m.tabManager.TabsData)
		m.saveRecentFiles()
		m.window.Close()
	})

	fileMenuItems := []*fyne.MenuItem{openItem}

	if len(m.recentFiles) > 0 {
		recentSubmenu := fyne.NewMenu("Recent", m.generateRecentMenuItems()...)
		fileMenuItems = append(fileMenuItems,
			fyne.NewMenuItemSeparator(),
			&fyne.MenuItem{
				Label:     "Recent Files",
				ChildMenu: recentSubmenu,
			},
		)
	}

	fileMenuItems = append(fileMenuItems,
		fyne.NewMenuItemSeparator(),
		quitItem,
	)

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File", fileMenuItems...),
	)
	m.window.SetMainMenu(mainMenu)
}

func (m *menuManager) generateRecentMenuItems() []*fyne.MenuItem {
	items := make([]*fyne.MenuItem, 0, len(m.recentFiles))
	for _, path := range m.recentFiles {
		p := path
		items = append(items, fyne.NewMenuItem(filepath.Base(p), func() {
			m.openedFileCallback(p)
		}))
	}
	return items
}

func (m *menuManager) openedFileCallback(path string) {
	content, err := readFileContent(path)
	if err != nil {
		dialog.ShowError(err, m.window)
		return
	}

	entry := m.tabManager.createEntry(filepath.Base(path), content)
	entry.Filepath = path
	tab := container.NewTabItem(entry.Title, container.NewStack(entry))

	m.tabManager.Tabs.Append(tab)
	m.tabManager.Tabs.Select(tab)
	m.tabManager.TabsData = append(m.tabManager.TabsData, &TabData{Entry: entry, Tab: tab})

	m.addRecentFile(path)
}

func (m *menuManager) addRecentFile(path string) {
	// Remove duplicate if exists
	m.recentFiles = slices.DeleteFunc(m.recentFiles, func(p string) bool {
		return p == path
	})

	// Prepend
	m.recentFiles = append([]string{path}, m.recentFiles...)
	if len(m.recentFiles) > 10 {
		m.recentFiles = m.recentFiles[:10]
	}
	m.buildMenu() // refresh menu
}

func (m *menuManager) saveRecentFiles() {
	data, err := json.MarshalIndent(m.recentFiles, "", "  ")
	if err != nil {
		fmt.Println("Failed to save recent files:", err)
		return
	}
	_ = os.WriteFile(recentFilePath, data, 0644)
}

func (m *menuManager) loadRecentFiles() {
	data, err := os.ReadFile(recentFilePath)
	if err != nil {
		if !os.IsNotExist(err) {
			fmt.Println("Failed to load recent files:", err)
		}
		return
	}
	_ = json.Unmarshal(data, &m.recentFiles)

	// Filter non-existing files
	m.recentFiles = slices.DeleteFunc(m.recentFiles, func(path string) bool {
		_, err := os.Stat(path)
		return err != nil && !os.IsNotExist(err)
	})
}
