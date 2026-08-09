package handlers

import (
	"fmt"
	"math"
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
		case "list":
			component, err := buildListTab(r, app, now)
			if err != nil {
				http.Error(w, "Failed to load expenses: "+err.Error(), http.StatusInternalServerError)
				return
			}
			tabComponent = component
		case "summary":
			component, err := buildSummaryTab(r, app, now)
			if err != nil {
				http.Error(w, "Failed to load summary: "+err.Error(), http.StatusInternalServerError)
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

// period is a user-selected date range, defaulting to the current month.
// from and to are inclusive; end() converts to the exclusive upper bound the
// repository expects.
type period struct {
	from, to           time.Time
	fromGiven, toGiven bool
}

func parsePeriod(r *http.Request, now time.Time) period {
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, -1)

	from, fromGiven := parseDateParam(r, "from", monthStart)
	to, toGiven := parseDateParam(r, "to", monthEnd)
	if to.Before(from) {
		from, to = to, from
		fromGiven, toGiven = toGiven, fromGiven
	}
	return period{from: from, to: to, fromGiven: fromGiven, toGiven: toGiven}
}

func (p period) active() bool { return p.fromGiven || p.toGiven }

func (p period) end() time.Time { return p.to.AddDate(0, 0, 1) }

func (p period) label(now time.Time) string {
	switch {
	case p.fromGiven && p.toGiven:
		return templates.DayLabel(p.from) + " – " + templates.DayLabel(p.to)
	case p.fromGiven:
		return "od " + templates.DayLabel(p.from)
	case p.toGiven:
		return "do " + templates.DayLabel(p.to)
	default:
		return templates.MonthTitle(now)
	}
}

func (p period) fromValue() string {
	if p.fromGiven {
		return p.from.Format(time.DateOnly)
	}
	return ""
}

func (p period) toValue() string {
	if p.toGiven {
		return p.to.Format(time.DateOnly)
	}
	return ""
}

func buildListTab(r *http.Request, app *controllers.AppController, now time.Time) (templ.Component, error) {
	p := parsePeriod(r, now)
	query := strings.TrimSpace(r.URL.Query().Get("q"))

	expenses, err := app.GetExpenses(r.Context(), p.from, p.end(), query)
	if err != nil {
		return nil, err
	}

	total := 0.0
	for _, e := range expenses {
		total += e.Amount
	}

	periodText := "wydane w tym miesiącu"
	if p.active() {
		periodText = "wydane w wybranym okresie"
	}
	countLabel := periodText + " · " + strconv.Itoa(len(expenses)) + " " + templates.PluralEntries(len(expenses))

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

	return templates.ListTab(p.label(now), templates.FormatAmount(total), countLabel, query, p.fromValue(), p.toValue(), p.active() || query != "", groups), nil
}

func buildSummaryTab(r *http.Request, app *controllers.AppController, now time.Time) (templ.Component, error) {
	p := parsePeriod(r, now)

	totals, err := app.GetCategoryTotals(r.Context(), p.from, p.end())
	if err != nil {
		return nil, err
	}

	// Category selection defaults to all. Submitted filter forms carry a
	// catsel=1 marker so "every pill unchecked" is distinguishable from
	// "no category filter in the URL at all".
	q := r.URL.Query()
	catFiltered := q.Get("catsel") == "1" || len(q["cat"]) > 0
	selected := make(map[string]bool, len(q["cat"]))
	for _, name := range q["cat"] {
		selected[name] = true
	}

	colors := make(map[string]string)
	for _, c := range app.GetCategories() {
		colors[c.Name] = c.Color
	}
	colorOf := func(name string) string {
		if color, ok := colors[name]; ok {
			return color
		}
		return domain.NeutralCategoryColor
	}

	options := make([]templates.SummaryOption, 0, len(totals))
	var included []domain.CategoryTotal
	total, entries := 0.0, 0
	for _, t := range totals {
		on := !catFiltered || selected[t.Category]
		options = append(options, templates.SummaryOption{Name: t.Category, Color: colorOf(t.Category), Checked: on})
		if on {
			included = append(included, t)
			total += t.Total
			entries += t.Count
		}
	}

	rows := make([]templates.SummaryRow, 0, len(included))
	slices := make([]templates.DonutSlice, 0, len(included))
	start := 0.0
	for _, t := range included {
		frac := t.Total / total
		rows = append(rows, templates.SummaryRow{
			Name:      t.Category,
			Color:     colorOf(t.Category),
			AmountFmt: templates.FormatAmount(t.Total),
			Percent:   formatPercent(frac),
		})
		slices = append(slices, donutSlice(colorOf(t.Category), frac, start, len(included)))
		start += frac
	}

	countLabel := strconv.Itoa(entries) + " " + templates.PluralEntries(entries) +
		" · " + strconv.Itoa(len(included)) + " " + templates.PluralCategories(len(included))

	return templates.SummaryTab(
		p.label(now), templates.FormatAmount(total), countLabel,
		p.fromValue(), p.toValue(), p.active() || catFiltered,
		options, slices, rows,
	), nil
}

// Donut geometry: slices are stroke dash segments on an r=40 circle in a
// 100×100 viewBox; the template rotates the group -90° so the first slice
// starts at 12 o'clock.
const donutCircumference = 2 * math.Pi * 40

// donutGap is the hairline gap between adjacent slices, in viewBox units.
const donutGap = 1.6

func donutSlice(color string, frac, start float64, sliceCount int) templates.DonutSlice {
	length := frac * donutCircumference
	gap := donutGap
	if sliceCount < 2 || length < 2*donutGap {
		gap = 0
	}
	visible := length - gap
	return templates.DonutSlice{
		Color:  color,
		Dash:   fmt.Sprintf("%.3f %.3f", visible, donutCircumference-visible),
		Offset: fmt.Sprintf("%.3f", -(start*donutCircumference + gap/2)),
	}
}

// formatPercent renders a fraction of the total as an integer percentage,
// with "<1%" for slivers that would round to nothing.
func formatPercent(frac float64) string {
	pct := frac * 100
	if pct < 1 {
		return "<1%"
	}
	return strconv.Itoa(int(math.Round(pct))) + "%"
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
