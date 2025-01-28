package migrations

import (
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

type Migration struct {
	ID        uint      `gorm:"primarykey"`
	Version   uint      `gorm:"uniqueIndex"`
	Name      string    `gorm:"not null"`
	AppliedAt time.Time `gorm:"not null"`
}

type MigrationStep struct {
	Version uint
	Name    string
	Up      func(*gorm.DB) error
	Down    func(*gorm.DB) error
}

var migrations = []MigrationStep{}

func RegisterMigration(version uint, name string, up, down func(*gorm.DB) error) {
	migrations = append(migrations, MigrationStep{
		Version: version,
		Name:    name,
		Up:      up,
		Down:    down,
	})
}

func RunMigrations(db *gorm.DB) error {
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	var appliedMigrations []Migration
	if err := db.Order("version asc").Find(&appliedMigrations).Error; err != nil {
		return fmt.Errorf("failed to fetch applied migrations: %w", err)
	}

	for _, m := range migrations {
		if !isMigrationApplied(appliedMigrations, m.Version) {
			if err := db.Transaction(func(tx *gorm.DB) error {
				if err := m.Up(tx); err != nil {
					return fmt.Errorf("failed to apply migration %d (%s): %w", m.Version, m.Name, err)
				}

				return tx.Create(&Migration{
					Version:   m.Version,
					Name:     m.Name,
					AppliedAt: time.Now(),
				}).Error
			}); err != nil {
				return err
			}
		}
	}

	return nil
}

func RollbackMigration(db *gorm.DB) error {
	var lastMigration Migration
	if err := db.Order("version desc").First(&lastMigration).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("no migrations to rollback")
		}
		return fmt.Errorf("failed to fetch last migration: %w", err)
	}

	for _, m := range migrations {
		if m.Version == lastMigration.Version {
			return db.Transaction(func(tx *gorm.DB) error {
				if err := m.Down(tx); err != nil {
					return fmt.Errorf("failed to rollback migration %d (%s): %w", m.Version, m.Name, err)
				}

				return tx.Delete(&lastMigration).Error
			})
		}
	}

	return fmt.Errorf("migration %d not found", lastMigration.Version)
}

func isMigrationApplied(applied []Migration, version uint) bool {
	for _, m := range applied {
		if m.Version == version {
			return true
		}
	}
	return false
}

func GetMigrationStatus(db *gorm.DB) ([]Migration, error) {
	var migrations []Migration
	err := db.Order("version asc").Find(&migrations).Error
	return migrations, err
}