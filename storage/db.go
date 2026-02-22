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

func (r *Repository) GetCategories(ctx context.Context) ([]domain.ExpenseCategory, error) {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()

	rows, err := r.db.Query(newCtx, "SELECT name FROM categories ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.ExpenseCategory
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, domain.ExpenseCategory(name))
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}
	return categories, nil
}

func (r *Repository) AddCategory(ctx context.Context, name string) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()
	_, err := r.db.Exec(newCtx, "INSERT INTO categories (name) VALUES ($1) ON CONFLICT DO NOTHING", name)
	if err != nil {
		return fmt.Errorf("failed to add category: %w", err)
	}
	return nil
}

func (r *Repository) DeleteCategory(ctx context.Context, name string) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()
	_, err := r.db.Exec(newCtx, "DELETE FROM categories WHERE name = $1", name)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}
	return nil
}

func (r *Repository) GetMonthlyTotal(ctx context.Context, start time.Time, end time.Time) (float64, error) {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()

	var total float64
	err := r.db.QueryRow(newCtx, "SELECT COALESCE(SUM(amount), 0) FROM expenses WHERE timestamp >= $1 AND timestamp < $2", start, end).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to query monthly total: %w", err)
	}
	return total, nil
}

func (r *Repository) GetMonthlyExpenses(ctx context.Context, start time.Time, end time.Time) ([]*domain.Expense, error) {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()

	rows, err := r.db.Query(newCtx, "SELECT timestamp, amount, category FROM expenses WHERE timestamp >= $1 AND timestamp < $2 ORDER BY timestamp DESC", start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query monthly expenses: %w", err)
	}
	defer rows.Close()

	var expenses []*domain.Expense
	for rows.Next() {
		exp := &domain.Expense{}
		var catStr string
		if err := rows.Scan(&exp.Timestamp, &exp.Amount, &catStr); err != nil {
			return nil, fmt.Errorf("failed to scan expense: %w", err)
		}
		exp.Category = domain.ExpenseCategory(catStr)
		expenses = append(expenses, exp)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating expenses: %w", err)
	}
	return expenses, nil
}
