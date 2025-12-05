package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DatabaseURL          string `mapstructure:"DATABASE_URL"`
	RedisURL             string `mapstructure:"REDIS_URL"`
	ClickHouseDSN        string `mapstructure:"CLICKHOUSE_DSN"`
	Port                 string `mapstructure:"PORT"`
	JWTSecret            string `mapstructure:"JWT_SECRET"`
	APISpecDir           string `mapstructure:"API_SPEC_DIR"`
	PasswordResetTimeout int    `mapstructure:"PASSWORD_RESET_TIMEOUT"` // in seconds, default 3 days
	RegionCode           string `mapstructure:"REGION_CODE"`
}

func LoadConfig() (*Config, error) {
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://user:password@localhost:5432/pulse?sslmode=disable")
	viper.SetDefault("REDIS_URL", "redis://localhost:6379")
	viper.SetDefault("CLICKHOUSE_DSN", "clickhouse://default@localhost:9000/default")
	viper.SetDefault("JWT_SECRET", "change-this-secret-in-production")
	viper.SetDefault("API_SPEC_DIR", "./api-specs")
	viper.SetDefault("PASSWORD_RESET_TIMEOUT", 259200) // 3 days in seconds
	viper.SetDefault("REGION_CODE", "apac")            // default region

	// Use .env file if it exists, otherwise rely on environment variables
	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
