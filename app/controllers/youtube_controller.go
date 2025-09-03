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
		return utils.ResponseError(c, err, "Bad Request Body!", fiber.StatusBadRequest)
	}

	// Youtube Link Validators
	if err := validators.YoutubeURLValidator(body.URL); err != nil {
		return utils.ResponseError(c, err, "Invalid youtube link!", fiber.StatusBadRequest)
	}

	return services.GetVideoResolutionListService(c, body)
}

func DownloadVideo(c *fiber.Ctx) error {
	// Get Request Body
	var body types.YoutubeDownloadRequest

	// General Validators
	err := utils.RequestBodyValidator(c, &body)
	if err != nil {
		return utils.ResponseError(c, err, "Bad Request Body!", fiber.StatusBadRequest)
	}

	// Youtube Link Validators
	if err := validators.YoutubeURLValidator(body.URL); err != nil {
		return utils.ResponseError(c, err, "Invalid youtube link!", fiber.StatusBadRequest)
	}

	return services.DownloadVideoService(c, body)
}

func DownloadVideoSection(c *fiber.Ctx) error {
	// Get Request Body
	var body types.YoutubeDownloadSectionRequest

	// General Validators
	err := utils.RequestBodyValidator(c, &body)
	if err != nil {
		return utils.ResponseError(c, err, "Bad Request Body!", fiber.StatusBadRequest)
	}

	// Youtube Link Validators
	if err := validators.YoutubeURLValidator(body.URL); err != nil {
		return utils.ResponseError(c, err, "Bad Request Body!", fiber.StatusBadRequest)
	}

	// Youtube Section Validators
	if err := validators.ValidateSectionTimes(body.StartTime, body.EndTime); err != nil {
		return utils.ResponseError(c, err, "Bad Request Body!", fiber.StatusBadRequest)
	}

	return services.DownloadVideoSectionService(c, body)
}
