package ui

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
	"github.com/cjbrigato/dofhunt/types"
)

type TravelHistoryEntry struct {
	From      types.MapPosition
	Direction types.ClueDirection
	ClueName  string
	To        types.MapPosition
}

const (
	MaxTravelHistoryEntries = 4
)

var (
	TravelHistory = NewTravelHistoryCollection(MaxTravelHistoryEntries)
)

type TravelHistoryCollection struct {
	entries         []TravelHistoryEntry
	maxEntries      int
	showHistory     bool
	lastShowHistory bool
}

func NewTravelHistoryCollection(max int) *TravelHistoryCollection {
	return &TravelHistoryCollection{
		entries:         make([]TravelHistoryEntry, 0),
		maxEntries:      max,
		showHistory:     true,
		lastShowHistory: true,
	}
}

func (t *TravelHistoryCollection) AddEntry(from types.MapPosition, dir types.ClueDirection, clue string, to types.MapPosition) {
	t.entries = append(t.entries, TravelHistoryEntry{
		From:      from,
		Direction: dir,
		ClueName:  clue,
		To:        to,
	})
	if len(t.entries) > t.maxEntries {
		t.entries = t.entries[len(t.entries)-t.maxEntries:]
	}
}

func (t *TravelHistoryCollection) GetEntries() []TravelHistoryEntry {
	return t.entries
}

func (t *TravelHistoryCollection) GenerateCurrentFileListTableRow() []*g.TableRowWidget {
	rows := make([]*g.TableRowWidget, 0)
	for i := len(t.entries) - 1; i >= 0; i-- {
		rows = append(rows, t.entries[i].TableRow())
	}
	return rows
}

func (t *TravelHistoryCollection) Table() *g.CustomWidget {
	return g.Custom(func() {
		imgui.PushStyleVarVec2(imgui.StyleVarSelectableTextAlign, imgui.Vec2{X: 1.0, Y: 0.0})
		g.Table().NoHeader(true).Freeze(0, 1).Flags(g.TableFlagsRowBg).
			Columns(
				g.TableColumn("F[").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("FX").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("F,").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("FY").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("F]").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("Dir").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("T[").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("TX").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("T,").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("TY").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("T]").Flags(g.TableColumnFlagsWidthFixed),
				g.TableColumn("Clue").Flags(g.TableColumnFlagsWidthStretch),
			).
			Rows(t.GenerateCurrentFileListTableRow()...,
			).Build()
		imgui.PopStyleVar()
	})
}

func (t *TravelHistoryCollection) HistoryLayout(w *g.MasterWindow) g.Widget {
	return g.Custom(func() {
		g.Row(g.Custom(func() {
			imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextAlign, imgui.Vec2{1.0, 1.0})
			imgui.PushStyleVarVec2(imgui.StyleVarSeparatorTextPadding, imgui.Vec2{20.0, 0.0})
			imgui.Checkbox("##shhist", &t.showHistory)
			imgui.SameLine()
			imgui.SeparatorText("History")
			imgui.PopStyleVarV(2)
		})).Build()
		if t.showHistory != t.lastShowHistory {
			ox, oy := w.GetSize()
			if t.lastShowHistory {
				w.SetSize(ox, oy-70)
			} else {
				w.SetSize(ox, oy+70)
			}
			t.lastShowHistory = t.showHistory
		}

		if t.showHistory {
			g.Custom(func() {
				if len(t.GetEntries()) > 0 {
					t.Table().Build()
				}
			}).Build()
		}
	})
}

func (te *TravelHistoryEntry) TableRow() *g.TableRowWidget {
	return g.TableRow(
		g.Label("["),
		g.Selectable(fmt.Sprintf("%d", te.From.X)).Flags(g.SelectableFlagsDisabled),
		g.Label(","),
		g.Selectable(fmt.Sprintf("%d", te.From.Y)).Flags(g.SelectableFlagsDisabled),
		g.Label("]"),
		g.Label(te.Direction.Arrow()),
		g.Label("["),
		g.Selectable(fmt.Sprintf("%d", te.To.X)).Flags(g.SelectableFlagsDisabled),
		g.Label(","),
		g.Selectable(fmt.Sprintf("%d", te.To.Y)).Flags(g.SelectableFlagsDisabled),
		g.Label("]"),
		g.Label(te.ClueName),
	)
}

func (te *TravelHistoryEntry) Row() *g.RowWidget {
	fromStr := fmt.Sprintf("%3d, %3d", te.From.X, te.From.Y)
	toStr := fmt.Sprintf("%3d, %3d", te.To.X, te.To.Y)
	return g.Row(g.Custom(func() {
		g.Label(fromStr).Build()
		imgui.SameLine()
		g.Label(te.Direction.Arrow()).Build()
		imgui.SameLine()
		g.Label(toStr).Build()
		imgui.SameLine()
		g.Label(te.ClueName).Build()

	}))
}
