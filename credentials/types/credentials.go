package types

import "time"

type Credential struct {
	ID           string          `json:"id" db:"id"`
	Title        string          `json:"title" db:"title"`
	Note         string          `json:"note" db:"note"`
	CreatedAt    *time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time      `json:"updated_at" db:"updated_at"`
	ExpiresAt    *time.Time      `json:"expires_at" db:"expires_at"`
	LastReadAt   *time.Time      `json:"last_read_at" db:"last_read_at"`
	CustomFields *map[string]any `json:"custom_fields" db:"custom_fields"`
}

type CardCredential struct {
	Credential
	CardAttributes
	UserIdentifierAttribute
}

type CardAttributes struct {
	OwnerName      string `json:"owner_name" db:"owner_name"`
	CVC            int    `json:"cvc" db:"cvc"`
	ExpirationDate string `json:"expiration_date" db:"expiration_date"`
	CardNumber     int64  `json:"card_number" db:"card_number"`
}

type PasswordCredential struct {
	Credential
	PasswordAttributes
	UserIdentifierAttribute
}

type PasswordAttributes struct {
	Password       string `json:"password" db:"password"`
	DomainName     string `json:"domain_name" db:"domain_name"`
}

type SSHKeyCredential struct {
	Credential
	SSHKeyAttributes
	UserIdentifierAttribute
}

type SSHKeyAttributes struct {
	PrivateKey string `json:"private_key" db:"private_key"`
	PublicKey  string `json:"public_key" db:"public_key"`
	Hostname   string `json:"hostname" db:"hostname"`
}

type UserIdentifierAttribute struct {
	UserIdentifier string `json:"user_identifier" db:"user_identifier"`
}

type CredentialType string

const (
	CredentialTypeCard     CredentialType = "card"
	CredentialTypePassword CredentialType = "password"
	CredentialTypeSSHKey   CredentialType = "ssh_key"
)

type CreateCredentialOpts struct {
	Title        string         `json:"title" db:"title"`
	Note         string         `json:"note" db:"note"`
	CustomFields map[string]any `json:"custom_fields" db:"custom_fields"`
	Type         CredentialType `json:"type" validate:"required"`
	SSHKeyAttributes
	PasswordAttributes
	CardAttributes
	UserIdentifierAttribute
}
