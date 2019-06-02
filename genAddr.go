// genAddr.go
// go run genAddr.go addr.go dbAddr.go encode.go 'passwd

// pgstart (WSL)  # Starting PostgreSQL 10 database server
// runpg (WSL)    # log into the psql prompt

package main

import (
	"fmt"
	// "log"
	"math/rand"
	"os"
)

// To reduce resources required for demonstration,
//   only 1/64th all address ranges are allocated with this version.
// Code is structured with randomization on which ranges would be used,
//   so that knowledge of the code base would not reduce security.
// This version can be upgraded to full addr range while the service is operational.

const NallocBits = 6 // allocate 1 out of 2^N bits (can allocate later)
var NrangeAlloc = Nrange >> NallocBits

// const addrTable = "addr"

func main() {
	err := test()
	if err != nil {
		fmt.Println(err)
	}
	// rand.Seed(time.Now().UnixNano()) // pick random seed
}

func test() error {
	passwd := os.Args[1]
	db, err := OpenDB(passwd)
	if err != nil {
		// fmt.Println("genAddr: error creating table 'addr'")
		return err
	}

	err = db.DropTable()
	if err != nil {
		fmt.Println(err)
	}

	err = db.CreateTable()
	if err != nil {
		return err
	}

	rand.Seed(0) // pick random seed
	fmt.Println("Test case - addresses not randomized")
	NrangeAlloc := Nrange >> NallocBits
	addrArr := rand.Perm(NrangeAlloc)[:10]
	err = db.SaveAddrArr(addrArr)
	if err != nil {
		return err
	}

	addrArrR, err := db.GetAddrArr()
	fmt.Println("   ", addrArrR[0])
	err = db.MarkAddrUsed(addrArrR[0])
	addrArrR1, err := db.GetAddrArr()
	fmt.Println(len(addrArrR1))
	if err != nil {
		return err
	}
	/* for _, addr := range addrArrR {
		fmt.Println(addr)
	} */
	return nil
}
