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
	MaxTravelHistoryEntries = 3
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

func (te *TravelHistoryEntry) Row() *g.RowWidget {
	fromStr := fmt.Sprintf("[%d, %d]", te.From.X, te.From.Y)
	toStr := fmt.Sprintf("[%d, %d]", te.To.X, te.To.Y)
	return g.Row(g.Custom(func() {
		g.Label(fromStr).Build()
		imgui.SameLine()
		te.Direction.Button().Build()
		imgui.SameLine()
		g.Label(te.ClueName).Build()
		imgui.SameLine()
		g.Label(toStr).Build()
	}))
}
