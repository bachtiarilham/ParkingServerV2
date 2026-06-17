package priority

import "strings"

func normalizePriority(severity string) string {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "high":
		return "high"
	case "medium":
		return "medium"
	default:
		return "low"
	}
}
