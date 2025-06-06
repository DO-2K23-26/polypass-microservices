package organization

type FolderCredential struct {
	IdFolder     string `gorm:"type:uuid;primaryKey"`
	IdCredential string `gorm:"type:uuid;primaryKey"`
}
