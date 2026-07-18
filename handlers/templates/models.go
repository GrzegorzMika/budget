package templates

import (
	"fmt"
	"strings"
	"time"
)

// FeedItem is one expense row in the list feed.
type FeedItem struct {
	ID        string
	Title     string
	Subtitle  string
	Color     string
	AmountFmt string
}

// DayGroup is a day's worth of feed rows under a shared date label.
type DayGroup struct {
	Label string
	Items []FeedItem
}

var monthsNominative = [...]string{
	"Styczeń", "Luty", "Marzec", "Kwiecień", "Maj", "Czerwiec",
	"Lipiec", "Sierpień", "Wrzesień", "Październik", "Listopad", "Grudzień",
}

var monthsGenitive = [...]string{
	"stycznia", "lutego", "marca", "kwietnia", "maja", "czerwca",
	"lipca", "sierpnia", "września", "października", "listopada", "grudnia",
}

// MonthTitle renders a month heading like "Lipiec 2026".
func MonthTitle(t time.Time) string {
	return fmt.Sprintf("%s %d", monthsNominative[t.Month()-1], t.Year())
}

// DayLabel renders a date like "16 lipca".
func DayLabel(t time.Time) string {
	return fmt.Sprintf("%d %s", t.Day(), monthsGenitive[t.Month()-1])
}

// FormatAmount renders 1142.86 as "1 142,86" (Polish convention,
// non-breaking-space thousands separator).
func FormatAmount(v float64) string {
	s := fmt.Sprintf("%.2f", v)
	intPart, decPart, _ := strings.Cut(s, ".")
	neg := strings.HasPrefix(intPart, "-")
	intPart = strings.TrimPrefix(intPart, "-")

	var groups []string
	for len(intPart) > 3 {
		groups = append([]string{intPart[len(intPart)-3:]}, groups...)
		intPart = intPart[:len(intPart)-3]
	}
	out := intPart
	for _, g := range groups {
		out += " " + g
	}
	if neg {
		out = "-" + out
	}
	return out + "," + decPart
}

// PluralEntries returns the Polish plural form of "wpis" for n.
func PluralEntries(n int) string {
	if n == 1 {
		return "wpis"
	}
	if m := n % 10; m >= 2 && m <= 4 && !(n%100 >= 12 && n%100 <= 14) {
		return "wpisy"
	}
	return "wpisów"
}
