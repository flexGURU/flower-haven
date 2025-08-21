package handlers

import (
	"net/http"

	"github.com/flexGURU/flower-haven/backend/pkg"
	"github.com/gin-gonic/gin"
)

func (s *Server) getDashboardDataHandler(ctx *gin.Context) {
	data, err := s.repo.ProductRepository.GetDashboardData(ctx)
	if err != nil {
		ctx.JSON(pkg.ErrorToStatusCode(err), errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data})
}
