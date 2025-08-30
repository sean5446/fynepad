package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.queso.gopad")
	w := a.NewWindow("GoPad")
	labelStatus := widget.NewLabel("")

	tabManager := newTabManager(a, w, labelStatus, 14)
	menuManager := newMenuManager(a, w, []string{})

	tabManager.menuManager = menuManager
	menuManager.tabManager = tabManager

	windowState := defaultWindowState

	// Load session if it exists
	if savedTabs, savedWindowState, recentFiles, err := loadSession(); err == nil {
		windowState = savedWindowState
		menuManager = newMenuManager(a, w, recentFiles)

		menuManager.tabManager = tabManager
		tabManager.menuManager = menuManager

		tabManager.newTabs(savedTabs, savedWindowState.TabSelected)
	} else {
		tabManager.newTab("", "") // default new tab
	}

	w.SetContent(container.NewBorder(nil, labelStatus, nil, nil, tabManager.tabs))
	w.Resize(fyne.NewSize(windowState.Width, windowState.Height))

	// Save session on close
	w.SetCloseIntercept(func() {
		saveSession(tabManager, menuManager, w)
		a.Quit()
	})

	w.ShowAndRun()
}
