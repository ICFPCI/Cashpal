package handlers

import (
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

func ListAccounts(w http.ResponseWriter, r *http.Request) {
	contextUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user user id cannot be loaded from the session data", http.StatusInternalServerError)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	accounts, err := query.ListAccountByUser(r.Context(), contextUserID)

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

	contextUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user user id cannot be loaded from the session data", http.StatusInternalServerError)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	userCheckParams := db.GetAccountWithUserCheckParams{
		ID:     int32(accountID),
		UserID: contextUserID,
	}

	account, err := query.GetAccountWithUserCheck(r.Context(), userCheckParams)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "this account does not exist", http.StatusNotFound)
		return
	}

	if account.IsMember != 1 {
		http.Error(w, "access denied", http.StatusForbidden)
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

func saveAccount(context context.Context, newAccount db.CreateAccountParams) (*db.Account, int, error) {
	userID, ok := context.Value(middleware.UserIDContextKey).(int32)

	if !ok {
		return nil, http.StatusInternalServerError, errors.New("user id cannot be loaded from the session")
	}

	query, connClose, tx, err := database.GetNewConnectionWithTransaction(context)

	if err != nil {
		log.Println(err.Error())
		return nil, http.StatusInternalServerError, errors.New("service unavailable")
	}

	defer connClose()
	defer tx.Rollback(context)

	qtx := query.WithTx(tx)

	account, err := qtx.CreateAccount(context, newAccount)

	if err != nil {
		log.Println(err.Error())
		return nil, http.StatusInternalServerError, errors.New("account creation failed")
	}

	member := db.CreateMemberParams{
		UserID:       userID,
		MemberRoleID: 1,
		AccountID:    account.ID,
	}

	if _, err := qtx.CreateMember(context, member); err != nil {
		log.Println(err.Error())
		return nil, http.StatusInternalServerError, errors.New("account creation failed")
	}

	tx.Commit(context)

	return &account, http.StatusOK, nil

}

func CreateAccount(w http.ResponseWriter, r *http.Request) {

	var newAccount db.CreateAccountParams

	if err := json.NewDecoder(r.Body).Decode(&newAccount); err != nil {
		log.Println(err.Error())
		http.Error(w, "error parsing json from request body", http.StatusBadRequest)
		return
	}

	account, statuscode, err := saveAccount(r.Context(), newAccount)

	if err != nil {
		http.Error(w, err.Error(), statuscode)
		return
	}

	serializedAccount, err := json.Marshal(account)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedAccount)

}

func compareData(updateData *db.UpdateAccountParams, currentData *db.GetAccountWithUserCheckRow) {
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

	contextUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user user id cannot be loaded from the session data", http.StatusInternalServerError)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	userCheckParams := db.GetAccountWithUserCheckParams{
		ID:     int32(accountID),
		UserID: contextUserID,
	}

	account, err := query.GetAccountWithUserCheck(r.Context(), userCheckParams)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "this account does not exist", http.StatusForbidden)
		return
	}

	if account.IsMember != 1 {
		http.Error(w, "access denied", http.StatusForbidden)
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
