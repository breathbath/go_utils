//Package url gives HTTP clients helper functions
package url

import (
	"strings"
)

//JoinURL joins url part into a single url
func JoinURL(parts ...string) string {
	preparedParts := make([]string, 0, len(parts))
	for k, part := range parts {
		if part == "" {
			continue
		}
		if k == 0 {
			preparedParts = append(preparedParts, strings.TrimRight(part, "/"))
		} else {
			preparedParts = append(preparedParts, strings.Trim(part, "/"))
		}
	}
	return strings.Join(preparedParts, "/")
}
