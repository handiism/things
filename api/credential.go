package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/handiism/smi/db/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go"
)

func (s *Server) setPicture() fiber.Handler {
	return func(c *fiber.Ctx) error {
		credentialID := c.Locals("CredentialID").(string)

		header, err := c.FormFile("picture")
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		s3, err := getMinioClient()
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		file, err := header.Open()
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		path := "dfdl0.png"
		_, err = s3.PutObject("smi", path, file, header.Size, minio.PutObjectOptions{})
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		credential, err := s.queries.SetPictureByCredentialId(c.Context(), sqlc.SetPictureByCredentialIdParams{
			ID:      uuid.FromStringOrNil(credentialID),
			Picture: pgtype.Text{Valid: true, String: path},
		})
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(200).JSON(map[string]any{"status": "success", "credential": credential})
	}
}

func getMinioClient() (*minio.Client, error) {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "TrgIxPibXMUlEjNSK7VK"
	secretAccessKey := "7PGnOTKlnyPWF72aZR9MxO8T28tRbxdtsmIPpIBW"
	useSSL := false

	minioClient, err := minio.New(endpoint, accessKeyID, secretAccessKey, useSSL)
	if err != nil {
		return nil, err
	}

	return minioClient, nil
}
