package database

import (
	"github.com/cuken/meeple-metrics/internal/config"
	"github.com/cuken/meeple-metrics/internal/database/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewConnection(cfg *config.Config) (*gorm.DB, error) {
	var dialector gorm.Dialector

	switch cfg.DB.Driver {
	case "postgres":
		dialector = postgres.Open(cfg.GetDSN())
	default:
		dialector = sqlite.Open(cfg.GetDSN())
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := migrations.RunMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}