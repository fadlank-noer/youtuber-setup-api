package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youtuber-setup-api/app/services"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/pkg/utils"
	"github.com/youtuber-setup-api/pkg/validators"
)

func GetVideoResolutionList(c *fiber.Ctx) error {
	// Get Request Body
	var body types.YoutubeGetResolutionRequest

	// General Validators
	err := utils.RequestBodyValidator(c, &body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request Body!",
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Youtube Link Validators
	if err := validators.YoutubeURLValidator(body.URL); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid youtube link!",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return services.GetVideoResolutionListService(c, body)
}

func DownloadVideo(c *fiber.Ctx) error {
	// Get Request Body
	var body types.YoutubeDownloadRequest

	// General Validators
	err := utils.RequestBodyValidator(c, &body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request Body!",
			"message": err.Error(),
			"data":    nil,
		})
	}

	// Youtube Link Validators
	if err := validators.YoutubeURLValidator(body.URL); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid youtube link!",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return services.DownloadVideoService(c, body)
}
