package deps

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/utils"
)

var (
	dbInstance *sql.DB
	dbOnce     sync.Once
	dbErr      error
)

// initDB initializes the singleton DB connection and configures connection pooling.
func initDB() {
	// Get Database Configuration
	databaseConfig, cfgErr := utils.ConnectionURLBuilder("database")
	if cfgErr != nil {
		dbErr = fmt.Errorf("error calling database config: %w", cfgErr)
		zerolog.Logger().Error().Err(dbErr).Msg("Database config error")
		return
	}

	// Connect DB (sql.DB is a pooled connection manager)
	conn, err := sql.Open("postgres", databaseConfig)
	if err != nil {
		dbErr = fmt.Errorf("error opening database: %w", err)
		zerolog.Logger().Error().Err(dbErr).Msg("Database open error")
		return
	}

	// Configure pool parameters (tune as needed or move to config/env)
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxIdleTime(5 * time.Minute)
	conn.SetConnMaxLifetime(60 * time.Minute)

	// Verify connection
	if pingErr := conn.Ping(); pingErr != nil {
		dbErr = fmt.Errorf("error pinging database: %w", pingErr)
		zerolog.Logger().Error().Err(dbErr).Msg("Database ping error")
		_ = conn.Close()
		return
	}

	dbInstance = conn
	zerolog.Logger().Info().Msg("Database connection initialized")
}

// PGConn returns a singleton pooled *sql.DB. The first call initializes the pool.
func PGConn() (*sql.DB, error) {
	dbOnce.Do(initDB)
	return dbInstance, dbErr
}

// GetDB provides a convenient accessor to retrieve the shared *sql.DB instance
// without re-opening a new connection each time.
func GetDB() (*sql.DB, error) {
	return PGConn()
}
