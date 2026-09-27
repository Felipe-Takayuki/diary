package usecase

import (
	"context"

	"diary/internal/domain"
)

// DailyGoalUseCase defines application use cases for daily goals management.
type DailyGoalUseCase interface {
	GetDailyGoals(ctx context.Context, rawDate string) (*domain.DailyGoal, error)
	SaveDailyGoals(ctx context.Context, rawDate string, items []domain.Item) (*domain.DailyGoal, error)
}

type dailyGoalUseCase struct {
	repo domain.GoalRepository
}

// NewDailyGoalUseCase creates a new DailyGoalUseCase instance.
func NewDailyGoalUseCase(repo domain.GoalRepository) DailyGoalUseCase {
	return &dailyGoalUseCase{
		repo: repo,
	}
}

// GetDailyGoals parses the raw date and retrieves the daily goals for that date.
func (uc *dailyGoalUseCase) GetDailyGoals(ctx context.Context, rawDate string) (*domain.DailyGoal, error) {
	date, err := domain.ParseDate(rawDate)
	if err != nil {
		return nil, err
	}

	return uc.repo.LoadByDate(ctx, date)
}

// SaveDailyGoals parses the raw date, sanitizes items, builds a DailyGoal aggregate, and persists it.
func (uc *dailyGoalUseCase) SaveDailyGoals(ctx context.Context, rawDate string, items []domain.Item) (*domain.DailyGoal, error) {
	date, err := domain.ParseDate(rawDate)
	if err != nil {
		return nil, err
	}

	sanitizedItems := make([]domain.Item, 0, len(items))
	for _, it := range items {
		sanitized, err := domain.NewItem(it.Text(), it.Done())
		if err != nil {
			// Skip empty items, preserving original app behavior
			continue
		}
		sanitizedItems = append(sanitizedItems, sanitized)
	}

	goal := domain.NewDailyGoal(date, sanitizedItems)
	if err := uc.repo.Save(ctx, goal); err != nil {
		return nil, err
	}

	return goal, nil
}
