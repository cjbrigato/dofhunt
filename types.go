package main

import (
	"fmt"
	"sort"

	g "github.com/AllenDang/giu"
)

type SupportedLanguage struct {
	countryCode  string
	FriendlyName string
}

var SupportedLanguages = []SupportedLanguage{
	SupportedLanguage{
		countryCode:  "fr",
		FriendlyName: "Francais",
	},
	SupportedLanguage{
		countryCode:  "en",
		FriendlyName: "English",
	},
	SupportedLanguage{
		countryCode:  "es",
		FriendlyName: "Espanol",
	},
	SupportedLanguage{
		countryCode:  "de",
		FriendlyName: "Deutsch",
	},
	SupportedLanguage{
		countryCode:  "pt",
		FriendlyName: "Portugues",
	},
}

type MapPosition struct {
	X int
	Y int
}

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

func (m *MapPosition) TravelCommand() string {
	return fmt.Sprintf("/travel %d %d", m.X, m.Y)
}

func (m *MapPosition) DirectedMapPositionsSet(dir ClueDirection) []MapPosition {
	return directedMapPositions(*m, dir, 10)
}

func (m *MapPosition) GetClueNames() []string {
	clues, ok := CluesPosMap[m.X][m.Y]
	if !ok {
		return nil
	}
	names := make([]string, 0)
	for _, clue := range clues {
		names = append(names, ClueNamesMap[clue])
	}
	return names
}

func (m *MapPosition) FindNextClue(dir ClueDirection) ClueResultSet {
	return getClueResultSet(*m, dir, 10)
}

func directedMapPositions(start MapPosition, dir ClueDirection, limit int) []MapPosition {
	if limit < 1 {
		return nil
	}
	results := make([]MapPosition, 0)
	switch dir {
	case ClueDirectionRight:
		for i := 1; i <= limit; i++ {
			results = append(results, MapPosition{
				X: start.X + i,
				Y: start.Y,
			})
		}
	case ClueDirectionLeft:
		for i := 1; i <= limit; i++ {
			results = append(results, MapPosition{
				X: start.X - i,
				Y: start.Y,
			})
		}
	case ClueDirectionUp:
		for i := 1; i <= limit; i++ {
			results = append(results, MapPosition{
				X: start.X,
				Y: start.Y - i,
			})
		}
	case ClueDirectionDown:
		for i := 1; i <= limit; i++ {
			results = append(results, MapPosition{
				X: start.X,
				Y: start.Y + i,
			})
		}
	}
	return results
}

func getClueResultSet(start MapPosition, dir ClueDirection, limit int) ClueResultSet {
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
