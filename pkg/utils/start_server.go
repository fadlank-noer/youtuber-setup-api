package utils

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	// Build Fiber connection URL.
	fiberConnURL, _ := ConnectionURLBuilder("fiber")
	log.Print("Server start with config", fiberConnURL)

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}
}
