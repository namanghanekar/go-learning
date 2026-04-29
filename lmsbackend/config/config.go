package config

import "os"

type Config struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string

	RedisAddr string

	JWTSecret string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", "naman"),
		DBName:     getEnv("DB_NAME", "lms"),
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),

		RedisAddr: getEnv("REDIS_ADDR", "localhost:6379"),

		JWTSecret: getEnv("JWT_SECRET", "supersecret"),
	}
}

func getEnv(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
