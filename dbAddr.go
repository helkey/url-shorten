// dbAddr.go
// go run dbAddr.go addr.go db.go encode.go 'passwd'

// pgstart (WSL)  # Starting PostgreSQL 10 database server
// runpg (WSL)    # log into the psql prompt

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/lib/pq"
)

// Select random address. Return err if can't access address DB.
func (dB DB) GetRandAddr() (addr uint64, err error) {
	for {
		addr = uint64(rand.Intn(Nrange))
		var count int
		fmt.Println(addr)
		count, err = dB.NumAddrRows(addr)
		if err != nil {
			fmt.Println("ERR dbAddr/NumAddrRows")
			return
		}
		fmt.Println(count, addr)
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
			const sleepSec = 1
			time.Sleep(sleepSec * time.Second)
			count, err = dB.NumAddrRows(addr)
			if err != nil {
				fmt.Println("ERR dbAddr/NumAddrRows - 2")
				return
			}
			if count == 1 {
				return // successfully selected rand addr
			}
		}
		// Otherwise try another rand addr

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
