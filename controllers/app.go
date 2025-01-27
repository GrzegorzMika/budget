package controllers

import (
	"context"

	"github.com/GrzegorzMika/budget/domain"
)

type Repository interface {
	SaveExpense(ctx context.Context, expense *domain.Expense) error
	Ping(ctx context.Context) error
}

type AppController struct {
	repository Repository
}

func NewAppController(repository Repository) *AppController {
	return &AppController{repository: repository}
}

func (c *AppController) SaveExpense(ctx context.Context, expense *domain.Expense) error {
	return c.repository.SaveExpense(ctx, expense)
}

func (c *AppController) Ping(ctx context.Context) error {
	return c.repository.Ping(ctx)
}
