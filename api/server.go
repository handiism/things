package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/handiism/things/db/sqlc"
	"github.com/handiism/things/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config  util.Config
	queries *sqlc.Queries
	pool    *pgxpool.Pool
	router  *fiber.App
}

func NewServer(config util.Config, pool *pgxpool.Pool) (*Server, error) {
	queries := sqlc.New(pool)
	server := &Server{
		config:  config,
		pool:    pool,
		queries: queries,
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := fiber.New()

	ability := router.Group("/ability")
	ability.Post("/", s.createAbility())
	ability.Get("/", s.getAbilities())
	ability.Put("/:id", s.updateAbility())
	ability.Delete("/:id", s.deleteAbility())

	role := router.Group("/role")
	role.Post("/", s.createRole())
	role.Get("/", s.getRoles())
	role.Put("/:id", s.updateRole())
	role.Put("/:id/abilities", s.setAbilities())
	role.Delete("/:id", s.deleteRole())

	auth := router.Group("/auth")
	auth.Post("/login", s.login())
	auth.Post("/register", s.register())
	auth.Get("/me", s.verify(), s.me())

	credential := router.Group("/credential")

	credential.Put("/picture", s.verify(), s.setPictureCredential())
	credential.Put("/profile", s.verify(), s.updateProfileCredential())

	s.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Listen(address)
}
