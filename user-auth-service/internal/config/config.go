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
	DBNAME       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBHOST       string
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
		DBPort:       getEnv("DB_PORT", "5432"),
		DBHOST:       getEnv("DBHOST", "localhost"),
		DBUser:       getEnv("DB_USER", "postgres"),
		DBPassword:   getEnv("DB_PASSWORD", "POSTGRES"),
		DBNAME:       getEnv("DB_NAME", "user_auth_db"),
	}
}

// getEnv - вспомогательная функция для получения значений из переменных окружения
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
