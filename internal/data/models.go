package data

import (
	"gorm.io/gorm"
)

type Models struct {
}

func NewModels(conn *gorm.DB) Models {
	return Models{}
}
