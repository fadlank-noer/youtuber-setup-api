package types

type YoutubeGetResolutionRequest struct {
	URL string `json:"url"`
}

type YoutubeDownloadRequest struct {
	YoutubeGetResolutionRequest
	ID string `json:"resolution_id"`
}
