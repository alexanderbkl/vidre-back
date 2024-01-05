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
		headersWritten := false

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
		timeParsed, err := time.Parse("2006-01-02T15:04:05.000Z", payload.Time)
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
			if workSchedule.EntryHour.IsZero() {
				workSchedule.EntryHour = timeParsed
			} else if (!workSchedule.RestStartHour.IsZero() || !workSchedule.RestEndHour.IsZero()) && !workSchedule.ExitHour.IsZero() {
				log.Errorf("Error: entry time cannot be set after exit time")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR LA ENTRADA DESPUÉS DE LA SALIDA"})
			} else {
				log.Errorf("Error: entry time cannot be set twice")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR LA ENTRADA DOS VECES"})
			}
		case "startRest":
			if workSchedule.RestStartHour.IsZero() && !workSchedule.EntryHour.IsZero() {
				workSchedule.RestStartHour = timeParsed
			} else if !workSchedule.RestStartHour.IsZero() && !workSchedule.RestEndHour.IsZero() {
				log.Errorf("Error: startRest time cannot be set twice")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR EL INICIO DEL DESCANSO DOS VECES"})
			} else if workSchedule.RestStartHour.IsZero() && workSchedule.EntryHour.IsZero() {
				log.Errorf("Error: startRest time cannot be set without entry time")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR EL INICIO DEL DESCANSO SIN REGISTRAR LA ENTRADA"})
			} else {
				log.Errorf("Error: startRest time cannot be set after endRest time")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR EL INICIO DEL DESCANSO DESPUÉS DE REGISTRAR EL FIN DEL DESCANSO"})
			}
		case "endRest":
			workSchedule.RestEndHour = timeParsed
			// if payload.Type is endRest and startRest is not set, return error
			if workSchedule.RestStartHour.IsZero() {
				log.Errorf("Error: endRest time cannot be set without startRest time")
				//time 30 minutes before the end of the rest in spanish timezones
				time := timeParsed.Add(-30 * time.Minute)
				//set the startRest time to the time calculated
				workSchedule.RestStartHour = time
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO HA INDICADO EL INICIO DEL DESCANSO, SE INDICA EL INICIO A LAS " + time.Format("15:04")})
			} else if !workSchedule.ExitHour.IsZero() {
				log.Errorf("Error: endRest time cannot be set after exit time")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR EL FIN DEL DESCANSO DESPUÉS DE LA SALIDA"})
			} else if workSchedule.EntryHour.IsZero() {
				log.Errorf("Error: endRest time cannot be set without entry time")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR EL FIN DEL DESCANSO SIN REGISTRAR LA ENTRADA"})
			}
		case "exit":
			if workSchedule.RestEndHour.IsZero() {
				log.Errorf("Error: exit time cannot be set without endRest time")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR LA SALIDA SIN REGISTRAR EL FIN DEL DESCANSO"})
			} else if workSchedule.RestStartHour.IsZero() {
				log.Errorf("Error: exit time cannot be set without startRest time")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR LA SALIDA SIN REGISTRAR EL INICIO DEL DESCANSO"})
			} else if !workSchedule.ExitHour.IsZero() {
				log.Errorf("Error: exit time cannot be set twice")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR LA SALIDA DOS VECES"})
			} else if workSchedule.EntryHour.IsZero() {
				log.Errorf("Error: exit time cannot be set without entry time")
				headersWritten = true
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "NO SE PUEDE REGISTRAR LA SALIDA SIN REGISTRAR LA ENTRADA"})
			} else {
				workSchedule.ExitHour = timeParsed
			}
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

		if !headersWritten {
			ctx.JSON(http.StatusOK, gin.H{"message": "Work schedule updated successfully", "work_schedule": workSchedule})
		}
	})
}

// AddWorkDay adds a workday by worker code and workday.
// POST /api/worker/work_day/add
func AddWorkDay(router *gin.RouterGroup) {
	router.POST("/worker/work_day/add", func(ctx *gin.Context) {
		var payload struct {
			WorkerCode    string `json:"worker_code"`
			Date          string `json:"date"` // Assuming the date comes in as a string like "2006-01-02"
			EnterHour     string `json:"enterHour"`
			ExitHour      string `json:"exitHour"`
			StartRestHour string `json:"startRestHour"`
			EndRestHour   string `json:"endRestHour"`
		}

		if err := ctx.ShouldBindJSON(&payload); err != nil {
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

		// Parse the date and time strings
		date, err := time.Parse("2006-01-02", payload.Date)
		if err != nil {
			log.Errorf("Error parsing date: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}
		enterHour, err := time.Parse("2006-01-02T15:04:05.000Z", payload.EnterHour)
		if err != nil {
			log.Errorf("Error parsing enterHour: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
			return
		}
		exitHour, err := time.Parse("2006-01-02T15:04:05.000Z", payload.ExitHour)
		if err != nil {
			log.Errorf("Error parsing exitHour: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
			return
		}
		startRestHour, err := time.Parse("2006-01-02T15:04:05.000Z", payload.StartRestHour)
		if err != nil {
			log.Errorf("Error parsing startRestHour: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
			return
		}
		endRestHour, err := time.Parse("2006-01-02T15:04:05.000Z", payload.EndRestHour)
		if err != nil {
			log.Errorf("Error parsing endRestHour: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
			return
		}

		// Find or initialize the work schedule for the worker on the given date
		var workSchedule entity.WorkSchedule
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
		workSchedule.EntryHour = enterHour
		workSchedule.ExitHour = exitHour
		workSchedule.RestStartHour = startRestHour
		workSchedule.RestEndHour = endRestHour

		// Save the updated work schedule
		if err := db.Db().Save(&workSchedule).Error; err != nil {
			log.Errorf("Error updating work schedule: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update work schedule"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Work schedule updated successfully", "work_schedule": workSchedule})
	})
}

// DeleteWorkday deletes a workday by worker code and workday.
// DELETE /api/worker/delete
func DeleteWorkDay(router *gin.RouterGroup) {
	router.DELETE("/worker/work_day/delete", func(ctx *gin.Context) {
		code := ctx.Query("worker_code") // assuming code is passed as a query parameter
		dateStr := ctx.Query("date")     // assuming date is passed as a query parameter

		if code == "" {
			log.Errorf("No code provided")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No code provided"})
			return
		}

		workerId, err := query.GetWorkerIDFromCode(code)
		if err != nil {
			log.Errorf("Error getting worker ID from code: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get worker ID from code"})
			return
		}

		if dateStr == "" {
			log.Errorf("No date provided")
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "No date provided"})
			return
		}

		// Parse the date and time strings
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Errorf("Error parsing date: %v", err)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
			return
		}

		// Delete the worker with the provided code
		if err := db.Db().Where("worker_id = ? AND date = ?", workerId, date).Delete(&entity.WorkSchedule{}).Error; err != nil {
			log.Errorf("cannot delete worker: %s", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete worker"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Workday deleted successfully"})
	})
}

// UpdateWorkDay updates a workday by worker code and workday.
// PUT /api/worker/work_day/update
func UpdateWorkDay(router *gin.RouterGroup) {
	router.PUT("/worker/work_day/update", func(ctx *gin.Context) {
		log.Printf("test")
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

		timeParsed, err := time.Parse("2006-01-02T15:04:05.000Z", payload.Time)
		if err != nil {
			timeParsed, err = time.Parse("2006-01-02", payload.Time)
			if err != nil {
				log.Errorf("Error parsing time date: %v", err)
				ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time format"})
				return
			}
		}

		// Find or initialize the work schedule for the worker on the given date
		var workSchedule entity.WorkSchedule
		result := db.Db().FirstOrCreate(&workSchedule, entity.WorkSchedule{
			WorkerID: workerId,
			Date:     date,
		})

		log.Printf("test6")
		log.Printf("Worker ID: %v", result.Error)
		if result.Error != nil {
			log.Printf("Error finding or creating work schedule: %v", result.Error)
			log.Errorf("Error finding or creating work schedule: %v", result.Error)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find or create work schedule"})
			return
		}

		// Update the relevant time based on the type of payload
		log.Printf("Type: %v", payload.Type)
		switch payload.Type {
		case "date":
			workSchedule.Date = timeParsed
		case "enterHour":
			workSchedule.EntryHour = timeParsed
		case "startRestHour":
			workSchedule.RestStartHour = timeParsed
		case "endRestHour":
			workSchedule.RestEndHour = timeParsed
		case "exitHour":
			workSchedule.ExitHour = timeParsed
		default:
			log.Errorf("Invalid type: %v", payload.Type)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type"})
			return
		}

		// Save the updated work schedule
		log.Printf("workSchedule date: %v", workSchedule.Date)
		if err := db.Db().Save(&workSchedule).Error; err != nil {
			log.Errorf("Error updating work schedule: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update work schedule"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"message": "Work schedule updated successfully", "work_schedule": workSchedule})
	})
}
