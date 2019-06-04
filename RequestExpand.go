// RequestExpand
// go run RequestExpand.go addr.go dbAddr.go encode.go 'passwd'
// localhost:8090/L6X000000bmG
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const UrlExpand = "localhost:8090" // 12.0.0.1 (IPv6 ::1)
const urlLen = NcharA + NcharR
const urlLongLen = NcharA + NcharRLong
const shardDefault = 1 << 31

var passwd = ""
var validLenShortUrl []int
var dbS DBS


func init() {
	validLenShortUrl = []int{urlLen, urlLongLen}
	dbS.shard = shardDefault
}

func main() {
	// http.HandleFunc("/create", expandHandler)
	http.HandleFunc("/", expandHandler)
	log.Fatal(http.ListenAndServe(UrlExpand, nil))
}

func expandHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[1:]
	fmt.Println("expand: ", shortUrl)
	fmt.Fprintf(w, "expand: ", shortUrl)

	// Check shortURL length matches all historically valid values
	fmt.Println(len(shortUrl), urlLen, urlLongLen)
	if (len(shortUrl) != urlLen) && (len(shortUrl) != urlLongLen) {
		// log.Fatal("RequestExpand err: invalid shortened URL length: ", shortUrl)
		fmt.Fprintf(w, "Error - invalid shortened URL")
		return
	}
	// Decode short URL components
	decodeA, decodeR, shard := DecodeURL("oxABCabs0123") // randSlice=1521
	if shard >= Nshard {
		log.Fatal("RequestExpand error: invalid DB shard", shortUrl)
		fmt.Fprintf(w, "Error - invalid DB shard")
		return
	}
	
	// Lookup randExt and fullURL (given database shard)
	passwd := os.Args[1]
	fullURL, randDB, err := dbS.ReadUrlDB(decodeA, shard, passwd)
	// Check saved randExt matches supplied value
	if err != nil {
		// log.Fatal("RequestExpand: error expanding URL: ", shortUrl)
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
	// Initialize default http.Request, http.ResponseWriter objects
	// expandHandler(w, r)
	// Check w
	return nil
}


