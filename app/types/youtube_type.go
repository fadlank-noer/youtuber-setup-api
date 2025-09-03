package types

import "github.com/youtuber-setup-api/lib/ytdlp"

type YoutubeGetResolutionRequest struct {
	URL string `json:"url"`
}

type YoutubeGetResolutionResponse struct {
	ytdlp.GetVideoMetadata
	ResolutionOption ytdlp.GetResolutionMediaFormats `json:"resolution_option"`
}

type YoutubeDownloadRequest struct {
	YoutubeGetResolutionRequest
	ID string `json:"resolution_id"`
}

type YoutubeDownloadSectionRequest struct {
	YoutubeDownloadRequest
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
