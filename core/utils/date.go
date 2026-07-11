package utils

import (
	"database/sql"
	"fmt"
	"time"
)

func FormatIndonesianDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	months := [...]string{
		"",
		"Januari",
		"Februari",
		"Maret",
		"April",
		"Mei",
		"Juni",
		"Juli",
		"Agustus",
		"September",
		"Oktober",
		"November",
		"Desember",
	}

	return t.Format("02") + " " + months[int(t.Month())] + " " + t.Format("2006")
}

func ParseISODate(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}

	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse date %q: %w", value, err)
	}

	return parsed, nil
}

func FormatRiwayatTime(values ...sql.NullTime) string {
	for _, v := range values {
		if t := NullTimeValue(v); !t.IsZero() {
			return t.Format("15:04")
		}
	}
	return ""
}
