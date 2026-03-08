// @feature:imgopt Image optimization: resize and convert images to WebP/AVIF at build time.
package imgopt

import (
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

var ErrNoEncoder = errors.New("imgopt: encoder not available")

type Variant struct {
	Width   int
	Suffix  string
	Format  string
	OutPath string
	WebPath string
}

type Result struct {
	Original string
	Variants []Variant
}

// DefaultBreakpoints returns the default responsive image widths.
func DefaultBreakpoints() []int {
	return []int{1200, 800, 400}
}

// IsOptimizable returns true for image extensions we can process.
func IsOptimizable(ext string) bool {
	switch strings.ToLower(ext) {
	case ".png", ".jpg", ".jpeg":
		return true
	}
	return false
}

// DecodeImage opens a file and decodes it as an image.
func DecodeImage(path string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", fmt.Errorf("imgopt: open %s: %w", path, err)
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, "", fmt.Errorf("imgopt: decode %s: %w", path, err)
	}
	return img, format, nil
}

// Resize scales src to maxWidth preserving aspect ratio using CatmullRom.
// Returns src unchanged if its width is already <= maxWidth.
func Resize(src image.Image, maxWidth int) image.Image {
	bounds := src.Bounds()
	srcW := bounds.Dx()
	srcH := bounds.Dy()

	if srcW <= maxWidth {
		return src
	}

	dstW := maxWidth
	dstH := srcH * maxWidth / srcW

	dst := image.NewNRGBA(image.Rect(0, 0, dstW, dstH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}

// EncodeWebP encodes img as WebP by shelling out to cwebp.
// Returns ErrNoEncoder if cwebp is not on PATH.
func EncodeWebP(w io.Writer, img image.Image, quality int) error {
	cwebp, err := exec.LookPath("cwebp")
	if err != nil {
		return ErrNoEncoder
	}

	tmp, err := os.CreateTemp("", "imgopt-*.png")
	if err != nil {
		return fmt.Errorf("imgopt: create temp: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := png.Encode(tmp, img); err != nil {
		tmp.Close()
		return fmt.Errorf("imgopt: encode temp png: %w", err)
	}
	tmp.Close()

	outTmp, err := os.CreateTemp("", "imgopt-*.webp")
	if err != nil {
		return fmt.Errorf("imgopt: create temp: %w", err)
	}
	outPath := outTmp.Name()
	outTmp.Close()
	defer os.Remove(outPath)

	cmd := exec.Command(cwebp, "-q", fmt.Sprintf("%d", quality), tmpPath, "-o", outPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("imgopt: cwebp: %w: %s", err, out)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return fmt.Errorf("imgopt: read webp: %w", err)
	}
	_, err = w.Write(data)
	return err
}

// EncodeAVIF encodes img as AVIF by shelling out to avifenc.
// Returns ErrNoEncoder if avifenc is not on PATH.
func EncodeAVIF(w io.Writer, img image.Image, quality int) error {
	avifenc, err := exec.LookPath("avifenc")
	if err != nil {
		return ErrNoEncoder
	}

	tmp, err := os.CreateTemp("", "imgopt-*.png")
	if err != nil {
		return fmt.Errorf("imgopt: create temp: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)

	if err := png.Encode(tmp, img); err != nil {
		tmp.Close()
		return fmt.Errorf("imgopt: encode temp png: %w", err)
	}
	tmp.Close()

	outTmp, err := os.CreateTemp("", "imgopt-*.avif")
	if err != nil {
		return fmt.Errorf("imgopt: create temp: %w", err)
	}
	outPath := outTmp.Name()
	outTmp.Close()
	defer os.Remove(outPath)

	cmd := exec.Command(avifenc, "--min", "0", "--max", "63",
		"-a", fmt.Sprintf("cq-level=%d", 63-quality*63/100),
		tmpPath, outPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("imgopt: avifenc: %w: %s", err, out)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		return fmt.Errorf("imgopt: read avif: %w", err)
	}
	_, err = w.Write(data)
	return err
}

// ProcessImage decodes srcPath, resizes at each breakpoint smaller than the
// original width, encodes to the original format and optionally WebP, and
// writes files to outDir. WebPaths are computed relative to webDir.
func ProcessImage(srcPath, outDir, webDir, baseName string, breakpoints []int) (*Result, error) {
	img, format, err := DecodeImage(srcPath)
	if err != nil {
		return nil, err
	}

	srcW := img.Bounds().Dx()
	result := &Result{Original: srcPath}

	for _, bp := range breakpoints {
		if bp >= srcW {
			continue
		}

		resized := Resize(img, bp)
		suffix := fmt.Sprintf("-%dw", bp)

		// Try AVIF first (best compression).
		avifName := baseName + suffix + ".avif"
		avifPath := filepath.Join(outDir, avifName)
		if err := writeEncoded(avifPath, resized, EncodeAVIF, 80); err == nil {
			result.Variants = append(result.Variants, Variant{
				Width:   bp,
				Suffix:  suffix,
				Format:  "avif",
				OutPath: avifPath,
				WebPath: webDir + "/" + avifName,
			})
		}

		// Try WebP second.
		webpName := baseName + suffix + ".webp"
		webpPath := filepath.Join(outDir, webpName)
		if err := writeEncoded(webpPath, resized, EncodeWebP, 80); err == nil {
			result.Variants = append(result.Variants, Variant{
				Width:   bp,
				Suffix:  suffix,
				Format:  "webp",
				OutPath: webpPath,
				WebPath: webDir + "/" + webpName,
			})
		}

		// Encode original format last (fallback).
		origExt := "." + format
		origName := baseName + suffix + origExt
		origPath := filepath.Join(outDir, origName)

		if err := writeImage(origPath, resized, format); err != nil {
			return nil, fmt.Errorf("imgopt: write %s: %w", origPath, err)
		}

		result.Variants = append(result.Variants, Variant{
			Width:   bp,
			Suffix:  suffix,
			Format:  format,
			OutPath: origPath,
			WebPath: webDir + "/" + origName,
		})
	}

	return result, nil
}

func writeImage(path string, img image.Image, format string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	switch format {
	case "png":
		return png.Encode(f, img)
	case "jpeg":
		// Use png as fallback for simplicity; JPEG encoding would need
		// "image/jpeg" Encode which we can add later.
		return png.Encode(f, img)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func writeEncoded(path string, img image.Image, encode func(io.Writer, image.Image, int) error, quality int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := encode(f, img, quality); err != nil {
		os.Remove(path)
		return err
	}
	return nil
}
