package types

import "github.com/elastic/go-elasticsearch/v9/typedapi/types"

type Tag struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	FolderId  string `json:"folder_id"`
}


var EsTag = map[string]types.Property{
	"id":        types.NewKeywordProperty(),
	"name":      types.NewSearchAsYouTypeProperty(),
	"folder_id": types.NewKeywordProperty(),
}

var TagIndex = "tags"
