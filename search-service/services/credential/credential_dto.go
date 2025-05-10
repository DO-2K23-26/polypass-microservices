package credential

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

// Request DTOs
type CreateCredentialRequest struct {
	ID       string `json:"id,omitempty"` // Optional, will be generated if not provided
	Title    string `json:"title"`
	FolderID string `json:"folder_id"`
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
type CredentialResponse struct {
	ID       string          `json:"id"`
	Title    string          `json:"title"`
	FolderID string          `json:"folder_id"`
	Tags     []types.Tag     `json:"tags,omitempty"`
	Folder   *FolderResponse `json:"folder,omitempty"`
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

// Conversion functions
func ConvertToTagResponse(tag types.Tag) types.Tag {
	return types.Tag{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func ConvertToTagResponses(tags []types.Tag) []types.Tag {
	if tags == nil {
		return nil
	}

	response := make([]types.Tag, len(tags))
	for i, tag := range tags {
		response[i] = ConvertToTagResponse(tag)
	}
	return response
}

func ConvertToFolderResponse(folder *types.Folder) *FolderResponse {
	if folder == nil {
		return nil
	}

	return &FolderResponse{
		ID:   folder.ID,
		Name: folder.Name,
	}
}

func ConvertToCredentialResponse(credential types.Credential) CredentialResponse {
	return CredentialResponse{
		ID:     credential.ID,
		Title:  credential.Title,
		Tags:   ConvertToTagResponses(credential.Tags),
		Folder: ConvertToFolderResponse(credential.Folder),
	}
}

func ConvertToCredentialResponses(credentials []types.Credential) []CredentialResponse {
	response := make([]CredentialResponse, len(credentials))
	for i, credential := range credentials {
		response[i] = ConvertToCredentialResponse(credential)
	}
	return response
}
