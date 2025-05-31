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
	fontLabel := widget.NewLabel(getLabelText())

	recentFiles := []string{
		"/home/user/notes1.txt",
		"/home/user/todo.md",
	}

	SetupMenu(w, recentFiles)

	// Add initial tab
	newTab(tabs, fontLabel, a)

	w.SetContent(container.NewBorder(nil, fontLabel, nil, nil, tabs))

	w.ShowAndRun()
}

func changeFontSize(a fyne.App, fontSize float32, entry *HotkeyEntry, fontLabel *widget.Label) {
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	ApplyTheme(a, fontSize)
	entry.Refresh()
	fontLabel.SetText(getLabelText())
}

func getLabelText() string {
	return "Font Size: " + strconv.Itoa(int(fontSize))
}

func newTab(tabs *container.AppTabs, fontLabel *widget.Label, a fyne.App) {
	var entry *HotkeyEntry

	entry = NewHotkeyEntry(func(shortcut fyne.Shortcut) {
		switch sc := shortcut.(type) {
		case *desktop.CustomShortcut:
			if sc.KeyName == fyne.KeyN && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+N
				newTab(tabs, fontLabel, a)
			} else if sc.KeyName == fyne.KeyW && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+W
				tabs.RemoveIndex(tabs.SelectedIndex())
			} else if sc.KeyName == fyne.KeyMinus && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Minus
				if fontSize > 6 {
					fontSize -= 2
				}
				changeFontSize(a, fontSize, entry, fontLabel)
			} else if sc.KeyName == fyne.KeyEqual && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Plus
				fontSize += 2
				changeFontSize(a, fontSize, entry, fontLabel)
			} else if sc.KeyName == fyne.Key0 && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+0
				fontSize = defaultFontSize
				changeFontSize(a, fontSize, entry, fontLabel)
			} else if sc.KeyName == fyne.KeyS && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+S
				println("implement save")
			} else if sc.KeyName == fyne.KeyO && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+O
				println("implemenet open")
			} else if sc.KeyName == fyne.KeyF && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+F
				println("implement find")
			} else if sc.KeyName == fyne.KeyQ && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Q
				a.Quit()
			}
		}
	})

	ApplyTheme(a, fontSize)

	tabName := "Untitled-" + strconv.Itoa(len(tabs.Items)+1)
	tab := container.NewTabItem(tabName, container.NewStack(entry))
	tabs.Append(tab)
	tabs.Select(tab)
}
