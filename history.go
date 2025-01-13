package main

import (
	"fmt"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
)

type TravelHistoryEntry struct {
	From      MapPosition
	Direction ClueDirection
	ClueName  string
	To        MapPosition
}

const (
	MaxTravelHistoryEntries = 4
)

var (
	TravelHistory = NewTravelHistoryCollection()
)

type TravelHistoryCollection struct {
	entries []TravelHistoryEntry
}

func NewTravelHistoryCollection() *TravelHistoryCollection {
	return &TravelHistoryCollection{
		entries: make([]TravelHistoryEntry, 0),
	}
}

func (t *TravelHistoryCollection) AddEntry(from MapPosition, dir ClueDirection, clue string, to MapPosition) {
	t.entries = append(t.entries, TravelHistoryEntry{
		From:      from,
		Direction: dir,
		ClueName:  clue,
		To:        to,
	})
	if len(t.entries) > MaxTravelHistoryEntries {
		t.entries = t.entries[len(t.entries)-MaxTravelHistoryEntries:]
	}
}

func (t *TravelHistoryCollection) GetEntries() []TravelHistoryEntry {
	return t.entries
}

func (t *TravelHistoryCollection) GenerateCurrentFileListTableRow() []*g.TableRowWidget {

	rows := make([]*g.TableRowWidget, 0)
	for _, entry := range t.GetEntries() {
		rows = append(rows, entry.TableRow())
	}
	return rows
}

func (t *TravelHistoryCollection) Table() *g.CustomWidget {
	return g.Custom(func() {
		imgui.PushStyleVarVec2(imgui.StyleVarSelectableTextAlign, imgui.Vec2{1.0, 0.0})
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
