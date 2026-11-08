package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName           string
	AppEnv            string
	GRPCPort          string
	HTTPPort          string
	JWTSecret         string
	StorageBaseURL    string
	StorageUploadDir  string
	DatabaseURL       string
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	MigrationsPath    string
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppName:           getEnv("APP_NAME", "content-service"),
		AppEnv:            getEnv("APP_ENV", "development"),
		GRPCPort:          getEnv("GRPC_PORT", "50054"),
		HTTPPort:          getEnv("HTTP_PORT", "50055"),
		JWTSecret:         getEnv("JWT_SECRET", ""),
		StorageBaseURL:    getEnv("STORAGE_BASE_URL", "http://localhost:50055/storage/uploads"),
		StorageUploadDir:  getEnv("STORAGE_UPLOAD_DIR", "storage/uploads"),
		DatabaseURL:       getEnv("DATABASE_URL", ""),
		DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		DBConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 30)) * time.Minute,
		MigrationsPath:    getEnv("MIGRATIONS_PATH", "db/migrations"),
	}

	if cfg.JWTSecret == "" {
		if cfg.AppEnv != "development" {
			return nil, errors.New("JWT_SECRET must be set when APP_ENV is not development")
		}
		log.Println("[WARN] JWT_SECRET is not set; tokens cannot be validated in development")
	}

	if cfg.DatabaseURL == "" {
		log.Println("[WARN] DATABASE_URL is not set")
	}

	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
