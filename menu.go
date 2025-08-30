package main

import (
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)


type MenuManager struct {
	window      fyne.Window
	tabManager  *TabManager
	recentFiles []string
}

func newMenuManager(w fyne.Window, tm *TabManager, rf []string) *MenuManager {
	m := &MenuManager{
		window:     w,
		tabManager: tm,
		recentFiles: rf,
	}
	m.loadRecentFiles()
	m.buildMenu()
	return m
}

func (m *MenuManager) buildMenu() {
	openItem := fyne.NewMenuItem("Open...", func() {
		m.tabManager.showOpenFileDialogWithCallback(m.openedFileCallback)
	})

	quitItem := fyne.NewMenuItem("Quit", func() {
		_ = saveSession(m.tabManager.TabsData, m.recentFiles, WindowState(m.window.Canvas().Size()))
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

func (m *MenuManager) generateRecentMenuItems() []*fyne.MenuItem {
	items := make([]*fyne.MenuItem, 0, len(m.recentFiles))
	for _, path := range m.recentFiles {
		p := path
		items = append(items, fyne.NewMenuItem(filepath.Base(p), func() {
			m.openedFileCallback(p)
		}))
	}
	return items
}

func (m *MenuManager) openedFileCallback(path string) {
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

func (m *MenuManager) addRecentFile(path string) {
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

func (m *MenuManager) loadRecentFiles() {
	_, _, recentFiles, err := loadSession()
	if err != nil {
		return
	}

	m.recentFiles = recentFiles

	// Filter non-existing files
	m.recentFiles = slices.DeleteFunc(m.recentFiles, func(path string) bool {
		_, err := readFileContent(path)
		return err != nil
	})
}
