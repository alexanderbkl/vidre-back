package db

import (

	"github.com/alexanderbkl/vidre-back/internal/event"
	"gorm.io/gorm"
)

var log = event.Log

// Db returns the default *gorm.DB connection.
func Db() *gorm.DB {
	if dbConn.db == nil {
		return nil
	}

	return dbConn.db
}

// UnscopedDb returns an unscoped *gorm.DB connection
// that returns all records including deleted records.
func UnscopedDb() *gorm.DB {
	return Db().Unscoped()
}
