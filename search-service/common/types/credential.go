package types

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

type Credential struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Tags   []Tag   `json:"tags"`
	Folder *Folder `json:"folder"`
}

var EsCredential = map[string]types.Property{
	"id":    types.NewKeywordProperty(),
	"title": types.NewSearchAsYouTypeProperty(),
	"tags": &types.NestedProperty{
		Properties: EsTag,
	},
	"folder": &types.ObjectProperty{
		Properties: EsFolder,
	},
}

var CredentialIndex = "credentials"