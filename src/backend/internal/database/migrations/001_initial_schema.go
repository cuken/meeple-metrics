package migrations

import (
	"time"
	
	"gorm.io/gorm"
)

func init() {
	RegisterMigration(1, "initial_schema", upInitialSchema, downInitialSchema)
}

func upInitialSchema(db *gorm.DB) error {
	type User struct {
		gorm.Model
		Email     string `gorm:"uniqueIndex;not null"`
		Username  string `gorm:"uniqueIndex;not null"`
		Password  string `gorm:"not null"`
		IsActive  bool   `gorm:"default:true"`
		TenantID  string `gorm:"index;not null"`
	}

	type Game struct {
		gorm.Model
		Name        string `gorm:"not null"`
		Description string
		MinPlayers  int
		MaxPlayers  int
		TenantID    string `gorm:"index;not null"`
	}

	type GameSession struct {
		gorm.Model
		GameID    uint      `gorm:"not null"`
		PlayedAt  time.Time `gorm:"not null"`
		Notes     string
		TenantID  string    `gorm:"index;not null"`
	}

	return db.AutoMigrate(&User{}, &Game{}, &GameSession{})
}

func downInitialSchema(db *gorm.DB) error {
	return db.Migrator().DropTable("users", "games", "game_sessions")
}