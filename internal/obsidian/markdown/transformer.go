package markdown

import (
	"fmt"
	"html/template"
	"regexp"
	"strings"
)

// Contains logic for transforming specific markdown extensions into rich HTML components.

// applyTransforms applies all the transform
func applyTransforms(htmlStr string) string {
	htmlStr = transformHighlights(htmlStr)
	htmlStr = transformMermaid(htmlStr)
	htmlStr = transformCallouts(htmlStr)
	htmlStr = transformTags(htmlStr)
	return htmlStr
}

// transformTags locates Obsidian-style tags (e.g. #tagname)
// and converts them into a stylized link.
func transformTags(htmlStr string) string {
	// Regex explanation:
	// (?m)       : Multi-line mode (though often redundant in simple string replaces, good for safety).
	// (\s|^)     : Group 1. Match either whitespace OR start of line. We need to capture this to put it back.
	// (#[a-zA-Z0-9_\-]+) : Group 2. The tag itself, starting with #.
	re := regexp.MustCompile(`(?m)(^|\s|>)(#[a-zA-Z0-9_\-]+)`)

	return re.ReplaceAllStringFunc(htmlStr, func(match string) string {
		submatches := re.FindStringSubmatch(match)

		// Safety check: ensure we have both groups (separator + tag)
		if len(submatches) < 3 {
			return match
		}

		separator := submatches[1] // The whitespace or empty string found before the tag
		fullTag := submatches[2]   // The tag, e.g., "#tagname"
		tagName := fullTag[1:]     // The tag without the hash, e.g., "tagname"

		// Construct the new HTML
		// We re-attach the 'separator' at the front to preserve sentence spacing.
		return fmt.Sprintf(
			`%s<span class="inline-tag"><a href="/tags/%s">%s</a></span>`,
			separator,
			tagName,
			fullTag,
		)
	})
}

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

		icon := getCalloutIcon(cType)

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
		sb.WriteString(`<div class="callout-icon">` + icon + `</div>`)
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

func getCalloutIcon(cType string) string {
	switch cType {
	case "abstract", "summary", "tldr":
		// Icon: clipboard-list
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-clipboard-list-icon lucide-clipboard-list h-5 w-5"><rect width="8" height="4" x="8" y="2" rx="1" ry="1"/><path d="M16 4h2a2 2 0 0 1 2 2v14a2 2 0 0 1-2 2H6a2 2 0 0 1-2-2V6a2 2 0 0 1 2-2h2"/><path d="M12 11h4"/><path d="M12 16h4"/><path d="M8 11h.01"/><path d="M8 16h.01"/></svg>`

	case "info", "new":
		// Icon: info
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-info-icon lucide-info h-5 w-5"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4"/><path d="M12 8h.01"/></svg>`

	case "todo":
		// Icon: check-circle-2 (provided as lucide-circle-check)
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-circle-check-icon lucide-circle-check h-5 w-5"><circle cx="12" cy="12" r="10"/><path d="m9 12 2 2 4-4"/></svg>`

	case "tip", "hint", "important":
		// Icon: lightbulb
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-lightbulb-icon lucide-lightbulb h-5 w-5"><path d="M15 14c.2-1 .7-1.7 1.5-2.5 1-.9 1.5-2.2 1.5-3.5A6 6 0 0 0 6 8c0 1 .2 2.2 1.5 3.5.7.7 1.3 1.5 1.5 2.5"/><path d="M9 18h6"/><path d="M10 22h4"/></svg>`

	case "success", "check", "done":
		// Icon: check
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-check-icon lucide-check h-5 w-5"><path d="M20 6 9 17l-5-5"/></svg>`

	case "question", "help", "faq":
		// Icon: help-circle (provided as message-circle-question-mark)
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-message-circle-question-mark-icon lucide-message-circle-question-mark h-5 w-5"><path d="M2.992 16.342a2 2 0 0 1 .094 1.167l-1.065 3.29a1 1 0 0 0 1.236 1.168l3.413-.998a2 2 0 0 1 1.099.092 10 10 0 1 0-4.777-4.719"/><path d="M9.09 9a3 3 0 0 1 5.83 1c0 2-3 3-3 3"/><path d="M12 17h.01"/></svg>`

	case "warning", "caution", "attention":
		// Icon: alert-triangle (provided as triangle-alert)
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-triangle-alert-icon lucide-triangle-alert h-5 w-5"><path d="m21.73 18-8-14a2 2 0 0 0-3.48 0l-8 14A2 2 0 0 0 4 21h16a2 2 0 0 0 1.73-3"/><path d="M12 9v4"/><path d="M12 17h.01"/></svg>`

	case "failure", "fail", "missing":
		// Icon: x
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-x-icon lucide-x h-5 w-5"><path d="M18 6 6 18"/><path d="m6 6 12 12"/></svg>`

	case "danger", "error":
		// Icon: zap
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-zap-icon lucide-zap h-5 w-5"><path d="M4 14a1 1 0 0 1-.78-1.63l9.9-10.2a.5.5 0 0 1 .86.46l-1.92 6.02A1 1 0 0 0 13 10h7a1 1 0 0 1 .78 1.63l-9.9 10.2a.5.5 0 0 1-.86-.46l1.92-6.02A1 1 0 0 0 11 14z"/></svg>`

	case "bug":
		// Icon: bug
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-bug-icon lucide-bug h-5 w-5"><path d="M12 20v-9"/><path d="M14 7a4 4 0 0 1 4 4v3a6 6 0 0 1-12 0v-3a4 4 0 0 1 4-4z"/><path d="M14.12 3.88 16 2"/><path d="M21 21a4 4 0 0 0-3.81-4"/><path d="M21 5a4 4 0 0 1-3.55 3.97"/><path d="M22 13h-4"/><path d="M3 21a4 4 0 0 1 3.81-4"/><path d="M3 5a4 4 0 0 0 3.55 3.97"/><path d="M6 13H2"/><path d="m8 2 1.88 1.88"/><path d="M9 7.13V6a3 3 0 1 1 6 0v1.13"/></svg>`

	case "example":
		// Icon: list
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-list-icon lucide-list h-5 w-5"><path d="M3 5h.01"/><path d="M3 12h.01"/><path d="M3 19h.01"/><path d="M8 5h13"/><path d="M8 12h13"/><path d="M8 19h13"/></svg>`

	case "quote", "cite":
		// Icon: quote
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-quote-icon lucide-quote h-5 w-5"><path d="M16 3a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2 1 1 0 0 1 1 1v1a2 2 0 0 1-2 2 1 1 0 0 0-1 1v2a1 1 0 0 0 1 1 6 6 0 0 0 6-6V5a2 2 0 0 0-2-2z"/><path d="M5 3a2 2 0 0 0-2 2v6a2 2 0 0 0 2 2 1 1 0 0 1 1 1v1a2 2 0 0 1-2 2 1 1 0 0 0-1 1v2a1 1 0 0 0 1 1 6 6 0 0 0 6-6V5a2 2 0 0 0-2-2z"/></svg>`

	default:
		// Icon: pencil
		return `<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="lucide lucide-pencil-icon lucide-pencil h-5 w-5"><path d="M21.174 6.812a1 1 0 0 0-3.986-3.987L3.842 16.174a2 2 0 0 0-.5.83l-1.321 4.352a.5.5 0 0 0 .623.622l4.353-1.32a2 2 0 0 0 .83-.497z"/><path d="m15 5 4 4"/></svg>`
	}
}
