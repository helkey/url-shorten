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

func addrHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "RequestAddr")
	return

	addrBase := addrBaseArr[iAddr]
	fmt.Fprintf(w, string(addrBase))

	// Extend addrBase with NallocBits random bits
	addr := (addrBase << NallocBits) | rand.Intn(Nalloc)
	addrShardStr := addrShardToStr(addr, shard)
	fmt.Fprintf(w, addrShardStr)

	shard = (shard + 1) % Nshard
	iAddr++
	// Mark addr in DB as allocated
}

func addrShardToStr(addr, iShard int) string {
	// Add iShard to address space
	addrShard := addr<<NshardBits | iShard
	addrShardStr := strconv.Itoa(addrShard)
	return addrShardStr
}

// Generate/save/load array of base addresses ranges
func loadAddrBase() {

	// dbAddr: SaveAddrBase(addrBaseArr)
}

var addrBaseArr []int
