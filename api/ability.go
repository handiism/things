package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) getAbilities() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		abilities, err := s.queries.GetAbilities(ctx)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "failed"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"abilities": abilities}})
	}
}
