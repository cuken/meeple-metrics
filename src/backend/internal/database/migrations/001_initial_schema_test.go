package migrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialSchemaMigration(t *testing.T) {
	db := setupTestDB(t)

	t.Run("Up Migration", func(t *testing.T) {
		err := upInitialSchema(db)
		assert.NoError(t, err)

		// Test User table
		var userCount int64
		err = db.Table("users").Count(&userCount).Error
		assert.NoError(t, err)

		// Test Game table
		var gameCount int64
		err = db.Table("games").Count(&gameCount).Error
		assert.NoError(t, err)

		// Test GameSession table
		var sessionCount int64
		err = db.Table("game_sessions").Count(&sessionCount).Error
		assert.NoError(t, err)
	})

	t.Run("Down Migration", func(t *testing.T) {
		err := downInitialSchema(db)
		assert.NoError(t, err)

		// Verify tables don't exist
		tables := []string{"users", "games", "game_sessions"}
		for _, table := range tables {
			exists := db.Migrator().HasTable(table)
			assert.False(t, exists, "Table %s should not exist", table)
		}
	})
}