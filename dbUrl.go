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

//
func OpenUrlDB(shard int, passwd string) (DB, error) {
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

/* Multiple tables to act as DB shards
func CreateTables(passwd string) error {
	const name = "url"
	if passwd == "" {
		if len(os.Args) <= 1 {
			e := "CreateTable: need Database password"
			return errors.New(e)
		}
		passwd = os.Args[1]
	}
	dB, err := OpenDB(0, passwd)
	if err != nil {
		fmt.Println("1:", err)
		return err
	}

	err = dB.CreateTable(name)
	err = dB.DropTable(name)
	return err
} */

func (dB DB) DropTable(name string) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't drop database table"
			err = errors.New(e)
		}
	}()

	sqlTbl := `DROP TABLE url;`
	_, err = dB.db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("DropTable:", err)
	}
	return err
}

// Create new DB shard / table
func (dB DB) CreateTable(name string) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't create database table"
			err = errors.New(e)
		}
	}()

	sqlTbl := fmt.Sprintf(`CREATE TABLE %s (addr INTEGER PRIMARY KEY, randext INT, fullurl TEXT);`, name)
	_, err = dB.db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("3:", err)
	}
	return
}

// Save URL mapping to DB
func (dB DB) SaveUrlDB(fullUrl string, addr uint64, randExt int) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "SaveUrl: can't save URL mapping in database"
			err = errors.New(e)
		}
	}()

	sqlIns := `INSERT INTO url (addr, randext, fullurl) VALUES ($1, $2, $3);`
	fmt.Println(sqlIns, addr, randExt, fullUrl)
	_, err = dB.db.Exec(sqlIns, addr, randExt, fullUrl)
	if err != nil {
		return errors.New("SaveUrl: error saving to 'url' DB")
	}
	return
}

//
func (dB DB) ReadUrlDB(addr uint64) (fullUrl string, randExt int, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "ReadUlr: can't read URL from database"
			err = errors.New(e)
		}
	}()

	sqlSel := fmt.Sprintf(`SELECT randext, fullurl FROM url WHERE addr = %d;`, addr)
	row := dB.db.QueryRow(sqlSel)
	err = row.Scan(&randExt, &fullUrl)
	if err != nil {
		err = errors.New("ReadUrlDB: URL not found")
	}
	fmt.Println(randExt, fullUrl)
	return
}

/* const FullUrl = "http://Full.Url"
const Addr, RandExt = uint64(0xaaaa), 0xcccc
const Shard = 3
func TestSaveurl() error {
	const tableName = "url"

	// encodeA, _ := encodeAddr(addr, NcharA)
	// randShard := (RandExt << NshardBits) | Shard
	// encodeR, _ := encodeAddr(randShard, charR)
	const isGrayList = false
	shortUrl, _ := encode(isGrayList, Addr, RandExt, Shard, NcharR)
	fmt.Println("TestSaveurl shortURL: ", shortUrl)

	dB, err := NewDBconn(Shard)
	if err != nil {
		return err
	}
	dB.DropTable(tableName)
	dB.CreateTable(tableName)
	if err != nil {
		return err
	}

	err = dB.SaveUrl(FullUrl, Addr, RandExt)
	if err != nil {
		return err
	}

	fullUrlR, randExtR, err := dB.ReadUrlDB(Addr)
	fmt.Println(fullUrlR, randExtR, err)
	return err
} */
