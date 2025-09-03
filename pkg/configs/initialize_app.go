package configs

import (
	"fmt"
	"os"

	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/deps"
)

func InitializeApp() error {
	// Ensure Dir
	if err := EnsureTMPdirs(); err != nil {
		zerolog.Logger().Panic().Msg("Error when creating tmp dirs!")
	}
	zerolog.Logger().Info().Msg("~~ TMP DIR Success ~~")

	// Ensure Database
	if _, err := deps.PGConn(); err != nil {
		zerolog.Logger().Panic().Msg("Database Error!")
	}
	zerolog.Logger().Info().Msg("~~ Database Success ~~")

	return nil
}

func EnsureTMPdirs() error {
	// Needed Dir
	dirs := []string{"./tmp_file", "./tmp_public"}

	for _, dir := range dirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			// Create folder with 0755 permission
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return fmt.Errorf("Error creating DIR %s: %w", dir, err)
			}
			fmt.Println("DIR created:", dir)
		} else {
			fmt.Println("DIR exist:", dir)
		}
	}
	return nil
}
