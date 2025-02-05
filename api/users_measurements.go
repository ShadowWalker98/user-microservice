package main

import (
	"net/http"
	"time"
	"user-microservice/internal/data"
)

func (app *application) AddUserMeasurementsHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("Adding measurements to the table for user id")

	var input struct {
		ID           int       `json:"id"`            // Corresponds to 'id' in the table
		Date         time.Time `json:"date"`          // Corresponds to 'date' in the table
		Weight       float64   `json:"weight"`        // Corresponds to 'weight' in the table
		Height       float64   `json:"height"`        // Corresponds to 'height' in the table
		Neck         int       `json:"neck"`          // Corresponds to 'neck' in the table
		Shoulders    int       `json:"shoulders"`     // Corresponds to 'shoulders' in the table
		Chest        int       `json:"chest"`         // Corresponds to 'chest' in the table
		LeftBicep    int       `json:"left_bicep"`    // Corresponds to 'left_bicep' in the table
		RightBicep   int       `json:"right_bicep"`   // Corresponds to 'right_bicep' in the table
		LeftForearm  int       `json:"left_forearm"`  // Corresponds to 'left_forearm' in the table
		RightForearm int       `json:"right_forearm"` // Corresponds to 'right_forearm' in the table
		Waist        int       `json:"waist"`         // Corresponds to 'waist' in the table
		Hips         int       `json:"hips"`          // Corresponds to 'hips' in the table
		LeftThigh    int       `json:"left_thigh"`    // Corresponds to 'left_thigh' in the table
		RightThigh   int       `json:"right_thigh"`   // Corresponds to 'right_thigh' in the table
		LeftCalf     int       `json:"left_calf"`     // Corresponds to 'left_calf' in the table
		RightCalf    int       `json:"right_calf"`    // Corresponds to 'right_calf' in the table
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	dummyUserMetadata := data.UsersMeasurementsMetadata{
		ID:           input.ID,
		Weight:       input.Weight,
		Height:       input.Height,
		Neck:         input.Neck,
		Shoulders:    input.Shoulders,
		Chest:        input.Chest,
		LeftBicep:    input.LeftBicep,
		RightBicep:   input.RightBicep,
		LeftForearm:  input.LeftForearm,
		RightForearm: input.RightForearm,
		Waist:        input.Waist,
		Hips:         input.Hips,
		LeftThigh:    input.LeftThigh,
		RightThigh:   input.RightThigh,
		LeftCalf:     input.LeftCalf,
		RightCalf:    input.RightCalf,
	}

	userMetadata, err := app.models.UsersMeasurementsMetadataModel.Insert(&dummyUserMetadata)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user_metadata": userMetadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
