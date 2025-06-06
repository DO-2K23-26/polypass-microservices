package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type Validator interface {
	Validate(c *fiber.Ctx) error
}

type BaseValidator struct{}


func (b *BaseValidator) Validate(c *fiber.Ctx, v any) error {
	validate := validator.New()
	if err := c.BodyParser(v); err != nil {
		return err
	}

	if err := validate.Struct(v); err != nil {
		return err
	}

	return nil
}


