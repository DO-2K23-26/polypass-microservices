package organization

type FolderCredential struct {
	IdFolder     string `gorm:"type:uuid;primaryKey"`
	IdCredential string `gorm:"type:uuid;primaryKey"`
}

// CredentialList represents the list response for credentials.
type CredentialList struct {
	Credentials []map[string]interface{} `json:"credentials"`
	Page        int                      `json:"page"`
	Limit       int                      `json:"limit"`
}
