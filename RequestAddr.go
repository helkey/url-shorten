// RequestAddr
// $ go run RequestAddr.go encode.go

package main

import (
	"fmt"
	"log"
	"net/http"
)

// flaviocopes.com/golang-tutorial-rest-api/#simple-http-response-handler

const UrlAddrServer = "127.0.0.1:8088" // (IPv6 ::1)

const MaxAddr = pow(Ncar, NcharA) - 1
const maxAddrBase = (MaxAddr >> NoffBit)

func main() {
	http.HandleFunc("/addr", addrShard)
	log.Fatal(http.ListenAndServe(UrlAddr, nil))
}

var shard int = 0 // Database shard assigned for address range

func addrShard(w http.ResponseWriter, r *http.Request) {
	// shard = (shard + 1) % Nshard
	fmt.Fprintf(w, "1024")
	// addr.go
	// addr, shard := <- ch
}
