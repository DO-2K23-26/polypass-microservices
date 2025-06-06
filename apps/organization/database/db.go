package database

import (
    "fmt"
    "log"
    "os"
    "gorm.io/driver/postgres"
	organization "github.com/DO-2K23-26/polypass-microservices/libs/interfaces/organization"
	"gorm.io/gorm"
)

func main() {
    user := os.Getenv("POSTGRES_USER")
    password := os.Getenv("POSTGRES_PASSWORD")
    dbname := os.Getenv("POSTGRES_DB")

    dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=5432 sslmode=disable",
        user, password, dbname)

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
        log.Fatal("Erreur de migration :", err)
    }

    log.Println("Tables créée ou à jour.")
}
