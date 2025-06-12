package credential

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

// Request DTOs
type CreateCredentialRequest struct {
	CredentialID string    `json:"credentialId" avro:"credentialId"`
	Name         *string   `json:"name" avro:"name"`
	FolderID     *string   `json:"folderId" avro:"folderId"`
	TagIDs       *[]string `json:"tagIds" avro:"tagIds"`
	Timestamp    int64     `json:"timestamp" avro:"timestamp"`
	CreatedAt    string    `json:"created_at" avro:"created_at"`
	UpdatedAt    string    `json:"updated_at" avro:"updated_at"`
}

type CreateCredentialResponse struct {
	Credential types.Credential `json:"credential"`
}

type UpdateCredentialRequest struct {
	CredentialID string    `json:"credentialId" avro:"credentialId"`
	Name         *string   `json:"name" avro:"name"`
	FolderID     *string   `json:"folderId" avro:"folderId"`
	TagIDs       *[]string `json:"tagIds" avro:"tagIds"`
	Timestamp    int64     `json:"timestamp" avro:"timestamp"`
	CreatedAt    string    `json:"created_at" avro:"created_at"`
	UpdatedAt    string    `json:"updated_at" avro:"updated_at"`
}

type GetCredentialRequest struct {
	ID string `json:"id"`
}

type DeleteCredentialRequest struct {
	CredentialID string    `json:"credentialId" avro:"credentialId"`
	Name         *string   `json:"name" avro:"name"`
	FolderID     *string   `json:"folderId" avro:"folderId"`
	TagIDs       *[]string `json:"tagIds" avro:"tagIds"`
	Timestamp    int64     `json:"timestamp" avro:"timestamp"`
	CreatedAt    string    `json:"created_at" avro:"created_at"`
	UpdatedAt    string    `json:"updated_at" avro:"updated_at"`
}

type SearchCredentialsRequest struct {
	SearchQuery string    `json:"search_query,omitempty"`
	FolderID    *string   `json:"folder_id,omitempty"`
	TagIDs      *[]string `json:"tag_ids,omitempty"`
	Limit       *int      `json:"limit,omitempty"`
	Page        *int      `json:"offset,omitempty"`
	UserID      string    `json:"user_id"` // Required to get user's folder access scope
}

// Response DTOs
type GetCredentialResponse struct {
	Credential types.Credential `json:"credential"`
}

type FolderResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SearchCredentialsResponse struct {
	Credentials []types.Credential `json:"credentials"`
	Total       int                `json:"total"`
}

type AddTagsToCredentialRequest struct {
	ID     string   `json:"id"`
	TagIds []string `json:"tag_ids,omitempty"`
}

type RemoveTagsFromCredentialRequest struct {
	ID     string   `json:"id"`
	TagIds []string `json:"tag_ids,omitempty"`
}
