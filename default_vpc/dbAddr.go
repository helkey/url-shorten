// dbAddr.go
// go run dbAddr.go addr.go db.go encode.go 'passwd'

// pgstart (WSL)  # Starting PostgreSQL 10 database server
// runpg (WSL)    # log into the psql prompt

package main

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq" // PostGres database
)

func OpenAddrDB(passwd string) (dB DB, err error) {
	// Recover from sql.Open() panic
	defer func() {
		if r := recover(); r != nil {
			e := "OpenAddrDB: can't open *ADDR* database table"
			err = errors.New(e)
		}
	}()
	
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		hostAddr, portAddr, dbUser, passwd, dbName)
	db, err := sql.Open(dbType, dbInfo)
	if err != nil {
		// e := "OpenAddrDB: not able to connect to database"
		return DB{}, err // errors.New(e)
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

	_, err = dB.db.Exec(`CREATE TABLE addrs (id BIGSERIAL PRIMARY KEY, addr BIGINT, avail BOOL);`)
	if err != nil {
		fmt.Println("dbAddr/createtable: ", err)
	}
	return
}

// Select random address. Return err if can't access address DB.
func (dB DB) GetRandAddr() (addr uint64, err error) {
	const SLEEPTIME = 1 // sec
	for {
		addr = uint64(rand.Intn(Nrange))
		var count int
		count, err = dB.NumRowsAddr(addr)
		if err != nil {
			// fmt.Println("ERR dbAddr/NumAddrRows", err)
			return
		}
		// If rand addr avail, save to DB, wait, check for concurrent selection
		if count == 0 {
			const NOTAVAIL = false
			err = dB.SaveAddrDB(addr, NOTAVAIL)
			if err != nil {
				fmt.Println("ERR dbAddr/SaveAddrDB")
				return
			}
			// Try again if addr not successfully inserted in DB,
			//   or two processes trying to use same random addr
			time.Sleep(SLEEPTIME * time.Second)
			count, err = dB.NumRowsAddr(addr)
			if err != nil {
				fmt.Println("ERR dbAddr/NumAddrRows - 2")
				return
			}
			if count == 1 {
				return // successfully selected rand addr
			}
		} else {
			fmt.Printf("dbAddr: addr=%v NOT AVAIL\n", addr)
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

// Number of rows with 'addr' in db
//   (should be max=1, except for simulaneous writes of same value)
func (dB DB) NumRowsAddr(addr uint64) (nAddr int, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			// e := "dbAddr: can't load addr array from database"
			// err = errors.New(e)
		}
	}()

	// Allocate addrArr in single step
	row := dB.db.QueryRow(`SELECT COUNT(*) FROM addrs WHERE addr = $1;`, addr)
	err = row.Scan(&nAddr)
	return
}
