package configs

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/youtuber-setup-api/lib/zerolog"
)

// FiberConfig func for configuration Fiber app.
// See: https://docs.gofiber.io/api/fiber#config
func FiberConfig() fiber.Config {
	// Check ENV File
	if err := godotenv.Load(); err != nil {
		zerolog.Logger().Panic().Msg("No ENV Loaded!")
	}

	// Define server settings.
	readTimeoutSecondsCount, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		zerolog.Logger().Panic().Msg("Invald SERVER_READ_TIMEOUT configuration!")
	}

	// Return Fiber configuration.
	return fiber.Config{
		BodyLimit:   305 * 1024 * 1024, // Up to 300 MB
		ReadTimeout: time.Second * time.Duration(readTimeoutSecondsCount),
	}
}
