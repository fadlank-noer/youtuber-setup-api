package zerolog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

func SetupZerolog() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
}

func Logger() *zerolog.Logger {
	// Timestamp YYYY-MM-DD hh:mm:ss
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}

	// Uppercase Level
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s", i))
	}

	// Message Handler
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("| MSG: %s", i)
	}

	log := zerolog.New(output).With().Timestamp().Logger()
	return &log
}
