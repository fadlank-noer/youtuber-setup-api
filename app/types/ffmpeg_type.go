package types

import "github.com/youtuber-setup-api/pkg/utils"

type FfmpegWriteTmcdRequest struct {
	VideoInput utils.RequestFileResult `json:"video_input"`
}
