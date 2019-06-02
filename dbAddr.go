// dbAddr.go
// go test dbAddr_test ***.go encode.go 'passwd'

// pgstart (WSL)  # Starting PostgreSQL 10 database server
// runpg (WSL)    # log into the psql prompt

package main

import (
	"database/sql"
	"errors"
	"fmt"

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
	_, err = dB.db.Exec(`DROP TABLE addr;`)
	fmt.Println("Table dropped?")
	if err != nil {
		fmt.Println("dbAddr: table 'addr' not dropped", err)
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

	_, err = dB.db.Exec(`CREATE TABLE addr (addr INTEGER PRIMARY KEY, avail BOOL);`)
	if err != nil {
		fmt.Println("dbAddr/createtable: ", err)
	}
	return
}

func (dB DB) SaveAddrArr(addrArr []int) error {
	const avail = true
	sqlInsert, err := dB.db.Prepare(`INSERT INTO addr (addr, avail) VALUES ($1, $2);`)
	if err != nil {
		return err
	}
	defer sqlInsert.Close()
	for _, addr := range addrArr {
		const avail = true
		err := dB.SaveAddrDB(sqlInsert, addr, avail)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dB DB) SaveAddrDB(sqlInsert *sql.Stmt, addr int, avail bool) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "SaveUrl: can't save URL mapping in database"
			err = errors.New(e)
		}
	}()

	_, err = dB.db.Exec(`INSERT INTO addr (addr, avail) VALUES ($1, $2);`, addr, avail)
	// _, err = sqlInsert.Query(addr, avail)
	if err != nil {
		return errors.New("dbAddr: error saving to 'addr' DB")
	}
	return
}

func (dB DB) GetAddrArr() (addrArr []uint64, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			// e := "dbAddr: can't load addr array from database"
			// err = errors.New(e)
		}
	}()

	// Allocate addrArr in single step
	row := dB.db.QueryRow(`SELECT COUNT(*) FROM addr WHERE avail = TRUE;`)
	var nAvail int
	err = row.Scan(&nAvail)
	if err != nil {
		return
	}
	addrArr = make([]uint64, nAvail)

	rows, err := dB.db.Query(`SELECT addr FROM addr WHERE avail = TRUE;`)
	var addr int
	// var avail bool
	indx := 0
	for rows.Next() {
		err = rows.Scan(&addr)
		fmt.Println(addr)
		addrArr[indx] = uint64(addr)
		indx++
	}
	return addrArr, nil
}

// Mark assigned address range as unavailable for additional assignment
func (dB DB) MarkAddrUsed(addr uint64) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "SaveUrl: can't save URL mapping in database"
			err = errors.New(e)
		}
	}()

	_, err = dB.db.Exec(`UPDATE addr SET avail = 'FALSE' WHERE addr = $1;`, addr)
	return

}
