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
	nRows, _ := dB.NumRowsDB()
	fmt.Printf("DB 'addr' has %v rows\n", nRows)
	dB.DropTable("addr")
	dB.CreateAddrTable()
}

func InitUrlTable() {
	dB, _ := OpenAddrDB(password())
	nRows, _ := dB.NumRowsDB()
	fmt.Printf("DB 'url' has %v rows\n", nRows)
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
	_, err = dB.db.Exec(fmt.Sprintf(`DROP TABLE %s;`, name))
	if err != nil {
		fmt.Println(e, ": ", err)
	}
	return err
}

