// @feature:ogimage Build-time Open Graph and Twitter Card image generator.
package ogimage

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	ogWidth       = 1200
	ogHeight      = 630
	twitterWidth  = 1200
	twitterHeight = 600

	accentBarHeight = 8
	padding         = 60
	titleY          = 120
	descY           = 260
	siteNameY       = 560
	maxTitleChars   = 80
	wrapWidth       = 50
)

type ImageConfig struct {
	Title       string
	Description string
	SiteName    string
	AccentColor string
	BgColor     string
	TextColor   string
	Face        font.Face
}

// GenerateOGImage creates a 1200x630 branded PNG for Open Graph.
func GenerateOGImage(cfg ImageConfig, outPath string) error {
	return generateImage(cfg, outPath, ogWidth, ogHeight)
}

// GenerateTwitterImage creates a 1200x600 branded PNG for Twitter Cards.
func GenerateTwitterImage(cfg ImageConfig, outPath string) error {
	return generateImage(cfg, outPath, twitterWidth, twitterHeight)
}

func generateImage(cfg ImageConfig, outPath string, width, height int) error {
	if cfg.Title == "" {
		cfg.Title = "Untitled"
	}

	bg := parseHex(cfg.BgColor)
	accent := parseHex(cfg.AccentColor)
	text := parseHex(cfg.TextColor)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Background
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)

	// Accent bar at top
	draw.Draw(
		img,
		image.Rect(0, 0, width, accentBarHeight),
		&image.Uniform{accent},
		image.Point{},
		draw.Src,
	)

	face := cfg.Face
	if face == nil {
		face = basicfont.Face7x13
	}

	// Title — wrap long text
	titleLines := wrapText(cfg.Title, wrapWidth)
	lineHeight := 42
	for i, line := range titleLines {
		// Draw each character twice offset for a faux-bold effect
		drawLabel(img, face, padding, titleY+i*lineHeight, line, text)
		drawLabel(img, face, padding+1, titleY+i*lineHeight, line, text)
	}

	// Description
	if cfg.Description != "" {
		descLines := wrapText(cfg.Description, wrapWidth+10)
		descStartY := descY + len(titleLines)*lineHeight
		for i, line := range descLines {
			drawLabel(img, face, padding, descStartY+i*lineHeight, line, text)
		}
	}

	// Site name at bottom
	if cfg.SiteName != "" {
		siteY := height - padding
		drawLabel(img, face, padding, siteY, cfg.SiteName, accent)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("ogimage: create file: %w", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		return fmt.Errorf("ogimage: encode png: %w", err)
	}
	return nil
}

func drawLabel(img *image.RGBA, face font.Face, x, y int, text string, c color.RGBA) {
	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(c),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)},
	}
	d.DrawString(text)
}

func parseHex(s string) color.RGBA {
	c := color.RGBA{A: 255}
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 6:
		fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 3:
		fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		return color.RGBA{R: 255, G: 0, B: 255, A: 255}
	}
	return c
}

func wrapText(text string, maxWidth int) []string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{text}
	}

	var lines []string
	current := words[0]

	for _, w := range words[1:] {
		if len(current)+1+len(w) > maxWidth {
			lines = append(lines, current)
			current = w
		} else {
			current += " " + w
		}
	}
	lines = append(lines, current)
	return lines
}
