package settings

import (
	"image/color"

	"github.com/cjbrigato/dofhunt/language"
)

type AppSettings struct {
	language    language.SupportedLanguage
	showHistory bool
	windowColor color.RGBA
}
