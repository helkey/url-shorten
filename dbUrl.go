// dbAddr.go
// go run dbUrl.go

// pgstart (WSL)  # Starting PostgreSQL 10 database server
// runpg (WSL)    # log into the psql prompt

package main

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

func OpenUrlDB(shard int, passwd string) (dB DB, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't create database table"
			err = errors.New(e)
		}
	}()

	// fmt.Println("password:", passwd)
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbName)
	db, err := sql.Open(dbType, dbInfo)
	if err != nil {
		e := "CreateTable: not able to connect to database"
		fmt.Println(e)
		// log.Fatal(e)
		return DB{}, errors.New(e)
	}
	return DB{db}, nil
}

// Create new DB shard / table
func (dB DB) CreateUrlTable() (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't create database table"
			err = errors.New(e)
		}
	}()

	_, err = dB.db.Exec(`CREATE TABLE url (addr INTEGER PRIMARY KEY, randext INT, nchar INT, fullurl TEXT);`)
	if err != nil {
		fmt.Println("3:", err)
	}
	return
}

// Save URL mapping to DB
func (dB DB) SaveUrlDB(fullUrl string, addr uint64, randExt, nChar int) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "SaveUrl: can't save URL mapping in database"
			err = errors.New(e)
		}
	}()

	sqlIns := `INSERT INTO url (addr, randext, nchar, fullurl) VALUES ($1, $2, $3, $4);`
	fmt.Printf("INSERT (addr=%v, randext=%v, nchar=%v, fullurl=%v)\n", addr, randExt, nChar, fullUrl)
	_, err = dB.db.Exec(sqlIns, addr, randExt, nChar, fullUrl)
	if err != nil {
		return errors.New("SaveUrl: error saving to 'url' DB")
	}
	return
}

// Read randExt, fullUrl given shortened address
func (dB DB) ReadUrlDB(addr uint64) (fullUrl string, randExt int, nChar int, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "ReadUlr: can't read URL from database"
			err = errors.New(e)
		}
	}()

	fmt.Println("ReadUrlDB:", addr)
	sqlSel := fmt.Sprintf(`SELECT randext, fullurl FROM url WHERE addr = %d;`, addr)
	row := dB.db.QueryRow(sqlSel)
	err = row.Scan(&randExt, &fullUrl)
	if err != nil {
		return // err = errors.New("ReadUrlDB: URL not found")
	}
	fmt.Println(randExt, fullUrl)
	return
}

// Check if long URL in database, return shortened URL
// func (dB DB) ExistsUrlDB(fullUrl string) (addr uint64, randExt int, nChar int, err error) {
func (dB DB) getShortUrl(fullUrl string, shard int) (shortUrl string, err error) {
	addr, randExt, nChar, err := dB.queryDBfullUrl(fullUrl)
	if err != nil {
		return
	}
	shortUrl, err = encode(addr, randExt, shard, nChar)
	return
}

func (dB DB) queryDBfullUrl(fullUrl string) (addr uint64, randExt int, nChar int, err error) {
	// Recover from database access panic
	defer func() {
		if r := recover(); r != nil {
			e := "ReadUlr: can't read URL from database"
			err = errors.New(e)
		}
	}()

	fmt.Println("search fullUrl:", fullUrl)
	sqlSel := fmt.Sprintf(`SELECT addr, randext, nchar FROM url WHERE fullurl = '%s';`, fullUrl)
	row := dB.db.QueryRow(sqlSel)
	fmt.Println("row", row)
	err = row.Scan(&addr, &randExt, &nChar)
	return
}
