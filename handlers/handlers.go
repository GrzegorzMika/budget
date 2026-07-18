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

func LandingPageHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		tab := r.URL.Query().Get("tab")

		var tabComponent templ.Component
		switch tab {
		case "list", "summary":
			tab = "list"
			component, err := buildListTab(r, app, now)
			if err != nil {
				http.Error(w, "Failed to load expenses: "+err.Error(), http.StatusInternalServerError)
				return
			}
			tabComponent = component
		case "categories":
			editing := strings.TrimSpace(r.URL.Query().Get("edit"))
			tabComponent = templates.CategoriesTab(app.GetCategories(), editing)
		default:
			tab = "add"
			saved := r.URL.Query().Get("saved") == "1"
			tabComponent = templates.AddTab(app.GetCategories(), now.Format(time.DateOnly), saved)
		}

		templates.Layout(tab, templates.MonthTitle(now)).Render(templ.WithChildren(r.Context(), tabComponent), w)
	}
}

func buildListTab(r *http.Request, app *controllers.AppController, now time.Time) (templ.Component, error) {
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)

	from, fromGiven := parseDateParam(r, "from", monthStart)
	to, toGiven := parseDateParam(r, "to", monthEnd)
	if to.Before(from) {
		from, to = to, from
		fromGiven, toGiven = toGiven, fromGiven
	}
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	rangeActive := fromGiven || toGiven

	// "to" is inclusive; the repository takes an exclusive upper bound
	expenses, err := app.GetExpenses(r.Context(), from, to.AddDate(0, 0, 1), query)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, e := range expenses {
		total += e.Amount
	}

	label := templates.MonthTitle(now)
	period := "wydane w tym miesiącu"
	if rangeActive {
		period = "wydane w wybranym okresie"
		switch {
		case fromGiven && toGiven:
			label = templates.DayLabel(from) + " – " + templates.DayLabel(to)
		case fromGiven:
			label = "od " + templates.DayLabel(from)
		default:
			label = "do " + templates.DayLabel(to)
		}
	}
	countLabel := period + " · " + strconv.Itoa(len(expenses)) + " " + templates.PluralEntries(len(expenses))

	colors := make(map[string]string)
	for _, c := range app.GetCategories() {
		colors[c.Name] = c.Color
	}

	var groups []templates.DayGroup
	lastDay := ""
	for _, e := range expenses {
		day := e.Timestamp.Format(time.DateOnly)
		if day != lastDay {
			groups = append(groups, templates.DayGroup{Label: templates.DayLabel(e.Timestamp)})
			lastDay = day
		}
		title, subtitle := e.Description, string(e.Category)
		if title == "" {
			title, subtitle = string(e.Category), ""
		}
		color, ok := colors[string(e.Category)]
		if !ok {
			color = domain.NeutralCategoryColor
		}
		g := &groups[len(groups)-1]
		g.Items = append(g.Items, templates.FeedItem{
			ID:        strconv.FormatInt(e.ID, 10),
			Title:     title,
			Subtitle:  subtitle,
			Color:     color,
			AmountFmt: templates.FormatAmount(e.Amount),
		})
	}

	fromValue, toValue := "", ""
	if fromGiven {
		fromValue = from.Format(time.DateOnly)
	}
	if toGiven {
		toValue = to.Format(time.DateOnly)
	}

	return templates.ListTab(label, templates.FormatAmount(total), countLabel, query, fromValue, toValue, rangeActive || query != "", groups), nil
}

func ExpensesHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			timestamp, err := time.Parse(time.DateOnly, r.FormValue("date"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			amount, err := parseAmount(r.FormValue("amount"))
			if err != nil {
				http.Error(w, "nieprawidłowa kwota", http.StatusBadRequest)
				return
			}
			category := strings.TrimSpace(r.FormValue("category"))
			if category == "" {
				http.Error(w, "kategoria jest wymagana", http.StatusBadRequest)
				return
			}
			description := strings.TrimSpace(r.FormValue("description"))
			if len(description) > 500 {
				http.Error(w, "opis może mieć najwyżej 500 znaków", http.StatusBadRequest)
				return
			}
			expense := &domain.Expense{
				Timestamp:   timestamp,
				Amount:      amount,
				Category:    domain.ExpenseCategory(category),
				Description: description,
			}
			if err := app.SaveExpense(r.Context(), expense); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Stay on the add form (rendered fresh) with a save confirmation.
			http.Redirect(w, r, "/?tab=add&saved=1", http.StatusSeeOther)
		case http.MethodDelete:
			id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
			if err != nil {
				http.Error(w, "nieprawidłowy identyfikator", http.StatusBadRequest)
				return
			}
			if err := app.DeleteExpense(r.Context(), id); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Totals and day groups change with the row, so ask htmx to
			// reload the whole page.
			w.Header().Set("HX-Refresh", "true")
			w.WriteHeader(http.StatusOK)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// parseAmount accepts Polish decimal notation ("12,50", "1 142,86").
func parseAmount(raw string) (float64, error) {
	s := strings.TrimSpace(raw)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, ",", ".")
	amount, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	if amount <= 0 {
		return 0, strconv.ErrRange
	}
	return amount, nil
}

// parseDateParam reads a YYYY-MM-DD query parameter, falling back to def when
// absent or malformed. The second result reports whether a valid value was given.
func parseDateParam(r *http.Request, name string, def time.Time) (time.Time, bool) {
	value := r.URL.Query().Get(name)
	if value == "" {
		return def, false
	}
	t, err := time.Parse(time.DateOnly, value)
	if err != nil {
		return def, false
	}
	return t, true
}

func AddCategoryHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		newCategory := strings.TrimSpace(r.FormValue("new_category"))
		color := r.FormValue("color")
		if !domain.IsValidCategoryColor(color) {
			color = "c1"
		}
		if newCategory != "" {
			if err := app.AddCategory(r.Context(), newCategory, color); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/?tab=categories", http.StatusSeeOther)
	}
}

func UpdateCategoryHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		oldName := strings.TrimSpace(r.FormValue("old"))
		newName := strings.TrimSpace(r.FormValue("name"))
		color := r.FormValue("color")
		if oldName == "" || newName == "" {
			http.Error(w, "nazwa kategorii jest wymagana", http.StatusBadRequest)
			return
		}
		if !domain.IsValidCategoryColor(color) {
			color = "c1"
		}
		if err := app.UpdateCategory(r.Context(), oldName, newName, color); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/?tab=categories", http.StatusSeeOther)
	}
}

func DeleteCategoryHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		category := strings.TrimSpace(r.URL.Query().Get("category"))
		if category == "" {
			http.Error(w, "category is required", http.StatusBadRequest)
			return
		}
		if err := app.DeleteCategory(r.Context(), category); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func StaticFileHandlerBuilder(_ *controllers.AppController) http.Handler {
	return http.FileServer(http.FS(templates.Static))
}

func HealthcheckHandlerBuilder(app *controllers.AppController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if err := app.Ping(r.Context()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
