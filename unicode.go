package main

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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
