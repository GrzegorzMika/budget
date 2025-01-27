package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/GrzegorzMika/budget/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

const DB_TIMEOUT = 5

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

var saveExpenseQuery = `INSERT INTO expenses (timestamp, amount, category) VALUES ($1, $2, $3)`

func (r *Repository) SaveExpense(ctx context.Context, expense *domain.Expense) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()
	_, err := r.db.Exec(newCtx, saveExpenseQuery, expense.Timestamp.Format(time.DateOnly), expense.Amount, string(expense.Category))
	if err != nil {
		return fmt.Errorf("failed to save expense: %w", err)
	}
	return nil
}

func (r *Repository) Ping(ctx context.Context) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()
	return r.db.Ping(newCtx)
}
