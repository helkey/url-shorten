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
	testAddr()
	/* addrShard, err := getBaseAddrServer(UrlAddr)
	if err == nil {
		fmt.Println(addrShard)
	} else {
		fmt.Println("err:", err)
	} */
}

func init() {
	rand.Seed(time.Now().UnixNano()) // initialize random seed
}

type BaseShard struct {
	addr  uint64
	shard uint32
}

const nAddrBit = 24 // # bits for address offset

const retryInterval time.Duration = 10   // interval to re-try address server (sec)
const requestTimeout = retryInterval / 2 // set timeouts < retryInterval
const UrlAddr = "http://127.0.0.1:8088/addr"

var maxAddrOff int  = (1 << NoffBit) - 1 // num addresses offset from each base address


func getAddr(urlAddr string, chAddr chan BaseShard) {
	const chDepth = 1 // channel queue depth to store lookahead base addresses
	chBase := make(chan uint64, chDepth)
	go getBaseAddr(urlAddr, chBase)
	shardMask := uint64(Nshard) - 1
	baseMask := ^shardMask
	// fmt.Printf("baseMask:%b  shardMask:%b", baseMask, shardMask)
	for {
		// Random array of 'maxAddrOff' addresses offset from
		//    each base address from server
		baseAddrShard := <-chBase // get base address from server
		baseAddr := baseAddrShard & baseMask
		baseShard := new(BaseShard)
		baseShard.shard = uint32(baseAddrShard & shardMask)
		addrs := rand.Perm(maxAddrOff + 1) // select addr offsets in random order
		for _, addr := range addrs {
			baseShard.addr = baseAddr + uint64(addr)
			chAddr <- *baseShard
		}
	}
}

// Request/retry base address, database shard.
//   Store using buffered go channel as queue.
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
	// https://stackoverflow.com/questions/11184336/how-to-convert-from-byte-to-int-in-go-programming
	// addrShard := binary.BigEndian.Uint64(body)
	addrShard := uint64(1024)
	return addrShard, nil
}

var requestAddr = getBaseAddrServer

// For automated testing. Overrides defaults, mocks address server
func MockServer() chan BaseShard {
	// Override operational defaults
	maxAddrOff = 9 // max value of address offset
	var base uint64 = 0
	const baseIncr uint64 = 1024
	rand.Seed(0) // const seed for repeatible test results
	requestAddr = func(url string) (uint64, error) {
		base += baseIncr
		return base, nil
	}

	chAddr := make(chan BaseShard)
	go getAddr(UrlAddr, chAddr)
	return chAddr
}

// Used for diagnostics when getting failures with 'go test'
func testAddr() {
	chAddr := MockServer()
	// Sleep time to avoid race conditions between channels
	const sleepMs = 10
	const nIter = 30
	for i := 0; i < nIter; i++ {
		addrShard := <-chAddr
		fmt.Print(addrShard.addr, ", ")
		time.Sleep(sleepMs * time.Millisecond)
	}
}
