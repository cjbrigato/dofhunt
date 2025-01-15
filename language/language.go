package language

import (
	"strings"

	"github.com/AllenDang/cimgui-go/imgui"
	g "github.com/AllenDang/giu"
	"github.com/cjbrigato/dofhunt/datas"
)

var AppSupportedLanguages = NewSupportedLanguagesCollection(SupportedLanguages)

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

type SupportedLanguage struct {
	countryCode  string
	FriendlyName string
}

type SupportedLanguagesCollection struct {
	supportedLanguages []SupportedLanguage
	countryCodeNameMap map[string]string
	nameCountryCodeMap map[string]string
	selectedIndex      int32
}

func (slc *SupportedLanguagesCollection) SelectedIndex() *int32 {
	return &slc.selectedIndex
}

func (slc *SupportedLanguagesCollection) CountryCode(lang string) string {
	return slc.nameCountryCodeMap[strings.ToLower(lang)]
}
func (slc *SupportedLanguagesCollection) Lang(countryCode string) string {
	return slc.countryCodeNameMap[strings.ToLower(countryCode)]
}

func (slc *SupportedLanguagesCollection) Langs() []string {
	langs := make([]string, len(slc.supportedLanguages))
	for i := range slc.supportedLanguages {
		langs[i] = slc.supportedLanguages[i].FriendlyName
	}
	return langs
}

func (slc *SupportedLanguagesCollection) CountryCodes() []string {
	ccs := make([]string, len(slc.supportedLanguages))
	for i := range slc.supportedLanguages {
		ccs[i] = slc.supportedLanguages[i].countryCode
	}
	return ccs
}

func (slc *SupportedLanguagesCollection) LangSetupLayout(initialized *bool) *g.RowWidget {
	return g.Row(g.Custom(func() {
		g.Dummy(-1, 5).Build()
		imgui.PushStyleVarVec2(imgui.StyleVarSelectableTextAlign, imgui.Vec2{0.5, 0.0})
		g.ListBox(slc.Langs()).Size(-1, 100).SelectedIndex(slc.SelectedIndex()).OnChange(func(idx int) {
			langs := slc.Langs()
			datas.GetDatas(slc.CountryCode(langs[idx]))
			*initialized = true
		}).Build()
		imgui.PopStyleVar()
	},
	))
}

func NewSupportedLanguagesCollection(languages []SupportedLanguage) *SupportedLanguagesCollection {
	col := &SupportedLanguagesCollection{
		supportedLanguages: languages,
		selectedIndex:      -1,
	}
	ccnmap := make(map[string]string)
	nccmap := make(map[string]string)
	for _, lang := range languages {
		ccnmap[strings.ToLower(lang.countryCode)] = lang.FriendlyName
		nccmap[strings.ToLower(lang.FriendlyName)] = lang.countryCode
	}
	col.countryCodeNameMap = ccnmap
	col.nameCountryCodeMap = nccmap
	return col
}
