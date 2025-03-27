package main

import (
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	app.Post("/secret", func(c *fiber.Ctx) error {

	})

	app.Listen(":3000")
}