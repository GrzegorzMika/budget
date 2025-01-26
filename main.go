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

	app := controllers.NewAppController(storage.NewRepository(pool))

	http.Handle("/", handlers.LandingPageHandlerBuilder(app))
	http.HandleFunc("/expenses", handlers.ExpensesHandlerBuilder(app))
	http.Handle("/static/", handlers.StaticFileHandlerBuilder(app))

	fmt.Println("Listening on :3000")
	fmt.Println(http.ListenAndServe(":3000", nil))
}
