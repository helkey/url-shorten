// RequestShorten
// go run RequestShorten.go addr.go db.go dbAddr.go dbDROP.go dbUrl.go encode.go network.go 'passwd

// localhost:8086/create  // (working WSL)
// localhost:8086/create/?source=&url=http://FullURL

// TODO
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var chAddrSh = make(chan AddrShard)

func main() {
	// Set up channel to supply channel addresses
	// fmt.Println("RequestShorten: go chan 'getAddr'")
	go getAddr(UrlAddrServer, chAddrSh)

	rand.Seed(time.Now().UnixNano())
	if INITIALIZEDB {
		rand.Seed(0)
		shard := 0
		InitUrlTable(shard)
		shortUrl, err := shortenUrl(testUrl)
		fmt.Println(shortUrl, err)
		return
	}

	fmt.Println("ReqShorten: listening")
	http.HandleFunc("/create/", shortenHandler)
	log.Fatal(http.ListenAndServe(UrlShorten, nil))
}

const shortUrlBase = "http://base.com/"

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// const header = "create/?source=&url=" // header before full URL
	keys, ok := r.URL.Query()["url"] // also: r.URL.Path
	if !ok {
		fmt.Println("ReqShort: key not found")
		// log.Println(e)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	fullUrl := keys[0]

	var shortUrl string
	var exists bool
	if shortUrl, exists = existsShort(fullUrl); !exists {
		fmt.Println("ReqShort: gen new short URL")
		shortNew, errMsg := shortenUrl(fullUrl)
		if errMsg != "" {
			fmt.Println("ReqShort err:", errMsg)
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}
		fmt.Fprintf(w, shortNew)
		return
	}
	fmt.Println("ReqShort: reuse existing shortened URL")
	fmt.Fprintf(w, shortUrl)
}

// Check if long URL already in database,
//   return false if can't access DB
func existsShort(fullUrl string) (shortUrl string, exists bool) {
	shard := 0
	dB, err := OpenUrlDB(shard, password())
	if err != nil {
		return "", false
	}
	defer dB.db.Close()

	addr, randExt, nChar, err := dB.CheckUrlDB(fullUrl)
	if err != nil {
		return "", false
	}
	shortUrl, err = encode(addr, randExt, shard, nChar)
	if err != nil {
		return "", false
	}
	return shortUrl, true
}

func shortenUrl(fullUrl string) (shortUrl string, errMsg string) {
	fmt.Println("ReqShort: *" + fullUrl + "*")
	// Get unique shortened address
	addrShard := <-chAddrSh
	addr := addrShard.addr
	shard := addrShard.shard

	// Generate shortened URL using address and database shard
	shortUrl, randExt, nChar, err := EncodeURL(fullUrl, addr, shard)
	if err != nil {
		return "", "Error shortening URL"
	}
	fmt.Printf("ReqShort: addr=%v;  shard=%v  shortUrl=%v\n", addr, shard, shortUrl)
	// "http://FullURL";  base=526058514; addr=533881127; shard=0  => shortUrl=ejK0000A86RV, randext=6880
	// (old) "http://FullURL"; base=619732968; addr=626419234; shard=0  => shortUrl=e4o0000Goog2

	dB, err := OpenUrlDB(shard, password())
	if err != nil {
		return "", "Error accessing URL database"
	}
	defer dB.db.Close()
	err = dB.SaveUrlDB(fullUrl, addr, randExt, nChar)
	if err != nil {
		return "", "Error storing shortened URL"
	}

	nRows, _ := dB.NumRowsDB("url")
	fmt.Printf("DB 'url' HAS %d rows\n", nRows)
	return shortUrl, ""
}

const testUrl = "http://FullURL"
