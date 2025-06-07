package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const defaultFontSize float32 = 14.0

var fontSize float32 = defaultFontSize
var tabsData []*TabEntryWithShortcut

func main() {
	a := app.New()
	w := a.NewWindow("GoPad")
	w.Resize(fyne.NewSize(800, 600))

	tabs := container.NewAppTabs()
	fontLabel := widget.NewLabel("")

	// load session data
	tabsData, _ = loadSessionData()
	for _, d := range tabsData {
		newTab(tabs, fontLabel, a, d.Title, d.Text)
	}

	// if no session data found, create a new tab
	if len(tabsData) == 0 {
		tab := newTab(tabs, fontLabel, a, "", "")
		tabsData = append(tabsData, tab)
	}

	// TODO implemenet last used files menu
	recentFiles := []string{
		"/home/user/notes1.txt",
		"/home/user/todo.md",
	}
	setupMenu(w, recentFiles)

	// setup the window content
	w.SetContent(container.NewBorder(nil, fontLabel, nil, nil, tabs))

	// save session data on close
	w.SetCloseIntercept(func() {
		println("Saving session data before closing ", tabsData)
		saveSessionData(tabsData)
		a.Quit()
	})

	w.ShowAndRun()
}

func changeFontSize(a fyne.App, fontSize float32, entry *TabEntryWithShortcut, fontLabel *widget.Label) {
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	applyTheme(a, fontSize)
	entry.Refresh()
	fontLabel.SetText(getLabelText(entry))
}
