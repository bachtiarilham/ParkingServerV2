package utils

import (
	"database/sql"
	"strconv"
	"time"
)

func NullStringValue(v sql.NullString) string {
	if !v.Valid {
		return ""
	}
	return v.String
}

func NullInt64Value(v sql.NullInt64) int64 {
	if !v.Valid {
		return 0
	}
	return v.Int64
}

func NullInt64StringValue(v sql.NullInt64) string {
	if !v.Valid {
		return ""
	}
	return strconv.FormatInt(v.Int64, 10)
}

func NullInt64Param(v int64) any {
	if v == 0 {
		return nil
	}
	return v
}

func NullTimeValue(v sql.NullTime) time.Time {
	if !v.Valid {
		return time.Time{}
	}
	return v.Time
}

func NullBoolValue(v sql.NullBool) bool {
	if !v.Valid {
		return false
	}
	return v.Bool
}

func NullFloat64Value(v sql.NullFloat64) float64 {
	if !v.Valid {
		return 0
	}
	return v.Float64
}
