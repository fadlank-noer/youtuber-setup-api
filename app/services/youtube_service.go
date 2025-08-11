package services

import (
	"bufio"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/lib/ytdlp"
	"github.com/youtuber-setup-api/lib/zerolog"
)

func GetVideoResolutionListService(c *fiber.Ctx, body types.YoutubeGetResolutionRequest) error {
	// Exec yt-dlp get resolution request
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}
	resolution_data, err := ytdlp_class.GetListResolution()
	if err != nil || len(resolution_data.VideoOnly) < 1 {
		// Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err,
			"message": "Unexpected Error, Please try again later",
			"data":    nil,
		})
	}

	// Return Data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   nil,
		"message": "Success",
		"data":    resolution_data,
	})
}

func DownloadVideoService(c *fiber.Ctx, body types.YoutubeDownloadRequest) error {
	// Exec yt-dlp get title
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}
	_, err := ytdlp_class.GetVideoTitle()
	if err != nil {
		// Error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err,
			"message": "Unexpected Error, Please try again later",
			"data":    nil,
		})
	}

	// Set Client Header
	filename := "video"
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
