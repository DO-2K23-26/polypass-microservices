package types

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

type Folder struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var EsFolder = map[string]types.Property{
	"id":   types.NewKeywordProperty(),
	"name": types.NewTextProperty(),
}

var FolderIndex = "folders"
