package main

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type TabData struct {
	Title    string `json:"title"`
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

const sessionFile = "session.json"
const defaultFontSize float32 = 14

var tabCount int = 1
var fontSize float32 = defaultFontSize


func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("GoPad")
	myWindow.Resize(fyne.NewSize(800, 600))

	tabs := container.NewAppTabs()
	tabMap := make(map[*container.TabItem]*widget.Entry)

	loadSession(tabs, tabMap)
	myWindow.SetContent(tabs)

	registerShortcuts(tabs, tabMap, myWindow, myApp)

	myWindow.SetCloseIntercept(func() {
		saveSession(tabs, tabMap)
		myWindow.Close()
	})

	myWindow.ShowAndRun()
}

func applyFontSize(app fyne.App, size float32) {
	app.Settings().SetTheme(&customTheme{fontSize: size})
}

func makeNewTab(title string, entry *widget.Entry, tab *container.TabItem) *container.TabItem {
	tab.Content = container.NewBorder(nil, nil, nil, nil, entry)
	tab = container.NewTabItem(title, entry)
	return tab
}

func addNewTab(tabs *container.AppTabs, tabMap map[*container.TabItem]*widget.Entry) {
	entry := widget.NewMultiLineEntry()
	entry.TextStyle = fyne.TextStyle{Monospace: true}
	entry.Wrapping = fyne.TextWrapWord // TODO
	entry.SetPlaceHolder("Type here...")
	
	title := "Untitled " + strconv.Itoa(tabCount)
	tabCount++

	dummyTab := &container.TabItem{}
	tab := makeNewTab(title, entry, dummyTab)
	*dummyTab = *tab // Assign real content back
	tabs.Append(tab)
	tabMap[tab] = entry
	tabs.Select(tab)
}

func closeCurrentTab(tabs *container.AppTabs, tabMap map[*container.TabItem]*widget.Entry) {
	current := tabs.Selected()
	if current == nil || len(tabs.Items) == 1 {
		return // Don't close the last tab
	}
	tabs.Remove(current)
	delete(tabMap, current)
}

func saveCurrentTab(tabs *container.AppTabs, tabMap map[*container.TabItem]*widget.Entry, myWindow fyne.Window) {
	current := tabs.Selected()
	if current == nil {
		return
	}
	entry := tabMap[current]
	content := entry.Text

	fileDialog := dialog.NewFileSave(func(uc fyne.URIWriteCloser, err error) {
		if uc != nil {
			uc.Write([]byte(content))
			uc.Close()
			current.Text = filepath.Base(uc.URI().Path())
			tabs.Refresh()
		}
	}, myWindow)
	fileDialog.SetFileName("untitled.txt")
	fileDialog.Show()
}

func loadSession(tabs *container.AppTabs, tabMap map[*container.TabItem]*widget.Entry) {
	if _, err := os.Stat(sessionFile); err == nil {
		data, err := os.ReadFile(sessionFile)
		if err == nil {
			var savedTabs []TabData
			if err := json.Unmarshal(data, &savedTabs); err == nil {
				for _, t := range savedTabs {
					entry := widget.NewMultiLineEntry()
					entry.SetText(t.Content)
					dummyTab := &container.TabItem{}
					tab := makeNewTab(t.Title, entry, dummyTab)
					*dummyTab = *tab
					tabMap[tab] = entry
					tabs.Append(tab)
				}
			}
		}
	}
	if len(tabs.Items) == 0 {
		addNewTab(tabs, tabMap)
	}
}

func saveSession(tabs *container.AppTabs, tabMap map[*container.TabItem]*widget.Entry) {
	var saveTabs []TabData
	for _, tab := range tabs.Items {
		entry := tabMap[tab]
		saveTabs = append(saveTabs, TabData{
			Title:    tab.Text,
			FilePath: "",
			Content:  entry.Text,
		})
	}
	data, _ := json.MarshalIndent(saveTabs, "", "  ")
	_ = os.WriteFile(sessionFile, data, 0644)
}

func openFile(tabs *container.AppTabs, tabMap map[*container.TabItem]*widget.Entry, w fyne.Window) {
	fd := dialog.NewFileOpen(func(uc fyne.URIReadCloser, err error) {
		if uc == nil {
			return
		}
		data, _ := io.ReadAll(uc)
		entry := widget.NewMultiLineEntry()
		entry.SetText(string(data))
		title := filepath.Base(uc.URI().Path())
		tab := makeNewTab(title, entry, &container.TabItem{})
		tabs.Append(tab)
		tabMap[tab] = entry
		tabs.Select(tab)
		uc.Close()
	}, w)
	fd.Show()
}

type customTheme struct {
	fontSize float32
}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(fyne.TextStyle{Monospace: true})
}

func (t *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	if name == theme.SizeNameText {
		return t.fontSize
	}
	return theme.DefaultTheme().Size(name)
}

func registerShortcuts(tabs *container.AppTabs, tabMap map[*container.TabItem]*widget.Entry, myWindow fyne.Window, myApp fyne.App) {
	ctrlN := &desktop.CustomShortcut{KeyName: fyne.KeyN, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlN, func(shotcut fyne.Shortcut) {
		println("New Tab Shortcut Triggered")
		addNewTab(tabs, tabMap)
	})

	ctrlT := &desktop.CustomShortcut{KeyName: fyne.KeyT, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlT, func(shotcut fyne.Shortcut) {
		println("New Tab Shortcut Triggered")
		addNewTab(tabs, tabMap)
	})

	ctrlS := &desktop.CustomShortcut{KeyName: fyne.KeyS, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlS, func(shotcut fyne.Shortcut) {
		println("Save Tab Shortcut Triggered")
		saveCurrentTab(tabs, tabMap, myWindow)
	})

	ctrlW := &desktop.CustomShortcut{KeyName: fyne.KeyW, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlW, func(shortcut fyne.Shortcut) {
		println("Close Tab Shortcut Triggered")
		closeCurrentTab(tabs, tabMap)
	})

	ctrlO := &desktop.CustomShortcut{KeyName: fyne.KeyO, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlO, func(shortcut fyne.Shortcut) {
		println("Opening File Shortcut Triggered")
		openFile(tabs, tabMap, myWindow)
	})

	ctrlZero := &desktop.CustomShortcut{KeyName: fyne.Key0, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlZero, func(shortcut fyne.Shortcut) {
		println("Font Reset Shortcut Triggered")
		fontSize = defaultFontSize
		applyFontSize(myApp, fontSize)
	})

	ctrlPlus := &desktop.CustomShortcut{KeyName: fyne.KeyEqual, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlPlus, func(sc fyne.Shortcut) {
		if fontSize < 30 {
			fontSize += 2
			println("Increased font size to", fontSize)
			applyFontSize(myApp, fontSize)
		}
	})

	ctrlMinus := &desktop.CustomShortcut{KeyName: fyne.KeyMinus, Modifier: fyne.KeyModifierControl}
	myWindow.Canvas().AddShortcut(ctrlMinus, func(sc fyne.Shortcut) {
		if fontSize > 8 {
			fontSize -= 2
			println("Decreased font size to", fontSize)
			applyFontSize(myApp, fontSize)
		}
	})
}
