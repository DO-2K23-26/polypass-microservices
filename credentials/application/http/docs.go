package http

import (
	"github.com/MarceloPetrucio/go-scalar-api-reference"
	"github.com/gofiber/fiber/v2"
)

type DocsController struct{}

func (d DocsController) Register(app *fiber.App) {
	app.Get("/docs", func(c *fiber.Ctx) error {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/swagger.json",
			Theme:   scalar.ThemeKepler,
			CustomOptions: scalar.CustomOptions{
				PageTitle: "Courses API",
			},
			DarkMode: true,
		})

		if err != nil {
			return c.Status(fiber.ErrInternalServerError.Code).JSON(map[string]string{
				"error": err.Error(),
			})
		}

		return c.Format(htmlContent)
	})
	app.Static("/config.schema.json", "./docs/swagger.json")
}

func NewDocsController() DocsController {
	return DocsController{}
}
