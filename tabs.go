package main

import (
	"fmt"
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
	Filepath   string
}

func (e *TabEntryWithShortcut) TypedShortcut(shortcut fyne.Shortcut) {
	// Handle shortcuts for the entry - allow default shortcuts like copy, paste, etc.
	switch shortcut.(type) {
	case *fyne.ShortcutCopy,
		*fyne.ShortcutPaste,
		*fyne.ShortcutCut,
		*fyne.ShortcutSelectAll,
		*fyne.ShortcutUndo,
		*fyne.ShortcutRedo:
		e.Entry.TypedShortcut(shortcut)
	default:
		if e.onShortcut != nil {
			e.onShortcut(shortcut)
		}
	}
}

func newTab(tabs *container.AppTabs, labelStatus *widget.Label, a fyne.App, tabTitle string, tabContent string) *TabEntryWithShortcut {
	if tabTitle == "" {
		tabTitle = "Untitled-" + strconv.Itoa(len(tabsData)+1)
	}

	var entry *TabEntryWithShortcut
	entry = assignShortcutAndData(func(shortcut fyne.Shortcut) {
		switch sc := shortcut.(type) {
		case *desktop.CustomShortcut:
			if sc.KeyName == fyne.KeyN && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+N
				tab := newTab(tabs, labelStatus, a, "", "")
				tabsData = append(tabsData, tab)
			} else if sc.KeyName == fyne.KeyT && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+T
				tab := newTab(tabs, labelStatus, a, "", "")
				tabsData = append(tabsData, tab)
			} else if sc.KeyName == fyne.KeyW && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+W
				tabs.RemoveIndex(tabs.SelectedIndex())
				// TODO do something better here - this crashes
				if len(tabsData) > 0 {
					tabsData = append(tabsData[:tabs.SelectedIndex()], tabsData[tabs.SelectedIndex()+1:]...)
				}
			} else if sc.KeyName == fyne.KeyMinus && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Minus
				if fontSize > 8 {
					fontSize -= 2
				}
				changeFontSize(a, fontSize, entry, labelStatus)
			} else if sc.KeyName == fyne.KeyEqual && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Plus
				fontSize += 2
				changeFontSize(a, fontSize, entry, labelStatus)
			} else if sc.KeyName == fyne.Key0 && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Zero
				fontSize = defaultFontSize
				changeFontSize(a, fontSize, entry, labelStatus)
			} else if sc.KeyName == fyne.KeyZ && sc.Modifier == fyne.KeyModifierAlt {
				// Alt+Z
				toggleWrap(entry)
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
				saveSessionData(tabsData)
				a.Quit()
			}
		}
	}, labelStatus, tabTitle, tabContent)

	applyTheme(a, fontSize)
	tab := container.NewTabItem(tabTitle, container.NewStack(entry))
	tabs.Append(tab)
	tabs.Select(tab)
	return entry
}

func assignShortcutAndData(onShortcut func(fyne.Shortcut), labelStatus *widget.Label, tabTitle string, tabContent string) *TabEntryWithShortcut {
	entry := &TabEntryWithShortcut{onShortcut: onShortcut}
	entry.MultiLine = true
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.Text = tabContent
	entry.ExtendBaseWidget(entry)
	entry.Title = tabTitle
	// more properties can be set here
	entry.OnChanged = func(s string) {
		labelStatus.SetText(getLabelText(entry))
	}
	entry.OnCursorChanged = func() {
		labelStatus.SetText(getLabelText(entry))
	}
	return entry
}

func getLabelText(entry *TabEntryWithShortcut) string {
	wrap := "on"
	if entry.Wrapping == fyne.TextWrapOff {
		wrap = "off"
	}
	return fmt.Sprintf("Ln: %d, Col: %d | %d characters | Font size: %.0fpx | Wrap: %s",
		entry.CursorRow+1, entry.CursorColumn+1, len(entry.Text), fontSize, wrap)
}

func toggleWrap(entry *TabEntryWithShortcut) {
	if entry.Wrapping == fyne.TextWrapOff {
		entry.Wrapping = fyne.TextWrapWord
	} else {
		entry.Wrapping = fyne.TextWrapOff
	}
	entry.Refresh()
}
