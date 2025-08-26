package utils

import "github.com/gofiber/fiber/v2"

func ResponseServerError(c *fiber.Ctx, err error) error {
	// Default Internal Server Error
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":   err.Error(),
		"message": "Internal Server Error!",
		"data":    nil,
	})

}

func ResponseError(c *fiber.Ctx, err error, msg string, status_code ...int) error {
	// Check Optional msg and status_code
	_msg := "Interval Server Error!"
	_status_code := fiber.StatusInternalServerError
	if msg != "" {
		_msg = msg
	}
	if len(status_code) > 0 {
		_status_code = status_code[0]
	}

	return c.Status(_status_code).JSON(fiber.Map{
		"error":   err.Error(),
		"message": _msg,
		"data":    nil,
	})
}

func ResponseSuccessJSON(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   nil,
		"message": "Success",
		"data":    data,
	})
}
