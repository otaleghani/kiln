// @feature:layouts Localization label registry for UI strings.
package i18n

// Labels holds all user-facing UI strings for a given language.
type Labels struct {
	SearchPlaceholder string
	ToggleTheme       string
	BackToTop         string
	OnThisPage        string
	TableOfContents   string
	LocalGraph        string
	Backlinks         string
	GeneratedWith     string
	PageNotFound      string
	GoBackHome        string
	Folder            string
	Tag               string
	Updated           string
	Created           string
	Words             string
	MinRead           string
	Copy              string
	NoResults         string
	Search            string
	Navbar            string
	Expand            string
	LastModified      string
}

var languages = map[string]*Labels{
	"en": {
		SearchPlaceholder: "Search notes...",
		ToggleTheme:       "Toggle theme",
		BackToTop:         "Back to top",
		OnThisPage:        "On this page",
		TableOfContents:   "Table of contents",
		LocalGraph:        "Local Graph",
		Backlinks:         "Backlinks",
		GeneratedWith:     "Generated with",
		PageNotFound:      "Page not found",
		GoBackHome:        "Go back home",
		Folder:            "Folder:",
		Tag:               "Tag:",
		Updated:           "Updated",
		Created:           "Created",
		Words:             "words",
		MinRead:           "%d min read",
		Copy:              "Copy",
		NoResults:         "No results found",
		Search:            "Search",
		Navbar:            "Navbar",
		Expand:            "Expand",
		LastModified:      "Last Modified",
	},
	"it": {
		SearchPlaceholder: "Cerca appunti...",
		ToggleTheme:       "Cambia tema",
		BackToTop:         "Torna su",
		OnThisPage:        "In questa pagina",
		TableOfContents:   "Indice",
		LocalGraph:        "Grafo locale",
		Backlinks:         "Backlinks",
		GeneratedWith:     "Generato con",
		PageNotFound:      "Pagina non trovata",
		GoBackHome:        "Torna alla home",
		Folder:            "Cartella:",
		Tag:               "Tag:",
		Updated:           "Aggiornato",
		Created:           "Creato",
		Words:             "parole",
		MinRead:           "%d min di lettura",
		Copy:              "Copia",
		NoResults:         "Nessun risultato trovato",
		Search:            "Cerca",
		Navbar:            "Navbar",
		Expand:            "Espandi",
		LastModified:      "Ultima modifica",
	},
}

// Resolve returns the Labels for the given language code, falling back to English.
func Resolve(lang string) *Labels {
	if l, ok := languages[lang]; ok {
		return l
	}
	return languages["en"]
}
