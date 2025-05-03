package types

import (
	"database/sql"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Folder struct {
	ID       string   `json:"id"  gorm:"primaryKey"`
	Name     string   `json:"name"`
	ParentID sql.NullString  `json:"parent_id"`
	Parent   *Folder  `gorm:"foreignKey:ParentID" json:"-"`
	Children []Folder `gorm:"foreignKey:ParentID" json:"children"`
	User     *[]User  `gorm:"many2many:user_folders;" json:"user"`
}

var EsFolder = map[string]types.Property{
	"id":        types.NewKeywordProperty(),
	"name":      types.NewSearchAsYouTypeProperty(),
	"parent_id": types.NewKeywordProperty(),
}

var FolderIndex = "folders"
