// RequestShorten
// go run RequestShorten.go addr.go db.go dbAddr.go dbDROP.go dbUrl.go encode.go network.go 'passwd
// localhost:8088/create/?source=&url=http://FullURL
// http://18.144.27.152:8088/create/?source=&url=http://FullURL

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var chAddrSh = make(chan AddrShard) // Go chan for buffering addr values

func main() {
	// Initialize url database
	if (len(os.Args) > 1) && os.Args[1] == init_str {
		fmt.Println("Initializing URL database")
		InitUrlTable(0)
		InitUrlTable(1)
		return
	}

	// Set up channel to supply channel addresses
	// fmt.Println("RequestShorten: go chan 'getAddr'")
	go getAddr(UrlAddrServer, chAddrSh)

	rand.Seed(time.Now().UnixNano())
	http.HandleFunc("/create/", shortenHandler)

	fmt.Println("ReqShorten: listening")
	log.Fatal(http.ListenAndServe(PortShorten, nil))
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// const header = "create/?source=&url=" // header before full URL
	keys, ok := r.URL.Query()["url"] // also: r.URL.Path
	if !ok {
		fmt.Println("ReqShort: key not found")
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	fullUrl := keys[0]

	shortUrl, err := getShortUrlAllDB(fullUrl)
	if err != nil {
		fmt.Println("ReqShort: Generate new short URL")
		errMsg := ""
		shortUrl, errMsg = shortenUrl(fullUrl)
		if errMsg != "" {
			fmt.Println("ReqShort err:", errMsg)
			http.Error(w, "404 not found.", http.StatusNotFound)
			return
		}
	} else {
		fmt.Println("ReqShort: reuse existing shortened URL")
	}
	_, _, shard := DecodeURL(shortUrl)
	instance := getShortenInstance() // Cloud instance identifier
	shortInfo := fmt.Sprintf("%s shard:%v; i:%s", shortUrl, shard, instance)
	fmt.Fprintf(w, shortInfo)
	fmt.Println(shortInfo)
}

func shortenUrl(fullUrl string) (shortUrl string, errMsg string) {
	// Get unique shortened address
	fmt.Println("Channel addresses in buffer:", len(chAddrSh))	
	addrShard := <-chAddrSh
	fmt.Println("Read channel address")
	addr := addrShard.addr
	shard := addrShard.shard

	// Generate shortened URL using address and database shard
	fmt.Println("ReqShort: *" + fullUrl + "*")
	shortUrl, randExt, nChar, err := EncodeURL(fullUrl, addr, shard)
	if err != nil {
		return "", "Error shortening URL"
	}
	fmt.Printf("ReqShort: addr=%v;  shard=%v  shortUrl=%v\n", addr, shard, shortUrl)

	dB, err := OpenUrlDB(shard, dbPassword())
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
