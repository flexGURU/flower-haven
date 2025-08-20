package pkg

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DATABASE_URL            string        `mapstructure:"DATABASE_URL"`
	MIGRATION_PATH          string        `mapstructure:"MIGRATION_PATH"`
	FRONTEND_URL            []string      `mapstructure:"FRONTEND_URL"`
	ENVIRONMENT             string        `mapstructure:"ENVIRONMENT"`
	SERVER_ADDRESS          string        `mapstructure:"SERVER_ADDRESS"`
	PASSWORD_COST           int           `mapstructure:"PASSWORD_COST"`
	PASSWORD_RESET_DURATION time.Duration `mapstructure:"PASSWORD_RESET_DURATION"`
	REFRESH_TOKEN_DURATION  time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TOKEN_DURATION          time.Duration `mapstructure:"TOKEN_DURATION"`
	TOKEN_SYMMETRIC_KEY     string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`
	TOKEN_ISSUER            string        `mapstructure:"TOKEN_ISSUER"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	setDefaults()

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("Config file not found, using environment variables")
		} else {
			return Config{}, Errorf(INTERNAL_ERROR, "failed to read config: %s", err.Error())
		}
	}

	var config Config

	return config, viper.Unmarshal(&config)
}

func setDefaults() {
	viper.SetDefault("DATABASE_URL", "")
	viper.SetDefault("MIGRATION_PATH", "")
	viper.SetDefault("FRONTEND_URL", []string{})
	viper.SetDefault("ENVIRONMENT", "")
	viper.SetDefault("SERVER_ADDRESS", "")
	viper.SetDefault("PASSWORD_COST", 0)
	viper.SetDefault("PASSWORD_RESET_DURATION", 0)
	viper.SetDefault("REFRESH_TOKEN_DURATION", 0)
	viper.SetDefault("TOKEN_DURATION", 0)
	viper.SetDefault("TOKEN_SYMMETRIC_KEY", "")
	viper.SetDefault("TOKEN_ISSUER", "")
}
