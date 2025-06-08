package core

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/DO-2K23-26/polypass-microservices/credentials/infrastructure/sql"
	"github.com/DO-2K23-26/polypass-microservices/credentials/types"
)

type CredentialsService interface {
	GetPasswordCredentials(ids []string) ([]types.PasswordCredential, error)
	GetCardCredentials(ids []string) ([]types.CardCredential, error)
	GetSSHKeyCredentials(ids []string) ([]types.SSHKeyCredential, error)

	CreateCredential(credential *types.CreateCredentialOpts) error
	CheckCredentialValidity(credential *types.CreateCredentialOpts) error
	CreatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error)
	CreateCardCredential(credential types.CardCredential) (types.CardCredential, error)
	CreateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error)

	UpdatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error)
	UpdateCardCredential(credential types.CardCredential) (types.CardCredential, error)
	UpdateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error)

	DeletePasswordCredentials(ids []string) error
	DeleteCardCredentials(ids []string) error
	DeleteSSHKeyCredentials(ids []string) error
}

type credentialService struct {
	sqlRepository sql.Sql
}

func NewCredentialService(sqlRepository sql.Sql) *credentialService {
	return &credentialService{
		sqlRepository: sqlRepository,
	}
}

func (c *credentialService) GetPasswordCredentials(ids []string) ([]types.PasswordCredential, error) {
	return c.sqlRepository.GetPasswordCredentials(ids)
}

func (c *credentialService) GetCardCredentials(ids []string) ([]types.CardCredential, error) {
	return c.sqlRepository.GetCardCredentials(ids)
}

func (c *credentialService) GetSSHKeyCredentials(ids []string) ([]types.SSHKeyCredential, error) {
	return c.sqlRepository.GetSSHKeyCredentials(ids)
}

func (c *credentialService) CreateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error) {
	return c.sqlRepository.CreateSSHKeyCredential(credential)
}

func (c *credentialService) UpdateSSHKeyCredential(credential types.SSHKeyCredential) (types.SSHKeyCredential, error) {
	return c.sqlRepository.UpdateSSHKeyCredential(credential)
}

func (c *credentialService) DeleteSSHKeyCredentials(ids []string) error {
	return c.sqlRepository.DeleteSSHKeyCredentials(ids)
}

func (c *credentialService) CreateCardCredential(credential types.CardCredential) (types.CardCredential, error) {
	return c.sqlRepository.CreateCardCredential(credential)
}

func (c *credentialService) UpdateCardCredential(credential types.CardCredential) (types.CardCredential, error) {
	return c.sqlRepository.UpdateCardCredential(credential)
}

func (c *credentialService) DeleteCardCredentials(ids []string) error {
	return c.sqlRepository.DeleteCardCredentials(ids)
}

func (c *credentialService) CreatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error) {
	return c.sqlRepository.CreatePasswordCredential(credential)
}

func (c *credentialService) UpdatePasswordCredential(credential types.PasswordCredential) (types.PasswordCredential, error) {
	return c.sqlRepository.UpdatePasswordCredential(credential)
}

func (c *credentialService) DeletePasswordCredentials(ids []string) error {
	return c.sqlRepository.DeletePasswordCredentials(ids)
}

var ERR_INVALID_CREDENTIAL_TYPE error = errors.New("invalid credential type")

// TODO it is not used yet
func (c *credentialService) CreateCredential(credentialOpts *types.CreateCredentialOpts) error {

	if err := c.CheckCredentialValidity(credentialOpts); err != nil {
		return err
	}

	switch credentialOpts.Type {
	case types.CredentialTypeCard:
		_, err := c.CreateCardCredential(types.CardCredential{
			Credential: types.Credential{
				Title:        credentialOpts.Title,
				Note:         credentialOpts.Note,
				CustomFields: &credentialOpts.CustomFields,
			},
			CardAttributes: credentialOpts.CardAttributes,
			UserIdentifierAttribute: types.UserIdentifierAttribute{
				UserIdentifier: credentialOpts.UserIdentifierAttribute.UserIdentifier,
			},
		})
		return err
	case types.CredentialTypePassword:
		_, err := c.CreatePasswordCredential(types.PasswordCredential{
			Credential: types.Credential{
				Title:        credentialOpts.Title,
				Note:         credentialOpts.Note,
				CustomFields: &credentialOpts.CustomFields,
			},
			PasswordAttributes: credentialOpts.PasswordAttributes,
			UserIdentifierAttribute: types.UserIdentifierAttribute{
				UserIdentifier: credentialOpts.UserIdentifierAttribute.UserIdentifier,
			},
		})
		return err
	case types.CredentialTypeSSHKey:
		_, err := c.CreateSSHKeyCredential(types.SSHKeyCredential{
			Credential: types.Credential{
				Title:        credentialOpts.Title,
				Note:         credentialOpts.Note,
				CustomFields: &credentialOpts.CustomFields,
			},
			SSHKeyAttributes: credentialOpts.SSHKeyAttributes,
			UserIdentifierAttribute: types.UserIdentifierAttribute{
				UserIdentifier: credentialOpts.UserIdentifierAttribute.UserIdentifier,
			},
		})
		return err
	default:
		return ERR_INVALID_CREDENTIAL_TYPE
	}
}

func (s *credentialService) CheckCredentialValidity(credentialOpts *types.CreateCredentialOpts) error {
	if credentialOpts == nil {
		return errors.New("credential options cannot be nil")
	}
	if credentialOpts.Title == "" {
		return errors.New("title cannot be empty")
	}

	switch credentialOpts.Type {
	case types.CredentialTypeCard:
		re := regexp.MustCompile(`^\d{16}$`)
		cardStr := strconv.FormatInt(credentialOpts.CardAttributes.CardNumber, 10)
		if !re.MatchString(cardStr) {
			return errors.New("invalid card number format: must be exactly 16 digits")
		}
		if credentialOpts.CardAttributes.CVC < 100 || credentialOpts.CardAttributes.CVC > 999 {
			return errors.New("CVC must be between 100 and 999")
		}
		expDate, err := time.Parse("2006-01-02", credentialOpts.CardAttributes.ExpirationDate)
		if err != nil {
			return errors.New("invalid expiry date format: expected YYYY-MM-DD")
		}
		if expDate.Before(time.Now()) {
			return errors.New("expired card")
		}
		if credentialOpts.CardAttributes.OwnerName == "" {
			return errors.New("owner name cannot be empty")
		}

	case types.CredentialTypePassword:
		if credentialOpts.PasswordAttributes.DomainName == "" {
			return errors.New("domain name cannot be empty")
		}
		if credentialOpts.PasswordAttributes.Password == "" {
			return errors.New("password cannot be empty")
		}

	case types.CredentialTypeSSHKey:
		if credentialOpts.SSHKeyAttributes.Hostname == "" {
			return errors.New("hostname cannot be empty")
		}
		if credentialOpts.SSHKeyAttributes.PrivateKey == "" {
			return errors.New("private key cannot be empty")
		}
		if credentialOpts.SSHKeyAttributes.PublicKey == "" {
			return errors.New("public key cannot be empty")
		}

	default:
		return ERR_INVALID_CREDENTIAL_TYPE
	}

	return nil
}
