package domain

import "errors"

var (
	// ErrInvalidDate is returned when a date string does not match any accepted date format.
	ErrInvalidDate = errors.New("invalid date format (use DD-MM-YYYY)")

	// ErrEmptyDate is returned when an empty date is provided.
	ErrEmptyDate = errors.New("date cannot be empty")

	// ErrEmptyItemText is returned when a goal item text is empty.
	ErrEmptyItemText = errors.New("goal item text cannot be empty")
)
