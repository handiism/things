package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/gosimple/slug"
	"github.com/handiism/things/db/sqlc"
	"github.com/handiism/things/util"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/minio/minio-go"
)

func (s *Server) setPictureCredential() fiber.Handler {
	return func(c *fiber.Ctx) error {
		credentialID := c.Locals("CredentialID").(string)

		header, err := c.FormFile("picture")
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		splited := strings.Split(header.Filename, ".")

		path := fmt.Sprintf("%v-%v.%v", time.Now().Unix(), slug.Make(splited[0]), splited[1])

		min, err := util.GetMinioClient()
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		file, err := header.Open()
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		if _, err = min.PutObject(s.config.MinioBucket, path, file, header.Size, minio.PutObjectOptions{}); err != nil {
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

type updateProfileRequest struct {
	Name     pgtype.Text `json:"name"`
	Email    pgtype.Text `json:"email"`
	Username pgtype.Text `json:"username"`
	RoleID   pgtype.Int4 `json:"roleId"`
}

func (u *updateProfileRequest) validate() error {
	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.By(func(value interface{}) error {
			if value, ok := value.(pgtype.Text); ok && value.Valid {
				return validation.Validate(value, is.Email)
			}

			return nil
		})),
	)
}

func (s *Server) updateProfileCredential() fiber.Handler {
	return func(c *fiber.Ctx) error {
		credentialID := c.Locals("CredentialID").(string)
		var req updateProfileRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		if err := req.validate(); err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		params := sqlc.UpdateCredentialParams{
			ID:       uuid.FromStringOrNil(credentialID),
			Name:     req.Name,
			Email:    req.Email,
			Username: req.Username,
			RoleID:   req.RoleID,
		}

		credential, err := s.queries.UpdateCredential(c.Context(), params)
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(200).JSON(map[string]any{"status": "success", "credential": credential})
	}
}
