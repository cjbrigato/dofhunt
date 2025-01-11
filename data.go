package main

import (
	"io"
	"log"
	"os"

	"github.com/tidwall/gjson"
)

var (
	clueMap   = make(map[int]map[int][]int)
	clueNames = make(map[int]string)
)

func readJson() ([]byte, error) {
	jsonFile, err := os.Open("clues.json")
	if err != nil {
		return nil, err
	}

	defer jsonFile.Close()

	return io.ReadAll(jsonFile)
}

func GetDatas() {
	log.Println("Reading Datas...")
	bytes, err := readJson()
	if err != nil {
		log.Fatalf("Cannot GetDatas: %v", err)
	}
	result := gjson.ParseBytes(bytes)
	log.Println("Loading ClueMap...")
	result.Get("maps").ForEach(func(key, value gjson.Result) bool {
		pos := value.Get("position")
		x := int(pos.Get("x").Int())
		y := int(pos.Get("y").Int())
		clues := make([]int, 0)
		c := value.Get("clues").Array()
		for _, clue := range c {
			clues = append(clues, int(clue.Int()))
		}
		_, ok := clueMap[x]
		if !ok {
			clueMap[x] = make(map[int][]int)
		}
		clueMap[x][y] = clues
		return true // keep iterating
	})
	log.Println("Loaded ClueMaps")
	log.Println("Loading ClueNames...")
	result.Get("clues").ForEach(func(key, value gjson.Result) bool {
		id := int(value.Get("clue-id").Int())
		name := value.Get("name-fr").String()
		name = NormalizeString("fr", name, true)
		clueNames[id] = name
		return true // keep iterating
	})
	log.Println("Loaded ClueNames")
}
