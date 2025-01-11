package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
)

const (
	SELECTED_CLUE_RESET       = "[SET Position -> Direction]"
	SELECTED_CLUE_TRAVELED    = "[Choose NEXT -> Direction]"
	SELECTED_CLUE_POS_CHANGED = "[Position Changed -> Set Direction]"
	SELECTED_CLUE_NOTFOUND    = "(X_x) No clues. You messed up"
)

var (
	curPosX         = int32(0)
	curPosY         = int32(0)
	curDir          = ClueDirectionNone
	curClues        = []string{}
	curSelectedClue = SELECTED_CLUE_RESET
	showTravel      = false
	curResultSet    = ClueResultSet{}
	lastPosX        = curPosX
	lastPosY        = curPosY
)

func loop() {
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	g.PushColorWindowBg(color.RGBA{50, 50, 70, 130})
	g.PushColorFrameBg(color.RGBA{30, 30, 60, 110})
	g.SingleWindow().Layout(
		g.Row(g.Custom(func() {
			imgui.PushItemWidth(40.0)
			g.DragInt("X", &curPosX, -100, 150).Build()
			imgui.SameLine()
			g.DragInt("Y", &curPosY, -100, 150).Build()
			imgui.PopItemWidth()
			imgui.SameLine()
			if imgui.BeginComboV("##dialogfilters", curSelectedClue, imgui.ComboFlags(imgui.ComboFlagsHeightRegular)) {
				for i, clue := range curClues {
					if imgui.SelectableBool(fmt.Sprintf("%s##%d", clue, i)) {
						curSelectedClue = clue
						showTravel = true
					}
				}
				imgui.EndCombo()
			}
		},
		),
		),
		g.Row(g.Custom(func() {
			g.Dummy(22.0, 0).Build()
			if curDir != ClueDirectionUp {
				imgui.SameLine()
				g.ArrowButton(g.DirectionUp).OnClick(func() {
					curDir = ClueDirectionUp
					UpdateClues()
				}).Build()
			} else {
				g.Label("").Build()
			}
		})),
		g.Row(g.Custom(func() {
			if curDir != ClueDirectionLeft {
				g.ArrowButton(g.DirectionLeft).OnClick(func() {
					curDir = ClueDirectionLeft
					UpdateClues()
				}).Build()
			} else {
				g.Dummy(22.0, 0).Build()
			}
			imgui.SameLine()
			if curDir != ClueDirectionNone {
				g.Button("    ").OnClick(func() {
					ResetClues(SELECTED_CLUE_RESET)
				}).Build()
			} else {
				g.Dummy(21.0, 0).Build()
			}
			imgui.SameLine()
			if curDir != ClueDirectionRight {
				g.ArrowButton(g.DirectionRight).OnClick(func() {
					curDir = ClueDirectionRight
					UpdateClues()
				}).Build()
			} else {
				g.Dummy(21.0, 0).Build()
			}
			imgui.SameLine()
			g.Label("  ").Build()
			if showTravel {
				imgui.SameLine()
				g.Button("Confirm Clue").OnClick(TravelNextClue).Build()
			}
		})),
		g.Row(g.Custom(func() {
			g.Dummy(22.0, 0).Build()
			if curDir != ClueDirectionDown {
				imgui.SameLine()
				g.ArrowButton(g.DirectionDown).OnClick(func() {
					curDir = ClueDirectionDown
					UpdateClues()
				}).Build()
			} else {
				g.Label("").Build()
			}
		})),
		g.Row(g.Custom(func() {
			imgui.SeparatorText("History")
		})),
		g.Custom(func() {
			for _, entry := range TravelHistory.GetEntries() {
				entry.Row().Build()
			}
		},
		),
	)
	g.PopStyleColor()
	g.PopStyleColor()
	imgui.PopStyleVar()
	if lastPosX != curPosX || lastPosY != curPosY {
		ResetClues(SELECTED_CLUE_POS_CHANGED)
	}
	lastPosX = curPosX
	lastPosY = curPosY
}

func UpdateClues() {
	curResultSet = getClueResultSet(MapPosition{
		X: int(curPosX),
		Y: int(curPosY),
	}, curDir, 10)
	curClues = curResultSet.Pois()
	if len(curClues) > 0 {
		curSelectedClue = curClues[0]
		showTravel = true
	} else {
		curSelectedClue = SELECTED_CLUE_NOTFOUND
		showTravel = false
	}
}

func ResetClues(message string) {
	curDir = ClueDirectionNone
	curClues = []string{}
	curSelectedClue = message
	curResultSet = ClueResultSet{}
	showTravel = false
}

func TravelNextClue() {
	poi := curSelectedClue
	pos, err := curResultSet.Pos(poi)
	if err != nil {
		log.Println(err)
		return
	}
	travel := pos.TravelCommand()
	imgui.LogToClipboard()
	imgui.LogText(travel)
	imgui.LogFinish()
	TravelHistory.AddEntry(MapPosition{
		X: int(curPosX),
		Y: int(curPosY),
	}, curDir, curSelectedClue, MapPosition{
		X: pos.X,
		Y: pos.Y,
	})
	curPosX = int32(pos.X)
	curPosY = int32(pos.Y)
	ResetClues(SELECTED_CLUE_TRAVELED)
}

func main() {
	GetDatas()
	wnd := g.NewMasterWindow("DofHunt", 400, 230, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent)
	wnd.SetTargetFPS(60)
	wnd.SetBgColor(color.RGBA{0, 0, 0, 0})
	wnd.Run(loop)

}
