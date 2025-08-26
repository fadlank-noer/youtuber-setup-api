package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/youtuber-setup-api/lib/zerolog"
)

// StartServer func for starting a simple server.
func StartServer(a *fiber.App) {
	// Build Fiber connection URL.
	fiberConnURL, _ := ConnectionURLBuilder("fiber")
	zerolog.Logger().Info().Msg(fmt.Sprintf("Server start with config ", fiberConnURL))

	// Run server.
	if err := a.Listen(fiberConnURL); err != nil {
		zerolog.Logger().Error().Msg(fmt.Sprintf("Oops... Server is not running! Reason: %v", err))
	}
}
