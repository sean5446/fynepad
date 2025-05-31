package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type HotkeyEntry struct {
	widget.Entry
	onShortcut func(shortcut fyne.Shortcut)
}

func NewHotkeyEntry(onShortcut func(fyne.Shortcut)) *HotkeyEntry {
	e := &HotkeyEntry{onShortcut: onShortcut}
	e.MultiLine = true
	e.Text = "some text"
	e.ExtendBaseWidget(e)
	// e.Wrapping = fyne.TextWrapWord
	return e
}

func (e *HotkeyEntry) TypedShortcut(shortcut fyne.Shortcut) {
	if e.onShortcut != nil {
		e.onShortcut(shortcut)
	}
}
