package main

import (
	"io"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/tidwall/gjson"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	CluesPosMap  = make(map[int]map[int][]int)
	ClueNamesMap = make(map[int]string)
)

func NormalizeString(lang string, s string, lower bool) string {
	asciiname := s
	if lang != "en" {
		t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
		asciiname, _, _ = transform.String(t, s)
	}
	if lower {
		asciiname = strings.ToLower(asciiname)
	}
	return asciiname
}

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
		_, ok := CluesPosMap[x]
		if !ok {
			CluesPosMap[x] = make(map[int][]int)
		}
		CluesPosMap[x][y] = clues
		return true
	})
	log.Println("Loaded ClueMaps")
	log.Println("Loading ClueNames...")
	result.Get("clues").ForEach(func(key, value gjson.Result) bool {
		id := int(value.Get("clue-id").Int())
		name := value.Get("name-fr").String()
		name = NormalizeString("fr", name, true)
		ClueNamesMap[id] = name
		return true // keep iterating
	})
	log.Println("Loaded ClueNames")
}
