package settings

import (
	"encoding/json"
	"fmt"
	"image/color"
	"os"

	g "github.com/AllenDang/giu"
)

const (
	settingsFile = "settings.json"
)

var (
	defaultWindowBgColor = color.RGBA{50, 50, 70, 130}
	defaultFrameBgColor  = color.RGBA{30, 30, 60, 110}
	defaultWindowPosXY   = 300
)

type AppSettings struct {
	GameLangCode   string
	ShowHistory    bool
	WindowColor    color.RGBA
	FrameColor     color.RGBA
	LastWindowPosX int
	LastWindowPosY int
}

func NewDefaultAppSettings() *AppSettings {
	return &AppSettings{
		GameLangCode:   "",
		ShowHistory:    true,
		WindowColor:    defaultWindowBgColor,
		FrameColor:     defaultFrameBgColor,
		LastWindowPosX: defaultWindowPosXY,
		LastWindowPosY: defaultWindowPosXY,
	}
}

func dirPath() string {
	return AppDataDir("DofHunt", false)
}

func filePath() string {
	dir := dirPath()
	return fmt.Sprintf("%s/%s", dir, settingsFile)
}

func ensureDirectory() {
	dir := dirPath()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
}

func fileExists() bool {
	file := filePath()
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return false
	}
	return true
}

// return appsettings and true if it was correctly loaded
func load() (*AppSettings, bool) {
	if !fileExists() {
		return NewDefaultAppSettings(), false
	}
	settingsJson, err := os.ReadFile(filePath())
	if err != nil {
		return NewDefaultAppSettings(), false
	}
	var savedSettings AppSettings
	if err := json.Unmarshal(settingsJson, &savedSettings); err != nil {
		return NewDefaultAppSettings(), false
	}
	return &savedSettings, true
}

func (a *AppSettings) ResetColors(from *AppSettings) {
	a.FrameColor = from.FrameColor
	a.WindowColor = from.WindowColor
}

func (a *AppSettings) RecallRefColors() {
	refSettings, _ := load()
	a.ResetColors(refSettings)
}

func (a *AppSettings) Save() {
	ensureDirectory()
	settingsJson, _ := json.Marshal(a)
	os.WriteFile(filePath(), settingsJson, 0666)
}

func (a *AppSettings) SaveColors() {
	refSettings, _ := load()
	refSettings.WindowColor = a.WindowColor
	refSettings.FrameColor = a.FrameColor
	refSettings.Save()
}

func (a *AppSettings) SaveHistory() {
	refSettings, _ := load()
	refSettings.ShowHistory = a.ShowHistory
	refSettings.Save()
}

func (a *AppSettings) SaveGameLangCode() {
	refSettings, _ := load()
	refSettings.GameLangCode = a.GameLangCode
	refSettings.Save()
}

func (a *AppSettings) SaveWindowPos(wnd *g.MasterWindow) {
	refSettings, _ := load()
	ox, oy := wnd.GetPos()
	refSettings.LastWindowPosX = ox
	refSettings.LastWindowPosY = oy
	refSettings.Save()
}

func InitSettings() *AppSettings {
	settings, loaded := load()
	if !loaded {
		ensureDirectory()
		settings.Save()
	}
	return settings
}
