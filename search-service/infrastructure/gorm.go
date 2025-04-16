package infrastructure

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormAdapter struct {
    // Define fields here
    Db *gorm.DB
}

func NewGormAdapter(host, user, password, dbname, port string) (*GormAdapter,error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil,err
	}
	return &GormAdapter{
		Db: db,
	},nil
}

func (a *GormAdapter) CheckHealth() bool {
	err := a.Db.Raw("SELECT 1").Row().Scan(&struct{}{})	
	if err != nil {
		return false
	}
	return true
}
