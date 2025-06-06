package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/optique-dev/core"
)

type HealthController interface {
	Health() fiber.Handler
	Register(app *fiber.App)
}

type healthController struct{}

func NewHealthController() HealthController {
	return healthController{}
}

func (h healthController) Health() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.SendString("Healthy")
	}
}

func (h healthController) Register(app *fiber.App) {
	core.Info("Registering health handler")
	app.Get("/health", h.Health())
}
