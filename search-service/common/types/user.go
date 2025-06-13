package types

type User struct {
	ID        string   `json:"id" gorm:"primaryKey;column:id"`
	Folders []Folder `json:"folders" gorm:"many2many:user_folders;"`
}

