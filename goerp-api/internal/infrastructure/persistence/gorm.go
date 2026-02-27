package persistence

import (
	"goerp-api/internal/infrastructure/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
