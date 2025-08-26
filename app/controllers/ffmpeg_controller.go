package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youtuber-setup-api/app/services"
	"github.com/youtuber-setup-api/app/types"
	"github.com/youtuber-setup-api/pkg/utils"
	"github.com/youtuber-setup-api/pkg/validators"
)

func WriteTmcdCompress(c *fiber.Ctx) error {
	// Define Body Type
	var body types.FfmpegTmcdCompressRequest

	// Check Uploaded File
	uploaded_file, err := utils.RequestBodyFileHandler(c, []string{"video_input"})
	if err != nil {
		return utils.ResponseError(c, err, "Bad Request File!", fiber.StatusBadRequest)
	}
	body = types.FfmpegTmcdCompressRequest{
		VideoInput: *uploaded_file["video_input"],
	}

	// General Validators
	if err = utils.RequestFormValidator(c, &body); err != nil {
		return utils.ResponseError(c, err, "Bad Request Body!", fiber.StatusBadRequest)
	}

	// Crf Code Validator
	if err := validators.CRFCodeValidator(body.CRFCode); err != nil {
		return utils.ResponseError(c, err, "Bad Request Body!", fiber.StatusBadRequest)
	}

	return services.WriteTmcdCompressService(c, &body)
}
