package search

import (
	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
)

type SearchCredentialQuery struct {
	Title      *string
	TagName    *string
	FolderName *string
}

type SearchCredentialResult struct {
	Credentials []types.Credential
}







