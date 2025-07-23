package config

import (
	"fmt"
	"pkg/config"
)

// Config holds all configuration for the business service
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Business BusinessConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type BusinessConfig struct {
	MaxProductsPerPage int
	DefaultCurrency    string
}

type LogConfig struct {
	Level string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: config.GetEnvAsInt("PORT", 8080),
			Host: config.GetEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     config.GetEnv("DB_HOST", "localhost"),
			Port:     config.GetEnvAsInt("DB_PORT", 5432),
			User:     config.GetEnv("DB_USER", "user"),
			Password: config.GetEnv("DB_PASSWORD", "password"),
			Name:     config.GetEnv("DB_NAME", "db"),
			SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
		},
		Business: BusinessConfig{
			MaxProductsPerPage: config.GetEnvAsInt("MAX_PRODUCTS_PER_PAGE", 50),
			DefaultCurrency:    config.GetEnv("DEFAULT_CURRENCY", "KRW"),
		},
		Log: LogConfig{
			Level: config.GetEnv("LOG_LEVEL", "info"),
		},
	}
}

// GetDatabaseURL returns formatted PostgreSQL connection string
func (dCfg *DatabaseConfig) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		dCfg.User,
		dCfg.Password,
		dCfg.Host,
		dCfg.Port,
		dCfg.Name,
		dCfg.SSLMode,
	)
}
