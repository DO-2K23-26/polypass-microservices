package organization

import "time"

type Tag struct {
	Id       string    `gorm:"column:id;type:uuid;primaryKey" json:"id" validate:"required,uuid"`
	Name      string    `gorm:"column:name;not null" json:"name" validate:"required"`
	Color     string    `gorm:"column:color;not null" json:"color" validate:"required,hexadecimal"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"created_at" validate:"required,datetime"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updated_at" validate:"required,datetime"`
	FolderID  string    `gorm:"column:folder_id;type:uuid;not null" json:"folder_id" validate:"required,uuid"`
	CreatedBy string    `gorm:"column:created_by;type:uuid;not null" json:"created_by" validate:"required,uuid"`

}

type CreateTagRequest struct {
	Name      string `json:"name" validate:"required"`
	Color     string `json:"color" validate:"required,hexadecimal"`
	FolderID  string `json:"folder_id" validate:"required,uuid"`
	CreatedBy string `json:"created_by" validate:"required,uuid"`
}

type UpdateTagRequest struct {
	Id       string `json:"id" validate:"required,uuid"`
	Name     string `json:"name" validate:"required"`
	Color    string `json:"color" validate:"required,hexadecimal"`
	FolderID string `json:"folder_id" validate:"required,uuid"`
}

type GetTagRequest struct {
	Page    int    `json:"page" validate:"required,min=1"`
	Limit  int    `json:"limit" validate:"required,min=10,max=100"`
}
