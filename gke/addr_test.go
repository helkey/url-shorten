// Addr_test
// go test addr_test.go addr.go encode.go network.go

package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var result = []uint64{1032, 1026, 1027, 1024, 1029, 1031, 1025, 1030, 1033, 1028,
	2053, 2056, 2055, 2054, 2057, 2051, 2048, 2052, 2049, 2050,
	3079, 3078, 3081, 3073, 3072, 3075, 3076, 3080, 3077, 3074}

func TestAddr(t *testing.T) {
	addrShard, _ := sepAddrShard("123456/7")
	assert.Equal(t, 123456, int(addrShard.addr))
	assert.Equal(t, 7, int(addrShard.shard))

	addrShard, err := baseAddrFromServer(UrlAddrServer)
	fmt.Println("addr", addrShard.addr, "shard", addrShard.shard, err)
	// LOOK at MockServer, run off real server instead

	rand.Seed(0) // const seed for repeatible test results
	chAddrM := MockServer(0)
	const sleepMs = 10 // avoid race condition between channels
	const nIter = 30
	for i := 0; i < nIter; i++ {
		addrShard := <-chAddrM
		assert.Equal(t, addrShard.addr, result[i])
		time.Sleep(sleepMs * time.Millisecond)
	}
}
