package dto

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/deps"
	"github.com/youtuber-setup-api/query"
)

// ServiceLogMetadata represents the comprehensive metadata for service logging
type ServiceLogMetadata struct {
	RequestID      string                 `json:"request_id"`
	TimestampStart time.Time              `json:"timestamp_start"`
	TimestampEnd   *time.Time             `json:"timestamp_end,omitempty"`
	UserID         *string                `json:"user_id,omitempty"`
	IPAddress      *string                `json:"ip_address,omitempty"`
	ServiceName    string                 `json:"service_name"`
	Status         string                 `json:"status"` // success, failed, canceled
	DurationMs     *int64                 `json:"duration_ms,omitempty"`
	ErrorMessage   *string                `json:"error_message,omitempty"`
	UserDevice     *string                `json:"user_device,omitempty"`
	GeoLocation    *string                `json:"geo_location,omitempty"`
	ConversionSource *string              `json:"conversion_source,omitempty"`

	// YouTube specific fields
	VideoURL             *string `json:"video_url,omitempty"`
	VideoID              *string `json:"video_id,omitempty"`
	VideoDuration        *int    `json:"video_duration,omitempty"`
	VideoFormat          *string `json:"video_format,omitempty"`
	VideoQualityRequested *string `json:"video_quality_requested,omitempty"`
	VideoSize            *int64  `json:"video_size,omitempty"`

	// FFmpeg specific fields
	InputSize        *int64   `json:"input_size,omitempty"`
	OutputSize       *int64   `json:"output_size,omitempty"`
	CompressionRatio *float64 `json:"compression_ratio,omitempty"`
	Codec            *string  `json:"codec,omitempty"`
	BitrateTarget    *string  `json:"bitrate_target,omitempty"`

	// Additional custom fields
	Custom map[string]interface{} `json:"custom,omitempty"`
}

// ServiceLogger manages service logging with timing
type ServiceLogger struct {
	metadata ServiceLogMetadata
	startTime time.Time
}

// NewServiceLogger creates a new service logger with initial metadata
func NewServiceLogger(serviceName, status string) *ServiceLogger {
	requestID := uuid.New().String()
	return &ServiceLogger{
		metadata: ServiceLogMetadata{
			RequestID:      requestID,
			TimestampStart: time.Now(),
			ServiceName:    serviceName,
			Status:         status,
			Custom:         make(map[string]interface{}),
		},
		startTime: time.Now(),
	}
}

// SetField sets a custom field in the metadata
func (sl *ServiceLogger) SetField(key string, value interface{}) *ServiceLogger {
	if sl.metadata.Custom == nil {
		sl.metadata.Custom = make(map[string]interface{})
	}
	sl.metadata.Custom[key] = value
	return sl
}

// SetYouTubeFields sets YouTube-specific metadata
func (sl *ServiceLogger) SetYouTubeFields(videoURL, videoID, format, quality *string, duration *int, size *int64) *ServiceLogger {
	sl.metadata.VideoURL = videoURL
	sl.metadata.VideoID = videoID
	sl.metadata.VideoFormat = format
	sl.metadata.VideoQualityRequested = quality
	sl.metadata.VideoDuration = duration
	sl.metadata.VideoSize = size
	return sl
}

// SetFFmpegFields sets FFmpeg-specific metadata
func (sl *ServiceLogger) SetFFmpegFields(inputSize, outputSize *int64, codec, bitrate *string) *ServiceLogger {
	sl.metadata.InputSize = inputSize
	sl.metadata.OutputSize = outputSize
	sl.metadata.Codec = codec
	sl.metadata.BitrateTarget = bitrate
	if inputSize != nil && outputSize != nil && *inputSize > 0 {
		ratio := float64(*outputSize) / float64(*inputSize)
		sl.metadata.CompressionRatio = &ratio
	}
	return sl
}

// SetUserInfo sets user-related metadata
func (sl *ServiceLogger) SetUserInfo(userID, ip, device, geo, source *string) *ServiceLogger {
	sl.metadata.UserID = userID
	sl.metadata.IPAddress = ip
	sl.metadata.UserDevice = device
	sl.metadata.GeoLocation = geo
	sl.metadata.ConversionSource = source
	return sl
}

// SetError sets error information and marks status as failed
func (sl *ServiceLogger) SetError(err error) *ServiceLogger {
	if err != nil {
		sl.metadata.Status = "failed"
		errMsg := err.Error()
		sl.metadata.ErrorMessage = &errMsg
	}
	return sl
}

// Complete marks the service as completed and logs it
func (sl *ServiceLogger) Complete() {
	endTime := time.Now()
	sl.metadata.TimestampEnd = &endTime
	duration := endTime.Sub(sl.startTime).Milliseconds()
	sl.metadata.DurationMs = &duration
	sl.Log()
}

// Log saves the service log to database
func (sl *ServiceLogger) Log() {
	ctx := context.Background()
	conn, err := deps.GetDB()
	if err != nil {
		zerolog.Logger().Error().Err(err).Msg("Failed to get DB connection for service log")
		return
	}

	queries := query.New(conn)
	metaBytes, err := json.Marshal(sl.metadata)
	if err != nil {
		zerolog.Logger().Error().Err(err).Msg("Failed to marshal service log metadata")
		return
	}

	err = queries.CreateServiceLog(ctx, query.CreateServiceLogParams{
		ServiceName: sl.metadata.ServiceName,
		Metadata: pqtype.NullRawMessage{
			RawMessage: metaBytes,
			Valid:      true,
		},
	})
	if err != nil {
		zerolog.Logger().Error().Err(err).Msg("Error inserting service log")
	}
}

// QuickLog provides a simple way to log with basic metadata
func QuickLog(serviceName, status string, metadata map[string]interface{}) {
	logger := NewServiceLogger(serviceName, status)
	for k, v := range metadata {
		logger.SetField(k, v)
	}
	logger.Complete()
}
