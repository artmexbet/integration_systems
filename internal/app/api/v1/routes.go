package v1

import (
	"github.com/gofiber/fiber/v2"

	"ris/internal/app/api/middleware"
)

// RegisterRoutes registers all API v1 routes
func RegisterRoutes(app *fiber.App, handler *Handler) {
	// API v1 group with authentication
	api := app.Group("/api/v1", middleware.AuthMiddleware())

	// Stats routes
	api.Get("/stats", handler.GetStats)
	api.Get("/stats/last-update", handler.GetLastUpdate)

	// Categories route
	api.Get("/categories", handler.GetCategories)

	// Laureates routes
	laureates := api.Group("/laureates")
	laureates.Get("/", handler.ListLaureates)
	laureates.Get("/:id", handler.GetLaureate)
	laureates.Post("/", handler.CreateLaureate)
	laureates.Put("/:id", handler.UpdateLaureate)
	laureates.Delete("/:id", handler.DeleteLaureate)

	// Prizes routes
	prizes := api.Group("/prizes")
	prizes.Get("/", handler.ListPrizes)
	prizes.Get("/category/:category", handler.GetPrizesByCategory)
	prizes.Get("/year/:year", handler.GetPrizesByYear)
	prizes.Get("/:id", handler.GetPrize)
	prizes.Post("/", handler.CreatePrize)
	prizes.Put("/:id", handler.UpdatePrize)
	prizes.Delete("/:id", handler.DeletePrize)
}
