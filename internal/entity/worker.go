package entity

import (
	"time"

	"github.com/alexanderbkl/vidre-back/internal/db"
	"gorm.io/gorm"
)

type Worker struct {
	ID        uint           `gorm:"primary_key" json:"id"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Code      string         `gorm:"type:varchar(255);unique" json:"code"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Worker) TableName() string {
	return "workers"
}

type Workers []Worker

func (worker *Worker) Create() error {
	return db.Db().Create(worker).Error
}

func (worker *Worker) TxCreate(tx *gorm.DB) error {
	return tx.Create(worker).Error
}

func (worker *Worker) Save() error {
	return db.Db().Session(&gorm.Session{FullSaveAssociations: true}).Save(worker).Error
}

func (worker *Worker) Count() (int64, error) {
	var count int64
	err := db.Db().Model(&Worker{}).Count(&count).Error
	return count, err
}


