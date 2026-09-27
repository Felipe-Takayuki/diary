package domain

import (
	"strings"
)

// Item represents a single daily goal item.
type Item struct {
	text string
	done bool
}

// NewItem creates a new sanitized goal item.
func NewItem(text string, done bool) (Item, error) {
	cleanText := SanitizeItemText(text)
	if cleanText == "" {
		return Item{}, ErrEmptyItemText
	}
	return Item{
		text: cleanText,
		done: done,
	}, nil
}

// SanitizeItemText removes line breaks and trims spaces from item text.
func SanitizeItemText(text string) string {
	clean := strings.ReplaceAll(text, "\r", " ")
	clean = strings.ReplaceAll(clean, "\n", " ")
	return strings.TrimSpace(clean)
}

// Text returns the goal text.
func (i Item) Text() string {
	return i.text
}

// Done returns whether the goal is completed.
func (i Item) Done() bool {
	return i.done
}

// Toggle returns a copy of the item with the done status inverted.
func (i Item) Toggle() Item {
	return Item{
		text: i.text,
		done: !i.done,
	}
}
