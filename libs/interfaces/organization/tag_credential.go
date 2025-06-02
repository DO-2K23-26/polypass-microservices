package organization

type TagCredential struct {
	IdTag        string `json:"id_tag" db:"id_tag" validate:"required,uuid"`
	IdCredential string `json:"id_credential" db:"id_credential" validate:"required,uuid"`
}
