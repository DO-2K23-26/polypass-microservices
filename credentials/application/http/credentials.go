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
//
//	@Summary		Get password credentials
//	@Description	Get a list of password credentials
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			ids	query		string	true	"Comma-separated list of credential IDs"
//	@Success		200	{object}	[]types.PasswordCredential
//	@Failure		404	{object}	fiber.Map
//	@Router			/credentials/password [get]
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
	Title        string         `json:"title" db:"title"`
	Note         string         `json:"note" db:"note"`
	CustomFields map[string]any `json:"custom_fields" db:"custom_fields"`
	types.PasswordAttributes
	types.UserIdentifierAttribute
}

func (c *CreatePasswordCredentialOpts) Validate(ctx *fiber.Ctx) error {
	return c.BaseValidator.Validate(ctx, c)
}

// CreatePasswordCredential godoc
//
//	@Summary		Create password credential
//	@Description	Create a password credential
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePasswordCredentialOpts	true	"Create password credential options"
//	@Success		201		{object}	types.PasswordCredential
//	@Failure		400		{object}	fiber.Map
//	@Failure		500		{object}	fiber.Map
//	@Router			/credentials/password [post]
func (c *CredentialsController) CreatePasswordCredential() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(CreatePasswordCredentialOpts)
		if err := payload.Validate(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := c.service.CheckCredentialValidity(&types.CreateCredentialOpts{
			Type:               types.CredentialTypePassword,
			Title:              payload.Title,
			Note:               payload.Note,
			CustomFields:       payload.CustomFields,
			PasswordAttributes: payload.PasswordAttributes,
			UserIdentifierAttribute: types.UserIdentifierAttribute{
				UserIdentifier: "",
			},
		}); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
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

		return ctx.Status(fiber.StatusCreated).JSON(cred)
	}
}

// UpdatePasswordCredential godoc
//
//	@Summary		Update password credential
//	@Description	Update a password credential
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreatePasswordCredentialOpts	true	"Update password credential options"
//	@Param			id		path		string							true	"Credential ID"
//	@Success		200		{object}	types.PasswordCredential
//	@Failure		400		{object}	fiber.Map
//	@Failure		500		{object}	fiber.Map
//	@Router			/credentials/password/:id [put]
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
//
//	@Summary		Delete password credentials
//	@Description	Delete a list of password credentials
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			ids	query		string	true	"Comma-separated list of credential IDs"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/credentials/password [delete]
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

// GetCardCredentials godoc
//
//	@Summary		Get card credentials
//	@Description	Get a list of card credentials
//	@Tags			card
//	@Accept			json
//	@Produce		json
//	@Param			ids	query		string	true	"Comma-separated list of card IDs"
//	@Success		200	{object}	[]types.CardCredential
//	@Failure		404	{object}	fiber.Map
//	@Router			/credentials/card [get]
func (c *CredentialsController) GetCardCredentials() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ids_query := ctx.Query("ids")
		if ids_query == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ids is required",
			})
		}

		ids := strings.Split(ids_query, ",")
		credentials, err := c.service.GetCardCredentials(ids)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(credentials)
	}
}

type CreateCardCredentialOpts struct {
	BaseValidator
	Title        string         `json:"title" db:"title"`
	Note         string         `json:"note" db:"note"`
	CustomFields map[string]any `json:"custom_fields" db:"custom_fields"`
	types.CardAttributes
}

func (c *CreateCardCredentialOpts) Validate(ctx *fiber.Ctx) error {
	return c.BaseValidator.Validate(ctx, c)
}

// CreateCardCredential godoc
//
//	@Summary		Create card credential
//	@Description	Create a card credential
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateCardCredentialOpts	true	"Create card credential options"
//	@Success		201		{object}	types.CardCredential
//	@Failure		400		{object}	fiber.Map
//	@Failure		500		{object}	fiber.Map
//	@Router			/credentials/card [post]
func (c *CredentialsController) CreateCardCredential() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(CreateCardCredentialOpts)
		if err := payload.Validate(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := c.service.CheckCredentialValidity(&types.CreateCredentialOpts{
			Type:           types.CredentialTypeCard,
			Title:          payload.Title,
			Note:           payload.Note,
			CustomFields:   payload.CustomFields,
			CardAttributes: payload.CardAttributes,
			UserIdentifierAttribute: types.UserIdentifierAttribute{
				UserIdentifier: "",
			},
		}); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		cred, err := c.service.CreateCardCredential(types.CardCredential{
			Credential: types.Credential{
				Title:        payload.Title,
				Note:         payload.Note,
				CustomFields: &payload.CustomFields,
			},
			CardAttributes: payload.CardAttributes,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusCreated).JSON(cred)
	}
}

// UpdateCardCredential godoc
//
//	@Summary		Update card credential
//	@Description	Update a card credential
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateCardCredentialOpts	true	"Update card credential options"
//	@Param			id		path		string						true	"Credential ID"
//	@Success		200		{object}	types.CardCredential
//	@Failure		400		{object}	fiber.Map
//	@Failure		500		{object}	fiber.Map
//	@Router			/credentials/card/:id [put]
func (c *CredentialsController) UpdateCardCredential() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(CreateCardCredentialOpts)
		if err := payload.Validate(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		cred, err := c.service.UpdateCardCredential(types.CardCredential{
			Credential: types.Credential{
				ID:           ctx.Params("id"),
				Title:        payload.Title,
				Note:         payload.Note,
				CustomFields: &payload.CustomFields,
			},
			CardAttributes: payload.CardAttributes,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(cred)
	}
}

// DeleteCardCredentials godoc
//
//	@Summary		Delete card credentials
//	@Description	Delete a list of Card credentials
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			ids	query		string	true	"Comma-separated list of credential IDs"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/credentials/card [delete]
func (c *CredentialsController) DeleteCardCredentials() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ids_query := ctx.Query("ids")
		if ids_query == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ids is required",
			})
		}

		ids := strings.Split(ids_query, ",")
		err := c.service.DeleteCardCredentials(ids)
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

// GetSSHKeyCredentials godoc
//
//	@Summary		Get SSHKey credentials
//	@Description	Get a list of SSHKey credentials
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			ids	query		string	true	"Comma-separated list of credential IDs"
//	@Success		200	{object}	[]types.SSHKeyCredential
//	@Failure		404	{object}	fiber.Map
//	@Router			/credentials/sshkey [get]
func (c *CredentialsController) GetSSHKeyCredentials() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ids_query := ctx.Query("ids")
		if ids_query == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ids is required",
			})
		}

		ids := strings.Split(ids_query, ",")
		credentials, err := c.service.GetSSHKeyCredentials(ids)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(credentials)
	}
}

type CreateSSHCredentialOpts struct {
	BaseValidator
	Title        string         `json:"title" db:"title"`
	Note         string         `json:"note" db:"note"`
	CustomFields map[string]any `json:"custom_fields" db:"custom_fields"`
	types.SSHKeyAttributes
}

func (c *CreateSSHCredentialOpts) Validate(ctx *fiber.Ctx) error {
	return c.BaseValidator.Validate(ctx, c)
}

// CreateSSHKeyCredential godoc
//
//	@Summary		Create SSHKey credential
//	@Description	Create a SSHKey credential
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateSSHCredentialOpts	true	"Create SSHKey credential options"
//	@Success		201		{object}	types.SSHKeyCredential
//	@Failure		400		{object}	fiber.Map
//	@Failure		500		{object}	fiber.Map
//	@Router			/credentials/sshkey [post]
func (c *CredentialsController) CreateSSHKeyCredential() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(CreateSSHCredentialOpts)
		if err := payload.Validate(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// ajout du check de validité
		if err := c.service.CheckCredentialValidity(&types.CreateCredentialOpts{
			Type:             types.CredentialTypeSSHKey,
			Title:            payload.Title,
			Note:             payload.Note,
			CustomFields:     payload.CustomFields,
			SSHKeyAttributes: payload.SSHKeyAttributes,
			UserIdentifierAttribute: types.UserIdentifierAttribute{
				UserIdentifier: "",
			},
		}); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		cred, err := c.service.CreateSSHKeyCredential(types.SSHKeyCredential{
			Credential: types.Credential{
				Title:        payload.Title,
				Note:         payload.Note,
				CustomFields: &payload.CustomFields,
			},
			SSHKeyAttributes: payload.SSHKeyAttributes,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusCreated).JSON(cred)
	}
}

// UpdateSSHKeyCredential godoc
//
//	@Summary		Update SSHKey credential
//	@Description	Update a SSHKey credential
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateSSHCredentialOpts	true	"Update SSHKey credential options"
//	@Param			id		path		string						true	"Credential ID"
//	@Success		200		{object}	types.SSHKeyCredential
//	@Failure		400		{object}	fiber.Map
//	@Failure		500		{object}	fiber.Map
//	@Router			/credentials/sshkey/:id [put]
func (c *CredentialsController) UpdateSSHKeyCredential() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		payload := new(CreateSSHCredentialOpts)
		if err := payload.Validate(ctx); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		cred, err := c.service.UpdateSSHKeyCredential(types.SSHKeyCredential{
			Credential: types.Credential{
				ID:           ctx.Params("id"),
				Title:        payload.Title,
				Note:         payload.Note,
				CustomFields: &payload.CustomFields,
			},
			SSHKeyAttributes: payload.SSHKeyAttributes,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(cred)
	}
}

// DeleteSSHKeyCredentials godoc
//
//	@Summary		Delete SSHKey credentials
//	@Description	Delete a list of SSHKey credentials
//	@Tags			credentials
//	@Accept			json
//	@Produce		json
//	@Param			ids	query		string	true	"Comma-separated list of credential IDs"
//	@Success		200	{object}	fiber.Map
//	@Failure		400	{object}	fiber.Map
//	@Failure		500	{object}	fiber.Map
//	@Router			/credentials/sshkey [delete]
func (c *CredentialsController) DeleteSSHKeyCredentials() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ids_query := ctx.Query("ids")
		if ids_query == "" {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ids is required",
			})
		}

		ids := strings.Split(ids_query, ",")
		err := c.service.DeleteSSHKeyCredentials(ids)
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
	app.Get("/credentials/card", c.GetCardCredentials())
	app.Post("/credentials/card", c.CreateCardCredential())
	app.Put("/credentials/card/:id", c.UpdateCardCredential())
	app.Delete("/credentials/card", c.DeleteCardCredentials())
	app.Get("/credentials/sshkey", c.GetSSHKeyCredentials())
	app.Post("/credentials/sshkey", c.CreateSSHKeyCredential())
	app.Put("/credentials/sshkey/:id", c.UpdateSSHKeyCredential())
	app.Delete("/credentials/sshkey", c.DeleteSSHKeyCredentials())
}
