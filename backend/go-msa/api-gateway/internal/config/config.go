package config

import (
	"pkg/config"
)

type Config struct {
	Server   ServerConfig
	Services ServicesConfig
	JWT      JWTConfig
	Log      LogConfig
}

type ServerConfig struct {
	Port int
	Host string
}

type ServicesConfig struct {
	AuthServiceURL     string
	BusinessServiceURL string
}

type JWTConfig struct {
	JWKSURL string
}

type LogConfig struct {
	Level string
}

func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: config.GetEnvAsInt("PORT", 8080),
			Host: config.GetEnv("HOST", "0.0.0.0"),
		},
		Services: ServicesConfig{
			AuthServiceURL:     config.GetEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
			BusinessServiceURL: config.GetEnv("BUSINESS_SERVICE_URL", "http://localhost:8082"),
		},
		JWT: JWTConfig{
			JWKSURL: config.GetEnv("JWKS_URL", "http://localhost:8081/jwks"),
		},
		Log: LogConfig{
			Level: config.GetEnv("LOG_LEVEL", "info"),
		},
	}
}
