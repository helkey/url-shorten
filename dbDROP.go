// dbDROP.go
// COULD DESTRY PRODUCTION DATABASE!!!
// DONT run this in same location as Prod DB

package main

import (
	"fmt"

	_ "github.com/lib/pq"
)

func init() {
}

func InitAddrTable() {
	dB, _ := OpenAddrDB(db_password())
	nRows, _ := dB.NumRowsDB("addrs")
	fmt.Printf("DB 'addr' HAD %d rows\n", nRows)
	dB.DropTable("addrs")
	dB.CreateAddrTable()
}

func InitUrlTable(shard int) (DB, error) {
	dB, err := OpenUrlDB(shard, db_password())
	if err != nil {
		fmt.Println("InitUrlTable: OpenUrlDB failed")
		return dB, err
	}

	// nRows, err := dB.NumRowsDB("url")
	if err != nil {
		fmt.Println("InitUrlTable: NumRowsDB failed")
		return dB, err
	}
	// fmt.Printf("DB 'url' HAD %d rows\n", nRows)

	err = dB.DropTable("url")
	if err != nil {
		fmt.Println("InitUrlTable: DropTable failed")
		return dB, err
	}

	err = dB.CreateUrlTable()
	if err != nil {
		fmt.Println("InitUrlTable: CreateUrlTable failed")
	}
	return dB, err
}

// Drop table of address assigned
func (dB DB) DropTable(name string) (err error) {
	e := fmt.Sprintf("ERR dbDROP: table '%s' not dropped", name)
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(e)
		}
	}()

	_, err = dB.db.Exec(`DROP TABLE ` + name + `;`)
	if err != nil {
		fmt.Println(e, ": ", err)
	}
	fmt.Printf("DROPPED TABLE %s; CAREFUL WITH PRODUCTION TABLES!\n", name)
	return
}
