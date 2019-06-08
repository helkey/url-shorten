// RequestShorten
// go run RequestShorten.go addr.go dbUrl.go encode.go 'passwd

// "localhost:8086/create" (working WSL)
// "localhost:8086/create/?source=&url=http://FullURL"

package main

import (
	"database/sql"
	"errors"
	"fmt"
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

func OpenAddrDB(passwd string) (DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbName)
	db, err := sql.Open(dbType, dbInfo)
	if err != nil {
		e := "dbAddr: not able to connect to database"
		return DB{}, errors.New(e)
	}
	return DB{db}, nil
}

func (dB DB) CreateAddrTable() (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "dbAddr: can't create database table"
			err = errors.New(e)
		}
	}()

	_, err = dB.db.Exec(`CREATE TABLE addrs (addr INTEGER PRIMARY KEY, avail BOOL);`)
	if err != nil {
		fmt.Println("dbAddr/createtable: ", err)
	}
	return
}

func password() string {
	if len(os.Args) <= 1 {
		log.Fatal("Supply DB password")
	}
	return os.Args[1]
}
