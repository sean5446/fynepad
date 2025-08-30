package main

import (
	"path/filepath"
	"slices"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
)

type MenuManager struct {
	app         fyne.App
	window      fyne.Window
	recentFiles []string
	tabManager  *TabManager
}

func newMenuManager(a fyne.App, w fyne.Window, rf []string) *MenuManager {
	m := &MenuManager{
		app:         a,
		window:      w,
		recentFiles: rf,
	}
	m.loadRecentFiles(rf)
	m.buildMenu()
	return m
}

func (m *MenuManager) buildMenu() {
	openItem := fyne.NewMenuItem("Open...", func() {
		m.showOpenFileDialogWithCallback(m.openedFileCallback)
	})

	quitItem := fyne.NewMenuItem("Quit", func() {
		saveSession(m.tabManager, m, m.window)
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

func (m *MenuManager) showOpenFileDialogWithCallback(callback func(path string)) {
	dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		if callback != nil {
			callback(reader.URI().Path())
		}
	}, m.app.Driver().AllWindows()[0]).Show()
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
	entry.filePath = path
	tab := container.NewTabItem(entry.title, container.NewStack(entry))

	m.tabManager.tabs.Append(tab)
	m.tabManager.tabs.Select(tab)
	m.tabManager.tabsData = append(m.tabManager.tabsData, &TabData{entry: entry, tab: tab})

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

func (m *MenuManager) loadRecentFiles(recentFiles []string) {
	m.recentFiles = recentFiles

	// Filter non-existing files
	m.recentFiles = slices.DeleteFunc(m.recentFiles, func(path string) bool {
		_, err := readFileContent(path)
		return err != nil
	})
}
