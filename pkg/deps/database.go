package deps

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/utils"
)

func PGConn() (*sql.DB, error) {
	// Get Database Configuration
	database_config, cfg_err := utils.ConnectionURLBuilder("database")
	if cfg_err != nil {
		return nil, fmt.Errorf("Error Calling Database Config!")
	}

	// Connect DB
	db, err := sql.Open("postgres", database_config)
	if err != nil {
		fmt.Printf("Database Error!: %w", err)
		zerolog.Logger().Error().Msg("Error when loading database!")
	}

	return db, err
}
