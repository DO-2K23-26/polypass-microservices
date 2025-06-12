package organization

type TagCredential struct {
	IdTag        string `gorm:"type:uuid;primaryKey"`
	IdCredential string `gorm:"type:uuid;primaryKey"`
}
