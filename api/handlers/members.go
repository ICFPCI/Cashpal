package handlers

import (
	"cashpal/database"
	db "cashpal/database/generated"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func ListMembers(w http.ResponseWriter, r *http.Request) {

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

	members, err := query.ListMemberByAccount(r.Context(), int32(accountID))

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusNotFound)
		return
	}

	if len(members) <= 0 {
		log.Println("no members found for this account")
		http.Error(w, "no members found for this account", http.StatusNotFound)
		return
	}

	serializedMembers, err := json.Marshal(members)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedMembers)
}

func GetMember(w http.ResponseWriter, r *http.Request) {

	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusBadRequest)
		return
	}

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

	memberData := db.GetMemberParams{
		AccountID: int32(accountID),
		UserID:    int32(userID),
	}

	member, err := query.GetMember(r.Context(), memberData)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "this member does not exist", http.StatusNotFound)
		return
	}

	serializedMember, err := json.Marshal(member)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedMember)
}

func AddMember(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusBadRequest)
		return
	}

	var newMember db.CreateMemberParams

	if err := json.NewDecoder(r.Body).Decode(&newMember); err != nil {
		log.Println(err.Error())
		http.Error(w, "error parsing json from request body", http.StatusBadRequest)
		return
	}

	newMember.AccountID = int32(accountID)

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	member, err := query.CreateMember(r.Context(), newMember)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "error adding this member", http.StatusNotFound)
		return
	}

	serializedMember, err := json.Marshal(member)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedMember)
}

func UpdateMember(w http.ResponseWriter, r *http.Request) {

	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(r.PathValue("userID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "user id is invalid or malformed", http.StatusBadRequest)
		return
	}

	var updatedMemberData db.UpdateMemberParams

	if err := json.NewDecoder(r.Body).Decode(&updatedMemberData); err != nil {
		log.Println(err.Error())
		http.Error(w, "error parsing json from request body", http.StatusBadRequest)
		return
	}

	updatedMemberData.AccountID = int32(accountID)
	updatedMemberData.UserID = int32(userID)

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()
	member, err := query.UpdateMember(r.Context(), updatedMemberData)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "user update failed", http.StatusNotFound)
		return
	}

	serializedMember, err := json.Marshal(member)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "data serialization failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(serializedMember)
}

func DeleteMember(w http.ResponseWriter, r *http.Request) {
	accountID, err := strconv.ParseInt(r.PathValue("accountID"), 10, 32)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "account id is invalid or malformed", http.StatusBadRequest)
		return
	}

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

	memberData := db.DeleteMemberParams{
		AccountID: int32(accountID),
		UserID:    int32(userID),
	}

	if err := query.DeleteMember(r.Context(), memberData); err != nil {
		fmt.Println(err.Error())
		http.Error(w, "error when deleting member", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("member deleted"))
}
