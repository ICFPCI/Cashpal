package handlers

import (
	"cashpal/database"
	db "cashpal/database/generated"
	"cashpal/middleware"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func ListAccounts(w http.ResponseWriter, r *http.Request) {
	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	accounts, err := query.ListAccount(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	serializedAccounts, err := json.Marshal(accounts)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedAccounts)
}

func GetAccount(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusBadRequest)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	account, err := query.GetAccount(r.Context(), int32(accountID))

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "this account does not exist", http.StatusInternalServerError)
		return
	}

	serializedAccount, err := json.Marshal(account)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedAccount)
}

func CreateAccount(w http.ResponseWriter, r *http.Request) {

	var newAccount db.CreateAccountParams

	if err := json.NewDecoder(r.Body).Decode(&newAccount); err != nil {
		log.Println(err.Error())
		http.Error(w, "error parsing json from request body", http.StatusBadRequest)
		return
	}

	userID, ok := r.Context().Value(middleware.UserContextKey).(int32)

	if !ok {
		http.Error(w, "user id cannot be loaded from the session", http.StatusInternalServerError)
		return
	}

	newAccount.UserID = userID

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	account, err := query.CreateAccount(r.Context(), newAccount)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account creation failed", http.StatusInternalServerError)
		return
	}

	serializedAccount, err := json.Marshal(account)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedAccount)

}

func compareData(updateData *db.UpdateAccountParams, currentData *db.Account) {
	if updateData.AccountName == "" {
		updateData.AccountName = currentData.AccountName
	}

	if updateData.AccountType == "" {
		updateData.AccountType = currentData.AccountType
	}
}

func UpdateAccount(w http.ResponseWriter, r *http.Request) {
	var accountUpdateData db.UpdateAccountParams

	if err := json.NewDecoder(r.Body).Decode(&accountUpdateData); err != nil {
		log.Println(err.Error())
		http.Error(w, "error parsing json from request body", http.StatusBadRequest)
		return
	}

	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusBadRequest)
		return
	}

	accountUpdateData.ID = int32(accountID)

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	account, err := query.GetAccount(r.Context(), int32(accountID))

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "this account does not exist", http.StatusNotFound)
		return
	}

	compareData(&accountUpdateData, &account)

	updatedAccount, err := query.UpdateAccount(r.Context(), accountUpdateData)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account update failed", http.StatusInternalServerError)
		return
	}

	serializedUpdatedAccount, err := json.Marshal(updatedAccount)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedUpdatedAccount)
}

func DeleteAccount(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("the accounts cannot be deleted yet"))
}
