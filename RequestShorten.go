// RequestShorten
// go run RequestShorten.go addr.go encode.go // RequestAddr.go

package main

import (
	"fmt"
	"log"
	"net/http"
)

const UrlShorten = "localhost:8086"    // 12.0.0.1 (IPv6 ::1)
const UrlAddrServer = "127.0.0.1:8088" // (IPv6 ::1)
var chAddr = make(chan BaseShard)

func main() {
	TestShorten()
	return
	
	// Set up channel to supply channel addresses
	go getAddr(UrlAddrServer, chAddr)

	http.HandleFunc("/create", shortenHandler)
	log.Fatal(http.ListenAndServe(UrlShorten, nil))
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Fprintf(w, "shorten", path)
	return
	if (r.Method != "GET") {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	const header = "?source=&url="
	if (len(path) <= len(header)) || (path[:len(header)] == header) {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	addrShard := <-chAddr
	addr := addrShard.addr
	fullURL := path[:len(header)]
	shortURL, randExt, err := EncodeURL(fullURL, addr, addrShard.shard) // encode.go
	randExt = 0
	if err != nil {
		fmt.Fprintf(w, "err shortening URL")
	} else {
		fmt.Fprintf(w, shortURL, randExt)
		// **PUT IN DBase SHARD** addr:?/fullURL; addr:?/randExt
	}
}

// go run RequestShorten.go addr.go encode.go // RequestAddr.go
func TestShorten() error {
	urls := []string{"Shorten This", "and THIS"}
	_, decodeA, _ := DecodeURL("8765431Kn")
	chAddrM := MockServer(decodeA)
	for _, url := range urls {
		addrShard := <-chAddrM
		shortURL, _, err := EncodeURL(url, addrShard.addr, addrShard.shard)
		if err != nil {
			return err
		}
		// Recover shard, compare to specification
		randURL, baseURL := shortURL[:NcharR], shortURL[NcharR:]
		fmt.Printf("'%s'  %s  %s  %s \n", url, shortURL, randURL, baseURL)
		dR, dA, iShard := DecodeURL(shortURL)
		fmt.Printf("rand:%b  base:%b  shard:%d \n", dR, dA, iShard) 
	}
	return nil
}

// e.g. tinyurl
// var testS = UrlShorten + "/create?source=&url=https%3A%2F%2Fwww.amazon.com%2Fhorsebattery"
