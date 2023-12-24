package api

import (
	"net/http"

	"github.com/alexanderbkl/vidre-back/internal/db"
	"github.com/alexanderbkl/vidre-back/internal/entity"
	"github.com/alexanderbkl/vidre-back/internal/form"
	"github.com/gin-gonic/gin"
)

// CreateWorker creates a worker.
//
// POST /api/worker/create
// - JSON body:
//   - name: string
//   - code: string
func CreateWorker(router *gin.RouterGroup) {
	router.POST("/worker/create", func(ctx *gin.Context) {
		tx := db.Db().Begin()
		var req form.CreateWorkerRequest

		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Errorf("cannot bind json: %s", err)
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		worker := entity.Worker{
			Name: req.Name,
			Code: req.Code,
		}

		if err := worker.TxCreate(tx); err != nil {
			log.Errorf("cannot create worker: %s", err)
			tx.Rollback()
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		tx.Commit()

		ctx.JSON(http.StatusOK, gin.H{
			"worker": worker,
		})

	})
}

// ModifyWorker modifies a worker.
//
// POST /api/worker/modify
// - JSON body:
//   - code: string
//   - name: string
func ModifyWorker(router *gin.RouterGroup) {
	router.POST("/worker/update", func(ctx *gin.Context) {
		var req form.ModifyWorkerRequest
		log.Printf("1")

		if err := ctx.ShouldBindJSON(&req); err != nil {
			log.Errorf("cannot bind json: %s", err)
			ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
			return
		}

		var worker entity.Worker
		if err := db.Db().Model(&worker).Where("code = ?", req.Code).First(&worker).Error; err != nil {
			log.Errorf("worker not founad: %s", err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		worker.Name = req.Name
		if err := worker.Save(); err != nil {
			log.Errorf("cannot save worker: %s", err)
			ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
			return
		}

		log.Printf("worker updated: %v", worker)

		ctx.JSON(http.StatusOK, gin.H{
			"worker": worker,
		})

	})
}

// DeleteWorker deletes a worker by code.
// DELETE /api/worker/delete
func DeleteWorker(router *gin.RouterGroup) {
	router.DELETE("/worker/delete", func(ctx *gin.Context) {
		code := ctx.Query("code") // assuming code is passed as a query parameter

		if code == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		// Delete the worker with the provided code
		if err := db.Db().Where("code = ?", code).Delete(&entity.Worker{}).Error; err != nil {
			log.Errorf("cannot delete worker: %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete worker"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Worker deleted successfully"})
	})
}
