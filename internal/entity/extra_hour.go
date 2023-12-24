package entity

import (
	"time"

	"github.com/alexanderbkl/vidre-back/internal/db"
	"gorm.io/gorm"
)

type ExtraHour struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	WorkerID  uint           `gorm:"type:integer" json:"worker_id"`
	Worker    Worker         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"worker"`
	DayType   string         `gorm:"type:varchar(255)" json:"day_type"`
	StartHour string         `gorm:"type:varchar(255)" json:"start_hour"`
	EndHour   string         `gorm:"type:varchar(255)" json:"end_hour"`
	//IsEntry is a bool (can be true or false)
	IsEntry   bool           `gorm:"type:boolean" json:"is_entry"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (ExtraHour) TableName() string {
	return "extra_hours"
}

type ExtraHours []ExtraHour

func (extraHour *ExtraHour) Create() error {
	return db.Db().Create(extraHour).Error
}

func (extraHour *ExtraHour) TxCreate(tx *gorm.DB) error {
	return tx.Create(extraHour).Error
}

func (extraHour *ExtraHour) Save() error {
	return db.Db().Session(&gorm.Session{FullSaveAssociations: true}).Save(extraHour).Error
}

func (extraHour *ExtraHour) Count() (int64, error) {
	var count int64
	err := db.Db().Model(&ExtraHour{}).Count(&count).Error
	return count, err
}
