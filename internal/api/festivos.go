package api

import (
	"net/http"

	"github.com/alexanderbkl/vidre-back/internal/db"
	"github.com/alexanderbkl/vidre-back/internal/entity"
	"github.com/gin-gonic/gin"
)

// GetFestivos returns the list of festivo dates
//
// GET /api/festivos
func GetFestivos(router *gin.RouterGroup) {
	router.GET("/festivos", func(c *gin.Context) {
		var festivos []entity.Festivo
		if err := db.Db().Find(&festivos).Error; err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, festivos)
		}
	})
}

// PostFestivo creates a new festivo date
//
// POST /api/festivos
func PostFestivo(router *gin.RouterGroup) {
	router.POST("/festivos", func(c *gin.Context) {
		var payload struct {
			Date string `json:"date"`
		}
		var festivo entity.Festivo
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		} else {
			festivo.Date = payload.Date
			if err := festivo.Create(db.Db()); err != nil {
				c.AbortWithStatus(http.StatusInternalServerError)
			} else {
				c.JSON(http.StatusOK, festivo)
			}
		}
	})
}

// DeleteFestivo deletes a festivo date
//
// DELETE /api/festivos/:date
func DeleteFestivo(router *gin.RouterGroup) {
	router.DELETE("/festivos/:date", func(c *gin.Context) {
		var festivo entity.Festivo
		if err := db.Db().Where("date = ?", c.Param("date")).Delete(&festivo).Error; err != nil {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			c.JSON(http.StatusOK, festivo)
		}
	})
}