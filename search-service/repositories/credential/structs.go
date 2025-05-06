package credential

import "github.com/DO-2K23-26/polypass-microservices/search-service/common/types"

type CreateCredentialQuery struct {
	ID     string        `json:"id"`
	Title  string        `json:"title"`
	Tags   []types.Tag   `json:"tags"`
	Folder *types.Folder `json:"folder"`
}

type CreateCredentialResult struct {
	Credential types.Credential `json:"credential"`
}

type UpdateCredentialQuery struct {
	ID     string        `json:"id"`
	Title  *string       `json:"title"`
	Folder *types.Folder `json:"folder"`
}

type UpdateCredentialResult struct {
	Credential types.Credential `json:"credential"`
}

type DeleteCredentialQuery struct {
	ID string `json:"id"`
}

type GetCredentialQuery struct {
	ID string `json:"id"`
}

type GetCredentialResult struct {
	Credential types.Credential `json:"credential"`
}

type SearchCredentialQuery struct {
	SearchQuery  string    `json:"search_query"`
	TagIds       []string `json:"tag_ids"`
	FoldersScope []string `json:"folders_scope"` // The name of the folders that the user making the request can access
	Limit        *int      `json:"limit"`         // The limit is the maximum number of credentials to return
	Offset       *int      `json:"offset"`        // The offset is the number of credentials to skip before returning results
}

type SearchCredentialResult struct {
	Credentials []types.Credential `json:"credentials"`
	Total       int                `json:"total"` // The total is the total number of credentials that match the query
	Limit       int                `json:"limit"`
	Offset      int                `json:"offset"`
}

type AddTagsToCredentialQuery struct {
	ID  string      `json:"id"`
	Tags []types.Tag `json:"tag"`
}

type AddTagsToCredentialResult struct {
	Credential types.Credential `json:"credential"`
}

type RemoveTagsFromCredentialQuery struct {
	ID     string   `json:"id"`
	TagIds []string `json:"tag_ids"`
}

type RemoveTagsFromCredentialResult struct {
	Credential types.Credential `json:"credential"`
}
