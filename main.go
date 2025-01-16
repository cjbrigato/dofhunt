package main

import (
	_ "embed"
	"image/color"

	g "github.com/AllenDang/giu"
	"github.com/cjbrigato/dofhunt/ui"
	"github.com/cjbrigato/dofhunt/winres"
)

const (
	SELECTED_CLUE_RESET       = "[SET Position -> Direction]"
	SELECTED_CLUE_TRAVELED    = "[Choose NEXT -> Direction]"
	SELECTED_CLUE_POS_CHANGED = "[Position Changed -> Set Direction]"
	SELECTED_CLUE_NOTFOUND    = "(X_x) No clues. You messed up"
	WND_BASE_WIDTH            = 380
	WND_BASE_HEIGHT           = 267
)

var uiState *ui.AppUIState

func main() {
	wnd := g.NewMasterWindow("DofHunt", 380, 273, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent) //g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent)
	wnd.SetTargetFPS(60)
	wnd.SetBgColor(color.RGBA{0, 0, 0, 0})
	winres.InitTextures()
	wnd.SetPos(300, 300)
	uiState = ui.NewAppUIState(wnd)
	wnd.Run(uiState.Loop)

}
