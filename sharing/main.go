package main

import (
	"os"
	"sharing/dto"
	"sharing/services"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

var (
	service    services.SecretService
)

func init() {
	service = services.NewSecretService()
	godotenv.Load()
}

func main() {
	app := fiber.New()

	app.Post("/secret", func(c *fiber.Ctx) error {
		var request dto.PostSecretRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		response, err := service.CreateSecret(request)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create secret"})
		}

		return c.JSON(response)
	})

	app.Get("/secret/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response, err := service.GetSecret(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Secret not found"})
		}
		return c.JSON(response)
	})

	app.Listen(":"+os.Getenv("PORT"))
}