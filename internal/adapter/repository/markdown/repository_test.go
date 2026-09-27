package markdown_test

import (
	"context"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"diary/internal/adapter/repository/markdown"
	"diary/internal/domain"
)

func TestRepository_SaveAndLoad(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "metas_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo, err := markdown.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	ctx := context.Background()
	date := domain.MustParseDate("24-09-2026")

	// 1. Loading non-existent date returns empty goals without error
	initialGoal, err := repo.LoadByDate(ctx, date)
	if err != nil {
		t.Fatalf("unexpected error loading non-existent date: %v", err)
	}
	if initialGoal.TotalCount() != 0 {
		t.Fatalf("expected 0 items, got %d", initialGoal.TotalCount())
	}

	// 2. Save goals
	item1, _ := domain.NewItem("Estudar Go", true)
	item2, _ := domain.NewItem("Praticar Clean Architecture", false)
	goalToSave := domain.NewDailyGoal(date, []domain.Item{item1, item2})

	if err := repo.Save(ctx, goalToSave); err != nil {
		t.Fatalf("unexpected error saving goal: %v", err)
	}

	// 3. Load saved goals
	loadedGoal, err := repo.LoadByDate(ctx, date)
	if err != nil {
		t.Fatalf("unexpected error loading saved goal: %v", err)
	}
	if loadedGoal.TotalCount() != 2 {
		t.Fatalf("expected 2 items, got %d", loadedGoal.TotalCount())
	}

	items := loadedGoal.Items()
	if items[0].Text() != "Estudar Go" || !items[0].Done() {
		t.Errorf("item 0 mismatch: text=%s, done=%v", items[0].Text(), items[0].Done())
	}
	if items[1].Text() != "Praticar Clean Architecture" || items[1].Done() {
		t.Errorf("item 1 mismatch: text=%s, done=%v", items[1].Text(), items[1].Done())
	}
}

func TestRepository_LegacyFileFallback(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "metas_legacy_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create legacy file format: 2026-09-24.md
	legacyContent := "# Metas - 2026-09-24\n\n- [x] Item Legado\n- [ ] Outro Item\n"
	legacyPath := filepath.Join(tempDir, "2026-09-24.md")
	if err := os.WriteFile(legacyPath, []byte(legacyContent), 0644); err != nil {
		t.Fatalf("failed to write legacy file: %v", err)
	}

	repo, err := markdown.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	ctx := context.Background()
	date := domain.MustParseDate("24-09-2026")

	loadedGoal, err := repo.LoadByDate(ctx, date)
	if err != nil {
		t.Fatalf("unexpected error loading legacy goal: %v", err)
	}

	if loadedGoal.TotalCount() != 2 {
		t.Fatalf("expected 2 items from legacy file, got %d", loadedGoal.TotalCount())
	}
	if loadedGoal.Items()[0].Text() != "Item Legado" || !loadedGoal.Items()[0].Done() {
		t.Errorf("unexpected legacy item 0: %v", loadedGoal.Items()[0])
	}
}

func TestRepository_Concurrency(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "metas_concurrent_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	repo, err := markdown.NewRepository(tempDir)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	ctx := context.Background()
	date := domain.MustParseDate("24-09-2026")

	var wg sync.WaitGroup
	workers := 10

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			item, _ := domain.NewItem("Item", idx%2 == 0)
			goal := domain.NewDailyGoal(date, []domain.Item{item})
			_ = repo.Save(ctx, goal)
			_, _ = repo.LoadByDate(ctx, date)
		}(i)
	}

	wg.Wait()
}
