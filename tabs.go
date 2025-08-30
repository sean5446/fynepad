package main

import (
	"fmt"
	"io"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type TabEntryWithShortcut struct {
	widget.Entry
	title      string
	filePath   string
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
	entry *TabEntryWithShortcut
	tab   *container.TabItem
}

type TabManager struct {
	app         fyne.App
	window      fyne.Window
	tabs        *container.AppTabs
	labelStatus *widget.Label
	tabsData    []*TabData
	fontSize    float32
	defaultSize float32
	menuManager *MenuManager
}

func newTabManager(app fyne.App, w fyne.Window, labelStatus *widget.Label, defaultFontSize float32) *TabManager {
	tabs := container.NewAppTabs()

	tm := &TabManager{
		app:         app,
		window:      w,
		tabs:        tabs,
		labelStatus: labelStatus,
		fontSize:    defaultFontSize,
		defaultSize: defaultFontSize,
	}

	tabs.OnSelected = func(tab *container.TabItem) {
		for _, data := range tm.tabsData {
			if data.tab == tab {
				tm.focusEntry(data.entry)
				break
			}
		}
	}

	return tm
}

func (tm *TabManager) newTab(title, content string) *TabEntryWithShortcut {
	if title == "" {
		title = "Untitled-" + strconv.Itoa(len(tm.tabsData)+1)
	}

	entry := tm.createEntry(title, content)
	tab := container.NewTabItem(title, container.NewStack(entry))

	tm.tabs.Append(tab)
	tm.tabs.Select(tab)

	tm.tabsData = append(tm.tabsData, &TabData{entry, tab})

	tm.applyFontSize(entry)
	tm.focusEntry(entry)
	return entry
}

func (tm *TabManager) newTabs(savedTabs []TabDetail, selectedIndex int) {
	for _, s := range savedTabs {
		entry := tm.createEntry(s.Title, s.Text)
		entry.Wrapping = fyne.TextWrap(s.Wrapping)
		entry.CursorRow = s.CursorRow
		entry.CursorColumn = s.CursorColumn
		entry.filePath = s.FilePath
		tab := container.NewTabItem(s.Title, container.NewStack(entry))
		tm.tabs.Append(tab)
		tm.tabsData = append(tm.tabsData, &TabData{entry, tab})
	}

	if len(tm.tabsData) == 0 {
		tm.newTab("", "") // default new tab
	}

	if selectedIndex >= 0 && selectedIndex < len(tm.tabsData) {
		tm.tabs.SelectIndex(selectedIndex)
		tm.focusEntry(tm.tabsData[selectedIndex].entry)
	}
}

func (tm *TabManager) createEntry(title, content string) *TabEntryWithShortcut {
	entry := &TabEntryWithShortcut{
		title: title,
	}
	entry.ExtendBaseWidget(entry)

	entry.MultiLine = true
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.Text = content

	entry.OnChanged = func(s string) {
		tm.labelStatus.SetText(tm.getLabelText(entry))
	}

	entry.OnCursorChanged = func() {
		tm.labelStatus.SetText(tm.getLabelText(entry))
	}

	entry.onShortcut = func(shortcut fyne.Shortcut) {
		tm.handleShortcut(entry, shortcut)
	}

	return entry
}

func (tm *TabManager) getCurrentTabIndex() (int, error) {
	index := tm.tabs.SelectedIndex()
	if index < 0 || index >= len(tm.tabsData) {
		return 0, fmt.Errorf("invalid tab index")
	}
	return index, nil
}

func (tm *TabManager) getCurrentEntry() (*TabEntryWithShortcut, error) {
	index, err := tm.getCurrentTabIndex()
	if err != nil {
		return nil, err
	}
	return tm.tabsData[index].entry, nil
}

func (tm *TabManager) closeCurrentTab() {
	index, err := tm.getCurrentTabIndex()
	if err != nil {
		fmt.Println("Error closing tab:", err)
		return
	}

	// if we don't have filepath set, ask to save
	entry := tm.tabsData[index].entry
	if entry.filePath == "" && entry.Text != "" {
		confirm := dialog.NewConfirm("Close Tab", "Do you want to save changes to "+entry.title+"?", func(save bool) {
			if save {
				tm.showSaveFileDialog(entry)
			} else {
				tm.tabs.RemoveIndex(index)
				tm.tabsData = append(tm.tabsData[:index], tm.tabsData[index+1:]...)
			}
		}, tm.window)
		confirm.SetDismissText("Don't Save")
		confirm.SetConfirmText("Save")
		confirm.Show()
	} else {
		tm.tabs.RemoveIndex(index)
		tm.tabsData = append(tm.tabsData[:index], tm.tabsData[index+1:]...)
	}
}

func (tm *TabManager) printCurrentTabText() {
	entry, err := tm.getCurrentEntry()
	if err != nil {
		fmt.Println("Error finding current tab:", err)
		return
	}
	fmt.Println(entry.Text)
}

func (tm *TabManager) toggleWrap(entry *TabEntryWithShortcut) {
	if entry.Wrapping == fyne.TextWrapOff {
		entry.Wrapping = fyne.TextWrapWord
	} else {
		entry.Wrapping = fyne.TextWrapOff
	}
	entry.Refresh()
}

func (tm *TabManager) applyFontSize(entry *TabEntryWithShortcut) {
	changeFontSize(tm.app, tm.fontSize, entry)
	tm.labelStatus.SetText(tm.getLabelText(entry))
}

func (tm *TabManager) getLabelText(entry *TabEntryWithShortcut) string {
	wrap := "on"
	if entry.Wrapping == fyne.TextWrapOff {
		wrap = "off"
	}
	return fmt.Sprintf("Ln: %d, Col: %d | %d characters | Font size: %.0fpx | Wrap: %s",
		entry.CursorRow+1, entry.CursorColumn+1, len(entry.Text), tm.fontSize, wrap)
}

func (tm *TabManager) showOpenFileDialog() {
	openDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil || reader == nil {
			return
		}
		defer reader.Close()

		data, err := io.ReadAll(reader)
		if err != nil {
			dialog.ShowError(err, tm.app.Driver().AllWindows()[0])
			return
		}

		entry := tm.newTab(reader.URI().Name(), string(data))
		entry.filePath = reader.URI().Path()

	}, tm.app.Driver().AllWindows()[0])

	openDialog.SetFilter(storage.NewExtensionFileFilter([]string{".txt", ".md", ""}))
	openDialog.Show()
}

func (tm *TabManager) saveCurrentFile() {
	entry, err := tm.getCurrentEntry()
	if err != nil {
		return
	}

	if entry.filePath != "" {
		err := writeFileContent(entry.filePath, entry.Text)
		if err != nil {
			dialog.ShowError(err, tm.app.Driver().AllWindows()[0])
		}
		return
	}

	// Fallback: show Save As
	tm.showSaveFileDialog(entry)
}

func (tm *TabManager) showSaveFileDialog(entry *TabEntryWithShortcut) {
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil || writer == nil {
			return
		}
		defer writer.Close()

		_, err = writer.Write([]byte(entry.Text))
		if err != nil {
			dialog.ShowError(err, tm.app.Driver().AllWindows()[0])
			return
		}

		entry.filePath = writer.URI().Path()
		entry.title = writer.URI().Name()

		// Update tab title
		index, _ := tm.getCurrentTabIndex()
		tm.tabs.Items[index].Text = entry.title
		tm.tabs.Refresh()
	}, tm.app.Driver().AllWindows()[0])

	saveDialog.SetFileName(entry.title + ".txt")
	saveDialog.Show()
}

func (tm *TabManager) handleShortcut(entry *TabEntryWithShortcut, shortcut fyne.Shortcut) {
	switch sc := shortcut.(type) {
	case *desktop.CustomShortcut:
		switch {
		case sc.KeyName == fyne.KeyN && sc.Modifier == fyne.KeyModifierControl:
			tm.newTab("", "")
		case sc.KeyName == fyne.KeyT && sc.Modifier == fyne.KeyModifierControl:
			tm.newTab("", "")
		case sc.KeyName == fyne.KeyW && sc.Modifier == fyne.KeyModifierControl:
			tm.closeCurrentTab()
		case sc.KeyName == fyne.KeyMinus && sc.Modifier == fyne.KeyModifierControl:
			if tm.fontSize > 8 {
				tm.fontSize -= 2
			}
			tm.applyFontSize(entry)
		case sc.KeyName == fyne.KeyEqual && sc.Modifier == fyne.KeyModifierControl:
			tm.fontSize += 2
			tm.applyFontSize(entry)
		case sc.KeyName == fyne.Key0 && sc.Modifier == fyne.KeyModifierControl:
			tm.fontSize = tm.defaultSize
			tm.applyFontSize(entry)
		case sc.KeyName == fyne.KeyZ && sc.Modifier == fyne.KeyModifierAlt:
			tm.toggleWrap(entry)
		case sc.KeyName == fyne.KeyD && sc.Modifier == fyne.KeyModifierControl:
			tm.printCurrentTabText()
		case sc.KeyName == fyne.KeyO && sc.Modifier == fyne.KeyModifierControl:
			tm.showOpenFileDialog()
		case sc.KeyName == fyne.KeyS && sc.Modifier == fyne.KeyModifierControl:
			tm.saveCurrentFile()
		case sc.KeyName == fyne.KeyQ && sc.Modifier == fyne.KeyModifierControl:
			saveSession(tm, tm.menuManager, tm.window)
			tm.app.Quit()
		default:
			fmt.Println("Unhandled shortcut:", sc)
		}
	}
}

func (tm *TabManager) focusEntry(entry *TabEntryWithShortcut) {
	// if this is called immediately after creating the entry, it may not focus correctly
	// https://github.com/fyne-io/fyne/wiki/Threading-migration
	time.AfterFunc(50*time.Millisecond, func() {
		fyne.Do(func() {
			canvas := tm.window.Canvas()
			if canvas != nil {
				canvas.Focus(entry)
			}
		})
	})
}
