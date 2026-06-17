package shift

import (
	"fmt"
	"modulegue/internal/domain/web/location"
	"strings"
)

func shiftCodeFromTemplate(shift location.ParkingShiftTemplate) string {
	label := strings.ToUpper(strings.TrimSpace(shift.Label))
	label = strings.ReplaceAll(label, " ", "_")
	label = strings.ReplaceAll(label, "-", "_")
	if label == "" {
		label = "SHIFT"
	}
	return fmt.Sprintf("%s_%s_%s", label, strings.ReplaceAll(shift.Start, ":", ""), strings.ReplaceAll(shift.End, ":", ""))
}
