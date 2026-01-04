package builder

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/otaleghani/kiln/assets"
)

// ThemeColors defines the specific color palette for a UI state (Light or Dark).
// It includes semantic colors (used for specific UI elements) and a generic palette.
type ThemeColors struct {
	// Semantic UI Colors
	Bg            string // Main background color
	Text          string // Main text color
	SidebarBg     string // Background color for navigation sidebars
	SidebarBorder string // Border color separating sidebar from content
	Accent        string // Primary accent color for links and active states
	Hover         string // Background color for hovered elements
	Comment       string // Text color for code comments or secondary text

	// Generic Palette (used for syntax highlighting or custom badges)
	Red    string
	Orange string
	Yellow string
	Green  string
	Blue   string
	Purple string
	Cyan   string
}

// Theme represents a complete visual style configuration.
// It bundles color schemas for both Light and Dark modes, along with typography settings.
type Theme struct {
	Light *ThemeColors
	Dark  *ThemeColors
	Font  *FontData
}

// FontData holds the metadata and CSS required to render a specific font family.
type FontData struct {
	Family   string   // The CSS font-family string (e.g., "'Inter', sans-serif")
	Files    []string // List of filenames (e.g., .woff2) that need to be extracted
	FontFace string   // The raw CSS @font-face declaration to inject into the stylesheet
}

// fonts is a registry of available font configurations supported by the builder.
// It maps a short string ID to the specific FontData.
var fonts = map[string]*FontData{
	"inter": {
		Family: "'Inter', sans-serif",
		Files:  []string{"Inter-Regular.woff2", "Inter-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'Inter';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('./Inter-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Inter';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('./Inter-Bold.woff2') format('woff2');
			}
		`,
	},
	"lato": {
		Family: "'Lato', sans-serif",
		Files:  []string{"Lato-Regular.woff2", "Lato-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'Lato';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('./Lato-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Lato';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('./Lato-Bold.woff2') format('woff2');
			}
		`,
	},
	"merriweather": {
		Family: "'Merriweather', serif",
		Files:  []string{"Merriweather-Regular.woff2", "Merriweather-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'Merriweather';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('./Merriweather-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Merriweather';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('./Merriweather-Bold.woff2') format('woff2');
			}
		`,
	},
	// System fonts rely on the OS font stack and require no external files.
	"system": {
		Family:   "-apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif",
		Files:    []string{},
		FontFace: "",
	},
}

// themes is a registry of predefined color schemes.
var themes = map[string]*Theme{
	"default": {
		Light: &ThemeColors{
			Bg: "#ffffff", Text: "#2e3338", SidebarBg: "#f5f6f8", SidebarBorder: "#e6e6e6", Accent: "#7e6df7", Hover: "#e8e8e8", Comment: "#9ca3af",
			Red: "#d14d41", Orange: "#d67e22", Yellow: "#dcb20a", Green: "#409c5a", Blue: "#4b8ad6", Purple: "#9463c6", Cyan: "#3aaeb4",
		},
		Dark: &ThemeColors{
			Bg: "#1e1e1e", Text: "#dcddde", SidebarBg: "#252526", SidebarBorder: "#2f2f31", Accent: "#7e6df7", Hover: "#333333", Comment: "#6b7280",
			Red: "#f07178", Orange: "#f78c6c", Yellow: "#ffcb6b", Green: "#c3e88d", Blue: "#82aaff", Purple: "#c792ea", Cyan: "#89ddff",
		},
	},
	"dracula": {
		Light: &ThemeColors{
			Bg: "#f8f8f2", Text: "#282a36", SidebarBg: "#e4e4db", SidebarBorder: "#bcc2cd", Accent: "#ff79c6", Hover: "#d6d6d6", Comment: "#6272a4",
			Red: "#ff5555", Orange: "#ffb86c", Yellow: "#f1fa8c", Green: "#50fa7b", Blue: "#6272a4", Purple: "#bd93f9", Cyan: "#8be9fd",
		},
		Dark: &ThemeColors{
			Bg: "#282a36", Text: "#f8f8f2", SidebarBg: "#21222c", SidebarBorder: "#44475a", Accent: "#ff79c6", Hover: "#44475a", Comment: "#6272a4",
			Red: "#ff5555", Orange: "#ffb86c", Yellow: "#f1fa8c", Green: "#50fa7b", Blue: "#8be9fd", Purple: "#bd93f9", Cyan: "#8be9fd",
		},
	},
	"catppuccin": {
		Light: &ThemeColors{
			Bg: "#eff1f5", Text: "#4c4f69", SidebarBg: "#e6e9ef", SidebarBorder: "#ccd0da", Accent: "#8839ef", Hover: "#acb0be", Comment: "#9ca0b0",
			Red: "#d20f39", Orange: "#fe640b", Yellow: "#df8e1d", Green: "#40a02b", Blue: "#1e66f5", Purple: "#8839ef", Cyan: "#04a5e5",
		},
		Dark: &ThemeColors{
			Bg: "#1e1e2e", Text: "#cdd6f4", SidebarBg: "#181825", SidebarBorder: "#313244", Accent: "#cba6f7", Hover: "#45475a", Comment: "#a6adc8",
			Red: "#f38ba8", Orange: "#fab387", Yellow: "#f9e2af", Green: "#a6e3a1", Blue: "#89b4fa", Purple: "#cba6f7", Cyan: "#89dceb",
		},
	},
	"nord": {
		Light: &ThemeColors{
			Bg: "#eceff4", Text: "#2e3440", SidebarBg: "#e5e9f0", SidebarBorder: "#d8dee9", Accent: "#5e81ac", Hover: "#d8dee9", Comment: "#4c566a",
			Red: "#bf616a", Orange: "#d08770", Yellow: "#ebcb8b", Green: "#a3be8c", Blue: "#5e81ac", Purple: "#b48ead", Cyan: "#88c0d0",
		},
		Dark: &ThemeColors{
			Bg: "#2e3440", Text: "#d8dee9", SidebarBg: "#242933", SidebarBorder: "#3b4252", Accent: "#88c0d0", Hover: "#434c5e", Comment: "#4c566a",
			Red: "#bf616a", Orange: "#d08770", Yellow: "#ebcb8b", Green: "#a3be8c", Blue: "#81a1c1", Purple: "#b48ead", Cyan: "#88c0d0",
		},
	},
}

// resolveTheme looks up a theme by name.
// If the theme is not found, it defaults to "default" and logs a warning.
// It also resolves the associated font using resolveFont.
func resolveTheme(themeName, fontName string) *Theme {
	theme, ok := themes[strings.ToLower(themeName)]
	if !ok {
		log.Printf("Theme '%s' not found, falling back to default.", themeName)
		theme = themes["default"]
	}
	theme.Font = resolveFont(fontName)
	return theme
}

// resolveFont looks up font data by name.
// If the font is not found, it defaults to "inter" and logs a warning.
func resolveFont(name string) *FontData {
	font, ok := fonts[strings.ToLower(name)]
	if !ok {
		log.Printf("Font '%s' not found, falling back to inter.", name)
		return fonts["inter"]
	}
	return font
}

// extractFonts writes the font files associated with the given FontData to disk.
// It reads the files from the embedded assets filesystem and writes them to fontsDir.
func (t *Theme) extractFonts(fontsDir string) {
	// If the font has no associated files (e.g., System fonts), return immediately.
	if len(t.Font.Files) == 0 {
		return
	}

	// Ensure the target directory exists.
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		log.Printf("Failed to create fonts directory: %s\n", err.Error())
	}

	for _, fileName := range t.Font.Files {
		// Read the binary content from the embedded assets FS.
		content, err := assets.TemplateFS.ReadFile(fileName)
		if err != nil {
			log.Printf("Failed to read embedded font %s: %s", fileName, err.Error())
		}

		// Write the binary content to the user's filesystem.
		destPath := filepath.Join(fontsDir, fileName)
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			log.Printf("Failed to write font file %s: %s", fileName, err.Error())
		}
	}
}
