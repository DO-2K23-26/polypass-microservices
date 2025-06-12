package tag

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

// Request DTOs
type CreateTagRequest struct {
	ID        string `json:"id" avro:"id"`
	Name      string `json:"name" avro:"name"`
	Color     string `json:"color" avro:"color"`
	CreatedAt string `json:"created_at" avro:"created_at"`
	UpdatedAt string `json:"updated_at" avro:"updated_at"`
	FolderID  string `json:"folder_id" avro:"folder_id"`
	CreatedBy string `json:"created_by" avro:"created_by"`
}

type UpdateTagRequest struct {
	ID        string `json:"id" avro:"id"`
	Name      string `json:"name" avro:"name"`
	Color     string `json:"color" avro:"color"`
	CreatedAt string `json:"created_at" avro:"created_at"`
	UpdatedAt string `json:"updated_at" avro:"updated_at"`
	FolderID  string `json:"folder_id" avro:"folder_id"`
	CreatedBy string `json:"created_by" avro:"created_by"`
}

type UpdateTagResponse struct {
	Tag types.Tag `json:"tag"`
}

type GetTagRequest struct {
	ID string `json:"id"`
}

type DeleteTagRequest struct {
	ID        string `json:"id" avro:"id"`
	Name      string `json:"name" avro:"name"`
	Color     string `json:"color" avro:"color"`
	CreatedAt string `json:"created_at" avro:"created_at"`
	UpdatedAt string `json:"updated_at" avro:"updated_at"`
	FolderID  string `json:"folder_id" avro:"folder_id"`
	CreatedBy string `json:"created_by" avro:"created_by"`
}

type SearchTagsRequest struct {
	SearchQuery string `json:"search_query,omitempty"`
	Limit       *int   `json:"limit,omitempty"`
	Page        *int   `json:"page,omitempty"`
	UserID      string `json:"user_id"` // Required to get the user's folder access scope
}

// Response DTOs

type CreateTagResponse struct {
	Tag types.Tag `json:"tag"`
}

type GetTagResponse struct {
	Tag types.Tag `json:"tag"`
}

type SearchTagsResponse struct {
	Tags   []types.Tag `json:"tags"`
	Limit  int         `json:"limit"`
	Offset int         `json:"page"`
	Total  int         `json:"total"`
}

type MGetTagRequest struct {
	IDs []string `json:"ids"`
}

type MGetTagResponse struct {
	Tags []types.Tag `json:"tags"`
}
