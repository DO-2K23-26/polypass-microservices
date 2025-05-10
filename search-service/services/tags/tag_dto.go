package tag

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

// Request DTOs
type CreateTagRequest struct {
	Name string `json:"name"`
}

type UpdateTagRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetTagRequest struct {
	ID string `json:"id"`
}

type DeleteTagRequest struct {
	ID string `json:"id"`
}

type SearchTagsRequest struct {
	SearchQuery string   `json:"search_query,omitempty"`
	FolderIDs   []string `json:"folder_id,omitempty"`
	Limit       *int     `json:"limit,omitempty"`
	Page        *int     `json:"page,omitempty"`
	UserID      string   `json:"user_id"` // Required to get the user's folder access scope
}

// Response DTOs
type TagResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

// Conversion functions
func ConvertToTagResponse(tag *types.Tag) TagResponse {
	return TagResponse{
		ID:   tag.ID,
		Name: tag.Name,
	}
}

func ConvertToTagsResponse(tags []types.Tag) []TagResponse {
	response := make([]TagResponse, len(tags))
	for i, tag := range tags {
		response[i] = TagResponse{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}
	return response
}
