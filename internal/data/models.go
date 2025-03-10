package data

import (
	"gorm.io/gorm"
)

type Models struct {
	UserModel                      UserModel
	UserTokenModel                 UserTokenModel
	UserVerificationModel          UserVerificationModel
	UsersMeasurementsMetadataModel UserMeasurementMetadataModel
}

func NewModels(conn *gorm.DB) Models {
	return Models{
		UserModel:                      UserModel{db: conn},
		UserTokenModel:                 UserTokenModel{conn: conn},
		UserVerificationModel:          UserVerificationModel{db: conn},
		UsersMeasurementsMetadataModel: UserMeasurementMetadataModel{db: conn},
	}
}
