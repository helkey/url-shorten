// dbDROP.go
// COULD DESTRY PRODUCTION DATABASE!!!
// DONT run this in same location as Prod DB

package main

import (
	"errors"
	"fmt"

	_ "github.com/lib/pq"
)

func init() {
}

// Drop table of address assigned
func (dB DB) DropAddrTable() (err error) {
	// Recover from db.Exec() panic
	defer func() {
		if r := recover(); r != nil {
			e := "CreateTable: can't drop database table"
			err = errors.New(e)
		}
	}()

	fmt.Println("DROPPED TABLE addrs; CAREFUL WITH PRODUCTION TABLES!")
	_, err = dB.db.Exec(`DROP TABLE addrs;`)
	if err != nil {
		fmt.Println("dbAddr: table 'addrs' not dropped", err)
	}
	return err
}
