// kiln-palette is a utility used to generate PNG preview of a theme palette for the documentation.
//
// Usage: go run ./cmd/kiln-palette/main.go --theme "dracula"
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log/slog"
	"os"
	"reflect"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/otaleghani/kiln/internal/builder"
	"github.com/spf13/cobra"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	FlagTheme        = "theme"
	FlagThemeShort   = "t"
	DefaultThemeName = "default"
)

var themeName string

func main() {
	var cmdGenerate = &cobra.Command{
		Use:   "gen-palette",
		Short: "Generate a PNG preview of a theme palette",
		Run: func(cmd *cobra.Command, args []string) {
			runGenerate()
		},
	}

	cmdGenerate.Flags().
		StringVarP(&themeName, FlagTheme, FlagThemeShort, DefaultThemeName, "Color theme (default, dracula, catppuccin, nord)")

	if err := cmdGenerate.Execute(); err != nil {
		log.Fatal(err)
	}
}

// Logic
func runGenerate() {
	handler := log.New(os.Stderr)
	handler.SetReportTimestamp(true)
	handler.SetFormatter(log.TextFormatter)
	theme := builder.ResolveTheme(themeName, "", slog.New(handler))
	handler.Info("Generating palette for theme", "name", themeName)

	// Setup Image Canvas
	// Width: 800 (400 Light | 400 Dark)
	// Height: Header + (Rows * RowHeight)
	const (
		width     = 500
		rowHeight = 60
		headerH   = 60
		colWidth  = width / 2
		padding   = 20
	)

	// Get list of fields to iterate in order
	fields := []string{
		"Bg", "Text", "SidebarBg", "SidebarBorder", "Accent", "Hover", "Comment",
		"Red", "Orange", "Yellow", "Green", "Blue", "Purple", "Cyan",
	}

	height := headerH + (len(fields) * rowHeight)
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Draw Backgrounds (Split Canvas)
	// Light Mode Background
	draw.Draw(
		img,
		image.Rect(0, 0, colWidth, height),
		&image.Uniform{parseHex(theme.Light.Bg)},
		image.Point{},
		draw.Src,
	)
	// Dark Mode Background
	draw.Draw(
		img,
		image.Rect(colWidth, 0, width, height),
		&image.Uniform{parseHex(theme.Dark.Bg)},
		image.Point{},
		draw.Src,
	)

	// Draw Content
	drawLabel := func(x, y int, text string, c color.Color) {
		point := fixed.Point26_6{X: fixed.I(x), Y: fixed.I(y)}
		d := &font.Drawer{
			Dst:  img,
			Src:  image.NewUniform(c),
			Face: basicfont.Face7x13,
			Dot:  point,
		}
		d.DrawString(text)
	}

	// Draw Headers
	drawLabel(
		padding,
		35,
		fmt.Sprintf("THEME: %s (Light)", strings.ToUpper(themeName)),
		parseHex(theme.Light.Text),
	)
	drawLabel(
		colWidth+padding,
		35,
		fmt.Sprintf("THEME: %s (Dark)", strings.ToUpper(themeName)),
		parseHex(theme.Dark.Text),
	)

	// Draw Rows
	v := reflect.ValueOf(*theme.Light)
	vDark := reflect.ValueOf(*theme.Dark)

	for i, fieldName := range fields {
		y := headerH + (i * rowHeight)

		// Light Column
		hexLight := v.FieldByName(fieldName).String()
		colLight := parseHex(hexLight)
		// Draw Swatch
		draw.Draw(
			img,
			image.Rect(padding, y, padding+40, y+40),
			&image.Uniform{colLight},
			image.Point{},
			draw.Src,
		)
		// Draw Border around swatch (for visibility on similar backgrounds)
		drawBorder(img, padding, y, 40, 40, color.Gray{Y: 128})
		// Draw Text
		labelColor := parseHex(theme.Light.Text)
		drawLabel(padding+50, y+25, fmt.Sprintf("%-14s %s", fieldName, hexLight), labelColor)

		// Dark Column
		hexDark := vDark.FieldByName(fieldName).String()
		colDark := parseHex(hexDark)
		// Draw Swatch
		draw.Draw(
			img,
			image.Rect(colWidth+padding, y, colWidth+padding+40, y+40),
			&image.Uniform{colDark},
			image.Point{},
			draw.Src,
		)
		// Draw Border
		drawBorder(img, colWidth+padding, y, 40, 40, color.Gray{Y: 128})
		// Draw Text
		labelColorDark := parseHex(theme.Dark.Text)
		drawLabel(
			colWidth+padding+50,
			y+25,
			fmt.Sprintf("%-14s %s", fieldName, hexDark),
			labelColorDark,
		)
	}

	// Save File
	filename := fmt.Sprintf("palette_%s.png", themeName)
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}

	handler.Info("Success!", "path", filename)
}

// Helpers
// parseHex converts a hex string (e.g. "#ffffff") to color.RGBA
func parseHex(s string) color.RGBA {
	c := color.RGBA{A: 255}
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 6:
		fmt.Sscanf(s, "%02x%02x%02x", &c.R, &c.G, &c.B)
	case 3:
		fmt.Sscanf(s, "%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the nibbles
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		// Return pink on error
		return color.RGBA{R: 255, G: 0, B: 255, A: 255}
	}
	return c
}

// drawBorder draws a simple 1px border
func drawBorder(img *image.RGBA, x, y, w, h int, c color.Color) {
	col := image.NewUniform(c)
	// Top
	draw.Draw(img, image.Rect(x, y, x+w, y+1), col, image.Point{}, draw.Src)
	// Bottom
	draw.Draw(img, image.Rect(x, y+h-1, x+w, y+h), col, image.Point{}, draw.Src)
	// Left
	draw.Draw(img, image.Rect(x, y, x+1, y+h), col, image.Point{}, draw.Src)
	// Right
	draw.Draw(img, image.Rect(x+w-1, y, x+w, y+h), col, image.Point{}, draw.Src)
}
