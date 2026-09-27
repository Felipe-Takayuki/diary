package domain

import "context"

// GoalRepository defines the output port for loading and persisting daily goals.
type GoalRepository interface {
	// LoadByDate loads the daily goals for the given date.
	// If no goals exist yet for the date, an empty DailyGoal aggregate is returned without error.
	LoadByDate(ctx context.Context, date Date) (*DailyGoal, error)

	// Save persists the daily goals for the given aggregate.
	Save(ctx context.Context, goal *DailyGoal) error
}
