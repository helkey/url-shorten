// dbAddr.go
// pgstart (WSL)  # Starting PostgreSQL 10 database server
// runpg (WSL)    # log into the psql prompt

package main

import (
	"database/sql"
	"errors"
	"fmt"
	// "log"
	"os"

	_ "github.com/lib/pq"
)

const (
	dbType = "postgres"
	host   = "localhost"
	port   = 5432
	user   = "helkey"
	dbName = "addrurlrand"
)

// Database object pointer, shard
type DBS struct {
	db     *sql.DB
	shard  uint32
	passwd string
}

func main() {
	passwd := os.Args[1]
	dbS := DBS{nil, 9999, passwd}
	fmt.Println(dbS.OpenDB(0, passwd))
	// CreateTables("")
}

//
func (dbS DBS) OpenDB(shard uint32, passwd string) error {
	if (passwd == dbS.passwd) && (shard == dbS.shard) {
		return nil
	}
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbName)
	db, err := sql.Open(dbType, dbInfo)
	if err != nil {
		e := "CreateTable: not able to connect to database"
		fmt.Println(e)
		// log.Fatal(e)
		return errors.New(e)
	}

	fmt.Println("password:", passwd)
	// sqlTbl := `CREATE DATABASE db0;`
	// _, err = db.Exec(sqlTbl)
	// fmt.Println("0:", err)

	sqlTbl := `CREATE TABLE randext (addr INT, fullurl TEXT);`
	_, err = dbS.db.Exec(sqlTbl)
	fmt.Println("1:", err)
	
	dbS = DBS{db, shard, passwd}
	return nil
}

func CreateTables(passwd string) {
	if passwd == "" {
		if len(os.Args) <= 1 {
			fmt.Println("CreateTables: need Database password")
			return
		}
		passwd = os.Args[1]
	}
	dbS := DBS{}
	dbS.OpenDB(uint32(0), passwd)
	fmt.Println(dbS.CreateTable(passwd))
	// if err != nil {
	// 	fmt.Println("CreateTables failed")
	//}
}

// defer db.Close()
// Create new DB shard / table
func (dbS DBS) CreateTable(passwd string) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't create database table"
			err = errors.New(e)

		}
	}()

	sqlTbl := `CREATE TABLE randext (addr INT, fullurl TEXT);`
	_, err = dbS.db.Exec(sqlTbl)
	// sqlTbl := `CREATE TABLE randext (addr INT, rand INT);`
	// _, err = db.Exec(sqlTbl) 
	return
}

	
// Save shortened URL to DB
func (dbS DBS) SaveAddr(fullUrl string, addr uint64, randExt uint32, shard uint32, passwd string) error {
	err := dbS.OpenDB(shard, passwd)
	if err != nil {
		return errors.New("SaveAddr: not able to connect to database")
	}
	db := dbS.db

	sqlIns := `
INSERT INTO fullUrl` + string(shard) + ` (addr, fullUrl)
VALUES ($1, $2);
INSERT INTO randExt` + string(shard) + ` (addr, randExt)
VALUES ($3, $4);`
	_, err = db.Exec(sqlIns, addr, fullUrl)
	if err != nil {
		return errors.New("")
	}
	return nil
}
