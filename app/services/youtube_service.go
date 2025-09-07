package services

import (
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/lib/ytdlp"
	"github.com/youtuber-setup-api/lib/zerolog"
	"github.com/youtuber-setup-api/pkg/dto"
	"github.com/youtuber-setup-api/pkg/utils"
)

func GetVideoResolutionListService(c *fiber.Ctx, body types.YoutubeGetResolutionRequest) error {
	// Initialize service logger
	logger := dto.NewServiceLogger("YT_RESOLUTION_LIST", "success")

	// Get user info
	ip := c.IP()
	logger.SetUserInfo(nil, &ip, nil, nil, nil)

	// Exec yt-dlp get resolution request
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}

	// Get Metadata
	ytdata, err := ytdlp_class.GetVideoMetadata()
	if err != nil {
		logger.SetError(err).Complete()
		return utils.ResponseError(c, err, "")
	}

	resolution_data, err := ytdlp_class.GetListResolution()
	if err != nil || len(resolution_data.VideoOnly) < 1 {
		logger.SetError(err).Complete()
		return utils.ResponseError(c, err, "")
	}

	// Set YouTube-specific metadata
	videoID := ytdata.ID
	logger.SetYouTubeFields(&body.URL, &videoID, nil, nil, nil, nil)

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

	// Log completion
	logger.SetField("resolution_count", len(resolution_data.VideoOnly)+len(resolution_data.VideoAudio)+len(resolution_data.AudioOnly))
	logger.Complete()

	// Return Data
	return utils.ResponseSuccessJSON(c, data)
}

func DownloadVideoService(c *fiber.Ctx, body types.YoutubeDownloadRequest) error {
	// Initialize service logger
	logger := dto.NewServiceLogger("YT_DOWNLOAD", "success")

	// Get user info
	ip := c.IP()
	logger.SetUserInfo(nil, &ip, nil, nil, nil)

	// Exec yt-dlp get video id
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}
	ytdata, err := ytdlp_class.GetVideoMetadata()
	if err != nil {
		logger.SetError(err).Complete()
		return utils.ResponseError(c, err, "")
	}

	// Set YouTube-specific metadata
	videoID := ytdata.ID
	format := "mp4"
	quality := body.ID
	logger.SetYouTubeFields(&body.URL, &videoID, &format, &quality, nil, nil)

	// Set Client Header
	filename := fmt.Sprintf("%s_youtube.mp4", ytdata.ID)
	c.Set("Content-Type", "video/mp4")
	c.Set("Content-Disposition", `attachment; filename="`+filename+`"`)

	// Stream Writer
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		if err := ytdlp_class.DownloadVideo(w, []string{"-f", body.ID, "-o", "-"}); err != nil {
			logger.SetError(err).Complete()
			zerolog.Logger().Error().Msg(fmt.Sprintln("Error streaming video:", err))
		} else {
			logger.SetField("filename", filename).Complete()
		}
	}))

	return nil
}

func DownloadVideoSectionService(c *fiber.Ctx, body types.YoutubeDownloadSectionRequest) error {
	// Initialize service logger
	logger := dto.NewServiceLogger("YT_DOWNLOAD_SECTION", "success")

	// Get user info
	ip := c.IP()
	logger.SetUserInfo(nil, &ip, nil, nil, nil)

	// Exec yt-dlp get video id
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}
	ytdata, err := ytdlp_class.GetVideoMetadata()
	if err != nil {
		logger.SetError(err).Complete()
		return utils.ResponseError(c, err, "")
	}

	// Set YouTube-specific metadata
	videoID := ytdata.ID
	format := "mp4"
	quality := body.ID
	logger.SetYouTubeFields(&body.URL, &videoID, &format, &quality, nil, nil)
	logger.SetField("start_time", body.StartTime).SetField("end_time", body.EndTime)

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
			logger.SetError(err).Complete()
			zerolog.Logger().Error().Msg(fmt.Sprintln("Error streaming video:", err))
		} else {
			logger.SetField("filename", filename).Complete()
		}
	}))

	return nil
}
