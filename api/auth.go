package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/handiism/smi/db/sqlc"
	"github.com/invopop/validation"
	"github.com/invopop/validation/is"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int32  `json:"roleId"`
}

func (r *registerRequest) Validate() error {
	return validation.ValidateStruct(r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(8, 0)),
		validation.Field(&r.RoleID, validation.Required),
	)
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *loginRequest) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required, validation.Length(8, 0)),
	)
}

func (s *Server) register() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req registerRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
		}

		if err, ok := req.Validate().(validation.Errors); ok && err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
		}

		safePassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
		}

		credential, err := s.queries.CreateCredential(c.Context(), sqlc.CreateCredentialParams{
			Email:    req.Email,
			Password: string(safePassword),
			RoleID:   req.RoleID,
		})

		if err != nil {
			return c.Status(http.StatusConflict).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]any{"status": "success", "credential": credential})
	}
}

type JwtClaims struct {
	jwt.RegisteredClaims
	CredentialID uuid.UUID `json:"credentialId"`
}

func (s *Server) login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req loginRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
		}

		if err := req.Validate(); err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})
		}

		credential, err := s.queries.GetCredentialByEmail(c.Context(), req.Email)
		if err != nil {
			return c.Status(http.StatusNotFound).JSON(map[string]any{"error": err.Error()})
		}

		if err := bcrypt.CompareHashAndPassword([]byte(credential.Password), []byte(req.Password)); err != nil {
			return c.Status(http.StatusConflict).JSON(map[string]any{"error": err.Error()})
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, JwtClaims{
			CredentialID: credential.ID,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: &jwt.NumericDate{
					Time: time.Now().Add(2 * time.Hour),
				},
			},
		})

		tokenString, err := token.SignedString([]byte("handiism"))
		if err != nil {
			return c.Status(http.StatusConflict).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(200).JSON(map[string]any{"status": "success", "credential": credential, "token": tokenString})
	}
}

func (s *Server) me() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get("CredentialID", "")
		if id == "" {
			return c.Status(404).JSON(map[string]any{"error": "id not found"})
		}

		credential, err := s.queries.GetCredentialById(c.Context(), uuid.FromStringOrNil(id))
		if err != nil {
			return c.Status(404).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(200).JSON(map[string]any{"status": "success", "credential": credential})
	}
}

func (s *Server) verify() fiber.Handler {
	return func(c *fiber.Ctx) error {
		headers := c.GetReqHeaders()
		auths := headers["Authorization"]
		if len(auths) == 0 {
			return c.Status(400).JSON(map[string]any{"error": "no auth"})
		}

		splitedAuth := strings.Split(auths[0], " ")
		if splitedAuth[0] != "Bearer" && len(splitedAuth[1]) <= 0 {
			return c.Status(400).JSON(map[string]any{"error": "false auth method"})
		}

		token, err := jwt.ParseWithClaims(splitedAuth[1], &JwtClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte("handiism"), nil
		})
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}
		claims, ok := token.Claims.(*JwtClaims)
		if !ok {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		c.Locals("CredentialID", claims.CredentialID.String())
		return c.Next()
	}
}
