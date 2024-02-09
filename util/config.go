package util

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Environment       string `mapstructure:"ENVIRONMENT"`
	PostgresURL       string `mapstructure:"POSTGRES_URL"`
	MigrationPath     string `mapstructure:"MIGRATION_PATH"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

func LoadConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	var config Config
	if _, err := os.Stat(".env"); err != nil {
		environment, found := os.LookupEnv("ENVIRONMENT")
		if !found {
			return config, fmt.Errorf("can't find ENVIRONMENT value")
		}
		config.Environment = environment

		postgresURL, found := os.LookupEnv("POSTGRES_URL")
		if !found {
			return config, fmt.Errorf("can't find POSTGRES_URL value")
		}
		config.PostgresURL = postgresURL

		migrationPath, found := os.LookupEnv("MIGRATION_PATH")
		if !found {
			return config, fmt.Errorf("can't find POSTGRES_URL value")
		}
		config.MigrationPath = migrationPath

		address, found := os.LookupEnv("HTTP_SERVER_ADDRESS")
		if !found {
			return config, fmt.Errorf("can't find HTTP_SERVER_ADDRESS value")
		}
		config.HTTPServerAddress = address

		return config, nil
	}

	if err := viper.ReadInConfig(); err != nil {
		return config, err
	}

	if err := viper.Unmarshal(&config); err != nil {
		return config, err
	}

	return config, nil
}
