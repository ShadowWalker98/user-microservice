package main

import (
	"errors"
	"fmt"
	"github.com/alexedwards/argon2id"
	"gorm.io/gorm"
	"net/http"
	"user-microservice/internal/data"
)

func (app *application) signupUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Printf("passed in data: %v\n", input)

	hash, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		app.logger.Println("error occurred while creating hash", err)
		app.serverErrorResponse(w, r, errors.New("error while processing your request"))
		return
	}

	fmt.Printf("the hash is %s\n", hash)

	user := data.User{
		Email:    input.Email,
		Password: hash,
	}

	err = app.models.UserModel.InsertUser(&user)
	if err != nil {
		app.logger.Println("error occurred while inserting into database: ", err)
		return
	}
}

func (app *application) loginUsersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Printf("passed in data: %v\n", input)

	user, err := app.models.UserModel.GetUserByEmail(input.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			app.logger.Println("error occurred while checking user details, user not found with email ", input.Email)
			app.IncorrectCredentialsResponse(w, r)
			return
		}
	}

	fmt.Printf("user found: %v\n", user)
	// checking if the password is correct by comparing hashes
	matched, err := argon2id.ComparePasswordAndHash(input.Password, user.Password)
	if err != nil {
		app.logger.Println("error occurred while comparing hashes", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	if !matched {
		app.IncorrectCredentialsResponse(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.logger.Println(err)
	}
}
