package folder

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

type CreateFolderQuery struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	ParentID string `json:"parent_id"`
}

type CreateFolderResult struct {
	Folder types.Folder `json:"folder"`
}

type UpdateFolderQuery struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}



type UpdateFolderResult struct {
	Folder types.Folder `json:"folder"`
}

type DeleteFolderQuery struct {
	ID string `json:"id"`
}

type GetFolderQuery struct {
	ID string `json:"id"`
}

type GetFolderResult struct {
	Folder types.Folder `json:"folder"`
}

type SearchFolderQuery struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Limit        *int      `json:"limit"`         // The limit is the maximum number of folders to return
	Offset       *int      `json:"offset"`        // The offset is the number of folders to skip before returning results
	FoldersScope *[]string `json:"folders_scope"` // The folders that the user making the request can access
}

type SearchFolderResult struct {
	Folders []types.Folder `json:"folders"`
	Total   int            `json:"total"` // The total is the total number of folders that match the query
	Limit   int            `json:"limit"`
	Offset  int            `json:"offset"`
}

type GetFolderHierarchyQuery struct {
	ID string `json:"id"`
}

type GetFolderHierarchyResult struct {
	Folders []types.Folder `json:"folders"`
}


