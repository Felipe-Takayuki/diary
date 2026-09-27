package domain

import "errors"

var (
	// ErrInvalidDate is returned when a date string does not match any accepted date format.
	ErrInvalidDate = errors.New("formato de data inválido (utilize dia-mês-ano: DD-MM-YYYY)")

	// ErrEmptyDate is returned when an empty date is provided.
	ErrEmptyDate = errors.New("data não informada")

	// ErrEmptyItemText is returned when a goal item text is empty.
	ErrEmptyItemText = errors.New("texto da meta não pode ser vazio")
)
