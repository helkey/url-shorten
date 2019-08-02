// RequestExpand
// go run RequestExpand.go addr.go db.go dbAddr.go dbUrl.go dbDrop.go encode.go network.go 'passwd'
//   localhost:8090/L6X000000bmG

package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

const urlLen = NcharA + NcharR
const urlLongLen = NcharA + NcharRLong

// var validLenShortUrl []int

func init() {
}

func main() {
	if INITIALIZEDB {
		longUrl, err := expandUrl("ejK0000A86RV")
		fmt.Println("ReqExp: long=", longUrl, err)
		return
	}

	shard := 0
	dB, _ := OpenUrlDB(shard, dbPassword())
	nRows, _ := dB.NumRowsDB("url")
	fmt.Printf("DB 'url' HAS %d rows\n", nRows)

	http.HandleFunc("/", expandHandler)
	log.Fatal(http.ListenAndServe(PortExpand, nil))
}

func expandHandler(w http.ResponseWriter, r *http.Request) {
	shortUrl := r.URL.Path[1:]
	fmt.Println("expand: ", shortUrl)
	longUrl, err := expandUrl(shortUrl)
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		// TODO: Redirect!!
		fmt.Fprint(w, longUrl)
	}
}

func expandUrl(shortUrl string) (longUrl string, err error) {
	// Check shortURL length matches all historically valid values
	fmt.Printf("len shortUrl:%v; min:%v;  max:%v\n", len(shortUrl), urlLen, urlLongLen)
	// validLenShortUrl = []int{urlLen, urlLongLen}
	if (len(shortUrl) != urlLen) && (len(shortUrl) != urlLongLen) {
		return "", errors.New("Error - invalid shortened URL length")
	}
	// Decode short URL components
	// decodeA, decodeR, shard := DecodeURL("oxABCabs0123") // randSlice=1521
	decodeA, decodeR, shard := DecodeURL(shortUrl)
	// decodeA, decodeR = 533881127, 6880
	fmt.Printf("decodeA:%v;  decodeR:%v;  shard:%v\n", decodeA, decodeR, shard)
	if shard >= Nshard {
		// log.Fatal("RequestExpand error: invalid DB shard", shortUrl)
		return "", errors.New("Error - invalid DB shard")
	}

	// Lookup randExt and fullURL (given database shard)
	dB, err := OpenUrlDB(shard, dbPassword())
	if err != nil {
		return "", errors.New("Error accessing URL database")
	}
	defer dB.db.Close()

	fullUrl, randDB, nChar, err := dB.ReadUrlDB(decodeA)
	fmt.Printf("fullUrl:%s, randDB:%v, nChar:%v\n", fullUrl, randDB, nChar)
	if err != nil {
		// log.Fatal("RequestExpand: error expanding URL: ", shortUrl)
		return "", errors.New("Error - shortened URL not found")
	}

	if randDB != decodeR {
		// log.Fatal("expandHandler: random extension not matched", shortUrl)
		return "", errors.New("Error decodeR - shortened URL not found")
	}
	return fullUrl, nil
}
