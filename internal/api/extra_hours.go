package api

import (
	"net/http"

	"github.com/alexanderbkl/vidre-back/internal/db"
	"github.com/alexanderbkl/vidre-back/internal/entity"
	"github.com/alexanderbkl/vidre-back/internal/query"
	"github.com/gin-gonic/gin"
)

// GetExtraHours returns the extra hours for the given user code.
//
// GET /api/worker/extra_hours/:code
func GetExtraHours(router *gin.RouterGroup) {
	router.GET("/worker/extra_hours/:code", func(ctx *gin.Context) {
		var extraHours entity.ExtraHours
		code := ctx.Param("code")
		worker_id, err := query.GetWorkerIDFromCode(code)
		if err != nil {
			log.Errorf("cannot get worker id from code: %s", err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}
		if err := db.Db().Where("worker_id = ?", worker_id).Find(&extraHours).Error; err != nil {
			log.Errorf("cannot find extra hours: %s", err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		ctx.JSON(
			http.StatusOK,
			extraHours,
		)
	})
}

// ToggleExtraHours toggles the extra hours for a worker.
//
// POST /api/worker/extra_hours/toggle
func ToggleExtraHours(router *gin.RouterGroup) {
	router.POST("/worker/extra_hours/toggle", func(ctx *gin.Context) {
		var payload struct {
			Code      string `json:"code"`
			DayType   string `json:"day_type"`
			StartHour string `json:"start_hour"`
			EndHour   string `json:"end_hour"`
			Enabled   bool   `json:"enabled"`
			IsEntry  bool   `json:"is_entry"`
		}

		if err := ctx.ShouldBindJSON(&payload); err != nil {
			log.Errorf("cannot bind json: %s", err)
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		worker_id, err := query.GetWorkerIDFromCode(payload.Code)
		if err != nil {
			log.Errorf("cannot get worker id from code: %s", err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		if payload.Enabled {
			// Create new extra hours record
			extraHour := entity.ExtraHour{
				WorkerID:  worker_id,
				DayType:   payload.DayType,
				StartHour: payload.StartHour,
				EndHour:   payload.EndHour,
				IsEntry:   payload.IsEntry,
			}
			if err := extraHour.Create(); err != nil {
				log.Errorf("cannot create extra hour: %s", err)
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}

			ctx.JSON(http.StatusOK, gin.H{
				"message": "Extra hour created",
			})
		} else {
			// Remove existing extra hours record
			if err := db.Db().Where("worker_id = ? AND day_type = ? AND start_hour = ? AND end_hour = ?", worker_id, payload.DayType, payload.StartHour, payload.EndHour).Delete(&entity.ExtraHour{}).Error; err != nil {
				log.Errorf("cannot delete extra hour: %s", err)
				ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
				return
			}
			ctx.JSON(http.StatusOK, gin.H{
				"message": "Extra hour deleted",
			})
		}
	})
}
