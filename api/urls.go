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
	// protected.HandleFunc("DELETE /users/{userID}", handlers.DeleteUser)

	// Accounts
	protected.HandleFunc("GET /accounts", handlers.ListAccounts)
	protected.HandleFunc("GET /accounts/{accountID}", handlers.GetAccount)
	protected.HandleFunc("POST /accounts", handlers.CreateAccount)
	protected.HandleFunc("PATCH /accounts/{accountID}", handlers.UpdateAccount)
	protected.HandleFunc("DELETE /accounts/{accountID}", handlers.DeleteAccount)

	// Members
	protected.HandleFunc("GET /accounts/{accountID}/members", handlers.ListMembers)
	protected.HandleFunc("GET /accounts/{accountID}/members/{userID}", handlers.GetMember)
	protected.HandleFunc("POST /accounts/{accountID}/members", handlers.AddMember)
	protected.HandleFunc("PATCH /accounts/{accountID}/members/{userID}", handlers.UpdateMember)
	protected.HandleFunc("DELETE /accounts/{accountID}/members/{userID}", handlers.DeleteMember)

	// Authentication
	router.HandleFunc("GET /login", handlers.Login)

}
