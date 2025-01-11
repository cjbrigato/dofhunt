package main

import (
	"fmt"
	"sort"
)

type MapPosition struct {
	X int
	Y int
}

type ClueDirection int

const (
	ClueDirectionRight ClueDirection = iota
	_
	ClueDirectionDown
	_
	ClueDirectionLeft
	_
	ClueDirectionUp
	ClueDirectionNone
)

func GetClueNames(pos MapPosition) []string {
	clues, ok := clueMap[pos.X][pos.Y]
	if !ok {
		return nil
	}
	names := make([]string, 0)
	for _, clue := range clues {
		names = append(names, clueNames[clue])
	}
	return names
}

func GetMapPositions(start MapPosition, dir ClueDirection, limit int) []MapPosition {
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

func GetClueResultSet(start MapPosition, dir ClueDirection, limit int) ClueResultSet {
	results := make(ClueResultSet)
	positions := GetMapPositions(start, dir, limit)
	for _, position := range positions {
		names := GetClueNames(position)
		for _, name := range names {
			if _, ok := results[name]; !ok {
				results[name] = position
			}
		}
	}
	return results
}

type ClueResultSet map[string]MapPosition

func (m *MapPosition) TravelCommand() string {
	return fmt.Sprintf("/travel %d %d", m.X, m.Y)
}

func (crs ClueResultSet) Pois() []string {
	pois := Keys(crs)
	sort.Strings(pois)
	return pois
}

func (crs ClueResultSet) Pos(p string) (*MapPosition, error) {
	pos, ok := crs[p]
	if !ok {
		return nil, fmt.Errorf("this clue/poi does not exists in result set")
	}
	return &pos, nil
}

func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	r := make([]K, 0, len(m))
	for k := range m {
		r = append(r, k)
	}
	return r
}
