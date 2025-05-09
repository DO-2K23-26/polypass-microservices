package folder

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

// Request DTOs
type CreateFolderRequest struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

type UpdateFolderRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetFolderRequest struct {
	ID string `json:"id"`
}

type DeleteFolderRequest struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"` // To check access permission
}

type SearchFoldersRequest struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Limit  *int   `json:"limit,omitempty"`
	Offset *int   `json:"offset,omitempty"`
	UserID string `json:"user_id"` // Required to get user's folder access scope
}

// Response DTOs
type FolderResponse struct {
	Folder types.Folder `json:"folder"`
}

type SearchFoldersResponse struct {
	Folders []types.Folder `json:"folders"`
	Total   int            `json:"total"`
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
}
