package data

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

func (UserToken) TableName() string {
	return "users_tokens"
}

// add tokens capability to the request

type UserTokenModel struct {
	conn *gorm.DB
}

type TokenType int

const (
	Session TokenType = iota
	CSRF
)

type UserToken struct {
	UserId    int
	TokenType int
	Token     string
	Expiry    time.Time
}

func (utm UserTokenModel) CreateToken(userId int, tokenType TokenType) (*UserToken, error) {
	generatedToken, err := GenerateRandomToken(32)
	if err != nil {
		return nil, err
	}

	userToken := UserToken{
		UserId:    userId,
		TokenType: int(tokenType),
		Token:     generatedToken,
		Expiry:    time.Now().AddDate(0, 0, 2).UTC(),
	}

	result := utm.conn.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "user_id"},
			{Name: "token_type"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"token", "expiry"}),
	}).Create(&userToken)
	if result.Error != nil {
		return nil, result.Error
	}

	return &userToken, nil
}

func (utm UserTokenModel) UpdateToken(userId int, tokenType TokenType) (*UserToken, error) {
	token, err := GenerateRandomToken(32)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("error while generating new token for userId %d", userId))
	}
	expiryTime := time.Now().AddDate(0, 0, 2).UTC()
	userToken := UserToken{
		UserId:    userId,
		TokenType: int(tokenType),
		Token:     token,
		Expiry:    expiryTime,
	}

	result := utm.conn.Save(&userToken)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userToken, nil
}

func (utm UserTokenModel) DeleteToken(userId int, tokenType TokenType) error {
	userToken := UserToken{
		UserId:    userId,
		TokenType: int(tokenType),
	}
	result := utm.conn.Omit("token", "expiry").Delete(&userToken)
	if result.Error != nil {
		fmt.Println(result.Error)
		return result.Error
	}
	return nil
}

func (utm UserTokenModel) DeleteAllTokens(userId int) error {
	result := utm.conn.Where("user_id = ?", userId).Delete(&UserToken{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GenerateRandomToken(length int) (string, error) {
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)
	return strings.Trim(token[:len(token)-1], " \n"), nil
}

func (utm UserTokenModel) CheckTokenValidityForUser(userId int, tokenType TokenType, token string) (bool, error) {
	var userToken UserToken
	result := utm.conn.Where("user_id = ?", userId).
		Where("token_type = ?", tokenType).
		Where("token = ?", token).First(&userToken)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, result.Error
	}

	if userToken.Expiry.Before(time.Now().UTC()) {
		return false, nil
	}

	if userToken.UserId != userId {
		return false, nil
	}

	return true, nil
}

// TODO: Add a method which checks for both tokens at once instead of one by one
func (utm UserTokenModel) CheckTokenValidity(tokenType TokenType, token string) (int, error) {
	// returns user id and nil if the token is valid
	// returns -1 if invalid

	var userToken UserToken

	result := utm.conn.Where("token_type = ?", tokenType).Where("token = ?", token).First(&userToken)
	if result.Error != nil {
		fmt.Println("no token found")
		return -1, result.Error
	}

	if userToken.Expiry.Before(time.Now()) {
		fmt.Println("token has expired")
		return -1, result.Error
	}

	return userToken.UserId, nil
}
