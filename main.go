package main

import (
	_ "embed"
	"image/color"
	"syscall"

	g "github.com/AllenDang/giu"
	"github.com/cjbrigato/dofhunt/ui"
	"github.com/cjbrigato/dofhunt/winres"
)

const (
	SELECTED_CLUE_RESET                  = "[SET Position -> Direction]"
	SELECTED_CLUE_TRAVELED               = "[Choose NEXT -> Direction]"
	SELECTED_CLUE_POS_CHANGED            = "[Position Changed -> Set Direction]"
	SELECTED_CLUE_NOTFOUND               = "(X_x) No clues. You messed up"
	WND_BASE_WIDTH                       = 380
	WND_BASE_HEIGHT                      = 273
	DPIAwarenessContextUnaware           = 16
	DPIAwarenessContextSystemAware       = 17
	DPIAwarenessContextPerMonitorAware   = 18
	DPIAwarenessContextPerMonitorAwareV2 = 34
)

var uiState *ui.AppUIState

func SetDpiAware() {
	user32 := syscall.NewLazyDLL("user32.dll")
	pSetProcessDpiAwarenessContext := user32.NewProc("SetProcessDpiAwarenessContext")
	pSetProcessDpiAwarenessContext.Call(DPIAwarenessContextPerMonitorAwareV2)
}

func main() {
	SetDpiAware()
	wnd := g.NewMasterWindow("DofHunt", 380, 273, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent) //g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent)
	wnd.SetTargetFPS(60)
	wnd.SetBgColor(color.RGBA{0, 0, 0, 0})
	winres.InitTextures()
	uiState = ui.NewAppUIState(wnd)
	wnd.SetPos(uiState.Settings.LastWindowPosX, uiState.Settings.LastWindowPosY)
	wnd.Run(uiState.Loop)
}
