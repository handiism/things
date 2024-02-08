package api

import (
	"github.com/gin-gonic/gin"
	"github.com/handiism/smi/db/sqlc"
	"github.com/handiism/smi/util"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	config  util.Config
	queries *sqlc.Queries
	pool    *pgxpool.Pool
	router  *gin.Engine
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
	router := gin.Default()

	ability := router.Group("/ability")

	ability.GET("/", s.getAbilities())

	s.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
