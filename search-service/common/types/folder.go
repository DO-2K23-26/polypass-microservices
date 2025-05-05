package types

import (

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Folder struct {
	ID       string         `json:"id"  gorm:"primaryKey"`
	Name     string         `json:"name"`
	ParentID *string `json:"parent_id"`
	Parent   *Folder        `gorm:"foreignKey:ParentID" json:"-"`
	Children []Folder       `gorm:"foreignKey:ParentID" json:"children"`
	User     *[]User        `gorm:"many2many:user_folders;" json:"user"`
}

var EsFolder = map[string]types.Property{
	"id":        types.NewKeywordProperty(),
	"name":      types.NewSearchAsYouTypeProperty(),
	"parent_id": types.NewKeywordProperty(),
}

// The parent id is stringified
type FolderWithStringifiedParentId struct {
	ID       string  `json:"id"  gorm:"primaryKey"`
	Name     string  `json:"name"`
	ParentID *string `json:"parent_id"`
}

var FolderIndex = "folders"

// func (f Folder) withParentIdToString() FolderWithStringifiedParentId {
// 	var parentId *string
// 	if f.ParentID.Valid {
// 		parentId = &f.Parent.ID
// 	}
// 	return FolderWithStringifiedParentId{
// 		ID:       f.ID,
// 		Name:     f.Name,
// 		ParentID: parentId,
// 	}
// }
