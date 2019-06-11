// db.go
// go run db.go

package main

import (
	"database/sql"
	"log"
	"os"
)

const (
	dbType = "postgres"
	host   = "localhost"
	port   = 5433
	user   = "postgres"
	dbName = "postgres"
)

type DB struct {
	db *sql.DB
}

func password() string {
	if len(os.Args) <= 1 {
		log.Fatal("Supply DB password")
	}
	return os.Args[1]
}
