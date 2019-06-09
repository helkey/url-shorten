// RequestAddr
// go run RequestAddr.go addr.go db.go dbAddr.go encode.go network.go 'passwd
//   {}: 127.0.0.1:8088/addr

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	// "os"
	"time"
)

var chAddr chan uint64

func main() {
	// rand.Seed(time.Now().UnixNano()) // initialize random seed
	rand.Seed(0) // initialize deterministic seed
	const gochanDepth = 1
	chAddr = make(chan uint64, gochanDepth)
	go sendBaseAddr(chAddr)

	http.HandleFunc("/addr", addrHandle)
	// -> ListenAndServeTLS for https
	log.Fatal(http.ListenAndServe(UrlAddrServer, nil))
}

var shard int = 0 // Database shard assigned for address range
var iAddr = 0     // pointer to current address range

// Get base address from go channel buffer
// Round-robin database shard allocation
func addrHandle(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "addrHandle")
	fmt.Println("RequestAddr: addrHandle")
	addr := <-chAddr
	fmt.Println("RequestAddr: addr=", addr)
	addrShard := AddrShardToStr(addr, shard)
	fmt.Fprint(w, addrShard)
	shard = (shard + 1) % Nshard
}

// Queue base addresses for assignment using
//   buffered go channel as queue.
func sendBaseAddr(chBase chan uint64) {
	for {
		const SLEEPSEC = 1

		// Open/close on each iteration to be
		// more robust to DB interruption
		passwd := password()
		dB, err := OpenAddrDB(passwd)
		if err != nil {
			fmt.Println("ERR RequestAddr: OpenDB")
			time.Sleep(SLEEPSEC * time.Second)
			continue
		}
		fmt.Println("ReqAddr: getRandAddr")
		addr, err := dB.GetRandAddr()
		if err != nil {
			fmt.Println("ERR ReqAddr: getRandAddr")
			time.Sleep(SLEEPSEC * time.Second)
			continue
		}
		fmt.Println("ReqAddr SEND: ", addr)
		dB.db.Close()

		if err != nil {
			fmt.Println("ERR RequestAddr: ", err)
			time.Sleep(SLEEPSEC * time.Second)
			continue
		}
		// Blocks when gochan buffer full
		fmt.Println("RequestAddr: addr", addr)
		chAddr <- addr
	}
}
