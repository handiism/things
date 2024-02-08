package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/handiism/smi/db/sqlc"
	"github.com/handiism/smi/util"
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

	s.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Listen(address)
}
