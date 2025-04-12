package types

import "github.com/elastic/go-elasticsearch/v8/typedapi/types"

type User struct {
	ID    string `json:"id"`
	FolderIds []string `json:"folder_ids"`
}

var EsUser = map[string]types.Property{
	"id":        types.NewKeywordProperty(),
	"folder_ids": types.NewKeywordProperty(),
}

var UserIndex = "users"
