package builder

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
	// "golang.org/x/text/cases"
)

// Contains logic for transforming specific tags

func transformHighlights(htmlStr string) string {
	re := regexp.MustCompile(`==([^=]+)==`)
	return re.ReplaceAllString(htmlStr, `<mark>$1</mark>`)
}

func transformMermaid(htmlStr string) string {
	re := regexp.MustCompile(
		`(?s)<pre[^>]*>\s*<code class="language-mermaid"[^>]*>(.*?)</code>\s*</pre>`,
	)
	return re.ReplaceAllStringFunc(htmlStr, func(match string) string {
		submatches := re.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}
		content := submatches[1]
		// Store original content in data-original so we can re-render on theme change
		encoded := template.HTMLEscapeString(content)
		return fmt.Sprintf(`<div class="mermaid" data-original="%s">%s</div>`, encoded, content)
	})
}

func transformCallouts(htmlStr string) string {
	// Updated regex to capture:
	// Group 1: Type (e.g., "faq")
	// Group 2: Fold Modifier ("+" or "-" or empty)
	// Group 3: Title (text remaining on the first line)
	// Group 4: Body (everything after the title line until blockquote ends)
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

		// Use title casing for the default title if none provided
		if cTitle == "" {
			// cTitle = cases.Title(cType)
			cTitle = strings.Title(cType)
		}

		// Cleanup body: If the regex split at a <br>, the closing </p> might still be in the body.
		// We ensure the body is reasonably wrapped or clean.
		if strings.HasSuffix(cBody, "</p>") && !strings.HasPrefix(cBody, "<p>") {
			// Prepend a p tag if we cut off the opening one via the regex split on <br>
			cBody = "<p>" + cBody
		}

		iconName := getCalloutIcon(cType)

		// Determine collapsible state
		// Obsidian: "-" means collapsed by default. "+" means open.
		// If neither, usually it's just a block, but we can treat it as open <details> for consistency.
		isCollapsible := foldParams == "-" || foldParams == "+"
		isOpen := foldParams != "-" // Open by default unless explicitly "-"

		var sb strings.Builder

		// We use HTML5 <details> for native collapse behavior
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

		// Title Content
		sb.WriteString(`<div class="callout-icon"><i data-lucide="` + iconName + `"></i></div>`)
		sb.WriteString(`<div class="callout-title-inner">` + cTitle + `</div>`)

		// Add fold icon if collapsible (optional visual indicator)
		if isCollapsible {
			sb.WriteString(
				`<div class="callout-fold-icon"><i data-lucide="chevron-down"></i></div>`,
			)
			sb.WriteString(`</summary>`)
		} else {
			sb.WriteString(`</div>`)
		}

		// Body Content
		sb.WriteString(`<div class="callout-content">`)
		sb.WriteString(cBody)
		sb.WriteString(`</div>`)

		if isCollapsible {
			sb.WriteString(`</details>`)
		} else {
			sb.WriteString(`</div>`)
		}

		return sb.String()
	})
}

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
