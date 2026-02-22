package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/GrzegorzMika/budget/controllers"
	"github.com/GrzegorzMika/budget/handlers"
	"github.com/GrzegorzMika/budget/migrations"
	"github.com/GrzegorzMika/budget/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := pgxpool.ParseConfig(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err.Error())
	}
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		panic(err.Error())
	}

	err = migrations.PerformMigrations(ctx)
	if err != nil {
		panic(err.Error())
	}

	repo := storage.NewRepository(pool)
	categories, err := repo.GetCategories(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to load categories: %w", err).Error())
	}

	app := controllers.NewAppController(repo, categories)

	http.HandleFunc("/", handlers.LandingPageHandlerBuilder(app))
	http.HandleFunc("/expenses", handlers.ExpensesHandlerBuilder(app))

	addCategory := handlers.AddCategoryHandlerBuilder(app)
	delCategory := handlers.DeleteCategoryHandlerBuilder(app)
	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "DELETE" {
			delCategory(w, r)
		} else {
			addCategory(w, r)
		}
	})

	http.Handle("/assets/", handlers.StaticFileHandlerBuilder(app))
	http.Handle("/healthz", handlers.HealthcheckHandlerBuilder(app))
	http.Handle("/readyz", handlers.HealthcheckHandlerBuilder(app))

	fmt.Println("Listening on :3000")
	fmt.Println(http.ListenAndServe(":3000", nil))
}
