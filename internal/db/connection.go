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

func CreateSchemas(dbPool *pgxpool.Pool) error {
	_, err := dbPool.Exec(context.Background(), `
	CREATE DATABASE IF NOT EXISTS passkey_auth;
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) UNIQUE NOT NULL,
		registration_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
  		last_updated_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
  		last_login_date TIMESTAMP WITH TIME ZONE, 
  		CONSTRAINT users_pkey PRIMARY KEY (id)
	);
	CREATE TYPE attestation_type AS ENUM ('direct', 'indirect', 'none');
	CREATE TABLE credentials (
	id SERIAL, 
	user_id INT NOT NULL, 
	credential_id VARCHAR(1023) NOT NULL UNIQUE, 
	public_key TEXT, 
	attestation_type ENUM('direct', 'indirect', 'none') NOT NULL, 
	aaguid CHAR(36) DEFAULT '00000000-0000-0000-0000-000000000000', 
	signature_count INT, 
	creation_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
	last_used_date TIMESTAMP WITH TIME ZONE, 
	last_updated_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, 
	type VARCHAR(25), 
	transports TEXT[], 
	backup_eligible BOOLEAN DEFAULT FALSE, 
	backup_state BOOLEAN DEFAULT FALSE, 
	PRIMARY KEY (id), 
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);	
`)
	if err != nil {
		fmt.Printf("Error while creating schemas: %v\n", err)
		return err
	}
	return nil
}
