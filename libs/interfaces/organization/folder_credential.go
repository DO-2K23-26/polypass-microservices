package organization

type CredentialType string

const (
	CredentialTypePassword CredentialType = "password"
	CredentialTypeCard     CredentialType = "card"
	CredentialTypeSSHKey   CredentialType = "sshkey"
)

type FolderCredential struct {
	IdFolder     string         `gorm:"type:uuid;primaryKey"`
	IdCredential string         `gorm:"type:uuid;primaryKey"`
	Type         CredentialType `gorm:"type:varchar(10);not null"`
}

// GetCredentialRequest represents the request for getting credentials.
type GetCredentialRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type GetCredentialResponse struct {
	Credentials []map[string]interface{} `json:"credentials"`
	Total       int                      `json:"total"`
	Page        int                      `json:"page"`
	Limit       int                      `json:"limit"`
}
