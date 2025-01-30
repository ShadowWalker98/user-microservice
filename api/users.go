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

	// inserting user into the users table

	err = app.models.UserModel.InsertUser(&user)
	if err != nil {
		app.logger.Println("error occurred while inserting into database: ", err)
		app.badRequestResponse(w, r, errors.New(fmt.Sprintf("user already exists with email %s", user.Email)))
		return
	}

	err = app.setCookies(w, r, &user)
	if err != nil {
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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
			app.logger.Println("error occurred while checking user details, user not found with email ",
				input.Email)
			app.IncorrectCredentialsResponse(w, r)
		}
		return
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

	// should I delete the pre-existing cookies and add new ones or
	// update the older ones?
	err = app.setCookies(w, r, user)
	if err != nil {
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.logger.Println(err)
	}
}

func (app *application) setCookies(w http.ResponseWriter, r *http.Request, user *data.User) error {
	// generating a new session token for the user if the login was successful

	sessionToken, err := app.models.UserTokenModel.CreateToken(user.UserId, data.Session)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
		return err
	}

	app.logger.Printf("the session token generated and stored was %s for user %d",
		sessionToken.Token,
		user.UserId)

	sessionCookie := &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken.Token,
		Path:     "/",
		Expires:  sessionToken.Expiry,
		HttpOnly: true,
	}

	// setting the session token on the response

	http.SetCookie(w, sessionCookie)

	// generating a CSRF token for the user

	csrfToken, err := app.models.UserTokenModel.CreateToken(user.UserId, data.CSRF)
	if err != nil {
		app.logger.Println("error occurred while generating CSRF token", err)
		app.serverErrorResponse(w, r, err)
		return err
	}

	csrfCookie := &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken.Token,
		Path:     "/",
		Expires:  csrfToken.Expiry,
		HttpOnly: false,
	}

	// setting the csrf token in the response
	http.SetCookie(w, csrfCookie)
	return nil
}

func (app *application) logoutHandler(w http.ResponseWriter, r *http.Request) {

	// we get the user id from the body

	var input struct {
		UserId int `json:"user_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// delete all tokens for this user from the users_tokens table

	err = app.models.UserTokenModel.DeleteAllTokens(input.UserId)
	if err != nil {
		app.logger.Println(err)
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w,
		http.StatusOK,
		envelope{"message": fmt.Sprintf("user with user id %d logged out",
			input.UserId)},
		nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) CheckIfLoggedInHandler(w http.ResponseWriter, r *http.Request) {
	// get the user id and the tokens that have been passed
	// Check if the token exists and has not expired
	// check if the tokens are matching the user id that has been passed
	// if it is all good then we return true to the caller, else we return false

	var input struct {
		UserId       int    `json:"user_id"`
		CSRFToken    string `json:"csrf_token"`
		SessionToken string `json:"session_token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.logger.Println("error occurred while parsing user input", err)
		app.badRequestResponse(w, r, err)
		return
	}

	// check for token validity for both
	valid, err := app.models.UserTokenModel.CheckTokenValidityForUser(input.UserId, data.Session, input.SessionToken)
	if err != nil {
		app.logger.Println("error occurred while checking token validity: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	if !valid {
		err = app.writeJSON(w, http.StatusOK, envelope{"validity": "false"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	valid, err = app.models.UserTokenModel.CheckTokenValidityForUser(input.UserId, data.CSRF, input.CSRFToken)
	if err != nil {
		app.logger.Println("error occurred while checking token validity: ", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	if !valid {
		err = app.writeJSON(w, http.StatusOK, envelope{"validity": "false"}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"validity": "true"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
