package main

import (
	_ "embed"
	"fmt"
	"log"
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

//go:embed clues.json
var jsonDatas []byte

func GetDatas(countryCode string) {
	langKey := fmt.Sprintf("name-%s", countryCode)
	log.Println("Reading Datas...")
	result := gjson.ParseBytes(jsonDatas)
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
		name := value.Get(langKey).String()
		name = NormalizeString(countryCode, name, true)
		ClueNamesMap[id] = name
		return true 
	})
	log.Println("Loaded ClueNames")
}
