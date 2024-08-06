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

func CreateTransactions(w http.ResponseWriter, r *http.Request) {}
func UpdateTransaction(w http.ResponseWriter, r *http.Request)  {}
func DeleteTransaction(w http.ResponseWriter, r *http.Request)  {}
