package organizations

type CreateFolderRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	ParentID    *string `json:"parent_id,omitempty"`
	CreatedBy   string  `json:"created_by"`
}

type CreateFolderResponse struct {
	ID          string   `json:"Id"`
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Icon        string   `json:"Icon"`
	CreatedAt   string   `json:"CreatedAt"`
	UpdatedAt   string   `json:"UpdatedAt"`
	ParentID    *string  `json:"ParentID"`
	CreatedBy   string   `json:"CreatedBy"`
	Members     []string `json:"members"`
}

type GetFolderResponse struct {
	ID          string   `json:"Id"`
	Name        string   `json:"Name"`
	Description string   `json:"Description"`
	Icon        string   `json:"Icon"`
	CreatedAt   string   `json:"CreatedAt"`
	UpdatedAt   string   `json:"UpdatedAt"`
	ParentID    *string  `json:"ParentID"`
	CreatedBy   string   `json:"CreatedBy"`
	Members     []string `json:"members"`
}

type GetFoldersResponse struct {
	Folders []GetFolderResponse `json:"folders"`
	Total   int                 `json:"total"`
	Page    int                 `json:"page"`
	Limit   int                 `json:"limit"`
}

type UpdateFolderRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type GetCredentialsResponse struct {
	Credentials []Credential `json:"credentials"`
}

type Credential struct {
	ID           string            `json:"id"`
	CustomFields map[string]string `json:"custom_fields,omitempty"`
	CardCredential
	PasswordCredential
	SSHCredential
	ExpiresAt      string `json:"expires_at"`
	LastReadAt     string `json:"last_read_at"`
	UpdatedAt      string `json:"updated_at"`
	UserIdentifier string `json:"user_identifier"`
	Title          string `json:"title"`
	Note           string `json:"note"`
}

type CardCredential struct {
	CardNumCardCredentialber string `json:"card_number"`
	CVC                      int    `json:"cvc"`
	ExpirationDate           string `json:"expiration_date"`
	OwnerName                string `json:"owner_name"`
}

type PasswordCredential struct {
	CreatedAt  string `json:"created_at"`
	DomainName string `json:"domain_name"`
	Password   string `json:"password"`
}

type SSHCredential struct {
	CreatedAt  string `json:"created_at"`
	Hostname   string `json:"hostname"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
}
