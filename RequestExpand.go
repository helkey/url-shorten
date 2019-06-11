// RequestExpand
// go run RequestExpand.go addr.go db.go dbAddr.go dbUrl.go dbDrop.go encode.go network.go 'passwd'
// localhost:8090/L6X000000bmG
package main

import (
	"fmt"
	"log"
	"net/http"
)

const urlLen = NcharA + NcharR
const urlLongLen = NcharA + NcharRLong

var validLenShortUrl []int

func init() {
	validLenShortUrl = []int{urlLen, urlLongLen}
}

func main() {
	http.HandleFunc("/", expandHandler)
	log.Fatal(http.ListenAndServe(UrlExpand, nil))
}

func expandHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[1:]
	fmt.Println("expand: ", shortUrl)
	fmt.Fprint(w, "expand: ", shortUrl)

	// Check shortURL length matches all historically valid values
	fmt.Println(len(shortUrl), urlLen, urlLongLen)
	if (len(shortUrl) != urlLen) && (len(shortUrl) != urlLongLen) {
		// log.Fatal("RequestExpand err: invalid shortened URL length: ", shortUrl)
		fmt.Fprint(w, "Error - invalid shortened URL")
		return
	}
	// Decode short URL components
	decodeA, decodeR, shard := DecodeURL("oxABCabs0123") // randSlice=1521
	if shard >= Nshard {
		log.Fatal("RequestExpand error: invalid DB shard", shortUrl)
		fmt.Fprint(w, "Error - invalid DB shard")
		return
	}

	// Lookup randExt and fullURL (given database shard)
	dB, err := OpenUrlDB(shard, password())
	if err != nil {
		fmt.Fprint(w, "Error accessing URL database")
		return
	}
	defer dB.db.Close()

	fullURL, randDB, err := dB.ReadUrlDB(decodeA)
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
	// **Redirect to decoded fullURL***
	return
}

// go run RequestShorten.go addr.go encode.go // RequestAddr.go
func TestExpand() error {
	// Initialize default http.Request, http.ResponseWriter objects
	// expandHandler(w, r)
	// Check w
	return nil
}
