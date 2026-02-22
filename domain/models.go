package domain

import "time"

type ExpenseCategory string

const (
	// Categories are now loaded from the database
)

type Expense struct {
	Timestamp time.Time
	Amount    float64
	Category  ExpenseCategory
}
