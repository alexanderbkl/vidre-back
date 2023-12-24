package api

import (
	"net/http"
	"time"

	"github.com/alexanderbkl/vidre-back/internal/db"
	"github.com/alexanderbkl/vidre-back/internal/entity"
	"github.com/alexanderbkl/vidre-back/internal/query"
	"github.com/gin-gonic/gin"
)

// GetWorkDay returns the work schedule for a worker on a given date frame.
//
// GET /api/worker/work_day
func GetWorkDay(router *gin.RouterGroup) {
	router.GET("/worker/work_days", func(ctx *gin.Context) {
		var payload struct {
			WorkerCode string `form:"worker_code"`
			StartDate  string `form:"start_date"` // Assuming the date comes in as a string like "2006-01-02"
			EndDate    string `form:"end_date"`   // Assuming the date comes in as a string like "2006-01-02"
		}


		if err := ctx.ShouldBindQuery(&payload); err != nil {
			log.Errorf("Error binding JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		workerId, err := query.GetWorkerIDFromCode(payload.WorkerCode)
		if err != nil {
			log.Errorf("Error getting worker ID from code: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get worker ID from code"})
			return
		}

		log.Printf("Start date: %v", payload.StartDate)
		log.Printf("Worker code: %v", payload.WorkerCode)
		// Parse the start and end dates
		startDate, err := time.Parse("2006-01-02", payload.StartDate)
		if err != nil {
			log.Errorf("Error parsing start date: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format"})
			return
		}
		endDate, err := time.Parse("2006-01-02", payload.EndDate)
		if err != nil {
			log.Errorf("Error parsing end date: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format"})
			return
		}

		// Find the work schedules for the worker on the given date frame
		var workSchedules []entity.WorkSchedule

		result := db.Db().Where("worker_id = ? AND date >= ? AND date <= ?", workerId, startDate, endDate).Find(&workSchedules)

		if result.Error != nil {
			log.Errorf("Error finding work schedules: %v", result.Error)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find work schedules"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"work_schedules": workSchedules})
	})
}

// PostWorkDay records the entry, exit and rest times for a worker's workday.
//
// POST /api/worker/work_day
func PostWorkDay(router *gin.RouterGroup) {
	router.POST("/worker/work_day", func(ctx *gin.Context) {
		var payload struct {
			WorkerCode string `json:"worker_code"`
			Date       string `json:"date"` // Assuming the date comes in as a string like "2006-01-02"
			Type       string `json:"type"` // Type can be "entry", "startRest", "endRest", "exit"
			Time       string `json:"time"` // Assuming time comes in as a string like "15:04"
		}


		if err := ctx.ShouldBindJSON(&payload); err != nil {
			log.Errorf("Error binding JSON: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		log.Printf("Payload: %v", payload)
		workerId, err := query.GetWorkerIDFromCode(payload.WorkerCode)
		if err != nil {
			log.Errorf("Error getting worker ID from code: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get worker ID from code"})
			return
		}

		// Parse the date and time strings
		date, err := time.Parse("2006-01-02", payload.Date)
		if err != nil {
			log.Errorf("Error parsing date: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		timeParsed, err := time.Parse("15:04", payload.Time)
		if err != nil {
			log.Errorf("Error parsing time: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
			return
		}

		// Find or initialize the work schedule for the worker on the given date
		var workSchedule entity.WorkSchedule
		log.Printf("Worker ID: %v", workerId)
		result := db.Db().FirstOrCreate(&workSchedule, entity.WorkSchedule{
			WorkerID: workerId,
			Date:     date,
		})

		if result.Error != nil {
			log.Errorf("Error finding or creating work schedule: %v", result.Error)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find or create work schedule"})
			return
		}

		// Update the relevant time based on the type of payload
		switch payload.Type {
		case "entry":
			workSchedule.EntryHour = timeParsed
		case "startRest":
			workSchedule.RestStartHour = timeParsed
		case "endRest":
			workSchedule.RestEndHour = timeParsed
		case "exit":
			workSchedule.ExitHour = timeParsed
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
			return
		}

		// Save the updated work schedule
		if err := db.Db().Save(&workSchedule).Error; err != nil {
			log.Errorf("Error updating work schedule: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update work schedule"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Work schedule updated successfully", "work_schedule": workSchedule})
	})
}
