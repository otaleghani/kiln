// @feature:wikilinks Tests for <picture> element rendering with image optimization variants.
package markdown

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/otaleghani/kiln/internal/imgopt"
)

func TestWriteImage_WithVariants(t *testing.T) {
	r := &IndexResolver{
		ImageResults: map[string]*imgopt.Result{
			"/img/photo": {
				Original: "/img/photo.png",
				Variants: []imgopt.Variant{
					{Width: 800, Format: "webp", WebPath: "/img/photo-800w.webp"},
					{Width: 400, Format: "webp", WebPath: "/img/photo-400w.webp"},
					{Width: 800, Format: "png", WebPath: "/img/photo-800w.png"},
					{Width: 400, Format: "png", WebPath: "/img/photo-400w.png"},
				},
			},
		},
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	r.writeImage(w, "/img/photo", []byte("a photo"), nil)
	w.Flush()
	out := buf.String()

	checks := []string{
		"<picture>",
		"</picture>",
		`<source type="image/webp"`,
		"srcset=",
		`loading="lazy"`,
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot: %s", want, out)
		}
	}
}

func TestWriteImage_NoVariants(t *testing.T) {
	r := &IndexResolver{
		ImageResults: map[string]*imgopt.Result{
			"/img/other": {Original: "/img/other.png"},
		},
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	r.writeImage(w, "/img/notfound", []byte("alt"), nil)
	w.Flush()
	out := buf.String()

	if !strings.Contains(out, `<img`) {
		t.Errorf("expected <img tag, got: %s", out)
	}
	if !strings.Contains(out, `loading="lazy"`) {
		t.Errorf("expected loading=lazy, got: %s", out)
	}
	if strings.Contains(out, "<picture>") {
		t.Errorf("should not contain <picture>, got: %s", out)
	}
}

func TestWriteImage_NilResults(t *testing.T) {
	r := &IndexResolver{
		ImageResults: nil,
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	r.writeImage(w, "/img/photo", []byte("alt"), nil)
	w.Flush()
	out := buf.String()

	if !strings.Contains(out, `<img`) {
		t.Errorf("expected <img tag, got: %s", out)
	}
	if !strings.Contains(out, `loading="lazy"`) {
		t.Errorf("expected loading=lazy, got: %s", out)
	}
	if strings.Contains(out, "<picture>") {
		t.Errorf("should not contain <picture>, got: %s", out)
	}
}
