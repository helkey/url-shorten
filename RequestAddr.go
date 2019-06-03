// RequestAddr
// $ go run RequestAddr.go addr.go encode.go genAddr.go
// {}: 127.0.0.1:8088/addr
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const UrlAddrServer = "127.0.0.1:8088" // (IPv6 ::1)

func main() {
	fmt.Println("START RequestAddr")
	addrArr, err := LoadAddrArr()
	
	http.HandleFunc("/addr", addrHandle)
	log.Fatal(http.ListenAndServe(UrlAddrServer, nil))
}

var shard int = 0 // Database shard assigned for address range
var iAddr = 0     // pointer to current address range

// Get base address from go channel buffer
// Round-robin database shard allocation
func addrHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RequestAddr")
	return

	addr = <-chAddr
	shard = (shard + 1) % Nshard
	addrShard := addrShartToStr(addr, shard)
	fmt.Fprintf(w, "RequestAddr")
}

func addrShardToStr(addr, shard int) string {
	// Add iShard to address space
	addrShard := addr<<NshardBits | shard
	addrShardStr := strconv.Itoa(addrShard)
	return addrShardStr
}

// Queue base addresses for assignment
//   using buffered go channel as queue.
func sendBaseAddr(urlAddrServer string, chBase chan AddrShard) {
	const SLEEPSEC = 1
	for {
		// Get one base address
		for {
			addr, err := GetRandAddr()
			if err == nil {
				// Block here when buffer full
				chAddr <- addr
				break
			}
			// Failed to get base address;
			//   wait, then re-try
			time.Sleep(SLEEPSEC * time.Second)
		}
	}
}




