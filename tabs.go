package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type HotkeyEntry struct {
	widget.Entry
	onShortcut func(shortcut fyne.Shortcut)
}

func (e *HotkeyEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if e.onShortcut != nil {
		e.onShortcut(shortcut)
	}
}

func newHotkeyEntry(onShortcut func(fyne.Shortcut)) *HotkeyEntry {
	e := &HotkeyEntry{onShortcut: onShortcut}
	e.MultiLine = true
	e.Text = "some text"
	e.ExtendBaseWidget(e)
	// e.Wrapping = fyne.TextWrapWord
	return e
}

func newTab(tabs *container.AppTabs, fontLabel *widget.Label, a fyne.App) {
	var entry *HotkeyEntry

	entry = newHotkeyEntry(func(shortcut fyne.Shortcut) {
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

	applyTheme(a, fontSize)

	tabName := "Untitled-" + strconv.Itoa(len(tabs.Items)+1)
	tab := container.NewTabItem(tabName, container.NewStack(entry))
	tabs.Append(tab)
	tabs.Select(tab)
}
