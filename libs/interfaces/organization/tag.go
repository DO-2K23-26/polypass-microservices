package organization

type Tag struct {
	Id        string `json:"id" db:"id" validate:"required,uuid"`
	Name      string `json:"name" db:"name" validate:"required"`
	Color     string `json:"color" db:"color" validate:"required,hexadecimal"`
	CreatedAt string `json:"created_at" db:"created_at" validate:"required,datetime"`
	UpdatedAt string `json:"updated_at" db:"updated_at" validate:"required,datetime"`
	FolderId  string `json:"folder_id" db:"folder_id" validate:"required,uuid"`
	CreatedBy string `json:"created_by" db:"created_by" validate:"required,uuid"`
}
