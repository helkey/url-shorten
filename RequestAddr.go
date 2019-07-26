// RequestAddr
// go run RequestAddr.go addr.go db.go dbAddr.go dbDrop.go dbUrl.go encode.go network.go 'passwd
//   {}: 127.0.0.1:8088/addr

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var chAddr chan uint64

func main() {
	// Initialize adddress database
	if (len(os.Args) > 1) && os.Args[1] == "InitAddr" {
		fmt.Println("Initializing addr database")
		InitAddrTable()
		return
	}

	const gochanDepth = 1
	rand.Seed(time.Now().UnixNano()) // initialize random seed
	chAddr = make(chan uint64, gochanDepth)
	go sendBaseAddr(chAddr)

	fmt.Println("ReqAddr: launched")
	http.HandleFunc("/addr", addrHandle)
	// http service:
	log.Fatal(http.ListenAndServe(UrlAddrServer, nil))
	// https service: USE ListenAndServeTLS()
}

var shard int = 0 // Database shard assigned for address range
var iAddr = 0     // pointer to current address range

// Get base address from go channel buffer
//   Uses round-robin database shard allocation
func addrHandle(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("RequestAddr: addrHandle")
	addr := <-chAddr
	addrShard := AddrShardToStr(addr, shard)
	shard = (shard + 1) % Nshard // cycle for next request
	fmt.Println("RequestAddr:   ", addrShard)
	fmt.Fprint(w, addrShard)
}

// Queue base addresses for assignment using buffered go channel
func sendBaseAddr(chBase chan uint64) {
	const SLEEPSEC = 1
	passwd := dbPassword()
	for {
		// Open/close DB on each iteration
		dB, err := OpenAddrDB(passwd)
		if err != nil {
			fmt.Println("ERR RequestAddr: OpenDB")
			time.Sleep(SLEEPSEC * time.Second)
			continue
		}
		// fmt.Println("ReqAddr: getRandAddr")
		addr, err := dB.GetRandAddr()
		if err != nil {
			fmt.Println("ERR ReqAddr: getRandAddr", err)
			time.Sleep(SLEEPSEC * time.Second)
			continue
		}
		dB.db.Close()

		if err != nil {
			fmt.Println("ERR RequestAddr: ", err)
			time.Sleep(SLEEPSEC * time.Second)
			continue
		}
		// Blocks when gochan buffer full
		fmt.Println("ReqAddr QUEUE: ", addr)
		chAddr <- addr
	}
}
