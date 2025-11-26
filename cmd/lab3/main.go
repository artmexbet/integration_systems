// Nobel Prize API - Lab3
package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	recoverer "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/jackc/pgx/v5/pgxpool"
	slogfiber "github.com/samber/slog-fiber"

	v1 "ris/internal/app/api/v1"
	"ris/pkg/postgres/queries"
)

func main() {
	// Setup logging
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
	})
	slog.SetDefault(slog.New(handler))

	// Get configuration from environment
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/ris"
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Connect to database
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Verify database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	slog.Info("Connected to database")

	// Initialize queries and service
	q := queries.New(pool)
	service := v1.NewNobelService(q)
	apiHandler := v1.NewHandler(service)

	// Setup Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Nobel Prize API v1.0",
	})

	// Middleware
	app.Use(slogfiber.New(slog.Default()))
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))
	app.Use(requestid.New())
	app.Use(recoverer.New())

	// Health check endpoint (no auth required)
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"service": "Nobel Prize API",
			"version": "1.0",
		})
	})

	// Swagger UI - serve static files
	app.Get("/swagger/*", swaggerHandler())

	// Register API routes
	v1.RegisterRoutes(app, apiHandler)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Starting server", "port", port)
		if err := app.Listen(":" + port); err != nil {
			slog.Error("Server error", "error", err)
		}
	}()

	<-quit
	slog.Info("Shutting down server...")

	if err := app.Shutdown(); err != nil {
		slog.Error("Server shutdown error", "error", err)
	}

	slog.Info("Server stopped")
}

// swaggerHandler returns a handler for Swagger UI
func swaggerHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Params("*")

		if path == "" || path == "/" {
			return c.Redirect("/swagger/index.html")
		}

		if path == "index.html" {
			c.Set("Content-Type", "text/html")
			return c.SendString(swaggerUIHTML)
		}

		if path == "doc.json" {
			c.Set("Content-Type", "application/json")
			return c.SendString(swaggerSpec)
		}

		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	}
}
