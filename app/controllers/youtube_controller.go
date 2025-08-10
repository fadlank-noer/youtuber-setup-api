package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youtuber-setup-api/app/services"
	"github.com/youtuber-setup-api/app/types"
)

func GetVideoResolutionList(c *fiber.Ctx) error {
	// Get Request Body
	var body types.YoutubeGetResolutionRequest

	// Parse body JSON ke struct
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "cannot parse JSON",
			"message": "error when parsing request body",
			"data":    nil,
		})
	}

	return services.GetVideoResolutionListService(c, body)
}
