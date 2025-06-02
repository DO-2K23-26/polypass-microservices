package organization

type FolderCredential struct {
	IdFolder     string `json:"id_folder" db:"id_folder" validate:"required,uuid"`
	IdCredential string `json:"id_credential" db:"id_credential" validate:"required,uuid"`
}
