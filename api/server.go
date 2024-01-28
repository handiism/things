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
	router  *gin.Engine
}

func NewServer(config util.Config, pool *pgxpool.Pool) (*Server, error) {
	queries := sqlc.New(pool)
	server := &Server{
		config:  config,
		queries: queries,
	}

	server.setupRouter()
	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()
	s.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
