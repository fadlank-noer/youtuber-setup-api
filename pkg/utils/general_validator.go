package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func RequestBodyValidator(c *fiber.Ctx, body interface{}) error {
	// Unknown Body Attribute Validator
	decoder := json.NewDecoder(bytes.NewReader(c.Body()))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(body); err != nil { // body harus pointer
		return errors.New(fmt.Sprintf("Invalid field detected: %v", err))
	}

	// Struct Validator
	if err := validate.Struct(body); err != nil {
		_errors := []string{}
		for _, err := range err.(validator.ValidationErrors) {
			_errors = append(_errors, fmt.Sprintf("Field '%s' failed on '%s' tag", err.Field(), err.Tag()))
		}

		return errors.New("Validation Error")
	}

	return nil
}
