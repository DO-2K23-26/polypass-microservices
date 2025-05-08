package tags

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

type GetTagQuery struct {
	ID string `json:"id"`
}

type GetTagResult struct {
	Tag *types.Tag `json:"tag"`
}

type CreateTagQuery struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	FolderID string `json:"folder_id"`
}

type CreateTagResult struct {
	Tag *types.Tag `json:"tag"`
}

type UpdateTagQuery struct {
	ID       string `json:"id"`
	Name     *string `json:"name"`
	FolderId *string
}

type UpdateTagResult struct {
	Tag *types.Tag `json:"tag"`
}

type DeleteTagQuery struct {
	ID string `json:"id"`
}

type SearchTagQuery struct {
	Name         string    `json:"name,omitempty"`
	FolderId     *string   `json:"folder_id,omitempty"`
	Limit        *int      `json:"limit,omitempty"`         // Limit is the maximum number of results to return.
	Offset       *int      `json:"offset,omitempty"`        // Offset is the number of results to skip.
	FoldersScope *[]string `json:"folders_scope,omitempty"` // The folders that the user making the request can access
}

type SearchTagResult struct {
	Tags   []types.Tag `json:"tags"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
	Total  int         `json:"total"`
}
