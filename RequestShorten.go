// RequestShorten
// go run RequestShorten.go addr.go db.go dbUrl.go encode.go network.go 'passwd

// "localhost:8086/create" (working WSL)
// "localhost:8086/create/?source=&url=http://FullURL"

package main

import (
	"fmt"
	"log"
	"net/http"
)

var chAddrSh = make(chan AddrShard)

func main() {
	// Set up channel to supply channel addresses
	fmt.Println("RequestShorten: go chan 'getAddr'")
	go getAddr(UrlAddrServer, chAddrSh)

	fmt.Println("RequestShorten/create")
	http.HandleFunc("/create/", shortenHandler)
	log.Fatal(http.ListenAndServe(UrlShorten, nil))
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	// const header = "create/?source=&url=" // header before full URL
	path := r.URL.Path

	if path == "/create/" {
		// Argument doesn't contain valid URL
		e := "RequestShorten - invalid request: " + path
		fmt.Println(e)
		// log.Println(e)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	fmt.Fprint(w, "Invalid API")
	fullUrl := "DEVO"
	return

	addrShard := <-chAddrSh
	addr := addrShard.addr
	shard := addrShard.shard
	shortUrl, randExt, err := EncodeURL(fullUrl, addr, shard) // encode.go
	if err != nil {
		log.Fatal("shortenHandler: error shortinging URL", fullUrl)
		fmt.Fprint(w, "error shortening URL")
		return
	}

	dB, err := OpenUrlDB(shard, password())
	if err != nil {
		fmt.Fprintf(w, "Error accessing URL database")
		return
	}
	defer dB.db.Close()
	err = dB.SaveUrlDB(fullUrl, addr, randExt)
	if err != nil {
		fmt.Fprintf(w, "Error storing shortened URL")
	}
	fmt.Fprintf(w, shortUrl, randExt)
	return
}
