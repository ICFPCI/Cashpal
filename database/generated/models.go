// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Account struct {
	ID          int32
	UserID      int32
	AccountName string
	AccountType string
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}

type AccountEvent struct {
	ID          int32
	AccountID   int32
	EventTypeID int32
	Description string
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}

type EventType struct {
	ID   int32
	Name string
}

type Member struct {
	ID           int32
	AccountID    int32
	UserID       int32
	MemberRoleID int32
	CreatedAt    pgtype.Timestamp
	UpdatedAt    pgtype.Timestamp
}

type MemberRole struct {
	ID   int32
	Name string
}

type Transaction struct {
	ID                int32
	AccountID         int32
	UserID            int32
	TransactionDate   pgtype.Date
	TransactionTypeID int32
	Amount            float64
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
	Description       string
}

type TransactionType struct {
	ID   int32
	Name string
}

type User struct {
	ID        int32
	Username  string
	Password  string
	CreatedAt pgtype.Timestamp
	UpdatedAt pgtype.Timestamp
}
