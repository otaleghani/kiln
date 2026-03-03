// @feature:layouts Template helper functions for templ components.
package templates

import (
	"fmt"
	"strings"
	"time"
)

// FormatDate formats a time.Time as "Jan 02, 2006".
func FormatDate(t time.Time) string {
	return t.Format("Jan 02, 2006")
}

// toStr safely converts an interface value to a string.
// Returns "" for nil values instead of panicking.
func toStr(v any) string {
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return fmt.Sprintf("%v", v)
	}
	return s
}

// buildThemeCSS builds the inline <style> block with CSS custom properties
// for light/dark themes and @font-face declarations.
func buildThemeCSS(theme *ThemeData) string {
	var b strings.Builder
	b.WriteString("<style>\n")

	// Font definitions
	b.WriteString("/* Font definitions */\n")
	b.WriteString(theme.FontFaceCSS)
	b.WriteString("\n\n")

	// Light theme (default)
	b.WriteString("/* Variable definitions */\n:root {\n")
	b.WriteString("/* --- LIGHT THEME (Default) --- */\n")
	writeColorVars(&b, theme.Light)
	b.WriteString("--sidebar-width: 280px;\n")
	b.WriteString(
		fmt.Sprintf(
			"--font-main: %s, -apple-system, BlinkMacSystemFont, \"Segoe UI\", Roboto, Helvetica, Arial, sans-serif;\n",
			theme.FontFamily,
		),
	)
	b.WriteString("}\n\n")

	// Dark theme (data attribute)
	b.WriteString("/* --- DARK THEME OVERRIDES (Data Attribute) --- */\n")
	b.WriteString(":root[data-theme=\"dark\"] {\n")
	writeColorVars(&b, theme.Dark)
	b.WriteString("}\n\n")

	// Dark theme (system preference)
	b.WriteString("/* --- DARK THEME OVERRIDES (System Preference) --- */\n")
	b.WriteString("@media (prefers-color-scheme: dark) {\n")
	b.WriteString(":root:not([data-theme=\"light\"]) {\n")
	writeColorVars(&b, theme.Dark)
	b.WriteString("}\n}\n")

	b.WriteString("</style>")
	return b.String()
}

func writeColorVars(b *strings.Builder, c *ThemeColors) {
	fmt.Fprintf(b, "--bg-color: %s;\n", c.Bg)
	fmt.Fprintf(b, "--text-color: %s;\n", c.Text)
	fmt.Fprintf(b, "--sidebar-bg: %s;\n", c.SidebarBg)
	fmt.Fprintf(b, "--sidebar-border: %s;\n", c.SidebarBorder)
	fmt.Fprintf(b, "--accent-color: %s;\n", c.Accent)
	fmt.Fprintf(b, "--hover-color: %s;\n", c.Hover)
	fmt.Fprintf(b, "/* Palette */\n")
	fmt.Fprintf(b, "--color-red: %s;\n", c.Red)
	fmt.Fprintf(b, "--color-orange: %s;\n", c.Orange)
	fmt.Fprintf(b, "--color-yellow: %s;\n", c.Yellow)
	fmt.Fprintf(b, "--color-green: %s;\n", c.Green)
	fmt.Fprintf(b, "--color-blue: %s;\n", c.Blue)
	fmt.Fprintf(b, "--color-purple: %s;\n", c.Purple)
	fmt.Fprintf(b, "--color-cyan: %s;\n", c.Cyan)
	fmt.Fprintf(b, "--color-comment: %s;\n", c.Comment)
}
