package config

import (
	"os"
	"strconv"
)

type AppConfig struct {
	Env  string
	Port int
}

type DbConfig struct {
	DbHost     string
	DbName     string
	DbPassword string
	DbPort     int
	DbSchema   string
	DbUsername string
}

type JwtConfig struct {
	RefreshExpirationMinutes int
	RefreshKey               string
	SecretExpirationMinutes  int
	SecretKey                string
}

type Config struct {
	App *AppConfig
	Db  *DbConfig
	Jwt *JwtConfig
}

func LoadConfig() (*Config, error) {
	return &Config{
		App: &AppConfig{
			Env:  os.Getenv("APP_ENV"),
			Port: getIntEnv("PORT", 8080),
		},
		Db: &DbConfig{
			DbHost:     os.Getenv("DB_HOST"),
			DbName:     os.Getenv("DB_NAME"),
			DbPassword: os.Getenv("DB_PASSWORD"),
			DbPort:     getIntEnv("DB_PORT", 5432),
			DbSchema:   os.Getenv("DB_SCHEMA"),
			DbUsername: os.Getenv("DB_USERNAME"),
		},
		Jwt: &JwtConfig{
			RefreshExpirationMinutes: getIntEnv(
				"JWT_REFRESH_EXPIRATION_MINUTES",
				60*24*7,
			),
			RefreshKey: os.Getenv("JWT_REFRESH_KEY"),
			SecretExpirationMinutes: getIntEnv(
				"JWT_SECRET_EXPIRATION_MINUTES",
				60*24,
			),
			SecretKey: os.Getenv("JWT_SECRET_KEY"),
		},
	}, nil
}

func getIntEnv(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil {
		return defaultValue
	}
	return value
}
