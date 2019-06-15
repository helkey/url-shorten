// shorten_test.go
// go test shorten_test.go addr.go encode.go // RequestAddr.go

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShorten(t *testing.T) {
	urls := []string{"Shorten This", "and THIS"}
	decodeA, _, _ := DecodeURL("8765431Kn")
	chAddrM := MockServer(decodeA)
	fmt.Println("shorten_test")
	// Pull addr from MockServer, as next two addr values are not sequential
	_ = <-chAddrM
	for _, url := range urls {
		addrShard := <-chAddrM
		addr, shard := addrShard.addr, addrShard.shard
		shortUrl, _, err := EncodeURL(url, addrShard.addr, addrShard.shard)
		assert.Equal(t, nil, err)

		// Recover shard, compare to specification
		// randUrl, baseUrl := shortUrl[:NcharR], shortUrl[NcharR:]
		fmt.Printf("'%s'  %s  %v\n", url, shortUrl, addr)
		dA, dR, shard := DecodeURL(shortUrl)
		fmt.Printf("rand:%v  base:%v  shard:%d \n", dR, dA, shard)
		assert.Equal(t, addr, uint64(dA))
	}
}
