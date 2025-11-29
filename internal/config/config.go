package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL   string `mapstructure:"DATABASE_URL"`
	RedisURL      string `mapstructure:"REDIS_URL"`
	ClickHouseDSN string `mapstructure:"CLICKHOUSE_DSN"`
	Port          string `mapstructure:"PORT"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/pulse?sslmode=disable")
	viper.SetDefault("REDIS_URL", "redis://localhost:6379")
	viper.SetDefault("CLICKHOUSE_DSN", "clickhouse://default@localhost:9000/default")
	viper.SetDefault("JWT_SECRET", "change-this-secret-in-production")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Printf("Error reading config file: %v", err)
			return nil, err
		}
		// Config file not found; ignore error if desired
		log.Println("No .env file found, using environment variables and defaults")
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
