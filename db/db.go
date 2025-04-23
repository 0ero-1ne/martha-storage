package db

import (
	"github.com/0ero-1ne/martha-storage/config"
	"github.com/0ero-1ne/martha-storage/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewDbConnection(config config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DbName), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.File{})
	return db, err
}
