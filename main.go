package main

import (
	"github.com/gofiber/fiber/v2"

	"github.com/youtuber-setup-api/pkg/configs"
	"github.com/youtuber-setup-api/pkg/utils"
)

func main() {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	utils.StartServer(app)
}
