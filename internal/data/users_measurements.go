package data

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserMeasurementMetadataModel struct {
	db *gorm.DB
}

var MeasurementMapping = map[string]string{
	"ID":           "id",
	"Date":         "date",
	"Weight":       "weight",
	"Height":       "height",
	"Neck":         "neck",
	"Shoulders":    "shoulders",
	"Chest":        "chest",
	"LeftBicep":    "left_bicep",
	"RightBicep":   "right_bicep",
	"LeftForearm":  "left_forearm",
	"RightForearm": "right_forearm",
	"Waist":        "waist",
	"Hips":         "hips",
	"LeftThigh":    "left_thigh",
	"RightThigh":   "right_thigh",
	"LeftCalf":     "left_calf",
	"RightCalf":    "right_calf",
}

type UsersMeasurementsMetadata struct {
	ID           int       // Corresponds to 'id' in the table
	Date         time.Time // Corresponds to 'date' in the table
	Weight       float64   // Corresponds to 'weight' in the table
	Height       float64   // Corresponds to 'height' in the table
	Neck         int       // Corresponds to 'neck' in the table
	Shoulders    int       // Corresponds to 'shoulders' in the table
	Chest        int       // Corresponds to 'chest' in the table
	LeftBicep    int       // Corresponds to 'left_bicep' in the table
	RightBicep   int       // Corresponds to 'right_bicep' in the table
	LeftForearm  int       // Corresponds to 'left_forearm' in the table
	RightForearm int       // Corresponds to 'right_forearm' in the table
	Waist        int       // Corresponds to 'waist' in the table
	Hips         int       // Corresponds to 'hips' in the table
	LeftThigh    int       // Corresponds to 'left_thigh' in the table
	RightThigh   int       // Corresponds to 'right_thigh' in the table
	LeftCalf     int       // Corresponds to 'left_calf' in the table
	RightCalf    int       // Corresponds to 'right_calf' in the table
}

func (u UserMeasurementMetadataModel) Insert(metadata *UsersMeasurementsMetadata) (*UsersMeasurementsMetadata, error) {
	var colUpdateList []string

	colUpdateList = getList(metadata)

	metadata.Date = time.Now().UTC()
	result := u.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
			{Name: "date"},
		},
		DoUpdates: clause.AssignmentColumns(colUpdateList),
	}).Select(colUpdateList, "id", "date").Create(metadata)
	if result.Error != nil {
		return nil, result.Error
	}
	return metadata, nil
}

func getList(metadata *UsersMeasurementsMetadata) []string {
	var cols []string
	if metadata.Weight > 0 {
		cols = append(cols, MeasurementMapping["Weight"])
	}
	if metadata.Height > 0 {
		cols = append(cols, MeasurementMapping["Height"])
	}
	if metadata.Neck > 0 {
		cols = append(cols, MeasurementMapping["Neck"])
	}
	if metadata.Shoulders > 0 {
		cols = append(cols, MeasurementMapping["Shoulders"])
	}
	if metadata.Chest > 0 {
		cols = append(cols, MeasurementMapping["Chest"])
	}
	if metadata.LeftBicep > 0 {
		cols = append(cols, MeasurementMapping["LeftBicep"])
	}
	if metadata.RightBicep > 0 {
		cols = append(cols, MeasurementMapping["RightBicep"])
	}
	if metadata.LeftForearm > 0 {
		cols = append(cols, MeasurementMapping["LeftForearm"])
	}
	if metadata.RightForearm > 0 {
		cols = append(cols, MeasurementMapping["RightForearm"])
	}
	if metadata.Waist > 0 {
		cols = append(cols, MeasurementMapping["Waist"])
	}
	if metadata.Hips > 0 {
		cols = append(cols, MeasurementMapping["Hips"])
	}
	if metadata.LeftThigh > 0 {
		cols = append(cols, MeasurementMapping["LeftThigh"])
	}
	if metadata.RightThigh > 0 {
		cols = append(cols, MeasurementMapping["RightThigh"])
	}
	if metadata.LeftCalf > 0 {
		cols = append(cols, MeasurementMapping["LeftCalf"])
	}
	if metadata.RightCalf > 0 {
		cols = append(cols, MeasurementMapping["RightCalf"])
	}

	fmt.Println(cols)

	return cols
}
