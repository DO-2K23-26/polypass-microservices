package credential

type RepoCredential struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Tags     []types.Tag  `json:"tags"`
	FolderId string `json:"folder_id"`
}

type ICredentialRepository interface {
	CreateCredential(query CreateCredentialQuery) (*CreateCredentialResult, error)
	UpdateCredential(query UpdateCredentialQuery) (*UpdateCredentialResult, error)
	DeleteCredential(query DeleteCredentialQuery) error
	GetCredential(query GetCredentialQuery) (*GetCredentialResult, error)
	SearchCredentials(query SearchCredentialQuery) (*SearchCredentialResult, error)
	AddTagsToCredential(query AddTagsToCredentialQuery) error
	RemoveTagsFromCredential(query RemoveTagsFromCredentialQuery) error
}


