package utils

import (
	"database/sql"
)

func NullString(s string) (ns sql.NullString) {
	if s != "" {
		ns.String = s
		ns.Valid = true
	}

	return ns
}

func NullInt64(f int64) (ns sql.NullInt64) {
	if f != 0 {
		ns.Int64 = f
		ns.Valid = true
	}
	return ns
}

func FormatNullTime(nt sql.NullTime, format string) string {
	if nt.Valid {
		return nt.Time.Format(format)
	}
	return ""
}
