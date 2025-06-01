package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const defaultFontSize float32 = 14.0

var fontSize float32 = defaultFontSize
var tabsData []TabData

func main() {
	a := app.New()
	w := a.NewWindow("GoPad")
	w.Resize(fyne.NewSize(800, 600))

	tabs := container.NewAppTabs()
	fontLabel := widget.NewLabel(getLabelText())

	// load session data
	tabsData, _ = loadSessionData()
	for _, d := range tabsData {
		newTab(tabs, fontLabel, a, d.Title, d.Content)
	}
	if len(tabs.Items) == 0 {
		newTab(tabs, fontLabel, a, "", "")
	}

	recentFiles := []string{
		"/home/user/notes1.txt",
		"/home/user/todo.md",
	}

	setupMenu(w, recentFiles)

	// setup the window content
	w.SetContent(container.NewBorder(nil, fontLabel, nil, nil, tabs))

	// save session data on close
	w.SetCloseIntercept(func() {
		saveSessionData(tabsData)
		a.Quit()
	})

	w.ShowAndRun()
}

func changeFontSize(a fyne.App, fontSize float32, entry *HotkeyEntry, fontLabel *widget.Label) {
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	applyTheme(a, fontSize)
	entry.Refresh()
	fontLabel.SetText(getLabelText())
}

func getLabelText() string {
	return "Font Size: " + strconv.Itoa(int(fontSize))
}
