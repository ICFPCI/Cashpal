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

func ListTransactions(w http.ResponseWriter, r *http.Request) {
	contextUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user user id cannot be loaded from the session data", http.StatusInternalServerError)
		return
	}

	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusInternalServerError)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	listTransactionParams := db.ListTransactionByAccountParams{
		AccountID: int32(accountID),
		UserID:    contextUserID,
	}

	transactions, err := query.ListTransactionByAccount(r.Context(), listTransactionParams)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	print(transactions)

	serializedTransactions, err := json.Marshal(transactions)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedTransactions)
}

func GetTransaction(w http.ResponseWriter, r *http.Request) {
	contextUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user user id cannot be loaded from the session data", http.StatusInternalServerError)
		return
	}

	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusInternalServerError)
		return
	}

	transactionID, err := strconv.ParseInt(r.PathValue("transactionID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "transaction id is invalid or malformed", http.StatusInternalServerError)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	getTransactionParams := db.GetTransactionWithCheckParams{
		AccountID: int32(accountID),
		ID:        int32(transactionID),
		UserID:    contextUserID,
	}

	transaction, err := query.GetTransactionWithCheck(r.Context(), getTransactionParams)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "access denied", http.StatusUnauthorized)
		return
	}

	serializedTransactions, err := json.Marshal(transaction)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedTransactions)
}

func verifyMembership(context context.Context, query *db.Queries, userID int32, accountID int32) (int, error) {

	getMemberParams := db.GetMemberParams{
		AccountID: accountID,
		UserID:    userID,
	}

	_, err := query.GetMember(context, getMemberParams)

	if err != nil {
		return http.StatusUnauthorized, errors.New("the user is not part of the list of members")
	}

	return http.StatusOK, nil
}

func CreateTransactions(w http.ResponseWriter, r *http.Request) {
	contextUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user user id cannot be loaded from the session data", http.StatusInternalServerError)
		return
	}

	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusInternalServerError)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	statusCode, err := verifyMembership(r.Context(), query, contextUserID, int32(accountID))

	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	var newTransaction db.CreateTransactionParams

	if err := json.NewDecoder(r.Body).Decode(&newTransaction); err != nil {
		log.Println(err.Error())
		http.Error(w, "error parsing json from request body", http.StatusBadRequest)
		return
	}

	newTransaction.AccountID = int32(accountID)
	newTransaction.UserID = contextUserID

	transaction, err := query.CreateTransaction(r.Context(), newTransaction)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "transaction creation failed", http.StatusInternalServerError)
		return
	}

	serializedTransaction, err := json.Marshal(transaction)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedTransaction)
}

func compareTransactionData(currentTransaction *db.Transaction, updatedTransaction *db.UpdateTransactionParams) {
	if updatedTransaction.Amount == 0 {
		updatedTransaction.Amount = currentTransaction.Amount
	}

	if updatedTransaction.Description == "" {
		updatedTransaction.Description = currentTransaction.Description
	}
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	contextUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user user id cannot be loaded from the session data", http.StatusInternalServerError)
		return
	}

	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusInternalServerError)
		return
	}

	transactionID, err := strconv.ParseInt(r.PathValue("transactionID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "transaction id is invalid or malformed", http.StatusInternalServerError)
		return
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	statusCode, err := verifyMembership(r.Context(), query, contextUserID, int32(accountID))

	if err != nil {
		http.Error(w, err.Error(), statusCode)
		return
	}

	var updatedData db.UpdateTransactionParams

	json.NewDecoder(r.Body).Decode(&updatedData)

	getTransactionParams := db.GetTransactionWithCheckParams{
		AccountID: int32(accountID),
		ID:        int32(transactionID),
		UserID:    contextUserID,
	}

	transaction, err := query.GetTransactionWithCheck(r.Context(), getTransactionParams)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "access denied", http.StatusUnauthorized)
		return
	}

	compareTransactionData(&transaction, &updatedData)

	updatedData.ID = int32(transactionID)

	updatedTransaction, err := query.UpdateTransaction(r.Context(), updatedData)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "transaction update failed", http.StatusInternalServerError)
		return
	}

	serializedUpdatedTransaction, err := json.Marshal(updatedTransaction)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization  failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedUpdatedTransaction)
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("the transactions cannot be deleted yet"))
}
