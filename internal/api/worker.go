package api

import (
	"net/http"

	"github.com/alexanderbkl/vidre-back/internal/db"
	"github.com/alexanderbkl/vidre-back/internal/entity"
	"github.com/gin-gonic/gin"
)

// GetWorkers returns all workers.
//
// GET /api/workers
func GetWorkers(router *gin.RouterGroup) {
	router.GET("/workers", func(ctx *gin.Context) {
		var workers entity.Workers
		if err := db.Db().Find(&workers).Error; err != nil {
			log.Errorf("cannot find workers: %s", err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		ctx.JSON(
			http.StatusOK,
			workers,
		)
	})
}
