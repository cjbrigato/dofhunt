package types

import (
	"fmt"

	"github.com/cjbrigato/dofhunt/datas"
)

type MapPosition struct {
	X int32
	Y int32
}

func (m *MapPosition) TravelCommand() string {
	return fmt.Sprintf("/travel %d %d", m.X, m.Y)
}

func (m *MapPosition) DirectedMapPositionsSet(dir ClueDirection) []MapPosition {
	return directedMapPositions(*m, dir, 10)
}

func (m *MapPosition) GetClueNames() []string {
	clues, ok := datas.CluesPosMap[int(m.X)][int(m.Y)]
	if !ok {
		return nil
	}
	names := make([]string, 0)
	for _, clue := range clues {
		names = append(names, datas.ClueNamesMap[clue])
	}
	return names
}

func (m *MapPosition) FindNextClue(dir ClueDirection) ClueResultSet {
	return GetClueResultSet(*m, dir, 10)
}

func directedMapPositions(start MapPosition, dir ClueDirection, limit int) []MapPosition {
	if limit < 1 {
		return nil
	}
	results := make([]MapPosition, 0)
	switch dir {
	case ClueDirectionRight:
		for i := int32(1); i <= int32(limit); i++ {
			results = append(results, MapPosition{
				X: start.X + i,
				Y: start.Y,
			})
		}
	case ClueDirectionLeft:
		for i := int32(1); i <= int32(limit); i++ {
			results = append(results, MapPosition{
				X: start.X - i,
				Y: start.Y,
			})
		}
	case ClueDirectionUp:
		for i := int32(1); i <= int32(limit); i++ {
			results = append(results, MapPosition{
				X: start.X,
				Y: start.Y - i,
			})
		}
	case ClueDirectionDown:
		for i := int32(1); i <= int32(limit); i++ {
			results = append(results, MapPosition{
				X: start.X,
				Y: start.Y + i,
			})
		}
	}
	return results
}
