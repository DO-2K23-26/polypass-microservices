package infrastructure

import (
	"fmt"

	"github.com/DO-2K23-26/polypass-microservices/search-service/common/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormAdapter struct {
	Db *gorm.DB
}

func NewGormAdapter(host, user, password, dbname, port string) (*GormAdapter, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			host, user, password, dbname, port),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &GormAdapter{
		Db: db,
	}, nil
}

func (a *GormAdapter) CheckHealth() bool {
	err := a.Db.Raw("SELECT 1").Row().Scan(new(int))
	if err != nil {
		return false
	}
return true
}

func (a *GormAdapter) Migrate() error {
	return a.Db.AutoMigrate(&types.Folder{})
}
