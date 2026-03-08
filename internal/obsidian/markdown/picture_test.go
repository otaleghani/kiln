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
		`sizes="min(65ch, 100vw)"`,
		`<figure class="img-figure">`,
		`</figure>`,
		`<button class="img-expand-btn"`,
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot: %s", want, out)
		}
	}

	// WebP source must appear before PNG source.
	webpIdx := strings.Index(out, `<source type="image/webp"`)
	pngIdx := strings.Index(out, `<source type="image/png"`)
	if webpIdx < 0 || pngIdx < 0 {
		t.Fatalf("expected both webp and png sources in output:\n%s", out)
	}
	if webpIdx >= pngIdx {
		t.Errorf("WebP source (at %d) must appear before PNG source (at %d)\ngot: %s", webpIdx, pngIdx, out)
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
	if strings.Contains(out, `sizes=`) {
		t.Errorf("plain <img> without srcset should not have sizes attribute, got: %s", out)
	}
	if !strings.Contains(out, `<figure class="img-figure">`) {
		t.Errorf("expected <figure> wrapper, got: %s", out)
	}
	if !strings.Contains(out, `<button class="img-expand-btn"`) {
		t.Errorf("expected expand button, got: %s", out)
	}
}

func TestWriteImage_WithAVIFAndWebP(t *testing.T) {
	r := &IndexResolver{
		ImageResults: map[string]*imgopt.Result{
			"/img/photo": {
				Original: "/img/photo.png",
				Variants: []imgopt.Variant{
					{Width: 800, Format: "avif", WebPath: "/img/photo-800w.avif"},
					{Width: 400, Format: "avif", WebPath: "/img/photo-400w.avif"},
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
		`<source type="image/avif"`,
		`<source type="image/webp"`,
		`loading="lazy"`,
		`sizes="min(65ch, 100vw)"`,
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\ngot: %s", want, out)
		}
	}

	// AVIF source must appear before WebP source.
	avifIdx := strings.Index(out, `<source type="image/avif"`)
	webpIdx := strings.Index(out, `<source type="image/webp"`)
	if avifIdx < 0 || webpIdx < 0 {
		t.Fatalf("expected both avif and webp sources in output:\n%s", out)
	}
	if avifIdx >= webpIdx {
		t.Errorf("AVIF source (at %d) must appear before WebP source (at %d)\ngot: %s", avifIdx, webpIdx, out)
	}

	// WebP source must appear before PNG source.
	pngIdx := strings.Index(out, `<source type="image/png"`)
	if pngIdx < 0 {
		t.Fatalf("expected png source in output:\n%s", out)
	}
	if pngIdx <= webpIdx {
		t.Errorf("PNG source (at %d) must appear after WebP source (at %d)\ngot: %s", pngIdx, webpIdx, out)
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
	if !strings.Contains(out, `<figure class="img-figure">`) {
		t.Errorf("expected <figure> wrapper, got: %s", out)
	}
	if !strings.Contains(out, `<button class="img-expand-btn"`) {
		t.Errorf("expected expand button, got: %s", out)
	}
}

func TestWriteImage_FullOutputStructure(t *testing.T) {
	r := &IndexResolver{
		ImageResults: map[string]*imgopt.Result{
			"/img/hero.png": {
				Original: "/img/hero.png",
				Variants: []imgopt.Variant{
					{Width: 1200, Format: "avif", WebPath: "/img/hero-1200w.avif"},
					{Width: 800, Format: "avif", WebPath: "/img/hero-800w.avif"},
					{Width: 1200, Format: "webp", WebPath: "/img/hero-1200w.webp"},
					{Width: 800, Format: "webp", WebPath: "/img/hero-800w.webp"},
					{Width: 1200, Format: "png", WebPath: "/img/hero-1200w.png"},
					{Width: 800, Format: "png", WebPath: "/img/hero-800w.png"},
				},
			},
		},
	}

	var buf bytes.Buffer
	w := bufio.NewWriter(&buf)
	r.writeImage(w, "/img/hero.png", []byte("Hero image"), nil)
	w.Flush()
	out := buf.String()

	// Verify structural order: figure > picture > sources > img > button > /figure
	checks := []struct{ name, substr string }{
		{"figure open", `<figure class="img-figure">`},
		{"picture open", `<picture>`},
		{"avif source with sizes", `<source type="image/avif" srcset="`},
		{"sizes on source", `sizes="min(65ch, 100vw)"`},
		{"webp source", `<source type="image/webp"`},
		{"png source", `<source type="image/png"`},
		{"fallback img", `<img src="/img/hero.png"`},
		{"picture close", `</picture>`},
		{"expand button", `<button class="img-expand-btn"`},
		{"figure close", `</figure>`},
	}
	for _, c := range checks {
		if !strings.Contains(out, c.substr) {
			t.Errorf("%s: missing %q\ngot: %s", c.name, c.substr, out)
		}
	}

	// Verify ordering: figure before picture before button before /figure
	figureIdx := strings.Index(out, `<figure`)
	pictureIdx := strings.Index(out, `<picture>`)
	buttonIdx := strings.Index(out, `<button class="img-expand-btn"`)
	closeFigureIdx := strings.Index(out, `</figure>`)

	if figureIdx >= pictureIdx {
		t.Error("figure must come before picture")
	}
	if pictureIdx >= buttonIdx {
		t.Error("picture must come before button")
	}
	if buttonIdx >= closeFigureIdx {
		t.Error("button must come before closing figure")
	}
}
