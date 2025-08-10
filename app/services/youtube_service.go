package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/lib/ytdlp"
)

func GetVideoResolutionListService(c *fiber.Ctx, body types.YoutubeGetResolutionRequest) error {
	// Exec yt-dlp get resolution request
	ytdlp_class := ytdlp.Ytdlp{
		Url: body.URL,
	}
	resolution_data, err := ytdlp_class.GetListResolution()
	if err != nil {
		// Error
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
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
