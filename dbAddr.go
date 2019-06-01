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
	db     *sql.DB
}

func OpenDB(passwd string) (DB, error) {
	dbInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, passwd, dbName)
	db, err := sql.Open(dbType, dbInfo)
	if err != nil {
		e := "dbAddr: not able to connect to database"
		fmt.Println(e)
		// log.Fatal(e)
		return DB{}, errors.New(e)
	}

	return DB{db}, nil
}


func (dB DB) DropTable(table string) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't drop database table"
			err = errors.New(e)
		}
	}()

	sqlDrop := fmt.Sprintf(`DROP TABLE %s;`, table)
	_, err = dB.db.Exec(sqlDrop)
	if err != nil {
		fmt.Println("DropTable:", err)
	}
	return err
}

func (dB DB) CreateTable(table string) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't create database table"
			err = errors.New(e)
		}
	}()

	sqlTbl := fmt.Sprintf(`CREATE TABLE %s (addr INTEGER PRIMARY KEY, avail BOOL);`, table)
	_, err = dB.db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("3:", err)
	}
	return
}

func (dB DB) SaveAddrArr(table string, addrArr []int) error {
	const assigned = false
	for _, addr := range addrArr {
		const avail = true
		err := dB.SaveAddrDB(table, addr, avail)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dB DB) SaveAddrDB(table string, addr int, avail bool) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "SaveUrl: can't save URL mapping in database"
			err = errors.New(e)
		}
	}()

	sqlIns := fmt.Sprintf(`INSERT INTO %s (addr, avail) VALUES ($1, $2);`, table)
	fmt.Println(sqlIns, addr, avail)
	_, err = dB.db.Exec(sqlIns, addr, avail)
	if err != nil {
		return errors.New("dbAddr: error saving to 'addr' DB")
	}
	return
}

func (dB DB) LoadAddrArr(table string) (addrArr []uint64, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			// e := "dbAddr: can't load addr array from database"
			// err = errors.New(e)
		}
	}()

	rows, err := dB.db.Query(`SELECT addr, avail FROM addr;`)
	addrArr = make([]uint64, 10)
	indx := 0
	var addr int
	var avail bool
	for rows.Next() {
		err = rows.Scan(&addr, &avail)
		// fmt.Println(addr, avail)
		if avail {
			fmt.Println(addr)
			addrArr[indx] = uint64(addr)
			indx++
		}
	}
	return addrArr, nil
}


// Mark assigned address range as unavailable for additional assignment
func (dB DB) MarkAddrUsed(table string, addr int) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "SaveUrl: can't save URL mapping in database"
			err = errors.New(e)
		}
	}()
	
	sqlUpdate := fmt.Sprintf(`UPDATE *s SET avail = 'FALSE' WHERE addr = %s;`, table, addr)
	_, err = dB.db.Exec(sqlUpdate)

	// READ valueQueryRow, make sure 'avail' now false
	// if err != nil {
	// err = errors.New("dbAddr: URL not found") }
	return

}


	
