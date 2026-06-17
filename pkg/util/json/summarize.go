package json

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

func summarizeThreshold(raw string) string {
	payload := map[string]any{}
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return raw
	}
	parts := make([]string, 0, len(payload))
	keys := make([]string, 0, len(payload))
	for key := range payload {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf("%s=%v", key, payload[key]))
	}
	return strings.Join(parts, " · ")
}
