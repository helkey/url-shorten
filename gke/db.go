// db.go
// go run db.go
//

package main

import (
	"database/sql"
	"fmt"
	"os"
)

const (
	dbType = "postgres"
	dbUser = "postgres"
	dbName = "postgres"
)

const INITIALIZEDB = false

type DB struct {
	db *sql.DB
}

const init_str = "init"

func dbPassword() (password string) {
	// db password from run-time argument
	if (len(os.Args) > 1) && (os.Args[1] != init_str) {
		fmt.Println("Password", os.Args[1])
		return os.Args[1]
	}

	/* db password from environment variable
	const password_env_variable = "TF_VAR_db_password"
	password = os.Getenv(password_env_variable)
	if password == "" {
		log.Fatal("DB: export " + password_env_variable + ": no password'")
	}
	fmt.Println(password)
	return password */

	// Specify db password as compile-time argument
	//   e.g. go build -ldflags "-X db.db_password=$TF_VAR_db_password"
	return db_password
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
