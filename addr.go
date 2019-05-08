// Addr: Get encoding addresses (incl base address from server)
// $ go run addr.go encode.go

package main

import (
	// "encoding/binary"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	// testAddr()
	addrShard, err := getBaseAddrServer(UrlAddr)
	if err == nil {
		fmt.Println(addrShard)
	} else {
		fmt.Println("err:", err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano()) // initialize random seed
}

type BaseShard struct {
	addr  uint64
	shard uint32
}

var nAddrs int = 1e6                   // number of addresses offset from each base address
var retryInterval time.Duration = 10   // interval to re-try address server (sec)
var requestTimeout = retryInterval / 2 // set timeouts < retryInterval
const UrlAddr = "http://127.0.0.1:8088/addr"

func getAddr(urlAddr string, chAddr chan BaseShard) {
	// const nAddrs = 1e6
	const chDepth = 1 // channel queue depth to store lookahead base addresses
	chBase := make(chan uint64, chDepth)
	go getBaseAddr(urlAddr, chBase)
	shardMask := uint64(Nshard) - 1
	baseMask := ^shardMask
	// fmt.Printf("baseMask:%b  shardMask:%b", baseMask, shardMask)
	for {
		// Random array of 'nAddrs' addresses offset from
		//    each base address from server
		baseAddrShard := <-chBase // get base address from server
		baseAddr := baseAddrShard & baseMask
		baseShard := new(BaseShard)
		baseShard.shard = uint32(baseAddrShard & shardMask)
		addrs := rand.Perm(nAddrs) // select addr in random order
		for _, addr := range addrs {
			baseShard.addr = baseAddr + uint64(addr)
			chAddr <- *baseShard
		}
	}
}

// Request/retry base address, database shard.
//   Store using buffered go channel as queue
func getBaseAddr(urlAddr string, chBase chan uint64) {
	var baseAddrShard uint64
	var err error
	for {
		// Request base address range from server
		baseAddrShard, err = requestAddr(urlAddr)
		if err != nil {
			// Retry address server until responds
			ticker := time.NewTicker(retryInterval * time.Second)
			for _ = range ticker.C {
				baseAddrShard, err = requestAddr(urlAddr)
				if err == nil {
					break
				}
			}
			ticker.Stop()
		}
		// Send base address; channel blocks when full
		// fmt.Print("\nNew base addr:", baseAddr, ": ")
		chBase <- baseAddrShard
		// time.Sleep(1 * time.Second)
	}
}

// Single request for base address, database shard from remote address server
//   Note: specify timeout (don't use default http request client)
//   TODO: re-write this with gRPC
func getBaseAddrServer(urlAddr string) (uint64, error) {
	var netClient = &http.Client{
		Timeout: time.Second * requestTimeout,
	}
	resp, err := netClient.Get(urlAddr)
	if err != nil {
		fmt.Println("Failed netClient.Get")
		return 0, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed ioutil.ReadAll")
		return 0, err
	}
	fmt.Println("body:", body, string(body[0]))
	// addrShard := binary.BigEndian.Uint64(body)
	addrShard := uint64(1024)
	return addrShard, nil
}

var requestAddr = getBaseAddrServer

// For automated testing. Overrides defaults, mocks address server
func Test() chan BaseShard {
	// Override operational defaults
	nAddrs = 10 // number of address per base address
	var base uint64 = 0
	var baseIncr uint64 = 1024
	rand.Seed(0) // const seed for repeatible test results
	requestAddr = func(url string) (uint64, error) {
		base += baseIncr
		return base, nil
	}

	chAddr := make(chan BaseShard)
	go getAddr(UrlAddr, chAddr)
	return chAddr
}

// For automated testing. Overrides defaults, mocks address server
func Testserver() chan BaseShard {
	// Override operational defaults
	nAddrs = 10 // number of address per base address
	rand.Seed(0) // const seed for repeatible test results

	chAddr := make(chan BaseShard)
	go getAddr(UrlAddr, chAddr) // UrlAddr defined in RequestAddr
	return chAddr
}

// Used for diagnostics when getting failures using
//  $ go test addr_test.go
func testAddr() {
	// chAddr := Test()
	chAddr := Testserver()
	const sleepMs = 10 // avoid race condition between channels
	const nIter = 30
	for i := 0; i < nIter; i++ {
		addrShard := <-chAddr
		fmt.Print(addrShard.addr, ", ")
		time.Sleep(sleepMs * time.Millisecond)
	}
}
