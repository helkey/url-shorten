// dbUrl_test.go
// go test dbUrl_test.go dbUrl.go db.go dbDROP.go dbAddr.go addr.go encode.go -args 'password
//  NOTE: go run does not use -args!!!
package main

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"testing"
)

const shard = 0
const 	fullUrl = "http://dbUrl_test.go"
const 	addr = uint64(123456789)
const 	randExt = 314159
const 	nChar = 3

// For testing with $go run  (instead of $go test)
func For_Go_run_mode() {
	dB, err := InitUrlTable(shard)
	if err != nil {
		fmt.Println("InitUrlTable:", err)
	}

	err = dB.SaveUrlDB(fullUrl, addr, randExt, nChar)
	if err != nil {
		fmt.Println("SaveUrlDB:", err)
	}

	fullUrl_new, randExt_new, nChar_new, err := dB.ReadUrlDB(addr)
	if err != nil {
		fmt.Println("ReadUrlDB:", err)
	} else {
		fmt.Printf("fullUrl=%s, randExt=%v, nChar=%v \n", fullUrl_new, randExt_new, nChar_new)
	}

	shortUrl, err := dB.getShortUrl(fullUrl_new, shard)
	if err != nil {
		fmt.Println("getShortUrl:", err)
	} else {
		fmt.Println("shortUrl:", shortUrl)
	}

	wrongUrl := "http://wrong.com"
	wrongShortUrl, err := dB.getShortUrl(wrongUrl, shard)
	if err != nil {
		fmt.Println("getShortUrl:", err)
	} else {
		fmt.Println("wrongShortUrl:", wrongShortUrl)
	}
}

func TestSaveUrl(t *testing.T) {
	dB, err := InitUrlTable(shard)
	assert.Equal(t, nil, err)

	err = dB.SaveUrlDB(fullUrl, addr, randExt, nChar)
	assert.Equal(t, nil, err)

	shortUrl, err := dB.getShortUrl(fullUrl, shard)
	assert.Equal(t, nil, err)
	assert.Equal(t, "xOE00008m0Kx", shortUrl)

	wrongUrl := "http://wrong.com"
	_, err = dB.getShortUrl(wrongUrl, shard)
	assert.NotEqual(t, nil, err)
}
