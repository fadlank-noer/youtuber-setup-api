package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

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

func RequestFormValidator(c *fiber.Ctx, form interface{}) error {
	// Ambil semua key dari form request
	mForm, err := c.MultipartForm()
	if err != nil {
		return fmt.Errorf("invalid multipart form: %v", err)
	}

	// Ambil semua field struct yang diizinkan (dari tag `form`)
	allowedFields := map[string]bool{}
	rt := reflect.TypeOf(form).Elem()
	for i := 0; i < rt.NumField(); i++ {
		tag := rt.Field(i).Tag.Get("form")
		if tag != "" {
			allowedFields[tag] = true
		}
	}

	// Cek apakah ada key yang tidak diizinkan
	for key := range mForm.Value {
		if !allowedFields[key] {
			return fmt.Errorf("invalid field detected: %s", key)
		}
	}

	// Bind string fields
	rv := reflect.ValueOf(form).Elem()
	for i := 0; i < rt.NumField(); i++ {
		field := rv.Field(i)
		tag := rt.Field(i).Tag.Get("form")

		if tag == "" {
			continue
		}

		// Reflect Form Value to Any Types
		if field.Kind() == reflect.String {
			field.SetString(c.FormValue(tag))
		}
	}

	// Validasi struct
	return validate.Struct(form)
}
