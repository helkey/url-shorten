// genAddr.go
// go run genAddr.go addr.go dbAddr.go encode.go 'passwd'

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

const NallocBits = 6           // allocate 1 out of 2^N bits (can allocate later)
var NrangeAlloc = Nrange >> NallocBits




const addrTable = "addr"

func main() {
        passwd := os.Args[1]
        db, err := OpenDB(passwd)
	if err != nil {
		fmt.Println("genAddr: error creating table 'addr'")
	}
	err = db.DropTable(addrTable)
	err = db.CreateTable(addrTable)

	NrangeAlloc := Nrange >> NallocBits
	// rand.Seed(time.Now().UnixNano()) // pick random seed
	rand.Seed(0) // pick random seed
	fmt.Println("Test result - addresses not randomized")
	addrArr := rand.Perm(NrangeAlloc)[:10]
	db.SaveAddrArr(addrTable, addrArr)

	addrArrR, err := db.LoadAddrArr(addrTable)
	fmt.Println("   ", addrArrR[0])
	// _, err = db.LoadAddrArr(addrTable)
	if err != nil {
		fmt.Println(err)
	}
	/* for _, addr := range addrArrR {
		fmt.Println(addr)
	} */
	
}






