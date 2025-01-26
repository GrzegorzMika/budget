package controllers

import (
	"context"

	"github.com/GrzegorzMika/budget/domain"
)

type Repository interface {
	SaveExpense(ctx context.Context, expense *domain.Expense) error
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
