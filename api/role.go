package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/handiism/things/db/sqlc"
	"github.com/invopop/validation"
	"github.com/jackc/pgx/v5"
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

type setAbilitiesRequest struct {
	Abilities []string `json:"abilities"`
}

func (s *setAbilitiesRequest) validate() error {
	return validation.ValidateStruct(s,
		validation.Field(&s.Abilities, validation.Required, validation.Each(validation.By(func(value interface{}) error {
			switch val := value.(type) {
			case string:
				ablity := sqlc.AbilityEnum(val)
				if !ablity.Valid() {
					return fmt.Errorf("is not valid enum")
				}
			default:
				return fmt.Errorf("is not valid enum")
			}
			return nil
		}))),
	)
}

func (s *Server) setAbilities() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		var req setAbilitiesRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		if err := req.validate(); err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		roleAbilities, err := func() ([]sqlc.RoleAbility, error) {
			tx, err := s.pool.BeginTx(c.Context(), pgx.TxOptions{IsoLevel: pgx.ReadUncommitted, AccessMode: pgx.ReadWrite})
			if err != nil {
				return []sqlc.RoleAbility{}, err
			}

			defer tx.Rollback(c.Context())

			wtx := s.queries.WithTx(tx)

			_, err = wtx.DeleteRoleAbilitiesByRoleId(c.Context(), int32(id))
			if err != nil {
				return []sqlc.RoleAbility{}, err
			}

			roleAbilities, err := wtx.SetRoleAbilities(c.Context(), sqlc.SetRoleAbilitiesParams{
				RoleID:        int32(id),
				RoleAbilities: req.Abilities,
			})
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return []sqlc.RoleAbility{}, err
			}

			if err := tx.Commit(c.Context()); err != nil {
				return []sqlc.RoleAbility{}, err
			}

			return roleAbilities, nil
		}()

		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(200).JSON(map[string]any{"status": "success", "roleAbilities": roleAbilities})
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
