package migrations

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	return db
}

func TestMigrations(t *testing.T) {
	db := setupTestDB(t)

	t.Run("RunMigrations", func(t *testing.T) {
		err := RunMigrations(db)
		assert.NoError(t, err)

		// Verify migrations table exists
		var count int64
		err = db.Model(&Migration{}).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// Verify tables were created
		tables := []string{"users", "games", "game_sessions"}
		for _, table := range tables {
			exists := db.Migrator().HasTable(table)
			assert.True(t, exists, "Table %s should exist", table)
		}
	})

	t.Run("RollbackMigration", func(t *testing.T) {
		err := RollbackMigration(db)
		assert.NoError(t, err)

		// Verify tables were dropped
		tables := []string{"users", "games", "game_sessions"}
		for _, table := range tables {
			exists := db.Migrator().HasTable(table)
			assert.False(t, exists, "Table %s should not exist", table)
		}

		// Verify migration record was removed
		var count int64
		err = db.Model(&Migration{}).Count(&count).Error
		assert.NoError(t, err)
		assert.Equal(t, int64(0), count)
	})

	t.Run("GetMigrationStatus", func(t *testing.T) {
		// Apply migration first
		err := RunMigrations(db)
		assert.NoError(t, err)

		// Check status
		migrations, err := GetMigrationStatus(db)
		assert.NoError(t, err)
		assert.Len(t, migrations, 1)
		assert.Equal(t, uint(1), migrations[0].Version)
		assert.Equal(t, "initial_schema", migrations[0].Name)
		assert.WithinDuration(t, time.Now(), migrations[0].AppliedAt, 2*time.Second)
	})
}