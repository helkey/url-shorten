// dbUrl_test.go
// go test dbUrl_test.go dbUrl.go db.go dbDROP.go dbAddr.go addr.go encode.go -args 'password
//  NOTE: go run does not use -args!!!
package main

import (
	"fmt"

	"github.com/stretchr/testify/assert"
	"testing"
)

func not_used() {
	shard := 0
	dB, err := InitUrlTable(shard)
	if err != nil {
		fmt.Println("InitUrlTable:", err)
	}

	fullUrl, addr, randExt, nChar := "http://dbUrl_test.go", uint64(123456789), 314159, 3
	err = dB.SaveUrlDB(fullUrl, addr, randExt, nChar)
	if err != nil {
		fmt.Println("SaveUrlDB:", err)
	}

	fullUrl, randExt, nChar, err = dB.ReadUrlDB(addr)
	if err != nil {
		fmt.Println("ReadUrlDB:", err)
	} else {
		fmt.Printf("fullUrl=%s, randExt=%v, nChar=%v \n", fullUrl, randExt, nChar)
	}

	shortUrl, err := dB.getShortUrl(fullUrl, shard)
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
	shard := 0
	dB, err := InitUrlTable(shard)
	assert.Equal(t, nil, err)

	fullUrl, addr, randExt, nChar := "http://dbUrl_test.go", uint64(1234567), 314159, 12
	err = dB.SaveUrlDB(fullUrl, addr, randExt, nChar)
	assert.Equal(t, nil, err)

	shortUrl, err := dB.getShortUrl(fullUrl, shard)
	assert.Equal(t, nil, err)
	assert.Equal(t, "axOE00008m0Kx", shortUrl)

	wrongUrl := "http://wrong.com"
	_, err = dB.getShortUrl(wrongUrl, shard)
	assert.NotEqual(t, nil, err)
}
