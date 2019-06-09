// shorten_test.go
// go run shorten_test.go addr.go encode.go // RequestAddr.go

package main

import (
	"fmt"
	"log"
	"net/http"
)

func TestShorten() error {
	urls := []string{"Shorten This", "and THIS"}
	decodeA, _, _ := DecodeURL("8765431Kn")
	chAddrM := MockServer(decodeA)
	for _, url := range urls {
		addrShard := <-chAddrM
		shortUrl, _, err := EncodeURL(url, addrShard.addr, addrShard.shard)
		if err != nil {
			return err
		}
		// Recover shard, compare to specification
		randUrl, baseUrl := shortUrl[:NcharR], shortUrl[NcharR:]
		fmt.Printf("'%s'  %s  %s  %s \n", url, shortUrl, randUrl, baseUrl)
		dA, dR, shard := DecodeURL(shortUrl)
		fmt.Printf("rand:%b  base:%b  shard:%d \n", dR, dA, shard)
	}
	return nil
}
