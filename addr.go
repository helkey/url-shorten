// Encode

package main

import (
	// "errors"
	// "fmt"
	"math/rand"
	"net/http"
	"time"
)


func main() {
	// TestAddr()
}

var nAddrs int = 1e6 // number of addresses offset from each base address
var retryInterval time.Duration = 10 // interval to re-try address server (sec)
var requestTimeout = retryInterval/2 // set timeouts to <retryInterval
const urlAddr = "localhost:8080"

func getAddr(urlAddr string, chAddr chan uint64) {
	// const nAddrs = 1e6
	const chDepth = 1 // channel queue depth to store lookahead base addresses
	chBase := make(chan uint64, chDepth)
	go getBaseAddr(urlAddr, chBase)
	for {
		// Random array of 'nAddrs' addresses offset from
		//    each base address from server
		baseAddr := <- chBase // get base address from server
		addrs := rand.Perm(nAddrs) // select addr in random order
		for addr := range addrs {
			chAddr <- baseAddr + uint64(addr)
		}
	}
}	


// Request/retry base address, database shard.
//   Store using buffered go channel as queue
func getBaseAddr(urlAddr string, chBase chan uint64)  {
	var baseAddr uint64
	var err error
	for {
		// Request base address range from server
		baseAddr, err = requestAddr(urlAddr)
		if err != nil {
			// Retry address server until responds
			ticker := time.NewTicker(retryInterval * time.Second)
			for _ = range ticker.C {
				baseAddr, err = requestAddr(urlAddr)
				if err == nil {
					break
				}
			}
			ticker.Stop()
		}
		// Send base address; channel blocks when full
		// fmt.Print("\nNew base addr:", baseAddr, ": ")
		chBase <- baseAddr
		// time.Sleep(1 * time.Second)
	}
}


// Single request for base address, database shard from remote address server
//   Note1: specify timeout (don't use default http request client)
//   TODO: re-write this with gRPC
func getBaseAddrServer(url string) (uint64, error) {
}
var requestAddr = getBaseAddrServer

func Test() (chan uint64, uint64, int)  {
	nAddrs = 10 // number of address per base address
	chAddr := make(chan uint64)
	go getAddr(urlAddr, chAddr)
	var base uint64 = 0
	var baseIncr uint64 = 1000
	requestAddr = func(url string) (uint64, error) {
		base += baseIncr
		return base, nil
	}
	return chAddr, baseIncr, nAddrs
}
