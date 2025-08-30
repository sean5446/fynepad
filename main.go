package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	app := app.NewWithID("com.queso.gopad")
	window := app.NewWindow("GoPad")
	labelStatus := widget.NewLabel("")

	tabManager := newTabManager(app, window, labelStatus, 14)
	menuManager := newMenuManager(app, window, []string{})

	tabManager.menuManager = menuManager
	menuManager.tabManager = tabManager

	windowState := defaultWindowState

	// load session if it exists
	if savedTabs, savedWindowState, recentFiles, err := loadSession(); err == nil {
		windowState = savedWindowState
		menuManager = newMenuManager(app, window, recentFiles)

		menuManager.tabManager = tabManager
		tabManager.menuManager = menuManager

		tabManager.newTabs(savedTabs, savedWindowState.TabSelected)
	} else {
		tabManager.newTab("", "") // default new tab
	}

	window.SetContent(container.NewBorder(nil, labelStatus, nil, nil, tabManager.tabs))
	window.Resize(fyne.NewSize(windowState.Width, windowState.Height))
	window.SetCloseIntercept(func() {
		saveSession(tabManager, menuManager, window)
		app.Quit()
	})

	window.ShowAndRun()
}
