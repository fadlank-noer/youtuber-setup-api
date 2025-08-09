package routes

import (
	"github.com/gofiber/fiber/v2"
)

// Public Routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")
}
