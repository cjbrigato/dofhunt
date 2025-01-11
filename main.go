package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
)

var (
	curPosX         = int32(0)
	curPosY         = int32(0)
	curDir          = ClueDirectionNone
	curClues        = []string{}
	curSelectedClue = ""
	curResultSet    = ClueResultSet{}
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
					}
				}
				imgui.EndCombo()
			}
		},
		),
		),
		g.Row(
			g.Label("       "),
			g.Button(" ^ ").OnClick(func() {
				curDir = ClueDirectionUp
				Do()
			}),
		),
		g.Row(
			g.Button(" < ").OnClick(func() {
				curDir = ClueDirectionLeft
				Do()
			}),
			g.Button(" X ").OnClick(func() {
				curDir = ClueDirectionNone
				curClues = []string{}
				curSelectedClue = ""
				curResultSet = ClueResultSet{}
			}),
			g.Button(" > ").OnClick(func() {
				curDir = ClueDirectionRight
				Do()
			}),
			g.Label("  "),
			g.Button("Travel").OnClick(GetTravelCommand),
		),
		g.Row(
			g.Label("       "),
			g.Button(" v ").OnClick(func() {
				curDir = ClueDirectionDown
				Do()
			}),
		),
	)
	g.PopStyleColor()
	g.PopStyleColor()
	imgui.PopStyleVar()
}

func Do() {
	curResultSet = GetClueResultSet(MapPosition{
		X: int(curPosX),
		Y: int(curPosY),
	}, curDir, 10)
	curClues = curResultSet.Pois()
}

func GetTravelCommand() {
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
	curPosX = int32(pos.X)
	curPosY = int32(pos.Y)
	curDir = ClueDirectionNone
	curClues = []string{}
	curSelectedClue = ""
	curResultSet = ClueResultSet{}
}

func main() {
	GetDatas()
	wnd := g.NewMasterWindow("DofHunt", 400, 230, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent)
	wnd.SetTargetFPS(60)
	wnd.SetBgColor(color.RGBA{0, 0, 0, 0})
	wnd.Run(loop)

}
