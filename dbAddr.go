// dbAddr.go
// go run dbAddr.go encode.go

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
// type DB struct {
//	*sql.DB
// }

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
	const name = "url"
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

	err = dbS.CreateTable(name)
	// err = DropTable(dbS.db)
	return err
}

func (dbS DBS) DropTable(name string) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't drop database table"
			err = errors.New(e)
		}
	}()

	sqlTbl := `DROP TABLE url;`
	_, err = dbS.db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("DropTable:", err)
	}
	return err
}

// Create new DB shard / table
func (dbS DBS) CreateTable(name string) (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't create database table"
			err = errors.New(e)
		}
	}()

	sqlTbl := fmt.Sprintf(`CREATE TABLE %s (addr INTEGER PRIMARY KEY, randext INT, fullurl TEXT);`, name)
	_, err = dbS.db.Exec(sqlTbl)
	if err != nil {
		fmt.Println("3:", err)
	}
	return
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
	sqlIns := `INSERT INTO url (addr, randext, fullurl) VALUES ($1, $2, $3);`
	fmt.Println(sqlIns, addr, randExt, fullUrl)
	_, err = dbS.db.Exec(sqlIns, addr, randExt, fullUrl)
	if err != nil {
		return errors.New("SaveUrl: error saving to 'url' DB")
	}
	return
}

//
func (dbS DBS) ReadUrlDB(addr uint64, shard uint32, passwd string) (fullUrl string, randE uint32, err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "ReadUlr: can't read URL from database"
			err = errors.New(e)
		}
	}()

	err = dbS.OpenDB(shard, passwd)
	if err != nil {
		err = errors.New("ReadUrlDB: database connection failed")
		return
	}

	sqlSel := fmt.Sprintf(`SELECT randext, fullurl FROM url WHERE addr = %d;`, addr)
	row := dbS.db.QueryRow(sqlSel)
	var randExt int
	err = row.Scan(&randExt, &fullUrl)
	randE = uint32(randExt)
	if err != nil {
		err = errors.New("ReadUrlDB: URL not found")
	}
	fmt.Println(randExt, fullUrl)
	return
}

func NewDBconn(shard uint32) (DBS, error) {
	const defaultShard, defaultPass = 9999, ""
	dbS := DBS{nil, defaultShard, defaultPass}
	if len(os.Args) <= 1 {
		return dbS, errors.New("NewDBconn error: password not set")
	}
	passwd := os.Args[1]
	err := dbS.OpenDB(shard, passwd)
	return dbS, err
}

const FullUrl = "http://Full.Url"
const Addr, RandExt = uint64(0xaaaa), uint32(0xcccc)
const Shard = 3
func TestSaveurl() error {
	const tableName = "url"

	// encodeA, _ := encodeAddr(addr, NcharA)
	// randShard := (RandExt << NshardBits) | Shard
	// encodeR, _ := encodeAddr(randShard, charR)
	const urlGrayList = false
	shortUrl, _ := EncodeURL(urlGrayList, Addr, RandExt, Shard)
	fmt.Println("TestSaveurl shortURL: ", shortUrl)
	
	dbS, err := NewDBconn(Shard)
	if err != nil {
		return err
	}
	dbS.DropTable(tableName)
	dbS.CreateTable(tableName)
	if err != nil {
		return err
	}

	err = dbS.SaveUrl(FullUrl, Addr, RandExt, Shard, dbS.passwd)
	if err != nil {
		return err
	}

	fullUrlR, randExtR, err := dbS.ReadUrlDB(Addr, Shard, dbS.passwd)
	fmt.Println(fullUrlR, randExtR, err)
	return err
}
