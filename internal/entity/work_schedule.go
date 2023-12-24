package entity

import (
	"time"

	"gorm.io/gorm"
)

type WorkSchedule struct {
	ID           uint           `gorm:"primary_key" json:"id"`
	WorkerID     uint           `gorm:"type:integer;not null" json:"worker_id"`
	Date         time.Time      `gorm:"type:date;not null" json:"date"`
	EntryHour    time.Time      `gorm:"type:time;not null" json:"entry_hour"`
	ExitHour     time.Time      `gorm:"type:time;not null" json:"exit_hour"`
	RestStartHour time.Time     `gorm:"type:time;not null" json:"rest_start_hour"`
	RestEndHour   time.Time     `gorm:"type:time;not null" json:"rest_end_hour"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (WorkSchedule) TableName() string {
	return "work_schedules"
}

func (schedule *WorkSchedule) Create(db *gorm.DB) error {
	return db.Create(schedule).Error
}

type WorkSchedules []WorkSchedule

func (schedule *WorkSchedule) TxCreate(tx *gorm.DB) error {
	return tx.Create(schedule).Error
}

func (schedule *WorkSchedule) Save(db *gorm.DB) error {
	return db.Session(&gorm.Session{FullSaveAssociations: true}).Save(schedule).Error
}

func (schedule *WorkSchedule) TxSave(tx *gorm.DB) error {
	return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(schedule).Error
}

func (schedule *WorkSchedule) Count(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&WorkSchedule{}).Count(&count).Error
	return count, err
}

