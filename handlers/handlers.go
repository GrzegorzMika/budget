package handlers

import (
	"net/http"
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

func LandingPageHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tab := r.URL.Query().Get("tab")
		if tab == "" {
			tab = "expense"
		}

		categories := app.GetCategories()
		categoryStrings := make([]string, len(categories))
		for i, category := range categories {
			categoryStrings[i] = string(category)
		}

		var tabComponent templ.Component
		switch tab {
		case "categories":
			tabComponent = templates.CategoriesTab(categoryStrings)
		case "summary":
			total, expenses, err := app.GetMonthlySummary(r.Context(), time.Now())
			if err != nil {
				http.Error(w, "Failed to load summary: " + err.Error(), http.StatusInternalServerError)
				return
			}
			tabComponent = templates.SummaryTab(total, expenses)
		default:
			// Default to 'expense'
			tab = "expense"
			tabComponent = templates.ExpenseTab(categoryStrings)
		}

		templates.Layout(tab).Render(templ.WithChildren(r.Context(), tabComponent), w)
	}
}

func AddCategoryHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			newCategory := strings.TrimSpace(r.FormValue("new_category"))
			if newCategory != "" {
				err := app.AddCategory(r.Context(), newCategory)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			http.Redirect(w, r, "/?tab=categories", http.StatusSeeOther)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func DeleteCategoryHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			category := strings.TrimSpace(r.URL.Query().Get("category"))
			if category == "" {
				http.Error(w, "category is required", http.StatusBadRequest)
				return
			}
			err := app.DeleteCategory(r.Context(), category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func StaticFileHandlerBuilder(_ *controllers.AppController) http.Handler {
	return http.FileServer(http.FS(templates.Static))
}

func HealthcheckHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			err := app.Ping(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
