package types

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

type Credential struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	FolderID  string   `json:"folder_id"`
	TagIDs    []string `json:"tag_ids"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

var EsCredential = map[string]types.Property{
	"id":         types.NewKeywordProperty(),
	"name":       types.NewSearchAsYouTypeProperty(),
	"folder_id":  types.NewKeywordProperty(),
	"tag_ids":    types.NewKeywordProperty(),
	"created_at": types.NewDateProperty(),
	"updated_at": types.NewDateProperty(),
}

var CredentialIndex = "credentials"
