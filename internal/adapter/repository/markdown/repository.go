package markdown

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"diary/internal/domain"
)

// Repository manages persistence of daily goals in Markdown files.
type Repository struct {
	dir string
	mu  sync.RWMutex
}

// NewRepository creates a new Markdown repository and ensures the target directory exists.
func NewRepository(dir string) (*Repository, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	return &Repository{dir: dir}, nil
}

func (r *Repository) filePath(date domain.Date) string {
	return filepath.Join(r.dir, fmt.Sprintf("%s.md", date.String()))
}

func (r *Repository) legacyFilePath(date domain.Date) string {
	legacy := date.LegacyISO()
	if legacy == "" {
		return ""
	}
	return filepath.Join(r.dir, fmt.Sprintf("%s.md", legacy))
}

// LoadByDate reads the Markdown file for the specified date.
// If the canonical file does not exist, it falls back to the legacy YYYY-MM-DD file name.
func (r *Repository) LoadByDate(_ context.Context, date domain.Date) (*domain.DailyGoal, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	path := r.filePath(date)
	file, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// Legacy fallback: YYYY-MM-DD.md
			legacyPath := r.legacyFilePath(date)
			if legacyPath != "" {
				if oldFile, oldErr := os.Open(legacyPath); oldErr == nil {
					defer oldFile.Close()
					items, parseErr := parseMarkdown(oldFile)
					if parseErr != nil {
						return nil, parseErr
					}
					return domain.NewDailyGoal(date, items), nil
				}
			}
			return domain.NewDailyGoal(date, []domain.Item{}), nil
		}
		return nil, fmt.Errorf("error opening file %s: %w", path, err)
	}
	defer file.Close()

	items, err := parseMarkdown(file)
	if err != nil {
		return nil, err
	}

	return domain.NewDailyGoal(date, items), nil
}

// Save writes the daily goals to the canonical Markdown file (DD-MM-YYYY.md).
func (r *Repository) Save(_ context.Context, goal *domain.DailyGoal) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	content := formatMarkdown(goal.Date(), goal.Items())
	path := r.filePath(goal.Date())

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("error saving file %s: %w", path, err)
	}

	return nil
}
