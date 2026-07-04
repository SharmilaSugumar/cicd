package config

import (
	"os"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	JWTSecret    string
	AppEnv       string
	SchedulerInt string
}

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		DBHost:       getEnv("DB_HOST", "localhost"),
		DBPort:       getEnv("DB_PORT", "5432"),
		DBUser:       getEnv("DB_USER", "forgeflow"),
		DBPassword:   getEnv("DB_PASSWORD", "forgeflow"),
		DBName:       getEnv("DB_NAME", "forgeflow"),
		JWTSecret:    getEnv("JWT_SECRET", "super_secret_key"),
		AppEnv:       getEnv("APP_ENV", "development"),
		SchedulerInt: getEnv("SCHEDULER_INTERVAL", "5s"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
