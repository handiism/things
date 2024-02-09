package util

import (
	"github.com/minio/minio-go"
)

func GetMinioClient() (*minio.Client, error) {
	cfg, err := LoadConfig(".")
	if err != nil {
		return nil, err
	}

	endpoint := cfg.MinioEndpoint
	accessKeyID := cfg.MinioAccessKey
	secretAccessKey := cfg.MinioSecretAccessKey
	useSSL := cfg.MinioSSL

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
