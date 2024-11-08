package config

import (
	"github.com/joho/godotenv" // пакет для работы с .env файлами
	"log"
	"os"
)

type Config struct {
	ServerPort   string
	JWTSecret    string
	RedisAddr    string
	SMTPHost     string
	SMTPPort     string
	SMTPUser     string
	SMTPPassword string
	DBUri        string
	DatabaseName string
}

// LoadConfiguration загружает конфигурацию из файла .env или среды
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. Loading environment variables.")
	}

	return &Config{
		ServerPort:   getEnv("PORT", "8080"),
		JWTSecret:    getEnv("JWT_SECRET", "supersecretkey"),
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", ""),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		DBUri:        getEnv("DB_URI", ""),
		DatabaseName: getEnv("DB_NAME", "user_auth_db"),
	}
}

// getEnv - вспомогательная функция для получения значений из переменных окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
