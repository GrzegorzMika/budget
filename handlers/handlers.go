package handlers

import (
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/GrzegorzMika/budget/controllers"
	"github.com/GrzegorzMika/budget/domain"
	"github.com/GrzegorzMika/budget/handlers/templates"
	"github.com/a-h/templ"
)

func ExpensesHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			timestamp, err := time.Parse(time.DateOnly, r.FormValue("date"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			amount, err := strconv.ParseFloat(strings.TrimSpace(r.FormValue("amount")), 64)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			category := r.FormValue("category")
			expense := &domain.Expense{
				Timestamp: timestamp,
				Amount:    amount,
				Category:  domain.ExpenseCategory(category),
			}
			err = app.SaveExpense(r.Context(), expense)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func LandingPageHandlerBuilder(_ *controllers.AppController) *templ.ComponentHandler {
	categories := make([]string, len(domain.ExpenseCategories))
	for i, category := range domain.ExpenseCategories {
		categories[i] = string(category)
	}
	slices.Sort(categories)
	component := templates.Index(categories)
	return templ.Handler(component)
}

func StaticFileHandlerBuilder(_ *controllers.AppController) http.Handler {
	return http.FileServer(http.FS(templates.Static))
}
