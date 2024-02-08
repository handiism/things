package api

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/handiism/smi/db/sqlc"
	"github.com/invopop/validation"
)

type upsertAbilityRequest struct {
	Name sqlc.AbilityEnum `json:"name"`
}

func (c *upsertAbilityRequest) Validate() error {
	return validation.ValidateStruct(c,
		validation.Field(&c.Name, validation.Required, validation.By(func(value interface{}) error {
			if value, ok := value.(sqlc.AbilityEnum); ok && value.Valid() {
				return nil
			} else {
				return fmt.Errorf("is not valid enum")
			}
		})),
	)
}

func (s *Server) createAbility() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req upsertAbilityRequest

		err := c.BodyParser(&req)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})

		}

		err = req.Validate()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})

		}

		ability, err := s.queries.CreateAbility(c.Context(), req.Name)
		if err != nil {
			return c.Status(http.StatusConflict).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]any{"status": "success", "ability": ability})
	}
}

func (s *Server) getAbilities() fiber.Handler {
	return func(c *fiber.Ctx) error {
		abilities, err := s.queries.GetAbilities(c.Context())
		if err != nil {
			return c.Status(200).JSON(map[string]any{"status": "failed"})

		}

		return c.JSON(map[string]any{"status": "success", "data": map[string]any{"abilities": abilities}})
	}
}

func (s *Server) updateAbility() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON("Parameter id harus berupa string numerik")
		}

		var req upsertAbilityRequest

		err = c.BodyParser(&req)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})

		}

		err = req.Validate()
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(map[string]any{"error": err.Error()})

		}

		ability, err := s.queries.UpdateAbility(c.Context(), sqlc.UpdateAbilityParams{ID: int32(id), Name: req.Name})
		if err != nil {
			return c.Status(http.StatusConflict).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(http.StatusOK).JSON(map[string]any{"status": "success", "ability": ability})
	}
}

func (s *Server) deleteAbility() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(map[string]any{"error": "Parameter id harus berupa string numerik"})
		}

		if err := s.queries.DeleteAbility(c.Context(), int32(id)); err != nil {
			return c.Status(409).JSON(map[string]any{"error": err.Error()})
		}

		return c.Status(200).JSON(map[string]any{"status": "success"})
	}
}
