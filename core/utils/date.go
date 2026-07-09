package utils

import "time"

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
