// RequestAddr
//$ go run RequestAddr.go addr.go dbAddr.go encode.go 'passwd
//   {}: 127.0.0.1:8088/addr
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const UrlAddrServer = "127.0.0.1:8088" // (IPv6 ::1)
var chAddr chan uint64

func main() {
	passwd := password()
	dB, err := OpenDB(passwd)
	if err != nil {
		fmt.Println("ERR RequestAddr:OpenDB")
		return
	}
	err = dB.CreateTable()
	fmt.Println("RequestAddr/CreateTable: ", err)
	
	const gochanDepth = 2
	chAddr = make(chan uint64, gochanDepth)
	go sendBaseAddr(chAddr)
	
	http.HandleFunc("/addr", addrHandle)
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
	shard = (shard + 1) % Nshard
	addrShard := addrShardToStr(addr, shard)
	fmt.Fprintf(w, addrShard)
}

func addrShardToStr(addr uint64, shard int) string {
	// Add shard to address space
	addrShard := addr<<NshardBits | uint64(shard)
	addrShardStr := strconv.Itoa(int(addrShard))
	return addrShardStr
}

// Queue base addresses for assignment using
//   buffered go channel as queue.
func sendBaseAddr(chBase chan uint64) {
	for {
		const SLEEPSEC = 1
		
		// Open/close on each iteration to be
		// more robust to DB interruption
		passwd := password()
		dB, err := OpenDB(passwd)
		if err != nil {
			fmt.Println("ERR sendBaseAddr: OpenDB")
			time.Sleep(SLEEPSEC * time.Second)
			continue
		}
		addr, err := dB.GetRandAddr()
		fmt.Println("RequestAddr: GEN addr", addr)
		dB.db.Close()
		
		if err != nil {
			fmt.Println("ERR sendBaseAddr: ", err)
			time.Sleep(SLEEPSEC * time.Second)
			continue
		}
		// Blocks when gochan buffer full
		fmt.Println("RequestAddr: send addr", addr)
		chAddr <- addr
	}
}

func password() string {
	if len(os.Args) <= 1 {
		log.Fatal("Supply DB password")
	}
	return os.Args[1]
}


