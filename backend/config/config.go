package config

import "os"

type Config struct {
    Port string
    DatabaseURL string
    AdminSecret string
}

func Load() *Config {
    return &Config{
        Port:           getEnv("PORT", "8088"),
        DatabaseURL:    getEnv("DATABASE_URL",""),
        AdminSecret:    getEnv("ADMIN_SECRET", "herenote-admin"),
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}