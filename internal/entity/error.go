package entity

import (
	"time"

	"github.com/alexanderbkl/vidre-back/internal/db"
	"github.com/alexanderbkl/vidre-back/internal/event"
	"github.com/sirupsen/logrus"
)

// Error represents an error message log.
type Error struct {
	ID           uint      `gorm:"primary_key"        json:"id"`
	ErrorTime    time.Time `                          json:"time"    sql:"index"`
	ErrorLevel   string    `gorm:"type:varchar(32)"   json:"level"`
	ErrorMessage string    `gorm:"type:varchar(2048)" json:"message"`
}

// Errors represents a list of error log messages.
type Errors []Error

// TableName returns the entity table name.
func (Error) TableName() string {
	return "errors"
}

// LogEvents logs published error events.
func (Error) LogEvents() {
	s := event.Subscribe("log.*")

	defer func() {
		event.Unsubscribe(s)
	}()

	for msg := range s.Receiver {
		level, ok := msg.Fields["level"]

		if !ok {
			continue
		}

		logLevel, err := logrus.ParseLevel(level.(string))

		if err != nil || logLevel >= logrus.InfoLevel {
			continue
		}

		newError := Error{ErrorLevel: logLevel.String()}

		if val, ok := msg.Fields["message"]; ok {
			newError.ErrorMessage = val.(string)
		}

		if val, ok := msg.Fields["time"]; ok {
			newError.ErrorTime = val.(time.Time)
		}

		db.Db().Create(&newError)
	}
}
