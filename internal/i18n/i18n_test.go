// @feature:layouts Tests for localization label registry.
package i18n

import (
	"reflect"
	"testing"
)

func TestResolve_English(t *testing.T) {
	l := Resolve("en")
	if l.SearchPlaceholder != "Search notes..." {
		t.Errorf("SearchPlaceholder = %q, want %q", l.SearchPlaceholder, "Search notes...")
	}
	if l.ToggleTheme != "Toggle theme" {
		t.Errorf("ToggleTheme = %q, want %q", l.ToggleTheme, "Toggle theme")
	}
	if l.BackToTop != "Back to top" {
		t.Errorf("BackToTop = %q, want %q", l.BackToTop, "Back to top")
	}
	if l.PageNotFound != "Page not found" {
		t.Errorf("PageNotFound = %q, want %q", l.PageNotFound, "Page not found")
	}
	if l.MinRead != "%d min read" {
		t.Errorf("MinRead = %q, want %q", l.MinRead, "%d min read")
	}
}

func TestResolve_Italian(t *testing.T) {
	l := Resolve("it")
	if l.SearchPlaceholder != "Cerca appunti..." {
		t.Errorf("SearchPlaceholder = %q, want %q", l.SearchPlaceholder, "Cerca appunti...")
	}
	if l.BackToTop != "Torna su" {
		t.Errorf("BackToTop = %q, want %q", l.BackToTop, "Torna su")
	}
	if l.PageNotFound != "Pagina non trovata" {
		t.Errorf("PageNotFound = %q, want %q", l.PageNotFound, "Pagina non trovata")
	}
	if l.MinRead != "%d min di lettura" {
		t.Errorf("MinRead = %q, want %q", l.MinRead, "%d min di lettura")
	}
}

func TestResolve_UnknownFallback(t *testing.T) {
	l := Resolve("xx")
	en := Resolve("en")
	if l != en {
		t.Error("unknown language code did not fall back to English")
	}
}

func TestResolve_AllFieldsPopulated(t *testing.T) {
	for code, labels := range languages {
		v := reflect.ValueOf(labels).Elem()
		typ := v.Type()
		for i := range v.NumField() {
			field := v.Field(i)
			if field.Kind() == reflect.String && field.String() == "" {
				t.Errorf("language %q: field %q is empty", code, typ.Field(i).Name)
			}
		}
	}
}
