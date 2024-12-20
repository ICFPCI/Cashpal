package database

import (
	"cashpal/config"
	database "cashpal/database/generated"
	"context"

	"github.com/jackc/pgx/v5"
)

func GetNewConnection(ctx context.Context) (*database.Queries, func(), error) {
	conn, err := pgx.Connect(ctx, config.GetSecret("DATABASE_URL"))

	if err != nil {
		return nil, nil, err
	}

	return database.New(conn), func() { conn.Close(ctx) }, nil
}

func GetNewConnectionWithTransaction(ctx context.Context) (*database.Queries, func(), pgx.Tx, error) {

	conn, err := pgx.Connect(ctx, config.GetSecret("DATABASE_URL"))

	if err != nil {
		return nil, nil, nil, err
	}

	tx, err := conn.Begin(ctx)

	if err != nil {
		return nil, nil, nil, err
	}

	return database.New(conn), func() { conn.Close(ctx) }, tx, nil
}
