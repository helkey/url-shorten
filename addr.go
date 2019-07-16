// Addr: Get encoding addresses (incl base address from server)
// $ go run addr.go encode.go

package main

import (
	// "encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/* func main() {
	testAddr()
}*/

type AddrShard struct {
	addr  uint64
	shard int
}

// Address space
const NaddrBit = 24 // # bits for address offset

const retryInterval time.Duration = 10   // interval to re-try address server (sec)
const requestTimeout = retryInterval / 2 // set timeouts < retryInterval

// var maxAddrOff int = (1 << NoffBit) - 1 // max offset from each base address
var maxAddrOff int = (1 << NaddrBit) - 1 // max offset from each base address

// Runs 'RequestShorten' context, full base addr + addr extension
func getAddr(urlAddrServer string, chAddrSh chan AddrShard) {
	const chDepth = 1 // channel queue depth to store lookahead base addresses
	chBase := make(chan AddrShard, chDepth)
	go getBaseAddr(urlAddrServer, chBase)
	for {
		// Random array of 'maxAddrOff' addresses offset from
		//    each base address from server
		baseShard := <-chBase // get base address, DB shard from server
		baseAddr := baseShard.addr
		addrs := rand.Perm(maxAddrOff + 1) // select address offsets in random order
		for _, addr := range addrs {
			addrSh := new(AddrShard)
			addrSh.shard = baseShard.shard
			addrSh.addr = baseAddr + uint64(addr)
			// fmt.Printf("getAddr: base=%v; addr=%v; shard=%v\n", baseAddr, addrSh.addr, addrSh.shard)
			chAddrSh <- *addrSh
		}
	}
}

// Request/retry base address, database shard.
//   Store using buffered go channel as queue.
func getBaseAddr(urlAddrServer string, chBase chan AddrShard) {
	for {
		// Request base address range from server
		// getBaseAddrServe == baseAddrFromServer (except during testing)
		baseShard, err := getBaseAddrServe(urlAddrServer)
		if err != nil {
			// Set timeout: Retry address server until responds
			ticker := time.NewTicker(retryInterval * time.Second)
			for _ = range ticker.C {
				baseShard, err = baseAddrFromServer(urlAddrServer)
				if err == nil {
					break
				}
			}
			ticker.Stop()
		}
		// Send base address; channel blocks when full
		// fmt.Print("\nNew base addr:", baseAddr, ": ")
		chBase <- baseShard
	}
}

// var shardMask = uint64(Nshard) - 1
// var baseMask = ^shardMask

// Single request for base address, database shard from remote address server
//   Note: specify timeout (don't use default http request client)
//   TODO: re-write with gRPC
func baseAddrFromServer(urlAddrServer string) (AddrShard, error) {
	var netClient = &http.Client{
		Timeout: time.Second * requestTimeout,
	}
	const addrAPI = "/addr"
	// -> https service
	addrApiUrl := fmt.Sprintf("http://%s/addr", urlAddrServer)
	resp, err := netClient.Get(addrApiUrl)
	if err != nil {
		// Write this to log file
		return AddrShard{}, err // errors.New("Failed netClient.Get")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AddrShard{}, errors.New("Failed ioutil.ReadAll")
	}
	addrShard, err := sepAddrShard(string(body))
	if err != nil {
		return AddrShard{}, errors.New("Failed to recognize address from server")
	}

	return addrShard, err
}

const shardDelimit = "/"

func AddrShardToStr(addr uint64, shard int) string {
	// Add shard to address space
	// addrShard := addr<<NshardBits | uint64(shard)
	// addrShardStr := strconv.Itoa(int(addrShard))
	addrShardStr := strconv.Itoa(int(addr)) + shardDelimit + strconv.Itoa(shard)
	return addrShardStr
}

// Convert addrShard into addr and shard
func sepAddrShard(addrShardStr string) (addrShard AddrShard, err error) {
	iSep := strings.Index(addrShardStr, shardDelimit)
	if iSep < 0 {
		return addrShard, errors.New("Err addr: delimiter not found")
	}
	addr, err := strconv.Atoi(addrShardStr[:iSep])
	if err != nil {
		return
	}
	addrShard.addr = uint64(addr)
	addrShard.shard, err = strconv.Atoi(addrShardStr[iSep+1:])
	return
}

var getBaseAddrServe = baseAddrFromServer

// For automated testing. Overrides defaults, mocks address server
func MockServer(baseAddr uint64) chan AddrShard {
	// Override operational defaults
	maxAddrOff = 9 // max value of address offset
	const baseIncr uint64 = 1024
	addrShard := new(AddrShard)
	addrShard.addr = baseAddr
	getBaseAddrServe = func(url string) (AddrShard, error) {
		addrShard.addr += baseIncr
		addrShard.shard = (addrShard.shard + 1) % Nshard
		return *addrShard, nil
	}
	chAddr := make(chan AddrShard)
	go getAddr("", chAddr)
	return chAddr
}

// Used for diagnostics when getting failures with 'go test'
func testAddr() {
	chAddr := MockServer(0)
	// Sleep time to avoid race conditions between channels
	const sleepMs = 10
	const nIter = 30
	for i := 0; i < nIter; i++ {
		addrShard := <-chAddr
		fmt.Print(addrShard.addr, ", ")
		time.Sleep(sleepMs * time.Millisecond)
	}
}
