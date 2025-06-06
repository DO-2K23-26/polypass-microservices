package main

import (
    "fmt"
    "log"
    "gorm.io/driver/postgres"
	organization "github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"gorm.io/gorm"
)

func main() {
    
    dsn := fmt.Sprintf("host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable")

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Erreur de connexion à la base de données :", err)
    }

    err = db.AutoMigrate(
		&organization.FolderCredential{}, 
		&organization.Folder{},
		&organization.TagCredential{},
		&organization.Tag{})
    if err != nil {
        fmt.Errorf("Failed to auto-migrate models: %w", err)
    }

    log.Println("Database connected and tables migrated.")
}
