package main

import (
	"context"
	"fmt"
	"os"

	"github.com/GrzegorzMika/budget/controllers"
	"github.com/GrzegorzMika/budget/handlers"
	"github.com/GrzegorzMika/budget/migrations"
	"github.com/GrzegorzMika/budget/storage"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/pprof"
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

	api := fiber.New()

	api.Use(cors.New(
		cors.Config{AllowOrigins: []string{os.Getenv("CORS_ALLOW_ORIGINS")}},
	))
	api.Use(logger.New())
	api.Use(pprof.New())
	api.Use(handlers.AuthMiddleware(os.Getenv("JWKS_URL")))

	api.Get("/healthz", handlers.HealthcheckHandlerBuilder(app))
	api.Get("/readyz", handlers.HealthcheckHandlerBuilder(app))

	api.Get("/categories", handlers.ExpenseCategoryHandlerBuilder(app))
	api.Post("/expenses", handlers.ExpensesHandlerBuilder(app))

	fmt.Println("Listening on :3000")
	fmt.Println(api.Listen("0.0.0.0:3000"))
}
