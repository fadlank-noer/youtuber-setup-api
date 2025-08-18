package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youtuber-setup-api/app/services"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/pkg/utils"
)

func WriteTmcd(c *fiber.Ctx) error {
	// Define Body Type
	var body types.FfmpegWriteTmcdRequest

	// Check Uploaded File
	uploaded_file, err := utils.RequestBodyFileHandler(c, []string{"video_input"})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request File!",
			"message": err.Error(),
			"data":    nil,
		})
	}
	body = types.FfmpegWriteTmcdRequest{
		VideoInput: *uploaded_file["video_input"],
	}

	// General Validators
	if err = utils.RequestFormValidator(c, &body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Bad Request Body!",
			"message": err.Error(),
			"data":    nil,
		})
	}

	return services.WriteTmcdService(c, &body)
}
