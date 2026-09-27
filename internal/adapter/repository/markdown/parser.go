package markdown

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"diary/internal/domain"
)

// parseMarkdown extracts checklist items from a Markdown file stream.
func parseMarkdown(r io.Reader) ([]domain.Item, error) {
	var items []domain.Item
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "- [ ] ") {
			text := strings.TrimSpace(line[6:])
			if item, err := domain.NewItem(text, false); err == nil {
				items = append(items, item)
			}
		} else if strings.HasPrefix(line, "- [x] ") || strings.HasPrefix(line, "- [X] ") {
			text := strings.TrimSpace(line[6:])
			if item, err := domain.NewItem(text, true); err == nil {
				items = append(items, item)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading markdown file: %w", err)
	}

	if items == nil {
		items = []domain.Item{}
	}

	return items, nil
}

// formatMarkdown serializes daily goals into Markdown checklist format.
func formatMarkdown(date domain.Date, items []domain.Item) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# Goals - %s\n\n", date.String()))

	for _, item := range items {
		status := " "
		if item.Done() {
			status = "x"
		}
		sb.WriteString(fmt.Sprintf("- [%s] %s\n", status, item.Text()))
	}

	return sb.String()
}
