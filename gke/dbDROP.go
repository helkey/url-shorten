// dbDROP.go
// COULD DESTRY PRODUCTION DATABASE!!!
// DONT run on Prod DB!!!

package main

import (
	"fmt"

	_ "github.com/lib/pq" // PostGres database
)

func InitAddrTable() {
	dB, err := OpenAddrDB(dbPassword())
	if err != nil {
		fmt.Println("InitAddrTable: OpenAddrDB failed")
		return
	}

	nRows, err := dB.NumRowsDB("addrs")
	if err == nil {
		fmt.Printf("DB 'addrs' HAD %d rows\n", nRows)
	}

	dB.DropTable("addrs")
	err = dB.CreateAddrTable()
	if err != nil {
		fmt.Println("InitUrlTable: CreateUrlTable failed")
	}
	dB.db.Close()
}

func InitUrlTable(shard int) {
	dB, err := OpenUrlDB(shard, dbPassword())
	if err != nil {
		fmt.Println("InitUrlTable: OpenUrlDB failed")
		return
	}

	nRows, err := dB.NumRowsDB("url")
	if err == nil {
		fmt.Printf("DB 'url' HAD %d rows\n", nRows)
	}


	dB.DropTable("url")
	err = dB.CreateUrlTable()
	if err != nil {
		fmt.Println("InitUrlTable: CreateUrlTable failed")
	}
	dB.db.Close()
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
