package credential

type Repository interface {
	CreateCredential(query CreateCredentialQuery) (*CreateCredentialResult, error)
	UpdateCredential(query UpdateCredentialQuery) (*UpdateCredentialResult, error)
	DeleteCredential(query DeleteCredentialQuery) error
	GetCredential(query GetCredentialQuery) (*GetCredentialResult, error)
	SearchCredentials(query SearchCredentialQuery) (*SearchCredentialResult, error)
}
