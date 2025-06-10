package main

import (
	"fmt"
	"log"
	"os"

	organization "github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_PORT"))

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
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	log.Println("Database connected and tables migrated.")
}
