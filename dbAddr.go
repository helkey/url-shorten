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
	port   = 5433
	user   = "postgres"
	dbName = "postgres"
)

// Database object pointer, shard
type DBS struct {
	db     *sql.DB
	shard  uint32
	passwd string
}

func main() {
	err := TestSaveurl()
	if err != nil {
		fmt.Println("main:", err)
	}
	// CreateTables(passwd)
}

//
func (dbS *DBS) OpenDB(shard uint32, passwd string) error {
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

	dbS.db = db
	dbS.shard = shard
	dbS.passwd = passwd
	return nil
}

func CreateTables(passwd string) error {
	if passwd == "" {
		if len(os.Args) <= 1 {
			e := "CreateTables: need Database password"
			return errors.New(e)
		}
		passwd = os.Args[1]
	}
	dbS := DBS{}
	err := dbS.OpenDB(uint32(0), passwd)
	if err != nil {
		fmt.Println("1:", err)
		return err
	}
	err = CreateTable(dbS.db)
	// err = DropTable(dbS.db)
	return err
}

func DropTable(db *sql.DB) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't drop database table"
			err = errors.New(e)
		}
	}()

	sqlTbl := `DROP TABLE url;`
	_, err = db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("3:", err)
	}

	sqlTbl = `DROP TABLE randext;`
	_, err = db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("4:", err)
		return err
	}
	return err
}

// Create new DB shard / table
func CreateTable(db *sql.DB) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't create database table"
			err = errors.New(e)
		}
	}()

	sqlTbl := `CREATE TABLE url (addr INTEGER PRIMARY KEY, fullurl TEXT);`
	_, err = db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("3:", err)
		return err
	}
	sqlTbl = `CREATE TABLE randext (addr INTEGER PRIMARY KEY, rand INT);`
	_, err = db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("4:", err)
	}
	return err
}

func NewDBconn(shard uint32) (DBS, error) {
	dbS := DBS{nil, 9999, ""}
	if len(os.Args) <= 1 {
		return dbS, errors.New("NewDBconn error: password not set")
	}
	passwd := os.Args[1]
	err := dbS.OpenDB(shard, passwd)
	return dbS, err
}

func TestSaveurl() error {
	const fullUrl, addr, randExt, shard = "http://Full.Url", uint64(0xaaaa), uint32(0xcccc), uint32(3)
	dbS, err := NewDBconn(shard)
	if err != nil {
		return err
	}
	dbS.SaveUrl(fullUrl, addr, randExt, shard, dbS.passwd)
	return nil
}

// Save URL mapping to DB
func (dbS DBS) SaveUrl(fullUrl string, addr uint64, randExt uint32, shard uint32, passwd string) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "SaveUrl: can't save URL mapping in database"
			err = errors.New(e)
		}
	}()

	err = dbS.OpenDB(shard, passwd)
	if err != nil {
		return errors.New("SaveUrl: database connection failed")
	}
	sqlIns := `INSERT INTO url (addr, fullUrl) VALUES ($1, $2);`
	fmt.Println(sqlIns, addr, fullUrl)
	// _, err = dbS.db.Exec(sqlIns, addr, fullUrl)
	if err != nil {
		return errors.New("SaveUrl: error saving to 'url' DB")
	}
	sqlIns = `INSERT INTO randext (addr, fullUrl) VALUES ($1, $2);`
	// _, err = dbS.db.Exec(sqlIns, addr, fullUrl)
	if err != nil {
		return errors.New("SaveUrl: error saving  to 'randext' DB table")
	}
	return nil
}

//
