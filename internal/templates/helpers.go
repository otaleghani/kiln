// @feature:layouts Template helper functions for templ components.
package templates

import "time"

// FormatDate formats a time.Time as "Jan 02, 2006".
func FormatDate(t time.Time) string {
	return t.Format("Jan 02, 2006")
}
