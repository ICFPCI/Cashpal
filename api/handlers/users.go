package handlers

import (
	"cashpal/api/utils"
	"cashpal/database"
	db "cashpal/database/generated"
	"cashpal/middleware"
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

func verifyUserOwnership(userID int32, context context.Context) (int, error) {

	contextUserID, ok := context.Value(middleware.UserIDContextKey).(int32)
	if !ok {
		return http.StatusInternalServerError, errors.New("user id cannot be loaded from the session data")
	}

	if contextUserID != int32(userID) {
		return http.StatusForbidden, errors.New("access denied")
	}

	return http.StatusOK, nil
}

// NOT USED: The URL for this handler is not enabled.
func ListUsers(w http.ResponseWriter, r *http.Request) {
	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	users, err := query.ListUsers(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	serializedUsers, err := json.Marshal(users)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedUsers)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser db.CreateUserParams

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		log.Println(err.Error())
		http.Error(w, "error parsing json from request body", http.StatusBadRequest)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	newUser.Password, err = utils.HashPassword(newUser.Password)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user, err := query.CreateUser(r.Context(), newUser)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "user creation failed", http.StatusInternalServerError)
		return
	}

	serializedUser, err := json.Marshal(user)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedUser)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "user id is invalid or malformed", http.StatusBadRequest)
		return
	}

	if statusCode, err := verifyUserOwnership(int32(userID), r.Context()); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	user, err := query.GetUser(r.Context(), int32(userID))

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "this user does not exist", http.StatusNotFound)
		return
	}

	serializedUser, err := json.Marshal(user)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "user id is invalid or malformed", http.StatusBadRequest)
		return
	}

	if statusCode, err := verifyUserOwnership(int32(userID), r.Context()); err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	var user db.UpdateUserParams
	user.ID = int32(userID)

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Println(err.Error())
		http.Error(w, "error parsing json from request body", http.StatusBadRequest)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	user.Password, err = utils.HashPassword(user.Password)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	updatedUser, err := query.UpdateUser(r.Context(), user)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "user update failed", http.StatusInternalServerError)
		return
	}

	serializedUpdatedUser, err := json.Marshal(updatedUser)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedUpdatedUser)
}
