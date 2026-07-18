package storage

import (
	"context"
	"fmt"
	"strings"
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

var saveExpenseQuery = `INSERT INTO expenses (timestamp, amount, category, description) VALUES ($1, $2, $3, NULLIF($4, ''))`

func (r *Repository) SaveExpense(ctx context.Context, expense *domain.Expense) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()
	_, err := r.db.Exec(newCtx, saveExpenseQuery, expense.Timestamp.Format(time.DateOnly), expense.Amount, string(expense.Category), expense.Description)
	if err != nil {
		return fmt.Errorf("failed to save expense: %w", err)
	}
	return nil
}

func (r *Repository) DeleteExpense(ctx context.Context, id int64) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()
	_, err := r.db.Exec(newCtx, "DELETE FROM expenses WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to delete expense: %w", err)
	}
	return nil
}

func (r *Repository) Ping(ctx context.Context) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()
	return r.db.Ping(newCtx)
}

func (r *Repository) GetCategories(ctx context.Context) ([]domain.Category, error) {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()

	rows, err := r.db.Query(newCtx, "SELECT name, color FROM categories ORDER BY name")
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.Name, &c.Color); err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}
	return categories, nil
}

func (r *Repository) AddCategory(ctx context.Context, name string, color string) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()
	_, err := r.db.Exec(newCtx, "INSERT INTO categories (name, color) VALUES ($1, $2) ON CONFLICT DO NOTHING", name, color)
	if err != nil {
		return fmt.Errorf("failed to add category: %w", err)
	}
	return nil
}

// UpdateCategory renames and/or recolors a category. Expenses reference
// categories by name, so a rename must cascade to them atomically.
func (r *Repository) UpdateCategory(ctx context.Context, oldName string, newName string, color string) error {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()

	tx, err := r.db.Begin(newCtx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(newCtx)

	_, err = tx.Exec(newCtx, "UPDATE categories SET name = $2, color = $3 WHERE name = $1", oldName, newName, color)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	if oldName != newName {
		_, err = tx.Exec(newCtx, "UPDATE expenses SET category = $2 WHERE category = $1", oldName, newName)
		if err != nil {
			return fmt.Errorf("failed to rename category on expenses: %w", err)
		}
	}
	if err := tx.Commit(newCtx); err != nil {
		return fmt.Errorf("failed to commit category update: %w", err)
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

var getExpensesQuery = `
SELECT id, timestamp, amount, category, COALESCE(description, '')
FROM expenses
WHERE timestamp >= $1 AND timestamp < $2
  AND ($3 = '' OR category ILIKE '%' || $3 || '%' OR COALESCE(description, '') ILIKE '%' || $3 || '%')
ORDER BY timestamp DESC, id DESC`

// GetExpenses returns expenses in [start, end), optionally filtered by a
// case-insensitive substring match on category or description.
func (r *Repository) GetExpenses(ctx context.Context, start time.Time, end time.Time, query string) ([]*domain.Expense, error) {
	newCtx, cancel := context.WithTimeout(ctx, DB_TIMEOUT*time.Second)
	defer cancel()

	rows, err := r.db.Query(newCtx, getExpensesQuery, start, end, escapeLike(query))
	if err != nil {
		return nil, fmt.Errorf("failed to query expenses: %w", err)
	}
	defer rows.Close()

	var expenses []*domain.Expense
	for rows.Next() {
		exp := &domain.Expense{}
		var catStr string
		if err := rows.Scan(&exp.ID, &exp.Timestamp, &exp.Amount, &catStr, &exp.Description); err != nil {
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

// escapeLike neutralizes LIKE wildcards in user-provided search text.
func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}
