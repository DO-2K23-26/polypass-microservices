package credential

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

// Request DTOs
type CreateCredentialRequest struct {
	ID       string `json:"id,omitempty"` // Optional, will be generated if not provided
	Title    string `json:"title"`
	FolderID string `json:"folder_id"`
}

type CreateCredentialResponse struct {
	Credential types.Credential `json:"credential"`
}

type UpdateCredentialRequest struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	FolderID string   `json:"folder_id"`
	TagIDs   []string `json:"tag_ids,omitempty"` // IDs of tags to associate with the credential
}

type GetCredentialRequest struct {
	ID string `json:"id"`
}

type DeleteCredentialRequest struct {
	ID string `json:"id"`
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

