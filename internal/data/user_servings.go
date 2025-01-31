package data

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserWaterServingModel struct {
	db *gorm.DB
}

type UserWaterServing struct {
	UserId           int
	WaterServingSize int
	Date             string
	RequiredServings int
	Consumed         int
}

func (uwsm UserWaterServingModel) AddServing(serving *UserWaterServing) {
	if serving.RequiredServings == 0 {
		serving.RequiredServings = 16
	}

	if serving.WaterServingSize == 0 {
		serving.WaterServingSize = 250
	}

	serving.Consumed++
	uwsm.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
			{Name: "date"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"consumed"}),
	})
}
