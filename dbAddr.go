// dbAddr.go
// go run dbAddr.go addr.go encode.go 'passwd'

// pgstart (WSL)  # Starting PostgreSQL 10 database server
// runpg (WSL)    # log into the psql prompt

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"
	
	_ "github.com/lib/pq"
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

func OpenDB(passwd string) (DB, error) {
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

func (dB DB) DropTable() (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't drop database table"
			err = errors.New(e)
		}
	}()

	fmt.Println("DropTable")
	_, err = dB.db.Exec(`DROP TABLE addrs;`)
	fmt.Println("Table dropped?")
	if err != nil {
		fmt.Println("dbAddr: table 'addrs' not dropped", err)
	}
	return err
}

func (dB DB) CreateTable() (err error) {
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


// Select random address. Return err if can't access address DB.
func (dB DB) GetRandAddr() (addr uint64, err error) {
	for addr = uint64(rand.Intn(Nrange)); ; {
		var count int
		count, err = dB.NumAddrRows(addr)
		if err != nil {
			fmt.Println("ERR dbAddr/NumAddrRows")
			return
		}
		// Select again if random addr not avail
		if count == 0 {
			const NOTAVAIL = false
			err = dB.SaveAddrDB(addr, NOTAVAIL)
			if err != nil {
				fmt.Println("ERR dbAddr/SaveAddrDB")
				return
			}
			// Try again if addr not successfully inserted in DB,
			//   or two processes trying to use same random addr
			const sleepSec = 1
			time.Sleep(sleepSec * time.Second)
			count, err = dB.NumAddrRows(addr)
			if err != nil {
				fmt.Println("ERR dbAddr/NumAddrRows - 2")
				return
			}
			if (count == 1) {
				return // successfully selected rand addr
			}
		}
			
	}

}

func (dB DB) SaveAddrDB(addr uint64, avail bool) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "SaveUrl: can't save URL mapping in database"
			err = errors.New(e)
		}
	}()

	_, err = dB.db.Exec(`INSERT INTO addrs (addr, avail) VALUES ($1, $2);`, addr, avail)
	if err != nil {
		return errors.New("dbAddr: error saving to 'addrs' DB")
	}
	return
}

func (dB DB) NumAddrRows(addr uint64) (nAvail int, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			// e := "dbAddr: can't load addr array from database"
			// err = errors.New(e)
		}
	}()

	// Allocate addrArr in single step
	row := dB.db.QueryRow(`SELECT COUNT(*) FROM addrs WHERE addr = $1;`, addr)
	err = row.Scan(&nAvail)
	return
}


