package core

import (
	"errors"

	"github.com/DO-2K23-26/polypass-microservices/credentials/infrastructure/sql"
	"github.com/DO-2K23-26/polypass-microservices/credentials/types"
)

type CredentialsService interface {
	GetPasswordCredentials(ids []string) ([]types.PasswordCredential, error)
	GetCardCredentials(ids []string) ([]types.CardCredential, error)
	GetSSHKeyCredentials(ids []string) ([]types.SSHKeyCredential, error)

	CreateCredential(credential types.CreateCredentialOpts) error
	CheckCredentialValidity(credential types.CreateCredentialOpts) error
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

func (c *credentialService) CreateCredential(credentialOpts *types.CreateCredentialOpts) error {
	switch credentialOpts.Type {
	case types.CredentialTypeCard:
		return c.CreateCardCredential(credentialOpts.CardAttributes)
	case types.CredentialTypePassword:
		return c.CreatePasswordCredential(credentialOpts.PasswordAttributes)
	default:
		return ERR_INVALID_CREDENTIAL_TYPE
	}
}
