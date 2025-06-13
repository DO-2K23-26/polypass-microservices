package main

import (
	"os"
	"sharing/dto"
	"sharing/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
	// Enable CORS for https://polypass.umontpeler.fr
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://polypass.umontpeler.fr",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "*",
		AllowCredentials: true,
	}))

	app.Post("/sharing", func(c *fiber.Ctx) error {
		var request dto.PostSecretRequest
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		useUserID := os.Getenv("USE_USER_ID")
		var userId string
		if useUserID == "false" {
			userId = ""
		} else {
			userId = ""
			authHeader := c.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				userId = authHeader[7:]
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing or invalid Authorization header"})
			}
		}

		response, err := service.CreateSecret(request, userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create secret"})
		}

		return c.JSON(response)
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/sharing/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		response, err := service.GetSecret(id)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Secret not found"})
		}
		return c.JSON(response)
	})

	app.Get("/sharing/history", func(c *fiber.Ctx) error {
		useUserID := os.Getenv("USE_USER_ID")

		if useUserID == "false" {
			return c.JSON([]dto.GetHistoryResponse{})
		}
		authHeader := c.Get("Authorization")
		var userId string
		if len(authHeader) > 7 && authHeader[:0] == "Bearer " {
			userId = authHeader[7:]
		} else if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header"})
		}
		response, err := service.GetHistory(userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve history"})
		}
		return c.JSON(response)
	})
	

	app.Listen(":"+os.Getenv("PORT"))
}