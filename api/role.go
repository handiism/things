package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/handiism/smi/db/sqlc"
	"github.com/invopop/validation"
)

type upsertRoleRequest struct {
	Name sqlc.RoleEnum `json:"name"`
}

func (c *upsertRoleRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.By(func(value interface{}) error {
			if value, ok := value.(sqlc.RoleEnum); ok && value.Valid() {
				return nil
			} else {
				return fmt.Errorf("is not valid enum")
			}
		})),
	)
}

func (s *Server) createRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req upsertRoleRequest

		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})

		}

		err = req.Validate()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})

		}

		role, err := s.queries.CreateRole(c.Context(), req.Name)
		if err != nil {
			return c.Status(http.StatusConflict).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]any{"status": "success", "role": role})
	}
}

func (s *Server) getRoles() fiber.Handler {
	return func(c *fiber.Ctx) error {
		roles, err := s.queries.GetRoles(c.Context())
		if err != nil {
			return c.Status(200).JSON(map[string]any{"status": "failed"})

		}

		return c.JSON(map[string]any{"status": "success", "data": map[string]any{"roles": roles}})
	}
}

func (s *Server) updateRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON("Parameter id harus berupa string numerik")
		}

		var req upsertRoleRequest

		err = c.BodyParser(&req)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})

		}

		err = req.Validate()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})

		}

		role, err := s.queries.UpdateRole(c.Context(), sqlc.UpdateRoleParams{
			ID:   int32(id),
			Name: req.Name,
		})
		if err != nil {
			return c.Status(http.StatusConflict).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]any{"status": "success", "role": role})
	}
}

func (s *Server) deleteRole() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": "Parameter id harus berupa string numerik"})
		}

		if err := s.queries.DeleteRole(c.Context(), int32(id)); err != nil {
			return c.Status(409).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(200).JSON(map[string]any{"status": "success"})
	}
}
