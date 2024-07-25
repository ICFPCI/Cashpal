package api

import (
	"cashpal/api/handlers"
	"cashpal/middleware"
	"net/http"
)

func SetupURLs(router *http.ServeMux) {

	protected := http.NewServeMux()
	router.Handle("/", middleware.ValidateJWT(protected))

	// Users
	protected.HandleFunc("GET /users", handlers.ListUsers)
	protected.HandleFunc("POST /users", handlers.CreateUser)
	protected.HandleFunc("GET /users/{userID}", handlers.GetUser)
	protected.HandleFunc("PATCH /users/{userID}", handlers.UpdateUser)

	// Authentication
	router.HandleFunc("GET /login", handlers.Login)

}
