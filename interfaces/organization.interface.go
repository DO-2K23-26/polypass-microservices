package interfaces

type Folder struct {
	Id          string   `json:"id" db:"id" validate:"required,uuid"`
	Name        string   `json:"name" db:"name" validate:"required"`
	Description *string  `json:"description" db:"description"`
	Icon        *string  `json:"icon" db:"icon"`
	CreatedAt   string   `json:"created_at" db:"created_at" validate:"required,datetime"`
	UpdatedAt   string   `json:"updated_at" db:"updated_at" validate:"required,datetime"`
	ParentId    *string  `json:"parent_id" db:"parent_id" validate:"uuid"`
	Members     []string `json:"members" db:"members" validate:"required,uuid"`
	CreatedBy   string   `json:"created_by" db:"created_by" validate:"required,uuid"`
}

type Tag struct {
	Id        string `json:"id" db:"id" validate:"required,uuid"`
	Name      string `json:"name" db:"name" validate:"required"`
	Color     string `json:"color" db:"color" validate:"required,hexadecimal"`
	CreatedAt string `json:"created_at" db:"created_at" validate:"required,datetime"`
	UpdatedAt string `json:"updated_at" db:"updated_at" validate:"required,datetime"`
	FolderId  string `json:"folder_id" db:"folder_id" validate:"required,uuid"`
	CreatedBy string `json:"created_by" db:"created_by" validate:"required,uuid"`
}

type TagCredential struct {
	IdTag        string `json:"id_tag" db:"id_tag" validate:"required,uuid"`
	IdCredential string `json:"id_credential" db:"id_credential" validate:"required,uuid"`
}

type FolderCredential struct {
	IdFolder     string `json:"id_folder" db:"id_folder" validate:"required,uuid"`
	IdCredential string `json:"id_credential" db:"id_credential" validate:"required,uuid"`
}
