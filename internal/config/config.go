package config

import (
    "fmt"
    "os"
    "strconv"

    "github.com/joho/godotenv"
    log "github.com/sirupsen/logrus"
)

type Config struct {
    DBHost     string
    DBPort     string
    DBUser     string
    DBPassword string
    DBName     string
    DBSSLMode  string
    ServerPort string
    LogLevel   string
}

func LoadConfig() (*Config, error) {
    // Загружаем .env файл если существует
    if err := godotenv.Load(); err != nil {
        log.Warn("No .env file found, using environment variables")
    }

    config := &Config{
        DBHost:     getEnv("DB_HOST", "localhost"),
        DBPort:     getEnv("DB_PORT", "5432"),
        DBUser:     getEnv("DB_USER", "postgres"),
        DBPassword: getEnv("DB_PASSWORD", "postgres"),
        DBName:     getEnv("DB_NAME", "subscriptions"),
        DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
        ServerPort: getEnv("SERVER_PORT", "8080"),
        LogLevel:   getEnv("LOG_LEVEL", "info"),
    }

    // Установка уровня логирования
    level, err := log.ParseLevel(config.LogLevel)
    if err != nil {
        log.SetLevel(log.InfoLevel)
    } else {
        log.SetLevel(level)
    }

    log.SetFormatter(&log.JSONFormatter{})
    log.SetOutput(os.Stdout)

    return config, nil
}

func (c *Config) GetDSN() string {
    return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSSLMode)
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
    if value := os.Getenv(key); value != "" {
        if intVal, err := strconv.Atoi(value); err == nil {
            return intVal
        }
    }
    return defaultValue
}