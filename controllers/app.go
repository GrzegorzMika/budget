package controllers

import (
	"context"
	"sync"
	"time"

	"github.com/GrzegorzMika/budget/domain"
)

type Repository interface {
	SaveExpense(ctx context.Context, expense *domain.Expense) error
	DeleteExpense(ctx context.Context, id int64) error
	GetCategories(ctx context.Context) ([]domain.Category, error)
	AddCategory(ctx context.Context, name string, color string) error
	UpdateCategory(ctx context.Context, oldName string, newName string, color string) error
	DeleteCategory(ctx context.Context, name string) error
	Ping(ctx context.Context) error
	GetExpenses(ctx context.Context, start time.Time, end time.Time, query string) ([]*domain.Expense, error)
}

type AppController struct {
	repository Repository
	mu         sync.RWMutex
	categories []domain.Category
}

func NewAppController(repository Repository, categories []domain.Category) *AppController {
	return &AppController{
		repository: repository,
		categories: categories,
	}
}

func (c *AppController) GetCategories() []domain.Category {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make([]domain.Category, len(c.categories))
	copy(result, c.categories)
	return result
}

func (c *AppController) SaveExpense(ctx context.Context, expense *domain.Expense) error {
	return c.repository.SaveExpense(ctx, expense)
}

func (c *AppController) DeleteExpense(ctx context.Context, id int64) error {
	return c.repository.DeleteExpense(ctx, id)
}

func (c *AppController) AddCategory(ctx context.Context, name string, color string) error {
	err := c.repository.AddCategory(ctx, name, color)
	if err != nil {
		return err
	}
	return c.reloadCategories(ctx)
}

func (c *AppController) UpdateCategory(ctx context.Context, oldName string, newName string, color string) error {
	err := c.repository.UpdateCategory(ctx, oldName, newName, color)
	if err != nil {
		return err
	}
	return c.reloadCategories(ctx)
}

func (c *AppController) DeleteCategory(ctx context.Context, name string) error {
	err := c.repository.DeleteCategory(ctx, name)
	if err != nil {
		return err
	}
	return c.reloadCategories(ctx)
}

func (c *AppController) reloadCategories(ctx context.Context) error {
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

// GetExpenses returns expenses in [start, end) matching the optional search
// query, newest first.
func (c *AppController) GetExpenses(ctx context.Context, start time.Time, end time.Time, query string) ([]*domain.Expense, error) {
	return c.repository.GetExpenses(ctx, start, end, query)
}
