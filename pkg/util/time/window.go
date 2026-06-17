package time

import "strings"

func splitHoursRange(raw string) (string, string) {
	parts := strings.Split(raw, "-")
	if len(parts) != 2 {
		return "", ""
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
}
