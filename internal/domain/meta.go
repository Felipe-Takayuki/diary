package domain

// DailyGoal represents the aggregate of goals for a specific date.
type DailyGoal struct {
	date  Date
	items []Item
}

// NewDailyGoal creates a DailyGoal aggregate.
func NewDailyGoal(date Date, items []Item) *DailyGoal {
	if items == nil {
		items = []Item{}
	}
	// Copy items slice to avoid external mutation
	copied := make([]Item, len(items))
	copy(copied, items)

	return &DailyGoal{
		date:  date,
		items: copied,
	}
}

// Date returns the date of this daily goal set.
func (g *DailyGoal) Date() Date {
	return g.date
}

// Items returns a copy of the items list.
func (g *DailyGoal) Items() []Item {
	copied := make([]Item, len(g.items))
	copy(copied, g.items)
	return copied
}

// AddItem appends a valid goal item.
func (g *DailyGoal) AddItem(item Item) {
	g.items = append(g.items, item)
}

// TotalCount returns the total number of goals for the day.
func (g *DailyGoal) TotalCount() int {
	return len(g.items)
}

// DoneCount returns the number of completed goals.
func (g *DailyGoal) DoneCount() int {
	count := 0
	for _, item := range g.items {
		if item.Done() {
			count++
		}
	}
	return count
}

// ProgressPercentage calculates the completion percentage from 0 to 100.
func (g *DailyGoal) ProgressPercentage() int {
	total := g.TotalCount()
	if total == 0 {
		return 0
	}
	done := g.DoneCount()
	return (done * 100) / total
}
