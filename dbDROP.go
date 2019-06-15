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
	dB, _ := OpenAddrDB(password())
	nRows, _ := dB.NumRowsDB("addrs")
	fmt.Printf("DB 'addr' HAD %d rows\n", nRows)
	dB.DropTable("addrs")
	dB.CreateAddrTable()
}

func InitUrlTable(shard int) {
	dB, _ := OpenUrlDB(shard, password())
	nRows, _ := dB.NumRowsDB("url")
	fmt.Printf("DB 'url' HAD %d rows\n", nRows)
	dB.DropTable("url")
	dB.CreateUrlTable()
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

	fmt.Printf("DROPPED TABLE %s; CAREFUL WITH PRODUCTION TABLES!\n", name)
	_, err = dB.db.Exec(`DROP TABLE ` + name + `;`)
	if err != nil {
		fmt.Println(e, ": ", err)
	}
	return err
}
