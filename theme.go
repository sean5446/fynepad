package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type customTheme struct {
	fyne.Theme
	textSize float32
}

func (m customTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return m.textSize
	}
	return theme.DefaultTheme().Size(name)
}

func (m customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (m customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (m customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func applyTheme(a fyne.App, size float32) {
	a.Settings().SetTheme(&customTheme{textSize: size})
}

func changeFontSize(a fyne.App, fontSize float32, entry *TabEntryWithShortcut) {
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	applyTheme(a, fontSize)
	entry.Refresh()
}
