package domain_test

import (
	"testing"

	"diary/internal/domain"
)

func TestDailyGoal_Progress(t *testing.T) {
	date := domain.MustParseDate("25-09-2026")

	t.Run("empty goals", func(t *testing.T) {
		dg := domain.NewDailyGoal(date, nil)
		if dg.TotalCount() != 0 {
			t.Errorf("got total %d, want 0", dg.TotalCount())
		}
		if dg.DoneCount() != 0 {
			t.Errorf("got done %d, want 0", dg.DoneCount())
		}
		if dg.ProgressPercentage() != 0 {
			t.Errorf("got progress %d, want 0", dg.ProgressPercentage())
		}
	})

	t.Run("mixed goals", func(t *testing.T) {
		item1, _ := domain.NewItem("Task 1", true)
		item2, _ := domain.NewItem("Task 2", false)
		item3, _ := domain.NewItem("Task 3", true)
		item4, _ := domain.NewItem("Task 4", false)

		dg := domain.NewDailyGoal(date, []domain.Item{item1, item2, item3, item4})

		if dg.TotalCount() != 4 {
			t.Errorf("got total %d, want 4", dg.TotalCount())
		}
		if dg.DoneCount() != 2 {
			t.Errorf("got done %d, want 2", dg.DoneCount())
		}
		if dg.ProgressPercentage() != 50 {
			t.Errorf("got progress %d, want 50", dg.ProgressPercentage())
		}
	})
}
