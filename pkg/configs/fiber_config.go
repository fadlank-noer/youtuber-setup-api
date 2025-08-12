package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	lib_zerolog "github.com/youtuber-setup-api/lib/zerolog"
)

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {
	// Define server settings.
	readTimeoutSecondsCount, _ := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	if readTimeoutSecondsCount == 0 {
		lib_zerolog.Logger().Panic().Msg("No ENV Loaded!")
		panic("")
	}

	// Return Fiber configuration.
	return fiber.Config{
		BodyLimit:   1024 * 1024 * 1024, // Up to 1 GB
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}
