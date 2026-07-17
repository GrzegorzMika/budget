package controllers

import (
	"context"
	"sync"
	"time"

	"github.com/GrzegorzMika/budget/domain"
)

type Repository interface {
	SaveExpense(ctx context.Context, expense *domain.Expense) error
	GetCategories(ctx context.Context) ([]domain.ExpenseCategory, error)
	AddCategory(ctx context.Context, name string) error
	DeleteCategory(ctx context.Context, name string) error
	Ping(ctx context.Context) error
	GetTotal(ctx context.Context, start time.Time, end time.Time) (float64, error)
	GetExpenses(ctx context.Context, start time.Time, end time.Time) ([]*domain.Expense, error)
}

type AppController struct {
	repository Repository
	mu         sync.RWMutex
	categories []domain.ExpenseCategory
}

func NewAppController(repository Repository, categories []domain.ExpenseCategory) *AppController {
	return &AppController{
		repository: repository,
		categories: categories,
	}
}

func (c *AppController) GetCategories() []domain.ExpenseCategory {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]domain.ExpenseCategory, len(c.categories))
	copy(result, c.categories)
	return result
}

func (c *AppController) SaveExpense(ctx context.Context, expense *domain.Expense) error {
	return c.repository.SaveExpense(ctx, expense)
}

func (c *AppController) AddCategory(ctx context.Context, name string) error {
	err := c.repository.AddCategory(ctx, name)
	if err != nil {
		return err
	}
	// Reload categories from DB after update
	categories, err := c.repository.GetCategories(ctx)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.categories = categories
	c.mu.Unlock()
	return nil
}

func (c *AppController) DeleteCategory(ctx context.Context, name string) error {
	err := c.repository.DeleteCategory(ctx, name)
	if err != nil {
		return err
	}
	// Reload categories from DB after update
	categories, err := c.repository.GetCategories(ctx)
	if err != nil {
		return err
	}
	c.mu.Lock()
	c.categories = categories
	c.mu.Unlock()
	return nil
}

func (c *AppController) Ping(ctx context.Context) error {
	return c.repository.Ping(ctx)
}

// GetSummary returns the total and expense list for [start, end).
func (c *AppController) GetSummary(ctx context.Context, start time.Time, end time.Time) (float64, []*domain.Expense, error) {
	total, err := c.repository.GetTotal(ctx, start, end)
	if err != nil {
		return 0, nil, err
	}

	expenses, err := c.repository.GetExpenses(ctx, start, end)
	if err != nil {
		return 0, nil, err
	}

	return total, expenses, nil
}
