package config

import (
	"fmt"
	"pkg/config"
	"time"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port string
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

type JWTConfig struct {
	Secret             string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	JWKSCacheDuration  time.Duration
}

type LogConfig struct {
	Level string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: config.GetEnv("PORT", "8081"),
			Host: config.GetEnv("HOST", "0.0.0.0"),
		},
		Database: DatabaseConfig{
			Host:     config.GetEnv("DB_HOST", "localhost"),
			Port:     config.GetEnvAsInt("DB_PORT", 5432),
			User:     config.GetEnv("DB_USER", "user"),
			Password: config.GetEnv("DB_PASSWORD", "password"),
			Name:     config.GetEnv("DB_NAME", "authdb"),
			SSLMode:  config.GetEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:             config.GetEnv("JWT_SECRET", "your-super-secret-jwt-key"),
			AccessTokenExpiry:  config.GetEnvAsDuration("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: config.GetEnvAsDuration("JWT_REFRESH_TOKEN_EXPIRY", 7*24*time.Hour),
			JWKSCacheDuration:  config.GetEnvAsDuration("JWKS_CACHE_DURATION", 1*time.Hour),
		},
		Log: LogConfig{
			Level: config.GetEnv("LOG_LEVEL", "info"),
		},
	}
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}
