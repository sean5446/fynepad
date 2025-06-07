package main

import (
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type TabEntryWithShortcut struct {
	widget.Entry
	onShortcut func(shortcut fyne.Shortcut)
	Title      string
	Filepath string
}

func (e *TabEntryWithShortcut) TypedShortcut(shortcut fyne.Shortcut) {
	if e.onShortcut != nil {
		e.onShortcut(shortcut)
	}
}

func newTab(tabs *container.AppTabs, fontLabel *widget.Label, a fyne.App, tabTitle string, tabData string) *TabEntryWithShortcut {
	if tabTitle == "" {
		tabTitle = "Untitled-" + strconv.Itoa(len(tabsData)+1)
	}

	var entry *TabEntryWithShortcut
	entry = assignShortcutAndData(func(shortcut fyne.Shortcut) {
		switch sc := shortcut.(type) {
		case *desktop.CustomShortcut:
			if sc.KeyName == fyne.KeyN && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+N
				tab := newTab(tabs, fontLabel, a, "", "")
				tabsData = append(tabsData, tab)
			} else if sc.KeyName == fyne.KeyT && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+T
				tab := newTab(tabs, fontLabel, a, "", "")
				tabsData = append(tabsData, tab)
			} else if sc.KeyName == fyne.KeyW && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+W
				tabs.RemoveIndex(tabs.SelectedIndex())
				// tabsData = append(tabsData[:tabs.SelectedIndex()], tabsData[tabs.SelectedIndex()+1:]...)
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
				// Ctrl+Zero
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
	}, tabData, tabTitle)

	applyTheme(a, fontSize)
	tab := container.NewTabItem(tabTitle, container.NewStack(entry))
	tabs.Append(tab)
	tabs.Select(tab)
	return entry
}

func assignShortcutAndData(onShortcut func(fyne.Shortcut), tabData string, tabTitle string) *TabEntryWithShortcut {
	e := &TabEntryWithShortcut{onShortcut: onShortcut}
	e.MultiLine = true
	e.TextStyle = fyne.TextStyle{Monospace: true}
	e.Text = tabData
	e.ExtendBaseWidget(e)
	e.Title = tabTitle
	// more properties can be set here
	// e.Wrapping = fyne.TextWrapWord
	return e
}
