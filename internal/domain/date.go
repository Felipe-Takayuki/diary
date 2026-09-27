package domain

import (
	"strings"
	"time"
)

const (
	CanonicalDateFormat = "02-01-2006" // DD-MM-YYYY
	slashDateFormat     = "02/01/2006" // DD/MM/YYYY
	isoDateFormat       = "2006-01-02" // YYYY-MM-DD
)

// Date represents a canonical date in DD-MM-YYYY format as a value object.
type Date struct {
	value string
	t     time.Time
}

// ParseDate parses, validates, and normalizes a date string into a canonical Date.
func ParseDate(dateStr string) (Date, error) {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return Date{}, ErrEmptyDate
	}

	// 1. Primary format: DD-MM-YYYY
	if t, err := time.Parse(CanonicalDateFormat, dateStr); err == nil {
		return Date{value: t.Format(CanonicalDateFormat), t: t}, nil
	}

	// 2. Slash format: DD/MM/YYYY
	if t, err := time.Parse(slashDateFormat, dateStr); err == nil {
		return Date{value: t.Format(CanonicalDateFormat), t: t}, nil
	}

	// 3. ISO format: YYYY-MM-DD
	if t, err := time.Parse(isoDateFormat, dateStr); err == nil {
		return Date{value: t.Format(CanonicalDateFormat), t: t}, nil
	}

	return Date{}, ErrInvalidDate
}

// MustParseDate parses a date string and panics if invalid (useful for constants or tests).
func MustParseDate(dateStr string) Date {
	d, err := ParseDate(dateStr)
	if err != nil {
		panic(err)
	}
	return d
}

// String returns the canonical DD-MM-YYYY string representation.
func (d Date) String() string {
	return d.value
}

// Time returns the underlying time.Time.
func (d Date) Time() time.Time {
	return d.t
}

// LegacyISO returns the date in YYYY-MM-DD format (used for backward compatibility with legacy filenames).
func (d Date) LegacyISO() string {
	if d.t.IsZero() {
		return ""
	}
	return d.t.Format(isoDateFormat)
}

// IsZero reports whether the date is unset.
func (d Date) IsZero() bool {
	return d.value == ""
}
