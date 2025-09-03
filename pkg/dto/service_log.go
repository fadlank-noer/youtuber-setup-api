package dto

import (
	"context"
	"encoding/json"

	"github.com/sqlc-dev/pqtype"
	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/deps"
	"github.com/youtuber-setup-api/query"
)

func ServiceLogging(service_name string, metadata interface{}) {
	// Context Definer
	ctx := context.Background()

	// Call DB
	conn, err := deps.PGConn()
	if err != nil {
		zerolog.Logger().Error().Msg("Error When Inserting Service Log!")
		return
	}
	defer conn.Close()

	// Query Call
	queries := query.New(conn)

	// Marshal metadata ke JSON
	metaBytes, err := json.Marshal(metadata)
	if err != nil {
		zerolog.Logger().Error().Err(err).Msg("Failed to marshal metadata")
		return
	}

	// Insert Query
	err = queries.CreateServiceLog(ctx, query.CreateServiceLogParams{
		ServiceName: service_name,
		Metadata: pqtype.NullRawMessage{
			RawMessage: metaBytes,
			Valid:      true,
		},
	})
	if err != nil {
		zerolog.Logger().Error().Err(err).Msg("Error inserting service log")
	}
}
