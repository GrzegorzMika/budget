package handlers

import (
	"log/slog"
	"slices"

	"github.com/GrzegorzMika/budget/controllers"
	"github.com/GrzegorzMika/budget/domain"
	"github.com/gofiber/fiber/v3"
)

func ExpensesHandlerBuilder(app *controllers.AppController) fiber.Handler {
	return func(c fiber.Ctx) error {
		expense := new(domain.Expense)
		if err := c.Bind().Body(expense); err != nil {
			slog.With("error", err).Error("failed to bind request body")
			return c.SendStatus(fiber.StatusBadRequest)
		}

		err := app.SaveExpense(c.Context(), expense)
		if err != nil {
			slog.With("error", err).Error("failed to save expense")
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusCreated)
	}
}

func ExpenseCategoryHandlerBuilder(app *controllers.AppController) fiber.Handler {
	categories := make([]string, len(domain.ExpenseCategories))
	for i, category := range domain.ExpenseCategories {
		categories[i] = string(category)
	}
	slices.Sort(categories)
	return func(c fiber.Ctx) error {
		return c.JSON(categories)
	}
}

func HealthcheckHandlerBuilder(app *controllers.AppController) fiber.Handler {
	return func(c fiber.Ctx) error {
		if err := app.Ping(c.Context()); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	}
}
