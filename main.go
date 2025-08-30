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
	menuManager := newMenuManager(w, tabManager, []string{})
	windowState := defaultWindowState

	// Load session if it exists
	if savedTabs, savedWindowState, recentFiles, err := loadSession(); err == nil && len(savedTabs) > 0 {
		windowState = savedWindowState
		menuManager = newMenuManager(w, tabManager, recentFiles)
		tabManager.newTabs(savedTabs)
	} else {
		tabManager.newTab("", "") // default new tab
	}

	w.SetContent(container.NewBorder(nil, labelStatus, nil, nil, tabManager.Tabs))
	w.Resize(fyne.NewSize(windowState.Width, windowState.Height))

	// Save session on close
	w.SetCloseIntercept(func() {
		saveSession(tabManager.TabsData, menuManager.recentFiles, WindowState(w.Canvas().Size()))
		a.Quit()
	})

	w.ShowAndRun()
}
