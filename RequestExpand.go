// RequestExpand
// go run RequestExpand.go addr.go dbAddr.go encode.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const UrlExpand = "localhost:8090" // 12.0.0.1 (IPv6 ::1)
const urlLen = NcharA + NcharR + Nshard
const urlLongLen = NcharA + NcharRLong + Nshard


var passwd = ""
var validLenShortUrl []int
var dbS DBS

func main() {
	fmt.Println(os.Args[1])

	validLenShortUrl = []int{urlLen, urlLongLen}
	dbS.shard = 1 << 31 // initialize to unused value
	
	// http.HandleFunc("/create", expandHandler)
	http.HandleFunc("/", expandHandler)
	log.Fatal(http.ListenAndServe(UrlExpand, nil))
}

func expandHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	fmt.Println("expand: ", path)
	fmt.Fprintf(w, "expand: ", path)
	return
	if r.Method != "GET" {
		http.Error(w, "405 method not allowed.", http.StatusMethodNotAllowed)
		return
	}
	const header = "?source=&url=" // header before long URL argument
	if (len(path) <= len(header)) || (path[:len(header)] == header) {
		// Argument doesn't contain URL
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}
	shortUrl := path[:len(header)]
	// Check shortURL length matches all historically valid values
	if (len(shortUrl) != urlLen) && (len(shortUrl) != urlLongLen) {
		log.Fatal("RequestExpand err: invalid shortened URL length", shortUrl)
		fmt.Fprintf(w, "Error - invalid shortened URL")
		return
	}

	// Lookup randExt and fullURL (given database shard)
	decodeA, decodeR, shard := DecodeURL("oxABCabs0123") // randSlice=1521
	// Check shard valid
	fullURL, randDB, err := dbS.ReadUrlDB(decodeA, shard, passwd)
	// Check saved randExt matches supplied rand extension
	if err != nil {
		log.Fatal("expandHandler: error expannding URL", shortUrl)
		fmt.Fprintf(w, "Error - full URL not found")
		return
	}
	if randDB != decodeR {
		log.Fatal("expandHandler: random extension not matched", shortUrl)
		fmt.Fprintf(w, "Error - full URL not found")
		return
	}
	fmt.Fprintf(w, fullURL)
	// Redirect to decoded fullURL
	return
}

// go run RequestShorten.go addr.go encode.go // RequestAddr.go
func TestExpand() error {
	return nil
}

// e.g. tinyurl
// var testS = UrlShorten + "/create?source=&url=https%3A%2F%2Fwww.amazon.com%2Fhorsebattery"
