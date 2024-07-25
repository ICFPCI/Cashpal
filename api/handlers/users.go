package handlers

import (
	"cashpal/api/utils"
	"cashpal/database"
	db "cashpal/database/generated"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

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

	w.Write(serializedUser)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "user id is invalid or malformed", http.StatusBadRequest)
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
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	serializedUser, err := json.Marshal(user)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.Write(serializedUser)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "user id is invalid or malformed", http.StatusBadRequest)
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
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := query.UpdateUser(r.Context(), user); err != nil {
		log.Println(err.Error())
		http.Error(w, "User update failed", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("user updated"))
}
