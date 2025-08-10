package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youtuber-setup-api/app/controllers"
)

// Public Routes.
func PublicRoutes(a *fiber.App) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Youtube
	route.Post("/youtube/resolution_list", controllers.GetVideoResolutionList)
}
