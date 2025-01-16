package types

import (
	"fmt"
	"sort"

	g "github.com/AllenDang/giu"
)

type ClueDirection int

const (
	ClueDirectionRight ClueDirection = iota
	ClueDirectionDown
	ClueDirectionLeft
	ClueDirectionUp
	ClueDirectionNone
)

func (cd ClueDirection) String() string {
	switch cd {
	case ClueDirectionRight:
		return "right"
	case ClueDirectionDown:
		return "down"
	case ClueDirectionLeft:
		return "left"
	case ClueDirectionUp:
		return "up"
	}
	return "none"
}

func (cd ClueDirection) Arrow() string {
	switch cd {
	case ClueDirectionRight:
		return "→"
	case ClueDirectionDown:
		return "↓"
	case ClueDirectionLeft:
		return "←"
	case ClueDirectionUp:
		return "↑"
	}
	return "none"
}

func (cd ClueDirection) Button() g.Widget {
	switch cd {
	case ClueDirectionRight:
		return g.ArrowButton(g.DirectionRight)
	case ClueDirectionDown:
		return g.ArrowButton(g.DirectionDown)
	case ClueDirectionLeft:
		return g.ArrowButton(g.DirectionLeft)
	case ClueDirectionUp:
		return g.ArrowButton(g.DirectionUp)
	}
	return g.Button("    ")
}

type ClueResultSet map[string]MapPosition

func (crs ClueResultSet) Pois() []string {
	r := make([]string, 0, len(crs))
	for k := range crs {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

func (crs ClueResultSet) Pos(p string) (*MapPosition, error) {
	pos, ok := crs[p]
	if !ok {
		return nil, fmt.Errorf("this clue/poi does not exists in result set")
	}
	return &pos, nil
}

func GetClueResultSet(start MapPosition, dir ClueDirection, limit int) ClueResultSet {
	results := make(ClueResultSet)
	positions := directedMapPositions(start, dir, limit)
	for _, position := range positions {
		names := position.GetClueNames()
		for _, name := range names {
			if _, ok := results[name]; !ok {
				results[name] = position
			}
		}
	}
	return results
}
