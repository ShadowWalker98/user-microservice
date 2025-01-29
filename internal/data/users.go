package data

import (
	"gorm.io/gorm"
)

type UserModel struct {
	db *gorm.DB
}

type User struct {
	UserId    int    `json:"user_id"`
	Email     string `json:"email"`
	Password  string `json:"-"`
	CreatedAt string `json:"created_at"`
}

func (u UserModel) InsertUser(user *User) error {

	result := u.db.Omit("user_id", "created_at").Create(user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u UserModel) GetUserByEmail(email string) (*User, error) {

	var user = &User{
		Email: email,
	}

	result := u.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}
