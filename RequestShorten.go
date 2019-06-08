// RequestShorten
// go run RequestShorten.go addr.go dbUrl.go encode.go network.go 'passwd

// "localhost:8086/create" (working WSL)
// "localhost:8086/create/?source=&url=http://FullURL"

package main

import (
	"fmt"
	"log"
	"net/http"
)

var chAddrSh = make(chan AddrShard)

/* func main() {
	// TestShorten()
	// return

	// Set up channel to supply channel addresses
	// fmt.Println("RequestShorten: go chan 'getAddr'")
	go getAddr(UrlAddrServer, chAddrSh)
	// dbS.shard = 1 << 31 // initialize to unused value

	fmt.Println("RequestShorten/create")
	http.HandleFunc("/", shortenHandler)
	log.Fatal(http.ListenAndServe(UrlShorten, nil))
} */

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	const header = "create/?source=&url=" // header before full URL
	path := r.URL.Path
	fmt.Fprint(w, "shorten", path)
	if (len(path) <= len(header)) || (path[:len(header)] != header) {
		// Argument too short to contain URL
		e := "RequestShorten - invalid request: " + path
		fmt.Println(e)
		// log.Println(e)
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	fullUrl := path[:len(header)]
	fmt.Fprint(w, "shorten ", path, fullUrl)
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

// go run RequestShorten.go addr.go encode.go // RequestAddr.go
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

// e.g. tinyurl
// var testS = UrlShorten + "/create?source=&url=https%3A%2F%2Fwww.amazon.com%2Fhorsebattery"
