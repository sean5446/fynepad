package main

import (
	"path/filepath"
	"fyne.io/fyne/v2"
)

func setupMenu(myWindow fyne.Window, recentFiles []string) {
	fileMenuItems := []*fyne.MenuItem{
		fyne.NewMenuItem("Open...", func() {
			println("Open clicked")
			// implement open logic
		}),
	}

	openRecent := func(file string) {
		println("Open recent file:", file)
	}

	// Create submenu for "Recent Files"
	if len(recentFiles) > 0 {
		recentSubmenu := fyne.NewMenu("Recent", generateRecentMenuItems(openRecent, recentFiles)...)
		fileMenuItems = append(fileMenuItems, fyne.NewMenuItemSeparator())
		fileMenuItems = append(fileMenuItems, &fyne.MenuItem{
			Label:     "Recent Files",
			Action:    nil,
			ChildMenu: recentSubmenu,
		})
	}

	fileMenuItems = append(fileMenuItems, fyne.NewMenuItemSeparator())
	fileMenuItems = append(fileMenuItems, fyne.NewMenuItem("Quit", func() {
		myWindow.Close()
	}))

	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("File", fileMenuItems...),
	)
	myWindow.SetMainMenu(mainMenu)
}

func generateRecentMenuItems(openRecent func(file string), files []string) []*fyne.MenuItem {
	items := make([]*fyne.MenuItem, 0, len(files))
	for _, file := range files {
		f := file // capture loop variable
		items = append(items, fyne.NewMenuItem(filepath.Base(f), func() {
			openRecent(f)
		}))
	}
	return items
}
