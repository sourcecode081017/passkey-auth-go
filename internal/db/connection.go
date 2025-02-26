package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// conection instance for postgres
// create connection instance once and use it for all db operations
const AUTH_DATABSE_URL = "postgres://postgres:postgres@localhost:5432/passkey_auth"

func GetConnectionPool() *pgxpool.Pool {
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("AUTH_DATABSE_URL"))
	if err != nil {
		fmt.Printf("Error while creating connection pool: %v\n", err)
		os.Exit(1)
	}
	return dbPool
}
