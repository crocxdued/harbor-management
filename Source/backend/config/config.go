package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	UseDB          bool
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	JWTSecret      string
	JWTExpiryHours int
	ServerPort     string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env не найден, используются переменные окружения")
	}
	useDB, _ := strconv.ParseBool(getEnv("USE_DB", "false"))
	expiry, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))
	return &Config{
		UseDB:          useDB,
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "harbor_db"),
		JWTSecret:      getEnv("JWT_SECRET", "secret"),
		JWTExpiryHours: expiry,
		ServerPort:     getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
