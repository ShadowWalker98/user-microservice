package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
	"user-microservice/internal/data"
)

func serveSwaggerJson(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, err := os.Open("docs/swagger.json") // Path to your swagger.json file
	if err != nil {
		http.Error(w, "Unable to open swagger.json", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Set the appropriate content-type header
	w.Header().Set("Content-Type", "application/json")

	// Serve the file content directly
	http.ServeFile(w, r, "docs/swagger.json")
}

func (app *application) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("session_token")
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// TODO: Change it to verify both the tokens at once in a single DB query instead of one by one

		sessionUser, err := app.models.UserTokenModel.CheckTokenValidity(data.Session, cookie.Value)
		if err != nil {
			err = app.writeJSON(w,
				http.StatusOK,
				envelope{"message": "please login to view this page, session token expired"},
				nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		csrfToken := r.Header.Get("x-csrf-token")
		csrfUser, err := app.models.UserTokenModel.CheckTokenValidity(data.CSRF, csrfToken)
		if err != nil {
			err = app.writeJSON(w,
				http.StatusOK,
				envelope{"message": "please login to view this page, csrf token expired"},
				nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		if sessionUser != csrfUser {
			err = app.writeJSON(w,
				http.StatusOK,
				envelope{"message": "malformed tokens, users do not match"},
				nil)
			if err != nil {
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("CORS Middleware Executing for:", r.Method, r.URL.Path)

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Expose-Headers", "x-csrf-token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			fmt.Println("Hello from preflight check")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) corsLogoutMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, err := r.Cookie("session_token")
		if err != nil {
			fmt.Println("No session token in request")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "x-csrf-token, Content-Type")

		if r.Method == http.MethodOptions {
			fmt.Println("Hello from preflight check")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) logoutOptionsHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Println("This is from the options request handler")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Headers", "x-csrf-token, Content-Type")

	w.WriteHeader(200)
}

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// addition of swagger
	router.HandlerFunc(http.MethodGet, "/swagger/*any", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))
	router.GET("/swagger.json", serveSwaggerJson)

	// TODO: Add reset password flow later
	// TODO: Add all the routes to swagger

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/signup", app.corsMiddleware(app.signupUsersHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users/login", app.corsMiddleware(app.loginUsersHandler))
	router.HandlerFunc(http.MethodOptions, "/v1/users/logout", app.logoutOptionsHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/logout", app.corsLogoutMiddleware(app.authMiddleware(app.logoutHandler)))
	router.HandlerFunc(http.MethodPost, "/v1/users/validate", app.CheckIfLoggedInHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/validate-session", app.CheckIfSessionIsValid)
	router.HandlerFunc(http.MethodPost, "/v1/users/reset-password", app.resetPasswordHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/validate-verification-code", app.ValidateVerificationCodeHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/validate-password", app.authMiddleware(app.ValidatePasswordHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users/users-measurements", app.authMiddleware(app.AddUserMeasurementsHandler))
	return router
}
