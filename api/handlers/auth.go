package handlers

import (
	"cashpal/api/utils"
	"cashpal/database"
	db "cashpal/database/generated"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
)

func parseUserRequestBody(body io.ReadCloser) (db.User, error) {
	var bodyData db.User
	if err := json.NewDecoder(body).Decode(&bodyData); err != nil {
		return bodyData, err
	}

	if bodyData.Username == "" {
		return bodyData, errors.New("username was not provided")
	}

	if bodyData.Password == "" {
		return bodyData, errors.New("password was not provided")
	}

	return bodyData, nil
}

func Login(w http.ResponseWriter, r *http.Request) {

	loginData, err := parseUserRequestBody(r.Body)

	if err != nil {
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

	user, err := query.GetUserByUsername(r.Context(), loginData.Username)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "invalid login credentials", http.StatusUnauthorized)
		return
	}

	if err := utils.VerifyPassword(loginData.Password, user.Password); err != nil {
		log.Println(err.Error())
		http.Error(w, "invalid login credentials", http.StatusUnauthorized)
		return
	}

	acessToken, err := utils.NewAccessToken(utils.GenerateClaims(user))

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "failed to create access token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"access_token": acessToken,
	}

	json.NewEncoder(w).Encode(response)
}
