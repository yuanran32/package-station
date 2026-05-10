package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port             string
	MySQLDSN         string
	JWTSecret        string
	TokenExpireHours int
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		MySQLDSN:         getEnv("MYSQL_DSN", "root:123456@tcp(127.0.0.1:3306)/parcel_station?charset=utf8mb4&parseTime=True&loc=Local"),
		JWTSecret:        getEnv("JWT_SECRET", "change-this-secret"),
		TokenExpireHours: getEnvInt("TOKEN_EXPIRE_HOURS", 72),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}
