package handlers

import (
	"cashpal/database"
	db "cashpal/database/generated"
	"cashpal/middleware"
	"context"
	"encoding/json"
	"errors"
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

	userID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user id cannot be loaded from the session", http.StatusInternalServerError)
	}

	ListMemberParams := db.ListMemberByAccountWithUserCheckParams{
		AccountID: int32(accountID),
		UserID:    int32(userID),
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	members, err := query.ListMemberByAccountWithUserCheck(r.Context(), ListMemberParams)

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

	contextUserID, ok := r.Context().Value(middleware.UserIDContextKey).(int32)

	if !ok {
		http.Error(w, "user id cannot be loaded from the session", http.StatusInternalServerError)
	}

	query, connClose, err := database.GetNewConnection(r.Context())

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	defer connClose()

	memberData := db.GetMemberWithUserCheckParams{
		AccountID: int32(accountID),
		UserID:    int32(userID),
		UserID_2:  int32(contextUserID), //UserID from the context
	}

	member, err := query.GetMemberWithUserCheck(r.Context(), memberData)

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

func verifyAdminRole(context context.Context, query db.Queries, accountID int32) error {

	contextUserID, ok := context.Value(middleware.UserIDContextKey).(int32)

	if !ok {
		return errors.New("user id cannot be loaded from the session")
	}

	logedMemberParams := db.GetMemberParams{
		AccountID: int32(accountID),
		UserID:    contextUserID,
	}

	logedMember, err := query.GetMember(context, logedMemberParams)

	if err != nil {
		return errors.New("error verifying account permissions")
	}

	if logedMember.MemberRoleID != 1 {
		return errors.New("administrator privileges are needed to modify the member list")
	}

	return nil
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

	if err := verifyAdminRole(r.Context(), *query, int32(accountID)); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	if err := verifyAdminRole(r.Context(), *query, int32(accountID)); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

	if err := verifyAdminRole(r.Context(), *query, int32(accountID)); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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
