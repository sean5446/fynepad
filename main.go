package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("com.queso.mytexteditor")
	w := a.NewWindow("GoPad")
	labelStatus := widget.NewLabel("")

	tabManager := NewTabManager(a, labelStatus, 14)

	menuManager := newMenuManager(w, tabManager)

	// Load session if it exists
	if savedTabs, err := loadSession(); err == nil && len(savedTabs) > 0 {
		for _, s := range savedTabs {
			entry := tabManager.createEntry(s.Title, s.Text)
			entry.Wrapping = s.Wrapping
			entry.CursorRow = s.CursorRow
			entry.CursorColumn = s.CursorColumn
			entry.Filepath = s.Filepath
			tab := container.NewTabItem(s.Title, container.NewStack(entry))
			tabManager.Tabs.Append(tab)
			tabManager.TabsData = append(tabManager.TabsData, &TabData{Entry: entry, Tab: tab})
		}
		tabManager.Tabs.SelectIndex(0)
	} else {
		tabManager.NewTab("", "") // default new tab
	}

	w.SetContent(container.NewBorder(nil, labelStatus, nil, nil, tabManager.Tabs))
	w.Resize(fyne.NewSize(800, 600))

	// Save session on close
	w.SetCloseIntercept(func() {
		saveSession(tabManager.TabsData)
		menuManager.saveRecentFiles()
		a.Quit()
	})

	w.ShowAndRun()
}
