// RequestAddr
// $ go run RequestAddr.go encode.go

package main

import (
	"fmt"
	// "log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const UrlAddrServer = "127.0.0.1:8088" // (IPv6 ::1)

const NallocBits = 6 // allocate 1 out of 2^N bits (can allocate later)
const Nalloc = 1 << NallocBits // ratio of unallocated to allocated addresses
var maxAddr = (pow(Nchar, NcharA) >> NoffBit)
var maxAddrAlloc = (maxAddr >> NallocBits)

func main() {
	fmt.Println("main")
	fmt.Println(maxAddr/1e6, maxAddrAlloc/1e6)
	return
	
	// http.HandleFunc("/addr", addrHandle)
	// log.Fatal(http.ListenAndServe(UrlAddrServer, nil))
}

var shard int = 0 // Database shard assigned for address range
var iAddr = 0 // pointer to current address range

func addrHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "1024")
	iShard := (shard + 1) % Nshard
	addrBase := addrBaseArr[iAddr]
	iAddr++
	// Extend addrBase with NallocBits random bits
	addr := (addrBase << NallocBits) | rand.Intn(Nalloc)
	addrShardStr := addrShardToStr(addr, iShard)
	fmt.Fprintf(w, addrShardStr)
	// Mark addr in DB as allocated 
}


func addrShardToStr(addr, iShard int) string {
	// Add iShard to address space
	addrShard := addr << NshardBits | iShard
	addrShardStr := strconv.Itoa(addrShard)
	return addrShardStr
}
	
// Generate and save an array of addresses to use
func saveAddrBase() []int{
	rand.Seed(time.Now().UnixNano()) // pick random seed
	randAddrs := rand.Perm(maxAddrAlloc)
	//for _ = range randAddrs {
	//	// write to DB
	//}
	return randAddrs
}
var addrBaseArr = saveAddrBase()






