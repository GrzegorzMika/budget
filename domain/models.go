package domain

import "time"

type ExpenseCategory string

type Expense struct {
	Timestamp time.Time
	Amount    float64
	Category  ExpenseCategory
}
