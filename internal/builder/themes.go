package builder

import (
	"html/template"
	"log/slog"
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
	Family           template.CSS // The CSS font-family string (e.g., "'Inter', sans-serif")
	Files            []string     // List of filenames (e.g., .woff2) that need to be extracted
	FontFace         template.CSS // The raw CSS @font-face declaration to inject into the stylesheet
	FontFaceReplaced template.CSS
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
				src: url('{{.Site.BaseURL}}/Inter-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Inter';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/Inter-Bold.woff2') format('woff2');
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
				src: url('{{.Site.BaseURL}}/Lato-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Lato';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/Lato-Bold.woff2') format('woff2');
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
				src: url('{{.Site.BaseURL}}/Merriweather-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Merriweather';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/Merriweather-Bold.woff2') format('woff2');
			}
		`,
	},
	"lora": {
		Family: "'Lora', serif",
		Files:  []string{"Lora-Regular.woff2", "Lora-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'Lora';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/Lora-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Lora';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/Lora-Bold.woff2') format('woff2');
			}
		`,
	},
	"libre-baskerville": {
		Family: "'Libre Baskerville', serif",
		Files:  []string{"LibreBaskerville-Regular.woff2", "LibreBaskerville-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'Libre Baskerville';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/LibreBaskerville-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Libre Baskerville';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/LibreBaskerville-Bold.woff2') format('woff2');
			}
		`,
	},
	"noto-serif": {
		Family: "'Noto Serif', serif",
		Files:  []string{"NotoSerif-Regular.woff2", "NotoSerif-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'Noto Serif';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/NotoSerif-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Noto Serif';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/NotoSerif-Bold.woff2') format('woff2');
			}
		`,
	},
	"ibm-plex-sans": {
		Family: "'IBM Plex Sans', sans-serif",
		Files:  []string{"IBMPlexSans-Regular.woff2", "IBMPlexSans-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'IBM Plex Sans';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/IBMPlexSans-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'IBM Plex Sans';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/IBMPlexSans-Bold.woff2') format('woff2');
			}
		`,
	},
	"google-sans": {
		Family: "'Google Sans', sans-serif",
		Files:  []string{"GoogleSans-Regular.woff2", "GoogleSans-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'Google Sans';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/GoogleSans-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Google Sans';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/GoogleSans-Bold.woff2') format('woff2');
			}
		`,
	},
	"roboto": {
		Family: "'Roboto', sans-serif",
		Files:  []string{"Roboto-Regular.woff2", "Roboto-Bold.woff2"},
		FontFace: `
			@font-face {
				font-family: 'Roboto';
				font-style: normal;
				font-weight: 400;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/Roboto-Regular.woff2') format('woff2');
			}
			@font-face {
				font-family: 'Roboto';
				font-style: normal;
				font-weight: 700;
				font-display: swap;
				src: url('{{.Site.BaseURL}}/Roboto-Bold.woff2') format('woff2');
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
	"tokyonight": {
		Light: &ThemeColors{
			Bg: "#d5d6db", Text: "#343b58", SidebarBg: "#cfc9c2", SidebarBorder: "#b4b5b9", Accent: "#7aa2f7", Hover: "#c0c0c0", Comment: "#565f89",
			Red: "#f7768e", Orange: "#ff9e64", Yellow: "#e0af68", Green: "#73daca", Blue: "#7aa2f7", Purple: "#bb9af7", Cyan: "#7dcfff",
		},
		Dark: &ThemeColors{
			Bg: "#1a1b26", Text: "#c0caf5", SidebarBg: "#16161e", SidebarBorder: "#414868", Accent: "#7aa2f7", Hover: "#292e42", Comment: "#565f89",
			Red: "#f7768e", Orange: "#ff9e64", Yellow: "#e0af68", Green: "#9ece6a", Blue: "#7aa2f7", Purple: "#bb9af7", Cyan: "#7dcfff",
		},
	},
	"rose-pine": {
		Light: &ThemeColors{ // Rosé Pine Dawn
			Bg: "#faf4ed", Text: "#575279", SidebarBg: "#fffaf3", SidebarBorder: "#cecacd", Accent: "#eb6f92", Hover: "#f2e9e1", Comment: "#9893a5",
			Red: "#b4637a", Orange: "#d7827e", Yellow: "#ea9d34", Green: "#286983", Blue: "#56949f", Purple: "#907aa9", Cyan: "#d7827e",
		},
		Dark: &ThemeColors{ // Rosé Pine Main
			Bg: "#191724", Text: "#e0def4", SidebarBg: "#1f1d2e", SidebarBorder: "#403d52", Accent: "#eb6f92", Hover: "#26233a", Comment: "#6e6a86",
			Red: "#eb6f92", Orange: "#ebbcba", Yellow: "#f6c177", Green: "#31748f", Blue: "#9ccfd8", Purple: "#c4a7e7", Cyan: "#ebbcba",
		},
	},
	"gruvbox": {
		Light: &ThemeColors{
			Bg: "#fbf1c7", Text: "#3c3836", SidebarBg: "#ebdbb2", SidebarBorder: "#d5c4a1", Accent: "#d65d0e", Hover: "#ebdbb2", Comment: "#928374",
			Red: "#cc241d", Orange: "#d65d0e", Yellow: "#d79921", Green: "#98971a", Blue: "#458588", Purple: "#b16286", Cyan: "#689d6a",
		},
		Dark: &ThemeColors{
			Bg: "#282828", Text: "#ebdbb2", SidebarBg: "#3c3836", SidebarBorder: "#504945", Accent: "#fe8019", Hover: "#3c3836", Comment: "#a89984",
			Red: "#fb4934", Orange: "#fe8019", Yellow: "#fabd2f", Green: "#b8bb26", Blue: "#83a598", Purple: "#d3869b", Cyan: "#8ec07c",
		},
	},
	"everforest": {
		Light: &ThemeColors{
			Bg: "#fdf6e3", Text: "#5c6a72", SidebarBg: "#f3efda", SidebarBorder: "#e6e2cc", Accent: "#f57d26", Hover: "#edf3e8", Comment: "#939f91",
			Red: "#f85552", Orange: "#f57d26", Yellow: "#dfa000", Green: "#8da101", Blue: "#3a94c5", Purple: "#df69ba", Cyan: "#35a77c",
		},
		Dark: &ThemeColors{
			Bg: "#2b3339", Text: "#d3c6aa", SidebarBg: "#323c41", SidebarBorder: "#4a555b", Accent: "#a7c080", Hover: "#374247", Comment: "#859289",
			Red: "#e67e80", Orange: "#e69875", Yellow: "#dbbc7f", Green: "#a7c080", Blue: "#7fbbb3", Purple: "#d699b6", Cyan: "#83c092",
		},
	},
	"cyberdream": {
		Light: &ThemeColors{
			Bg: "#ffffff", Text: "#16181a", SidebarBg: "#eff1f5", SidebarBorder: "#7b8496", Accent: "#5ea1ff", Hover: "#e4e5e8", Comment: "#7b8496",
			Red: "#ff6e5e", Orange: "#ffbd5e", Yellow: "#f2cdcd", Green: "#5eff6c", Blue: "#5ea1ff", Purple: "#bd5eff", Cyan: "#5ef1ff",
		},
		Dark: &ThemeColors{
			Bg: "#16181a", Text: "#ffffff", SidebarBg: "#1e2124", SidebarBorder: "#3c4048", Accent: "#5ea1ff", Hover: "#272a2d", Comment: "#7b8496",
			Red: "#ff6e5e", Orange: "#ffbd5e", Yellow: "#f2cdcd", Green: "#5eff6c", Blue: "#5ea1ff", Purple: "#bd5eff", Cyan: "#5ef1ff",
		},
	},
}

// ResolveTheme looks up a theme by name.
// If the theme is not found, it defaults to "default" and logs a warning.
// It also resolves the associated font using resolveFont.
func ResolveTheme(themeName, fontName string, log *slog.Logger) *Theme {
	theme, ok := themes[strings.ToLower(themeName)]
	if !ok {
		log.Warn("Theme not found. Using default theme.", "name", themeName)
		theme = themes["default"]
	}
	theme.Font = resolveFont(fontName, log)
	return theme
}

// resolveFont looks up font data by name.
// If the font is not found, it defaults to "inter" and logs a warning.
func resolveFont(name string, log *slog.Logger) *FontData {
	font, ok := fonts[strings.ToLower(name)]
	if !ok {
		log.Warn("Font not found. Using default theme.", "name", name)
		return fonts["inter"]
	}
	return font
}

// extractFonts writes the font files associated with the given FontData to disk.
// It reads the files from the embedded assets filesystem and writes them to fontsDir.
func (t *Theme) extractFonts(fontsDir string, log *slog.Logger) {
	// If the font has no associated files (e.g., System fonts), return immediately.
	if len(t.Font.Files) == 0 {
		return
	}

	// Ensure the target directory exists.
	if err := os.MkdirAll(fontsDir, 0755); err != nil {
		log.Error("Failed to create fonts directory", "error", err)
	}

	for _, fileName := range t.Font.Files {
		// Read the binary content from the embedded assets FS.
		content, err := assets.TemplateFS.ReadFile(fileName)
		if err != nil {
			log.Error("Failed to read embed font", "error", err)
		}

		// Write the binary content to the user's filesystem.
		destPath := filepath.Join(fontsDir, fileName)
		if err := os.WriteFile(destPath, content, 0644); err != nil {
			log.Error("Failed to write font", "error", err)
		}
	}
}
