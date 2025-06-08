package main

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

// -----------------------------
// Structs
// -----------------------------

type TabEntryWithShortcut struct {
	widget.Entry
	Title      string
	Filepath   string
	onShortcut func(fyne.Shortcut)
}

func (e *TabEntryWithShortcut) TypedShortcut(shortcut fyne.Shortcut) {
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

type TabData struct {
	Entry *TabEntryWithShortcut
	Tab   *container.TabItem
}

type TabManager struct {
	App         fyne.App
	Tabs        *container.AppTabs
	LabelStatus *widget.Label
	TabsData    []*TabData
	FontSize    float32
	DefaultSize float32
}

func NewTabManager(app fyne.App, labelStatus *widget.Label, defaultFontSize float32) *TabManager {
	return &TabManager{
		App:         app,
		Tabs:        container.NewAppTabs(),
		LabelStatus: labelStatus,
		FontSize:    defaultFontSize,
		DefaultSize: defaultFontSize,
	}
}

// -----------------------------
// Core Tab Creation
// -----------------------------

func (tm *TabManager) NewTab(title, content string) {
	if title == "" {
		title = "Untitled-" + strconv.Itoa(len(tm.TabsData)+1)
	}

	entry := tm.createEntry(title, content)
	tab := container.NewTabItem(title, container.NewStack(entry))

	tm.Tabs.Append(tab)
	tm.Tabs.Select(tab)

	tm.TabsData = append(tm.TabsData, &TabData{
		Entry: entry,
		Tab:   tab,
	})

	tm.applyFontSize(entry)
}

func (tm *TabManager) createEntry(title, content string) *TabEntryWithShortcut {
	entry := &TabEntryWithShortcut{
		Title: title,
	}
	entry.ExtendBaseWidget(entry)

	entry.MultiLine = true
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.Text = content

	entry.OnChanged = func(s string) {
		tm.LabelStatus.SetText(tm.getLabelText(entry))
	}
	entry.OnCursorChanged = func() {
		tm.LabelStatus.SetText(tm.getLabelText(entry))
	}

	entry.onShortcut = func(shortcut fyne.Shortcut) {
		tm.handleShortcut(entry, shortcut)
	}

	return entry
}

// -----------------------------
// Helper Methods
// -----------------------------

func (tm *TabManager) getCurrentTabIndex() (int, error) {
	index := tm.Tabs.SelectedIndex()
	if index < 0 || index >= len(tm.TabsData) {
		return 0, fmt.Errorf("invalid tab index")
	}
	return index, nil
}

func (tm *TabManager) getCurrentEntry() (*TabEntryWithShortcut, error) {
	index, err := tm.getCurrentTabIndex()
	if err != nil {
		return nil, err
	}
	return tm.TabsData[index].Entry, nil
}

func (tm *TabManager) CloseCurrentTab() {
	index, err := tm.getCurrentTabIndex()
	if err != nil {
		fmt.Println("Error closing tab:", err)
		return
	}

	tm.Tabs.RemoveIndex(index)
	tm.TabsData = append(tm.TabsData[:index], tm.TabsData[index+1:]...)
}

func (tm *TabManager) PrintCurrentTabText() {
	entry, err := tm.getCurrentEntry()
	if err != nil {
		fmt.Println("Error finding current tab:", err)
		return
	}
	fmt.Println("Text:", entry.Text)
}

func (tm *TabManager) ToggleWrap(entry *TabEntryWithShortcut) {
	if entry.Wrapping == fyne.TextWrapOff {
		entry.Wrapping = fyne.TextWrapWord
	} else {
		entry.Wrapping = fyne.TextWrapOff
	}
	entry.Refresh()
}

func (tm *TabManager) applyFontSize(entry *TabEntryWithShortcut) {
	changeFontSize(tm.App, tm.FontSize, entry)
	tm.LabelStatus.SetText(tm.getLabelText(entry))
}

func (tm *TabManager) getLabelText(entry *TabEntryWithShortcut) string {
	wrap := "on"
	if entry.Wrapping == fyne.TextWrapOff {
		wrap = "off"
	}
	return fmt.Sprintf("Ln: %d, Col: %d | %d characters | Font size: %.0fpx | Wrap: %s",
		entry.CursorRow+1, entry.CursorColumn+1, len(entry.Text), tm.FontSize, wrap)
}

// -----------------------------
// Shortcut Handling
// -----------------------------

func (tm *TabManager) handleShortcut(entry *TabEntryWithShortcut, shortcut fyne.Shortcut) {
	switch sc := shortcut.(type) {
	case *desktop.CustomShortcut:
		switch {
		case sc.KeyName == fyne.KeyN && sc.Modifier == fyne.KeyModifierControl:
			tm.NewTab("", "")
		case sc.KeyName == fyne.KeyW && sc.Modifier == fyne.KeyModifierControl:
			tm.CloseCurrentTab()
		case sc.KeyName == fyne.KeyMinus && sc.Modifier == fyne.KeyModifierControl:
			if tm.FontSize > 8 {
				tm.FontSize -= 2
			}
			tm.applyFontSize(entry)
		case sc.KeyName == fyne.KeyEqual && sc.Modifier == fyne.KeyModifierControl:
			tm.FontSize += 2
			tm.applyFontSize(entry)
		case sc.KeyName == fyne.Key0 && sc.Modifier == fyne.KeyModifierControl:
			tm.FontSize = tm.DefaultSize
			tm.applyFontSize(entry)
		case sc.KeyName == fyne.KeyZ && sc.Modifier == fyne.KeyModifierAlt:
			tm.ToggleWrap(entry)
		case sc.KeyName == fyne.KeyS && sc.Modifier == fyne.KeyModifierControl:
			fmt.Println("Save triggered")
			tm.PrintCurrentTabText()
		case sc.KeyName == fyne.KeyQ && sc.Modifier == fyne.KeyModifierControl:
			SaveSession(tm.TabsData)
			tm.App.Quit()
		default:
			fmt.Println("Unhandled shortcut:", sc)
		}
	}
}
