package ui

import (
	"log"
	"strings"
	"unsafe"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
	"github.com/cjbrigato/dofhunt/datas"
	"github.com/cjbrigato/dofhunt/datas/types"
	"github.com/cjbrigato/dofhunt/language"
	"github.com/cjbrigato/dofhunt/winres"
)

type CluesSelectablesState struct {
	curClues         []string
	curSelectedClue  string
	curSelectedIndex int32
	curFilteredClues []string
	filterText       string
}

func (cs *CluesSelectablesState) filterClues(gameLangCode string) {
	if cs.filterText != "" && len(cs.curClues) > 0 {
		cs.curFilteredClues = []string{}
		for _, clue := range cs.curClues {
			if strings.Contains(strings.ToLower(clue), datas.NormalizeString(gameLangCode, cs.filterText, true)) {
				cs.curFilteredClues = append(cs.curFilteredClues, clue)
			}
		}
	} else {
		cs.curFilteredClues = cs.curClues
	}
}

type AppWindowState struct {
	wnd                *g.MasterWindow
	shouldFilterFocus  bool
	shouldListboxFocus bool
	canConfirm         bool
	isMovingFrame      bool
}

type AppUIState struct {
	initialized  bool
	gameLangCode string
	windowState  *AppWindowState

	CurrentMapPosition types.MapPosition
	LastMapPosition    types.MapPosition
	CurrentDirection   types.ClueDirection

	CurrentClues types.ClueResultSet
	CluesState   *CluesSelectablesState
}

func NewAppUIState(w *g.MasterWindow) *AppUIState {
	return &AppUIState{

		gameLangCode: "fr",
		windowState: &AppWindowState{
			wnd: w,
		},

		CurrentMapPosition: types.MapPosition{0, 0},
		LastMapPosition:    types.MapPosition{0, 0},
		CurrentDirection:   types.ClueDirectionNone,

		CurrentClues: types.ClueResultSet{},
		CluesState: &CluesSelectablesState{
			curClues:         []string{},
			curFilteredClues: []string{},
			curSelectedIndex: int32(-1),
			curSelectedClue:  "",
			filterText:       "",
		},
	}
}

func (s *AppUIState) titleBarLayout() *g.CustomWidget {
	return FramelessWindowMoveWidget(g.Custom(func() {
		winres.Icon16Texture.ToImageWidget().Scale(0.75, 0.75).Build()
		imgui.SameLine()
		imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextAlign, imgui.Vec2{0.0, 1.0})
		imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextPadding, imgui.Vec2{20.0, 2.0})
		imgui.SeparatorText("DofHunt")
		imgui.PopStyleVarV(2)
	}), &s.windowState.isMovingFrame, s.windowState.wnd)
}

func (s *AppUIState) inputsLineLayout() *g.RowWidget {
	return g.Row(g.Custom(func() {
		imgui.PushItemWidth(40.0)
		g.DragInt("X", (*int32)(unsafe.Pointer(&s.CurrentMapPosition.X)), -100, 150).Build()
		imgui.SameLine()
		g.DragInt("Y", (*int32)(unsafe.Pointer(&s.CurrentMapPosition.Y)), -100, 150).Build()
		imgui.PopItemWidth()
		imgui.SameLine()
		if s.windowState.shouldFilterFocus {
			g.SetKeyboardFocusHere()
			s.windowState.shouldFilterFocus = false
		}
		g.InputText(&s.CluesState.filterText).Flags(g.InputTextFlagsEnterReturnsTrue).OnChange(func() {
			s.windowState.shouldListboxFocus = true
		}).Build()
		s.CluesState.filterClues(s.gameLangCode)
	},
	),
	)
}

func (s *AppUIState) setupPageWindowLayout() []g.Widget {
	return append(make([]g.Widget, 0),
		g.Dummy(-1, 5),
		FramelessWindowMoveWidget(winres.SplashTexture.ToImageWidget(), &s.windowState.isMovingFrame, s.windowState.wnd),
		g.Custom(func() {
			imgui.SeparatorText("Hunt Smarter")
		}),
		language.AppSupportedLanguages.LangSetupLayout(&s.initialized),
	)
}

func (s *AppUIState) UpdateClues() {
	s.CurrentClues = types.GetClueResultSet(s.CurrentMapPosition, s.CurrentDirection, 10)
	s.CluesState.curClues = s.CurrentClues.Pois()
	if len(s.CluesState.curClues) > 0 {
		s.windowState.shouldFilterFocus = true
		s.CluesState.curSelectedClue = s.CluesState.curClues[0]
		s.windowState.canConfirm = true
	} else {
		s.CluesState.curSelectedClue = "SELECTED_CLUE_NOTFOUND"
		s.windowState.canConfirm = false
	}
}

func (s *AppUIState) ResetClues(message string) {
	s.CurrentDirection = types.ClueDirectionNone
	s.CluesState.curClues = []string{}
	s.CluesState.curSelectedClue = message
	s.CurrentClues = types.ClueResultSet{}
	s.windowState.canConfirm = false
}

func (s *AppUIState) TravelNextClue() {
	poi := s.CluesState.curSelectedClue
	pos, err := s.CurrentClues.Pos(poi)
	if err != nil {
		log.Println(err)
		return
	}
	travel := pos.TravelCommand()
	imgui.LogToClipboard()
	imgui.LogText(travel)
	imgui.LogFinish()
	TravelHistory.AddEntry(s.CurrentMapPosition, s.CurrentDirection, s.CluesState.curSelectedClue, types.MapPosition{
		X: pos.X,
		Y: pos.Y,
	})
	s.CurrentMapPosition.X = pos.X
	s.CurrentMapPosition.Y = pos.Y
	s.CluesState.filterText = ""
	s.ResetClues("SELECTED_CLUE_TRAVELED")
}

func (s *AppUIState) clueResultsListboxLayout() *g.CustomWidget {
	return g.Custom(func() {
		if s.windowState.shouldListboxFocus {
			imgui.SetNextWindowFocus()
			s.windowState.shouldListboxFocus = false
		} else {
			if g.IsKeyPressed(g.KeyEscape) {
				s.windowState.shouldFilterFocus = true
			}
		}
		onChange := func(selectedIndex int) {
			if g.IsKeyPressed(g.KeyEnter) {
				s.CluesState.curSelectedIndex = int32(selectedIndex)
				if len(s.CluesState.curFilteredClues) > int(selectedIndex) {
					s.CluesState.curSelectedClue = s.CluesState.curFilteredClues[selectedIndex]
					s.TravelNextClue()
				}
			}
		}
		onDclick := func(selectedIndex int) {
			s.CluesState.curSelectedIndex = int32(selectedIndex)
			if len(s.CluesState.curFilteredClues) > int(selectedIndex) {
				s.CluesState.curSelectedClue = s.CluesState.curFilteredClues[selectedIndex]
				s.TravelNextClue()
			}
		}
		g.ListBox(s.CluesState.curFilteredClues).Size(-1, 100).OnChange(onChange).SelectedIndex(&s.CluesState.curSelectedIndex).OnDClick(onDclick).Build()
		if int(s.CluesState.curSelectedIndex) >= 0 && len(s.CluesState.curFilteredClues) > int(s.CluesState.curSelectedIndex) {
			s.CluesState.curSelectedClue = s.CluesState.curFilteredClues[s.CluesState.curSelectedIndex]
		} else {
			s.CluesState.curSelectedIndex = -1
		}
	})
}

func (s *AppUIState) directionPadChildLayout() *g.ChildWidget {
	return g.Child().Flags(g.WindowFlagsNoNav).Size(115, 100).Layout(
		g.Row(g.Custom(func() {
			g.Dummy(22.0, 0).Build()
			if s.CurrentDirection != types.ClueDirectionUp {
				imgui.SameLine()
				g.ArrowButton(g.DirectionUp).OnClick(func() {
					s.CurrentDirection = types.ClueDirectionUp
					s.UpdateClues()
				}).Build()
			} else {
				g.Label("").Build()
			}
		})),
		g.Row(g.Custom(func() {
			if s.CurrentDirection != types.ClueDirectionLeft {
				g.ArrowButton(g.DirectionLeft).OnClick(func() {
					s.CurrentDirection = types.ClueDirectionLeft
					s.UpdateClues()
				}).Build()
			} else {
				g.Dummy(22.0, 0).Build()
			}
			imgui.SameLine()
			if s.CurrentDirection != types.ClueDirectionNone {
				g.Button("    ").OnClick(func() {
					s.ResetClues("SELECTED_CLUE_RESET")
				}).Build()
			} else {
				g.Dummy(21.0, 0).Build()
			}
			imgui.SameLine()
			if s.CurrentDirection != types.ClueDirectionRight {
				g.ArrowButton(g.DirectionRight).OnClick(func() {
					s.CurrentDirection = types.ClueDirectionRight
					s.UpdateClues()
				}).Build()
			} else {
				g.Dummy(21.0, 0).Build()
			}
		})),
		g.Row(g.Custom(func() {
			g.Dummy(22.0, 0).Build()
			if s.CurrentDirection != types.ClueDirectionDown {
				imgui.SameLine()
				g.ArrowButton(g.DirectionDown).OnClick(func() {
					s.CurrentDirection = types.ClueDirectionDown
					s.UpdateClues()
				}).Build()
			} else {
				g.Label("").Build()
			}
		})),
		g.Row(g.Custom(func() {
			var label string
			if s.windowState.canConfirm {
				label = "*Double-Click :"
			}
			g.Label(label).Build()
		})),
	)
}

func (s *AppUIState) Loop() {
	WithUIStyle(func() {
		if !s.initialized {
			g.SingleWindow().Layout(
				s.setupPageWindowLayout()...,
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
				s.titleBarLayout(),
				s.inputsLineLayout(),
				g.Row(
					s.directionPadChildLayout(),
					s.clueResultsListboxLayout(),
				),
				TravelHistory.HistoryLayout(s.windowState.wnd),
			)
		}
		if s.LastMapPosition.X != s.CurrentMapPosition.X || s.LastMapPosition.Y != s.CurrentMapPosition.Y {
			s.ResetClues("SELECTED_CLUE_POS_CHANGED")
		}
		s.LastMapPosition.X = s.CurrentMapPosition.X
		s.LastMapPosition.Y = s.CurrentMapPosition.Y
	})
}
