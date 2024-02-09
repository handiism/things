package util

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string `mapstructure:"ENVIRONMENT"`
	PostgresURL          string `mapstructure:"POSTGRES_URL"`
	MigrationPath        string `mapstructure:"MIGRATION_PATH"`
	HTTPServerAddress    string `mapstructure:"HTTP_SERVER_ADDRESS"`
	MinioEndpoint        string `mapstructure:"MINIO_ENDPOINT"`
	MinioBucket          string `mapstructure:"MINIO_BUCKET"`
	MinioAccessKey       string `mapstructure:"MINIO_ACCESS_KEY"`
	MinioSecretAccessKey string `mapstructure:"MINIO_SECRET_ACCESS_KEY"`
	MinioSSL             bool   `mapstructure:"MINIO_SSL"`
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

		minioEndpoint, found := os.LookupEnv("MINIO_ENDPOINT")
		if !found {
			return config, fmt.Errorf("can't find MINIO_ENDPOINT value")
		}
		config.MinioEndpoint = minioEndpoint

		minioAccessKey, found := os.LookupEnv("MINIO_ACCESS_KEY")
		if !found {
			return config, fmt.Errorf("can't find MINIO_ACCESS_KEY value")
		}
		config.MinioAccessKey = minioAccessKey

		minioSecretAccessKey, found := os.LookupEnv("MINIO_SECRET_ACCESS_KEY")
		if !found {
			return config, fmt.Errorf("can't find MINIO_ACCESS_KEY value")
		}
		config.MinioSecretAccessKey = minioSecretAccessKey

		minioBucket, found := os.LookupEnv("MINIO_BUCKET")
		if !found {
			return config, fmt.Errorf("can't find MINIO_BUCKET value")
		}
		config.MinioBucket = minioBucket

		minioSSLString, found := os.LookupEnv("MINIO_SSL")
		if !found {
			return config, fmt.Errorf("can't find MINIO_SSL value")
		}
		minioSSL, err := strconv.ParseBool(minioSSLString)
		if err != nil {
			return config, fmt.Errorf("invalid MINIO_SSL value")
		}
		config.MinioSSL = minioSSL

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
