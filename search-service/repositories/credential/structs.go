package credential

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

type CreateCredentialQuery struct {
	ID       string      `json:"id"`
	Title    string      `json:"title"`
	Tags     []types.Tag `json:"tags"`
	FolderId string      `json:"folder_id"`
}

type CreateCredentialResult struct {
	Credential RepoCredential `json:"credential"`
}

type UpdateCredentialQuery struct {
	ID       *string     `json:"id"`
	Title    *string     `json:"title"`
	Tags     []types.Tag `json:"tags"`
	FolderId *string     `json:"folder_id"`
}

type UpdateCredentialResult struct {
	Credential RepoCredential `json:"credential"`
}

type DeleteCredentialQuery struct {
	ID string `json:"id"`
}

type GetCredentialQuery struct {
	ID string `json:"id"`
}

type GetCredentialResult struct {
	Credential RepoCredential `json:"credential"`
}

type SearchCredentialQuery struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	FolderId     *string   `json:"folder_id"`
	FolderName   *string   `json:"folder_name"`
	TagIds       *[]string `json:"tag_ids"`
	TagName      *string   `json:"tag_name"`
	FoldersScope *[]string `json:"folders_scope"` // The folders that the user making the request can access
	Limit        *int      `json:"limit"`         // The limit is the maximum number of credentials to return
	Offset       *int      `json:"offset"`        // The offset is the number of credentials to skip before returning results
}

type SearchCredentialResult struct {
	Credentials []RepoCredential `json:"credentials"`
	Total       int              `json:"total"` // The total is the total number of credentials that match the query
	Limit       int              `json:"limit"`
	Offset      int              `json:"offset"`
}

type AddTagsToCredentialQuery struct {
	ID     string   `json:"id"`
	TagIds []string `json:"tag_ids"`
}

type AddTagsToCredentialResult struct {
	Credential RepoCredential `json:"credential"`
}

type RemoveTagsFromCredentialQuery struct {
	ID     string   `json:"id"`
	TagIds []string `json:"tag_ids"`
}

type RemoveTagsFromCredentialResult struct {
	Credential RepoCredential `json:"credential"`
}
