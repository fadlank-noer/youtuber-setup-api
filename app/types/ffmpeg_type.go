package types

import "github.com/youtuber-setup-api/pkg/utils"

type FfmpegTmcdCompressRequest struct {
	VideoInput utils.RequestFileResult `form:"video_input"`
	CRFCode    string                  `form:"crf_code"`
}
