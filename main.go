package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

const defaultFontSize float32 = 14.0
var fontSize float32 = defaultFontSize


func main() {
	a := app.New()
	w := a.NewWindow("GoPad")
	w.Resize(fyne.NewSize(800, 600))

	tabs := container.NewAppTabs()
	fontLabel := widget.NewLabel("Font Size: " + strconv.Itoa(int(fontSize)))

	recentFiles := []string{
		"/home/user/notes1.txt",
		"/home/user/todo.md",
	}

	SetupMenu(w, func(file string) {
		println("Opening recent file:", file)
		// loadFile(file)
	}, recentFiles)

	// Add initial tab
	newTab(tabs, fontLabel, a)

	w.SetContent(container.NewBorder(nil, fontLabel, nil, nil, tabs))

	w.ShowAndRun()
}

func changeFontSize(a fyne.App, fontSize float32, entry *HotkeyEntry, fontLabel *widget.Label) {
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	ApplyTheme(a, fontSize)
	entry.Refresh()
	fontLabel.SetText("Font Size: " + strconv.Itoa(int(fontSize)))
}

func newTab(tabs *container.AppTabs, fontLabel *widget.Label, a fyne.App) {
	var entry *HotkeyEntry

	entry = NewHotkeyEntry(func(shortcut fyne.Shortcut) {
		switch sc := shortcut.(type) {
		case *desktop.CustomShortcut:
			// Check for combinations
			if sc.KeyName == fyne.KeyN && sc.Modifier == fyne.KeyModifierControl {
				newTab(tabs, fontLabel, a)
			} else if sc.KeyName == fyne.KeyW && sc.Modifier == fyne.KeyModifierControl {
				tabs.RemoveIndex(tabs.SelectedIndex())
			} else if sc.KeyName == fyne.KeyMinus && sc.Modifier == fyne.KeyModifierControl {
				if fontSize > 6 {
					fontSize -= 2
				}
				changeFontSize(a, fontSize, entry, fontLabel)
			} else if sc.KeyName == fyne.KeyEqual && sc.Modifier == fyne.KeyModifierControl {
				fontSize += 2
				changeFontSize(a, fontSize, entry, fontLabel)
			} else if sc.KeyName == fyne.Key0 && sc.Modifier == fyne.KeyModifierControl {
				fontSize = defaultFontSize
				changeFontSize(a, fontSize, entry, fontLabel)
			}
		}
	})

	ApplyTheme(a, fontSize)

	tab := container.NewTabItem("Tab", container.NewStack(entry))
	tabs.Append(tab)
	tabs.Select(tab)
}
