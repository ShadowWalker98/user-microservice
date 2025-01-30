package main

import (
	"github.com/julienschmidt/httprouter"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"os"
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

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	// addition of swagger
	router.HandlerFunc(http.MethodGet, "/swagger/*any", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))
	router.GET("/swagger.json", serveSwaggerJson)

	// TODO: Add reset password flow later

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/signup", app.signupUsersHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/login", app.loginUsersHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/logout", app.logoutHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/validate", app.CheckIfLoggedInHandler)
	return router
}
