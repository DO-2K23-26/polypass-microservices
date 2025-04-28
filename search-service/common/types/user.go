package types

import "github.com/elastic/go-elasticsearch/v9/typedapi/types"

type User struct {
	ID        string   `json:"id" gorm:"primaryKey;column:id"`
	Folders []Folder `json:"folders" gorm:"many2many:user_folders;"`
}

