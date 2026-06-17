package format

import (
	"fmt"
	"strings"
)

func formatIDR(value int64) string {
	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}
	raw := fmt.Sprintf("%d", value)
	if len(raw) <= 3 {
		return sign + "Rp " + raw
	}
	parts := make([]string, 0, (len(raw)+2)/3)
	for len(raw) > 3 {
		parts = append([]string{raw[len(raw)-3:]}, parts...)
		raw = raw[:len(raw)-3]
	}
	parts = append([]string{raw}, parts...)
	return sign + "Rp " + strings.Join(parts, ".")
}
