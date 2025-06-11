package folder

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

// Request DTOs
type CreateFolderRequest struct {
	ID          string   `json:"id" avro:"id"`
	Name        string   `json:"name" avro:"name"`
	Description string   `json:"description" avro:"description"`
	Icon        string   `json:"icon" avro:"icon"`
	ParentID    string   `json:"parent_id" avro:"parent_id"`
	Members     []string `json:"members" avro:"members"`
	CreatedBy   string   `json:"created_by" avro:"created_by"`
	CreatedAt   string   `json:"created_at" avro:"created_at"`
	UpdatedAt   string   `json:"updated_at" avro:"updated_at"`
}

type UpdateFolderRequest struct {
	ID          string   `json:"id" avro:"id"`
	Name        string   `json:"name" avro:"name"`
	Description string   `json:"description" avro:"description"`
	Icon        string   `json:"icon" avro:"icon"`
	ParentID    string   `json:"parent_id" avro:"parent_id"`
	Members     []string `json:"members" avro:"members"`
	CreatedAt   string   `json:"created_at" avro:"created_at"`
	UpdatedAt   string   `json:"updated_at" avro:"updated_at"`
}

type GetFolderRequest struct {
	ID string `json:"id" avro:"id"`
}

type DeleteFolderRequest struct {
	ID          string   `json:"id" avro:"id"`
	Name        string   `json:"name" avro:"name"`
	Description string   `json:"description" avro:"description"`
	Icon        string   `json:"icon" avro:"icon"`
	ParentID    string   `json:"parent_id" avro:"parent_id"`
	Members     []string `json:"members" avro:"members"`
	CreatedAt   string   `json:"created_at" avro:"created_at"`
	UpdatedAt   string   `json:"updated_at" avro:"updated_at"`
}

type SearchFoldersRequest struct {
	SearchQuery string `json:"search_query,omitempty" avro:"search_query"`
	Limit       *int   `json:"limit,omitempty" avro:"limit"`
	Page        *int   `json:"offset,omitempty" avro:"offset"`
	UserID      string `json:"user_id" avro:"user_id"` // Required to get user's folder access scope
}

// Response DTOs
type FolderResponse struct {
	Folder types.Folder `json:"folder" avro:"folder"`
}

type SearchFoldersResponse struct {
	Folders []types.Folder `json:"folders" avro:"folders"`
	Total   int            `json:"total" avro:"total"`
}

type GetUserFoldersRequest struct {
	UserID string `json:"user_id" avro:"user_id"`
}

type GetUserFoldersResponse struct {
	Folders []types.Folder `json:"folders" avro:"folders"`
}
