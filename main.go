package main

import (
	"github.com/gofiber/fiber/v2"

	_ "github.com/joho/godotenv/autoload" // load .env file automatically

	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/configs"
	"github.com/youtuber-setup-api/pkg/routes"
	"github.com/youtuber-setup-api/pkg/utils"
)

func main() {
	// Configure Zerolog Timestamp
	zerolog.SetupZerolog()

	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Define Routes
	routes.PublicRoutes(app)

	// Start Server
	utils.StartServer(app)
}
