package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/GrzegorzMika/budget/templates"
	"github.com/a-h/templ"
	_ "github.com/mattn/go-sqlite3"
)

var categories = []string{"Środki czystości", "Spożywcze", "Odzież i obuwie", "Rozrywka", "Gastronomia", "Kosmetyki", "Edukacja", "Transport", "Zdrowie"}

var createTable = `CREATE TABLE IF NOT EXISTS expenses (
	timestamp TIMESTAMP WITH TIMEZONE,
	amount FLOAT,
	category TEXT
);`

var insertionQuery = `INSERT INTO expenses (timestamp, amount, category) VALUES (?,?,?)`

func expensesHandlerBuilder(db *sql.DB) http.HandlerFunc {
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
			_, err = db.Exec(insertionQuery, timestamp, amount, category)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	}
}

func main() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(1)

	_, err = db.Exec(createTable)
	if err != nil {
		panic(err.Error())
	}
	slices.Sort(categories)
	component := templates.Index(categories)
	http.Handle("/", templ.Handler(component))
	http.HandleFunc("/expenses", expensesHandlerBuilder(db))
	http.Handle("/static/", http.FileServer(http.FS(templates.Static)))

	fmt.Println("Listening on :3000")
	fmt.Println(http.ListenAndServe(":3000", nil))
}
