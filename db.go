// db.go
// go run db.go
// 

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


func db_password() (password string) {
	const db_password_env_variable = "TF_VAR_db_password"
	if len(os.Args) >1 {
		return os.Args[1]
	}
	password = os.Getenv(db_password_env_variable)
	if password == "" {
		log.Fatal("DB: export TF_VAR_db_password='password'")
	}
	return password
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
