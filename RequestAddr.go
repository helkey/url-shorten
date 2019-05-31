// RequestAddr
// $ go run RequestAddr.go addr.go encode.go
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

// To reduce resources required for demonstration,
//   only 1/64th all address ranges can be allocated with this code.
// Code is structured so that it could be upgraded to the full range
//   in service, with randomization on which rangs would be use,
//   so that access to the code base would not reduce security.
const NallocBits = 6           // allocate 1 out of 2^N bits (can allocate later)
const Nalloc = 1 << NallocBits // ratio of unallocated to allocated addresses
var Nrange = pow(Nchar, NcharA) >> NaddrBit
var NrangeAlloc = Nrange >> NallocBits

func main() {
	fmt.Println("START RequestAddr")
	loadAddrBase()
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
	rand.Seed(time.Now().UnixNano()) // pick random seed
	addrBaseArr = rand.Perm(NrangeAlloc)
	// dbAddr: SaveAddrBase(addrBaseArr)
}

var addrBaseArr []int
