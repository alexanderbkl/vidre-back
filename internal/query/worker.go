package query

import (
	"github.com/alexanderbkl/vidre-back/internal/db"
	"github.com/alexanderbkl/vidre-back/internal/entity"
)

func GetWorkerIDFromCode(code string) (uint, error) {
	var worker entity.Worker
	if err := db.Db().Where("code = ?", code).Find(&worker).Error; err != nil {
		return 0, err
	}
	return worker.ID, nil
}