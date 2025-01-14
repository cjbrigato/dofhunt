package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
)

func DecodeEmbedded(data []byte) (*image.RGBA, error) {
	r := bytes.NewReader(data)
	img, err := png.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("LoadImage: error decoding png image: %w", err)
	}
	return g.ImageToRgba(img), nil
}

//go:embed winres/splash.png
var splashHeaderLogo []byte

func DecodeSplashHeaderLogo() (*image.RGBA, error) {
	return DecodeEmbedded(splashHeaderLogo)
}

//go:embed winres/icon16.png
var appIcon16 []byte

func DecodeAppIcon16() (*image.RGBA, error) {
	return DecodeEmbedded(appIcon16)
}

//go:embed winres/icon.png
var appIcon []byte

func DecodeAppIcon() (*image.RGBA, error) {
	return DecodeEmbedded(appIcon)
}

const (
	SELECTED_CLUE_RESET       = "[SET Position -> Direction]"
	SELECTED_CLUE_TRAVELED    = "[Choose NEXT -> Direction]"
	SELECTED_CLUE_POS_CHANGED = "[Position Changed -> Set Direction]"
	SELECTED_CLUE_NOTFOUND    = "(X_x) No clues. You messed up"
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
	rgbaIcon16         *image.RGBA
	rgbaIcon           *image.RGBA
	headerSplashRgba   *image.RGBA
	splashTexture      = &g.ReflectiveBoundTexture{}
	icon16Texture      = &g.ReflectiveBoundTexture{}
	curSelectedIndex   = int32(-1)
	filterText         = ""
	wnd                *g.MasterWindow
	isMovingFrame      = false
	language           = "fr"
	initialized        = false
	shouldFilterFocus  = false
	shouldListboxFocus = false
)

func framelessWindowMoveWidget(widget g.Widget) *g.CustomWidget {
	return g.Custom(func() {
		if isMovingFrame && !g.IsMouseDown(g.MouseButtonLeft) {
			isMovingFrame = false
			return
		}

		widget.Build()

		if g.IsItemHovered() {
			if g.IsMouseDown(g.MouseButtonLeft) {
				isMovingFrame = true
			}
		}

		if isMovingFrame {
			delta := imgui.CurrentIO().MouseDelta()
			dx := int(delta.X)
			dy := int(delta.Y)
			if dx != 0 || dy != 0 {
				ox, oy := wnd.GetPos()
				wnd.SetPos(ox+dx, oy+dy)
			}
		}
	})
}

func titleBarLayout() *g.CustomWidget {
	return framelessWindowMoveWidget(g.Custom(func() {
		icon16Texture.ToImageWidget().Scale(0.75, 0.75).Build()
		imgui.SameLine()
		imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextAlign, imgui.Vec2{0.0, 1.0})
		imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextPadding, imgui.Vec2{20.0, 2.0})
		imgui.SeparatorText("DofHunt")
		imgui.PopStyleVarV(2)
	}))
}

func langSetupLayout() *g.RowWidget {
	return g.Row(g.Custom(func() {
		g.Dummy(-1, 5).Build()
		imgui.PushStyleVarVec2(imgui.StyleVarSelectableTextAlign, imgui.Vec2{0.5, 0.0})
		g.ListBox(AppSupportedLanguages.Langs()).Size(-1, 100).SelectedIndex(AppSupportedLanguages.SelectedIndex()).OnChange(func(idx int) {
			langs := AppSupportedLanguages.Langs()
			GetDatas(AppSupportedLanguages.CountryCode(langs[idx]))
			initialized = true
		}).Build()
		imgui.PopStyleVar()
	},
	))
}

func onChange() {
	shouldListboxFocus = true
}

func headerLayout() *g.RowWidget {
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
		g.InputText(&filterText).Flags(g.InputTextFlagsEnterReturnsTrue).OnChange(onChange).Build()
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

func loop() {

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
	if !initialized {
		g.SingleWindow().Layout(
			g.Dummy(-1, 5),
			framelessWindowMoveWidget(splashTexture.ToImageWidget()),
			g.Custom(func() {
				imgui.SeparatorText("Hunt Smarter")
			}),
			langSetupLayout(),
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
			headerLayout(),
			g.Row(
				g.Child().Flags(g.WindowFlagsNoNav).Size(115, 100).Layout(
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
						if canConfirm {
							g.Button("Confirm Clue").OnClick(TravelNextClue).Build()
						} else {
							g.Label("").Build()
						}
					})),
				),
				g.Custom(func() {
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
				}),
			),
			g.Row(g.Custom(func() {
				imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextAlign, imgui.Vec2{1.0, 1.0})
				imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextPadding, imgui.Vec2{20.0, 0.0})
				imgui.SeparatorText("History")
				imgui.PopStyleVarV(2)
			})),
			g.Custom(func() {
				if len(TravelHistory.GetEntries()) > 0 {
					TravelHistory.Table().Build()
				}
			}),
		)
	}
	g.PopStyleColor()
	g.PopStyleColor()
	imgui.PopStyleVar()
	imgui.PopStyleVar()
	imgui.PopStyleVar()
	imgui.PopStyleVar()
	imgui.PopStyleVar()
	imgui.PopStyleVar()
	imgui.PopStyleColor()
	imgui.PopStyleColor()

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
	wnd = g.NewMasterWindow("DofHunt", 380, 263, g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFrameless|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent) //g.MasterWindowFlagsNotResizable|g.MasterWindowFlagsFloating|g.MasterWindowFlagsTransparent)
	wnd.SetTargetFPS(60)
	wnd.SetBgColor(color.RGBA{0, 0, 0, 0})
	rgbaIcon, _ = DecodeAppIcon()
	rgbaIcon16, _ = DecodeAppIcon16()
	headerSplashRgba, _ := DecodeSplashHeaderLogo()
	splashTexture.SetSurfaceFromRGBA(headerSplashRgba, false)
	icon16Texture.SetSurfaceFromRGBA(rgbaIcon16, false)
	wnd.SetPos(300, 300)
	wnd.Run(loop)

}
