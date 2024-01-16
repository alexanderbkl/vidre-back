package entity

import (
	"gorm.io/gorm"
)


type Festivo struct {
	// unique date of the holiday
	Date string `gorm:"primary_key" json:"date"`
}

func (Festivo) TableName() string {
	return "festivos"
}

type Festivos []Festivo

func (festivo *Festivo) Create(db *gorm.DB) error {
	return db.Create(festivo).Error
}

func (festivo *Festivo) TxCreate(tx *gorm.DB) error {
	return tx.Create(festivo).Error
}

func (festivo *Festivo) Save(db *gorm.DB) error {
	return db.Session(&gorm.Session{FullSaveAssociations: true}).Save(festivo).Error
}

func (festivo *Festivo) TxSave(tx *gorm.DB) error {
	return tx.Session(&gorm.Session{FullSaveAssociations: true}).Save(festivo).Error
}

func (festivo *Festivo) Count(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&Festivo{}).Count(&count).Error
	return count, err
}