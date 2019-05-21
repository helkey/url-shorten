package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	dbType = "postgres"
	host   = "localhost"
	port   = 5430 // plus shard #
	user   = "postgres"
	dbName = "addrUrlRand"
)

// Database object pointer, shard
type DBS struct {
	db    *sql.DB
	shard uint32
}

func main() {
}

//
func (dbS DBS) OpenDB(passwd string, shard uint32) error {
	if shard == dbS.shard {
		return nil
	}
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbName)
	db, err := sql.Open(dbType, dbInfo)
	if err != nil {
		log.Fatal("CreateTable: not able to connect to database")
		return err
	}
	dbS.db = db
	dbS.shard = shard
	return nil
}

// defer db.Close()
// Create new DB shard / table
func (dbS DBS) CreateTable(name string) error {
	db := dbS.db
	sqlTbl := `
CREATE TABLE addr (
  addr INT,
  fullUrl TEXT,
);
CREATE TABLE rand (
  addr INT,
  randExt TEXT,
);`
	_, err := db.Exec(sqlTbl)
	if err != nil {
		return errors.New("")
	}
	return nil
}

// Save shortened URL to DB
func (dbS DBS) SaveAddr(fullUrl string, addr uint64, randExt uint32, passwd string, shard uint32) error {
	err := dbS.OpenDB(passwd, shard)
	if err != nil {
		return errors.New("SaveAddr: not able to connect to database")
	}
	sqlIns := `
INSERT INTO fullUrl (addr, fullUrl)
VALUES ($1, $2)
INSERT INTO randExt (addr, randExt)
VALUES ($3, $4)`
	db := dbS.db
	_, err = db.Exec(sqlIns, addr, fullUrl)
	if err != nil {
		return errors.New("")
	}
	return nil
}
