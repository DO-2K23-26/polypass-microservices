package organization

import "time"

type Folder struct {
	Id          string `gorm:"primaryKey;type:uuid"`
	Name        string `gorm:"not null"`
	Description *string
	Icon        *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	ParentID    *string
	CreatedBy   string
	User        *[]User `gorm:"many2many:user_folders;" json:"user"`
}

type CreateFolderRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	ParentID    *string `json:"parent_id"`
	CreatedBy   string  `json:"created_by" binding:"required"`
}

type UpdateFolderRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Icon        *string `json:"icon"`
	ParentID    *string `json:"parent_id"`
}

type GetFolderRequest struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

type GetFolderResponse struct {
	Folders []Folder `json:"folders"`
	Total   int      `json:"total"`
	Page    int      `json:"page"`
	Limit   int      `json:"limit"`
}

type GetUsersInFolderRequest struct {
	Page  int `json:"page" form:"page"`
	Limit int `json:"limit" form:"limit"`
}

type GetUsersInFolderResponse struct {
	Users []string `json:"users"`
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
}
