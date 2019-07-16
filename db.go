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
	dbUser   = "postgres"
	dbName = "postgres"
)

const INITIALIZEDB = false

type DB struct {
	db *sql.DB
}

func password() string {
	if len(os.Args) <= 1 {
		log.Fatal("Supply DB password")
	}
	return os.Args[1]
}

// Number of rows in db
func (dB DB) NumRowsDB(name string) (nRows int, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			// e := "dbAddr: can't load addr array from database"
			// err = errors.New(e)
		}
	}()

	// Allocate addrArr in single step
	row := dB.db.QueryRow(`SELECT COUNT(*) FROM ` + name + `;`)
	err = row.Scan(&nRows)
	return
}
