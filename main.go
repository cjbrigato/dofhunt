package main

import (
	_ "embed"
	"image/color"
	"log"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
)

const (
	SELECTED_CLUE_RESET       = "[SET Position -> Direction]"
	SELECTED_CLUE_TRAVELED    = "[Choose NEXT -> Direction]"
	SELECTED_CLUE_POS_CHANGED = "[Position Changed -> Set Direction]"
	SELECTED_CLUE_NOTFOUND    = "(X_x) No clues. You messed up"
	WND_BASE_WIDTH            = 380
	WND_BASE_HEIGHT           = 267
)

var (
	curPosX            = int32(0)
	curPosY            = int32(0)
	curDir             = ClueDirectionNone
	curClues           = []string{}
	curFilteredClues   = []string{}
	curSelectedClue    = SELECTED_CLUE_RESET
	canConfirm         = false
	curResultSet       = ClueResultSet{}
	lastPosX           = curPosX
	lastPosY           = curPosY
	curSelectedIndex   = int32(-1)
	filterText         = ""
	wnd                *g.MasterWindow
	isMovingFrame      = false
	language           = "fr"
	initialized        = false
	shouldFilterFocus  = false
	shouldListboxFocus = false
	showHistory        = true
)

func titleBarLayout() *g.CustomWidget {
	return framelessWindowMoveWidget(g.Custom(func() {
		icon16Texture.ToImageWidget().Scale(0.75, 0.75).Build()
		imgui.SameLine()
		imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextAlign, imgui.Vec2{0.0, 1.0})
		imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextPadding, imgui.Vec2{20.0, 2.0})
		imgui.SeparatorText("DofHunt")
		imgui.PopStyleVarV(2)
	}), &isMovingFrame, wnd)
}

func inputsLineLayout() *g.RowWidget {
	return g.Row(g.Custom(func() {
		imgui.PushItemWidth(40.0)
		g.DragInt("X", &curPosX, -100, 150).Build()
		imgui.SameLine()
		g.DragInt("Y", &curPosY, -100, 150).Build()
		imgui.PopItemWidth()
		imgui.SameLine()
		if shouldFilterFocus {
			g.SetKeyboardFocusHere()
			shouldFilterFocus = false
		}
		g.InputText(&filterText).Flags(g.InputTextFlagsEnterReturnsTrue).OnChange(func() {
			shouldListboxFocus = true
		}).Build()
		filterClues(&filterText)
	},
	),
	)
}

func filterClues(filter *string) {
	if *filter != "" && len(curClues) > 0 {
		curFilteredClues = []string{}
		for _, clue := range curClues {
			if strings.Contains(strings.ToLower(clue), NormalizeString(language, *filter, true)) {
				curFilteredClues = append(curFilteredClues, clue)
			}
		}
	} else {
		curFilteredClues = curClues
	}
}

func setupPageWindowLayout() []g.Widget {
	return append(make([]g.Widget, 0),
		g.Dummy(-1, 5),
		framelessWindowMoveWidget(splashTexture.ToImageWidget(), &isMovingFrame, wnd),
		g.Custom(func() {
			imgui.SeparatorText("Hunt Smarter")
		}),
		AppSupportedLanguages.LangSetupLayout(&initialized),
	)
}

func clueResultsListboxLayout() *g.CustomWidget {
	return g.Custom(func() {
		if shouldListboxFocus {
			imgui.SetNextWindowFocus()
			shouldListboxFocus = false
		} else {
			if g.IsKeyPressed(g.KeyEscape) {
				shouldFilterFocus = true
			}
		}
		onChange := func(selectedIndex int) {
			if g.IsKeyPressed(g.KeyEnter) {
				curSelectedIndex = int32(selectedIndex)
				if len(curFilteredClues) > int(selectedIndex) {
					curSelectedClue = curFilteredClues[selectedIndex]
					TravelNextClue()
				}
			}
		}
		onDclick := func(selectedIndex int) {
			curSelectedIndex = int32(selectedIndex)
			if len(curFilteredClues) > int(selectedIndex) {
				curSelectedClue = curFilteredClues[selectedIndex]
				TravelNextClue()
			}
		}
		g.ListBox(curFilteredClues).Size(-1, 100).OnChange(onChange).SelectedIndex(&curSelectedIndex).OnDClick(onDclick).Build()
		if int(curSelectedIndex) >= 0 && len(curFilteredClues) > int(curSelectedIndex) {
			curSelectedClue = curFilteredClues[curSelectedIndex]
		} else {
			curSelectedIndex = -1
		}
	})
}

func directionPadChildLayout() *g.ChildWidget {
	return g.Child().Flags(g.WindowFlagsNoNav).Size(115, 100).Layout(
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
			var label string
			if canConfirm {
				label = "*Double-Click :"
			}
			g.Label(label).Build()
		})),
	)
}

func WithUIStyle(fn func()) {
	imgui.PushStyleVarVec2(imgui.StyleVarCellPadding, imgui.Vec2{1.0, 1.0})
	imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextAlign, imgui.Vec2{1.0, 1.0})
	imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextPadding, imgui.Vec2{20.0, 0.0})
	imgui.PushStyleVarFloat(imgui.StyleVarWindowBorderSize, 0)
	imgui.PushStyleVarFloat(imgui.StyleVarWindowRounding, 6.0)
	imgui.PushStyleVarFloat(imgui.StyleVarChildBorderSize, 0)
	imgui.PushStyleColorVec4(imgui.ColChildBg, g.ToVec4Color(color.RGBA{50, 50, 70, 0}))
	imgui.PushStyleColorVec4(imgui.ColButton, g.ToVec4Color(color.RGBA{50, 50, 70, 130}))
	g.PushColorWindowBg(color.RGBA{50, 50, 70, 130})
	g.PushColorFrameBg(color.RGBA{30, 30, 60, 110})

	fn()

	g.PopStyleColorV(4)
	imgui.PopStyleVarV(6)
}

func loop() {
	WithUIStyle(func() {
		if !initialized {
			g.SingleWindow().Layout(
				setupPageWindowLayout()...,
			)
		} else {
			g.SingleWindow().Flags(
				g.WindowFlags(imgui.WindowFlagsNoTitleBar)|
					g.WindowFlags(imgui.WindowFlagsNoCollapse)|
					g.WindowFlags(imgui.WindowFlagsNoScrollbar)|
					g.WindowFlags(imgui.WindowFlagsNoMove)|
					g.WindowFlags(imgui.WindowFlagsNoResize)|
					g.WindowFlags(imgui.WindowFlagsNoNav),
			).Layout(
				titleBarLayout(),
				inputsLineLayout(),
				g.Row(
					directionPadChildLayout(),
					clueResultsListboxLayout(),
				),
				TravelHistory.HistoryLayout(wnd),
			)
		}
		if lastPosX != curPosX || lastPosY != curPosY {
			ResetClues(SELECTED_CLUE_POS_CHANGED)
		}
		lastPosX = curPosX
		lastPosY = curPosY
	})
}

func UpdateClues() {
	curResultSet = getClueResultSet(MapPosition{
		X: int(curPosX),
		Y: int(curPosY),
	}, curDir, 10)
	curClues = curResultSet.Pois()
	if len(curClues) > 0 {
		shouldFilterFocus = true
		curSelectedClue = curClues[0]
		canConfirm = true
	} else {
		curSelectedClue = SELECTED_CLUE_NOTFOUND
		canConfirm = false
	}
}

func ResetClues(message string) {
	curDir = ClueDirectionNone
	curClues = []string{}
	curSelectedClue = message
	curResultSet = ClueResultSet{}
	canConfirm = false
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
	filterText = ""
	ResetClues(SELECTED_CLUE_TRAVELED)
}

func main() {
	wnd = g.NewMasterWindow("DofHunt", 380, 273, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent) //g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent)
	wnd.SetTargetFPS(60)
	wnd.SetBgColor(color.RGBA{0, 0, 0, 0})
	initTextures()
	wnd.SetPos(300, 300)
	wnd.Run(loop)

}
