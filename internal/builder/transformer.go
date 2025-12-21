package builder

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	// "golang.org/x/text/cases"
)

// Contains logic for transforming specific markdown extensions into rich HTML components.

// transformHighlights converts Obsidian-style highlight syntax (==text==)
// into standard HTML <mark> tags.
func transformHighlights(htmlStr string) string {
	re := regexp.MustCompile(`==([^=]+)==`)
	return re.ReplaceAllString(htmlStr, `<mark>$1</mark>`)
}

// transformMermaid locates code blocks designated as "mermaid" diagrams
// and converts them into a container div suitable for client-side rendering (e.g., via mermaid.js).
func transformMermaid(htmlStr string) string {
	// Regex looks for: <pre><code class="language-mermaid"> ...content... </code></pre>
	re := regexp.MustCompile(
		`(?s)<pre[^>]*>\s*<code class="language-mermaid"[^>]*>(.*?)</code>\s*</pre>`,
	)
	return re.ReplaceAllStringFunc(htmlStr, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		content := submatches[1]

		// We store the original un-rendered Mermaid syntax in 'data-original'.
		// This is crucial for themes that support dynamic re-rendering (e.g., switching between light/dark mode diagrams).
		encoded := template.HTMLEscapeString(content)
		return fmt.Sprintf(`<div class="mermaid" data-original="%s">%s</div>`, encoded, content)
	})
}

// transformCallouts processes Obsidian-style "admonition" or "callout" blocks.
// Syntax: > [!type]+/- Title
// It supports collapsible sections (<details>) and static boxes (<div>).
func transformCallouts(htmlStr string) string {
	// Regex breakdown:
	// Group 1: Type (e.g., "info", "warning") -> `\[!([\w-]+)\]`
	// Group 2: Fold Modifier ("+" for open, "-" for closed, or empty) -> `([-+]?)`
	// Group 3: Title (text remaining on the first line) -> `(.*?)`
	// Group 4: Body (everything after the title line until blockquote ends) -> `(.*?)`
	re := regexp.MustCompile(
		`(?s)<blockquote[^>]*>\s*<p>\s*\[!([\w-]+)\]([-+]?)\s*(.*?)(?:</p>|<br\s*/?>|\n)(.*?)</blockquote>`,
	)

	return re.ReplaceAllStringFunc(htmlStr, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 5 {
			return match
		}

		cType := strings.ToLower(submatches[1])
		foldParams := submatches[2] // "+" or "-" or ""
		cTitle := strings.TrimSpace(submatches[3])
		cBody := strings.TrimSpace(submatches[4])

		// Fallback: Use the callout type as the title if no specific title is provided.
		// e.g., > [!info] -> Title becomes "Info"
		if cTitle == "" {
			// Note: strings.Title is deprecated in favor of cases.Title, but used here for simplicity/legacy support.
			cTitle = strings.Title(cType)
		}

		// HTML Cleanup:
		// Because the regex splits the content at the first newline/<br>, we might break the surrounding <p> tags.
		// This check ensures the body starts with a <p> if it ends with </p> but lost its opener.
		if strings.HasSuffix(cBody, "</p>") && !strings.HasPrefix(cBody, "<p>") {
			cBody = "<p>" + cBody
		}

		iconName := getCalloutIcon(cType)

		// Determine collapsible state based on modifiers:
		// "-" : Collapsed by default
		// "+" : Expanded by default
		// ""  : Static block (or expanded details depending on styling preference)
		isCollapsible := foldParams == "-" || foldParams == "+"
		isOpen := foldParams != "-" // Open unless explicitly set to collapsed ("-")

		var sb strings.Builder

		// Construct the container.
		// We use HTML5 <details> for native collapse behavior if modifiers are present.
		if isCollapsible {
			sb.WriteString(`<details class="callout" data-callout="` + cType + `"`)
			if isOpen {
				sb.WriteString(` open`)
			}
			sb.WriteString(`>`)
			sb.WriteString(`<summary class="callout-title">`)
		} else {
			sb.WriteString(`<div class="callout" data-callout="` + cType + `">`)
			sb.WriteString(`<div class="callout-title">`)
		}

		// Render Title with Icon
		sb.WriteString(`<div class="callout-icon"><i data-lucide="` + iconName + `"></i></div>`)
		sb.WriteString(`<div class="callout-title-inner">` + cTitle + `</div>`)

		// Add a chevron icon if the element is collapsible
		if isCollapsible {
			sb.WriteString(
				`<div class="callout-fold-icon"><i data-lucide="chevron-down"></i></div>`,
			)
			sb.WriteString(`</summary>`)
		} else {
			sb.WriteString(`</div>`)
		}

		// Render Body Content
		sb.WriteString(`<div class="callout-content">`)
		sb.WriteString(cBody)
		sb.WriteString(`</div>`)

		// Close tags
		if isCollapsible {
			sb.WriteString(`</details>`)
		} else {
			sb.WriteString(`</div>`)
		}

		return sb.String()
	})
}

// getCalloutIcon maps a callout type string (e.g., "info", "bug") to a corresponding Lucide icon name.
// It supports various aliases common in markdown ecosystems.
func getCalloutIcon(cType string) string {
	switch cType {
	case "abstract", "summary", "tldr":
		return "clipboard-list"
	case "info":
		return "info"
	case "todo":
		return "check-circle-2"
	case "tip", "hint", "important":
		return "lightbulb"
	case "success", "check", "done":
		return "check"
	case "question", "help", "faq":
		return "help-circle"
	case "warning", "caution", "attention":
		return "alert-triangle"
	case "failure", "fail", "missing":
		return "x"
	case "danger", "error":
		return "zap"
	case "bug":
		return "bug"
	case "example":
		return "list"
	case "quote", "cite":
		return "quote"
	default:
		return "pencil"
	}
}
