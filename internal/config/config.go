// filepath: internal/config/config.go
package config

import (
    "fmt"
    "github.com/spf13/viper"
)

type Config struct {
    Port        string
    DatabaseURL string
    APIToken    string
}

func Load() (*Config, error) {
    // 1) Carga .env (solo en dev)
    viper.SetConfigFile(".env")
    _ = viper.ReadInConfig()

    // 2) Env vars automáticas y defaults
    viper.AutomaticEnv()
    viper.SetDefault("PORT", "3000")
    viper.SetDefault("DATABASE_URL",
        "host=timescaledb user=goapp password=secret dbname=sensors port=5432 sslmode=disable TimeZone=UTC",
    )
    // API_TOKEN no tiene default

    cfg := &Config{
        Port:        viper.GetString("PORT"),
        DatabaseURL: viper.GetString("DATABASE_URL"),
        APIToken:    viper.GetString("API_TOKEN"),
    }

    if cfg.DatabaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL no puede estar vacío")
    }
    return cfg, nil
}