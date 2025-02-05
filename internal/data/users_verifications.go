package data

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserVerificationModel struct {
	db *gorm.DB
}

type UsersVerifications struct {
	UserId           int
	VerificationCode int
	Expiry           time.Time
}

func (uvm UserVerificationModel) Insert(verification *UsersVerifications) error {
	fmt.Println("current time when adding code to db: " + time.Now().UTC().String())
	verification.Expiry = time.Now().UTC().Add(time.Minute * 30)
	result := uvm.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"verification_code", "expiry"}),
	}).Create(verification)
	if result.Error != nil {
		fmt.Println("error while inserting verification code for user: ", verification.UserId)
		return result.Error
	}
	return nil
}

func (uvm UserVerificationModel) Delete(userId int) error {

	result := uvm.db.Where("user_id = ?", userId).Delete(&UsersVerifications{})
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}

	return nil
}

func (uvm UserVerificationModel) GetByUserIdAndVerificationCode(userId int, verificationCode int) (*UsersVerifications, error) {
	var userVerification UsersVerifications
	result := uvm.db.
		Where("user_id = ?", userId).
		Where("verification_code = ?", verificationCode).
		First(&userVerification)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			fmt.Printf("there was no record found with user id %d and code %d", userId, verificationCode)
		}
		return nil, result.Error
	}

	fmt.Println("code expiry: ", userVerification.Expiry.UTC())
	fmt.Println("current time when checking expiry: " + time.Now().UTC().String())

	if userVerification.Expiry.Before(time.Now()) {
		return nil, errors.New("verification code expired, please try again")
	}

	return &userVerification, nil
}
