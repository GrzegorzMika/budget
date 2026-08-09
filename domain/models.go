package domain

import (
	"slices"
	"time"
)

type ExpenseCategory string

type Expense struct {
	ID          int64
	Timestamp   time.Time
	Amount      float64
	Category    ExpenseCategory
	Description string
}

type Category struct {
	Name  string
	Color string
}

// CategoryTotal is an aggregate of expenses for one category over a period.
type CategoryTotal struct {
	Category string
	Total    float64
	Count    int
}

// CategoryColorIDs is the fixed palette from the design system; categories
// store one of these ids and the UI maps them to swatch colors via CSS.
var CategoryColorIDs = []string{
	"c1", "c2", "c3", "c4", "c5", "c6", "c7", "c8",
	"c9", "c10", "c11", "c12", "c13", "c14", "c15", "c16",
}

// NeutralCategoryColor is the graphite swatch, used as the fallback for
// expenses whose category no longer exists.
const NeutralCategoryColor = "c16"

func IsValidCategoryColor(id string) bool {
	return slices.Contains(CategoryColorIDs, id)
}
