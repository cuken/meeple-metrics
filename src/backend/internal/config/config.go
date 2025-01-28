package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Server ServerConfig
	DB     DBConfig
}

type ServerConfig struct {
	Port string
}

type DBConfig struct {
	Driver     string
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	SQLitePath string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		DB: DBConfig{
			Driver:     getEnv("DB_DRIVER", "sqlite"),
			Host:       getEnv("DB_HOST", "localhost"),
			Port:       getEnv("DB_PORT", "5432"),
			User:       getEnv("DB_USER", ""),
			Password:   getEnv("DB_PASSWORD", ""),
			DBName:     getEnv("DB_NAME", "meeple_metrics"),
			SQLitePath: getEnv("DB_SQLITE_PATH", "meeple_metrics.db"),
		},
	}

	return config, nil
}

func (c *Config) GetDSN() string {
	switch c.DB.Driver {
	case "postgres":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			c.DB.Host, c.DB.Port, c.DB.User, c.DB.Password, c.DB.DBName)
	case "sqlite":
		return c.DB.SQLitePath
	default:
		return c.DB.SQLitePath
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}