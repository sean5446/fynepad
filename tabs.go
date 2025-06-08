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

func newTab(tabs *container.AppTabs, tabsData []*TabEntryWithShortcut, labelStatus *widget.Label, a fyne.App, tabTitle string, tabContent string) *TabEntryWithShortcut {
	if tabTitle == "" {
		tabTitle = "Untitled-" + strconv.Itoa(len(tabsData)+1)
	}

	var entry *TabEntryWithShortcut
	entry = assignShortcutAndData(func(shortcut fyne.Shortcut) {
		switch sc := shortcut.(type) {
		case *desktop.CustomShortcut:
			if sc.KeyName == fyne.KeyN && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+N
				tab := newTab(tabs, tabsData, labelStatus, a, "", "")
				tabsData = append(tabsData, tab)
			} else if sc.KeyName == fyne.KeyT && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+T
				tab := newTab(tabs, tabsData, labelStatus, a, "", "")
				tabsData = append(tabsData, tab)
			} else if sc.KeyName == fyne.KeyW && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+W
				closeCurrentTab(tabs, tabsData)
			} else if sc.KeyName == fyne.KeyMinus && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Minus
				if currentFontSize > 8 {
					currentFontSize -= 2
				}
				changeFontSize(a, currentFontSize, entry, labelStatus)
			} else if sc.KeyName == fyne.KeyEqual && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Plus
				currentFontSize += 2
				changeFontSize(a, currentFontSize, entry, labelStatus)
			} else if sc.KeyName == fyne.Key0 && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+Zero
				currentFontSize = defaultFontSize
				changeFontSize(a, currentFontSize, entry, labelStatus)
			} else if sc.KeyName == fyne.KeyZ && sc.Modifier == fyne.KeyModifierAlt {
				// Alt+Z
				toggleWrap(entry)
			} else if sc.KeyName == fyne.KeyS && sc.Modifier == fyne.KeyModifierControl {
				// Ctrl+S
				println("implement save")

				printCurrentTabText(tabs, tabsData) // try to debug why saving does not save recent changes

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

	applyTheme(a, currentFontSize)
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
	entry.Title = tabTitle
	entry.CursorColumn = 0 // TODO
	entry.CursorRow = 0    // TODO
	// entry.Wrapping = // TODO
	// TODO somehow set focus to the text
	entry.OnChanged = func(s string) {
		labelStatus.SetText(getLabelText(entry))
	}
	entry.OnCursorChanged = func() {
		labelStatus.SetText(getLabelText(entry))
	}
	entry.ExtendBaseWidget(entry)
	return entry
}

func getLabelText(entry *TabEntryWithShortcut) string {
	wrap := "on"
	if entry.Wrapping == fyne.TextWrapOff {
		wrap = "off"
	}
	return fmt.Sprintf("Ln: %d, Col: %d | %d characters | Font size: %.0fpx | Wrap: %s",
		entry.CursorRow+1, entry.CursorColumn+1, len(entry.Text), currentFontSize, wrap)
}

func toggleWrap(entry *TabEntryWithShortcut) {
	if entry.Wrapping == fyne.TextWrapOff {
		entry.Wrapping = fyne.TextWrapWord
	} else {
		entry.Wrapping = fyne.TextWrapOff
	}
	entry.Refresh()
}

func closeCurrentTab(tabs *container.AppTabs, tabsData []*TabEntryWithShortcut) {
	index, err := findCurrentTab(tabs, tabsData)
	if err != nil {
		println("Error finding current tab:", err.Error())
		return
	}
	tabs.RemoveIndex(index)
	tabsData = append(tabsData[:index], tabsData[index+1:]...)
}

// print text of current tab
func printCurrentTabText(tabs *container.AppTabs, tabsData []*TabEntryWithShortcut) {
	index, err := findCurrentTab(tabs, tabsData)
	if err != nil {
		println("Error finding current tab:", err.Error())
		return
	}
	if index < 0 || index >= len(tabsData) {
		println("No current tab selected")
		return
	}
	println("Current tab text:", tabsData[index].Entry.Text)
}

func findCurrentTab(tabs *container.AppTabs, tabsData []*TabEntryWithShortcut) (int, error) {
	index := tabs.SelectedIndex()
	if index < 0 || index >= len(tabsData) || index >= len(tabs.Items) {
		return 0, fmt.Errorf("no current tab selected")
	}
	return index, nil
}
