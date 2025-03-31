package folder

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

// Request DTOs
type CreateFolderRequest struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type UpdateFolderRequest struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type GetFolderRequest struct {
    ID     string `json:"id"`
    UserID string `json:"user_id"` // To check access permission
}

type DeleteFolderRequest struct {
    ID     string `json:"id"`
    UserID string `json:"user_id"` // To check access permission
}

type SearchFoldersRequest struct {
    ID       string  `json:"id,omitempty"`
    Name     string  `json:"name,omitempty"`
    Limit    *int    `json:"limit,omitempty"`
    Offset   *int    `json:"offset,omitempty"`
    UserID   string  `json:"user_id"` // Required to get user's folder access scope
}

// Response DTOs
type FolderResponse struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

type SearchFoldersResponse struct {
    Folders []FolderResponse `json:"folders"`
    Total   int              `json:"total"`
    Limit   int              `json:"limit"`
    Offset  int              `json:"offset"`
}

// Conversion functions
func ConvertToFolderResponse(folder types.Folder) FolderResponse {
    return FolderResponse{
        ID:   folder.ID,
        Name: folder.Name,
    }
}

func ConvertToFoldersResponse(folders []types.Folder) []FolderResponse {
    response := make([]FolderResponse, len(folders))
    for i, folder := range folders {
        response[i] = FolderResponse{
            ID:   folder.ID,
            Name: folder.Name,
        }
    }
    return response
}