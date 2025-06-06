package http

import (
	"strings"

	"github.com/DO-2K23-26/polypass-microservices/credentials/core"
	"github.com/DO-2K23-26/polypass-microservices/credentials/types"
	"github.com/gofiber/fiber/v2"
)

type CredentialsController struct {
	service core.CredentialsService
}

func NewCredentialsController(service core.CredentialsService) *CredentialsController {
	return &CredentialsController{
		service: service,
	}
}

// GetPasswordCredentials godoc
// @Summary Get password credentials
// @Description Get a list of password credentials
// @Tags credentials
// @Accept json
// @Produce json
// @Param ids query string true "Comma-separated list of credential IDs"
// @Success 200 {object} []types.PasswordCredential
// @Failure 404 {object} fiber.Map
// @Router /credentials/password [get]
func (c *CredentialsController) GetPasswordCredentials() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ids_query := ctx.Query("ids")
		if ids_query == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ids is required",
			})
		}

		ids := strings.Split(ids_query, ",")
		credentials, err := c.service.GetPasswordCredentials(ids)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(credentials)
	}
}

type CreatePasswordCredentialOpts struct {
	BaseValidator
	Title        string               `json:"title" db:"title"`
	Note         string               `json:"note" db:"note"`
	CustomFields map[string]any       `json:"custom_fields" db:"custom_fields"`
	types.PasswordAttributes
}

func (c *CreatePasswordCredentialOpts) Validate(ctx *fiber.Ctx) error {
	return c.BaseValidator.Validate(ctx, c)
}

// CreatePasswordCredential godoc
// @Summary Create password credential
// @Description Create a password credential
// @Tags credentials
// @Accept json
// @Produce json
// @Param payload body CreatePasswordCredentialOpts true "Create password credential options"
// @Success 200 {object} types.PasswordCredential
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /credentials/password [post]
func (c *CredentialsController) CreatePasswordCredential() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(CreatePasswordCredentialOpts)
		if err := payload.Validate(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		cred, err := c.service.CreatePasswordCredential(types.PasswordCredential{
			Credential: types.Credential{
				Title:        payload.Title,
				Note:         payload.Note,
				CustomFields: &payload.CustomFields,
			},
			PasswordAttributes: payload.PasswordAttributes,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(cred)
	}
}

// UpdatePasswordCredential godoc
// @Summary Update password credential
// @Description Update a password credential
// @Tags credentials
// @Accept json
// @Produce json
// @Param payload body CreatePasswordCredentialOpts true "Update password credential options"
// @Param id path string true "Credential ID"
// @Success 200 {object} types.PasswordCredential
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /credentials/password/:id [put]
func (c *CredentialsController) UpdatePasswordCredential() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(CreatePasswordCredentialOpts)
		if err := payload.Validate(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		cred, err := c.service.UpdatePasswordCredential(types.PasswordCredential{
			Credential: types.Credential{
				ID:           ctx.Params("id"),
				Title:        payload.Title,
				Note:         payload.Note,
				CustomFields: &payload.CustomFields,
			},
			PasswordAttributes: payload.PasswordAttributes,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(cred)
	}
}

// DeletePasswordCredentials godoc
// @Summary Delete password credentials
// @Description Delete a list of password credentials
// @Tags credentials
// @Accept json
// @Produce json
// @Param ids query string true "Comma-separated list of credential IDs"
// @Success 200 {object} fiber.Map
// @Failure 400 {object} fiber.Map
// @Failure 500 {object} fiber.Map
// @Router /credentials/password [delete]
func (c *CredentialsController) DeletePasswordCredentials() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ids_query := ctx.Query("ids")
		if ids_query == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ids is required",
			})
		}

		ids := strings.Split(ids_query, ",")
		err := c.service.DeletePasswordCredentials(ids)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "success",
		})
	}
}

func (c *CredentialsController) Register(app *fiber.App) {
	app.Get("/credentials/password", c.GetPasswordCredentials())
	app.Post("/credentials/password", c.CreatePasswordCredential())
	app.Put("/credentials/password/:id", c.UpdatePasswordCredential())
	app.Delete("/credentials/password", c.DeletePasswordCredentials())
}
