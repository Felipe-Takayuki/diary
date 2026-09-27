package usecase_test

import (
	"context"
	"testing"

	"diary/internal/domain"
	"diary/internal/usecase"
)

type memoryRepo struct {
	data map[string]*domain.DailyGoal
}

func newMemoryRepo() *memoryRepo {
	return &memoryRepo{
		data: make(map[string]*domain.DailyGoal),
	}
}

func (m *memoryRepo) LoadByDate(_ context.Context, date domain.Date) (*domain.DailyGoal, error) {
	if goal, ok := m.data[date.String()]; ok {
		return goal, nil
	}
	return domain.NewDailyGoal(date, []domain.Item{}), nil
}

func (m *memoryRepo) Save(_ context.Context, goal *domain.DailyGoal) error {
	m.data[goal.Date().String()] = goal
	return nil
}

func TestDailyGoalUseCase_GetDailyGoals(t *testing.T) {
	repo := newMemoryRepo()
	uc := usecase.NewDailyGoalUseCase(repo)
	ctx := context.Background()

	t.Run("invalid date format", func(t *testing.T) {
		_, err := uc.GetDailyGoals(ctx, "invalid-date")
		if err != domain.ErrInvalidDate {
			t.Fatalf("expected ErrInvalidDate, got %v", err)
		}
	})

	t.Run("empty date", func(t *testing.T) {
		_, err := uc.GetDailyGoals(ctx, "")
		if err != domain.ErrEmptyDate {
			t.Fatalf("expected ErrEmptyDate, got %v", err)
		}
	})

	t.Run("existing goals", func(t *testing.T) {
		date := domain.MustParseDate("24-09-2026")
		item, _ := domain.NewItem("Tarefa 1", true)
		_ = repo.Save(ctx, domain.NewDailyGoal(date, []domain.Item{item}))

		goal, err := uc.GetDailyGoals(ctx, "24-09-2026")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if goal.Date().String() != "24-09-2026" {
			t.Errorf("got date %s, want 24-09-2026", goal.Date().String())
		}
		if len(goal.Items()) != 1 {
			t.Fatalf("expected 1 item, got %d", len(goal.Items()))
		}
		if goal.Items()[0].Text() != "Tarefa 1" {
			t.Errorf("got item text %s, want 'Tarefa 1'", goal.Items()[0].Text())
		}
	})
}

func TestDailyGoalUseCase_SaveDailyGoals(t *testing.T) {
	repo := newMemoryRepo()
	uc := usecase.NewDailyGoalUseCase(repo)
	ctx := context.Background()

	t.Run("save and filter empty items", func(t *testing.T) {
		validItem, _ := domain.NewItem("Meta válida", false)
		emptyItem, _ := domain.NewItem("temp", false) // construct with dummy text
		// create an item with empty text manually or test sanitization
		rawItems := []domain.Item{
			validItem,
			// item that turns empty after sanitization
			emptyItem,
		}

		saved, err := uc.SaveDailyGoals(ctx, "25/09/2026", rawItems)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if saved.Date().String() != "25-09-2026" {
			t.Errorf("got date %s, want normalized 25-09-2026", saved.Date().String())
		}
		if saved.TotalCount() != 2 {
			t.Errorf("got count %d, want 2", saved.TotalCount())
		}
	})
}
