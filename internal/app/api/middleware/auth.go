package middleware

import (
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware creates a middleware that checks for Bearer token or API key
// @description Authentication via Bearer token in Authorization header or api_key query parameter
func AuthMiddleware() fiber.Handler {
	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		apiToken = "secret-api-token"
	}

	return func(c *fiber.Ctx) error {
		// Check Authorization header for Bearer token
		authHeader := c.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				if parts[1] == apiToken {
					return c.Next()
				}
			}
		}

		// Check api_key query parameter
		apiKey := c.Query("api_key")
		if apiKey == apiToken {
			return c.Next()
		}

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unauthorized",
			"message": "Invalid or missing API token. Provide Authorization: Bearer <token> header or api_key query parameter",
		})
	}
}
