package services

import (
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/lib/ytdlp"
	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/utils"
)

func GetVideoResolutionListService(c *fiber.Ctx, body types.YoutubeGetResolutionRequest) error {
	// Exec yt-dlp get resolution request
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}

	// Get Metadata
	ytdata, err := ytdlp_class.GetVideoMetadata()
	if err != nil {
		return utils.ResponseError(c, err, "")
	}

	resolution_data, err := ytdlp_class.GetListResolution()
	if err != nil || len(resolution_data.VideoOnly) < 1 {
		return utils.ResponseError(c, err, "")
	}

	// Compose Data
	data := types.YoutubeGetResolutionResponse{
		GetVideoMetadata: ytdlp.GetVideoMetadata{
			ID:           ytdata.ID,
			ThumbnailURL: ytdata.ThumbnailURL,
			Title:        ytdata.Title,
			Duration:     ytdata.Duration,
		},
		ResolutionOption: ytdlp.GetResolutionMediaFormats{
			VideoAudio: resolution_data.VideoAudio,
			VideoOnly:  resolution_data.VideoOnly,
			AudioOnly:  resolution_data.AudioOnly,
		},
	}

	// Return Data
	return utils.ResponseSuccessJSON(c, data)
}

func DownloadVideoService(c *fiber.Ctx, body types.YoutubeDownloadRequest) error {
	// Exec yt-dlp get video id
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}
	ytdata, err := ytdlp_class.GetVideoMetadata()
	if err != nil {
		return utils.ResponseError(c, err, "")
	}

	// Set Client Header
	filename := fmt.Sprintf("%s_youtube.mp4", ytdata.ID)
	c.Set("Content-Type", "video/mp4")
	c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)

	// Stream Writer
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if err := ytdlp_class.DownloadVideo(w, []string{"-f", body.ID, "-o", "-"}); err != nil {
			zerolog.Logger().Error().Msg(fmt.Sprintln("Error streaming video:", err))
		}
	}))

	return nil
}

func DownloadVideoSectionService(c *fiber.Ctx, body types.YoutubeDownloadSectionRequest) error {
	// Exec yt-dlp get video id
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}
	ytdata, err := ytdlp_class.GetVideoMetadata()
	if err != nil {
		return utils.ResponseError(c, err, "")
	}

	// Set Client Header
	filename := fmt.Sprintf("%s_youtube.mp4", ytdata.ID)
	c.Set("Content-Type", "video/mp4")
	c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)

	// Section Check
	section := fmt.Sprintf("\"*%s-%s\"", body.StartTime, body.EndTime)

	// Args Builder
	args := []string{
		"-f", body.ID,
		"--download-sections", section, "--force-keyframes-at-cuts",
		"-o", "-",
	}

	// Stream Writer
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if err := ytdlp_class.DownloadVideo(w, args); err != nil {
			zerolog.Logger().Error().Msg(fmt.Sprintln("Error streaming video:", err))
		}
	}))

	return nil
}
