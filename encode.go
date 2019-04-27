// Encode

package main

import (
	"errors"
	"fmt"
	"net"
	"net/url"
)

// Number of bits for addr field and number of bits for database shard
//   (may be increased in stages if service becomes popular and
//   larger URL space and more database shards is needed)
const bitsEncode = 12 // num bits for encoded address
const bitsEncodeLong = 15 // num bits for long encoded address
const bitShift = uint(bitsEncodeLong - bitsEncode)
const bitsShard = 3 // number database sharding bits
const charEncode = 10 // number chars for encoded address
const charEncodeLong = 12 // number chars for long encoded address

const numDig = 10 // number of digits
const numLett = 26 // number of letters
const numChar = 2 * numLett + numDig // number of characters used for encoding


func main() {
	// invertEncode("eNcOdEtHiS")
	invertEncode("aazzAAZZ0099")
}

// Encode ULR string by with counter and database shard
func Encode(longURL string, address, nAddr, iShard uint) (string, error) {
	lengthen, err := urlGrayListed(longURL)
	if err != nil {
		return "", err
	}
	urlEncode := address
	nChars := charEncode
	if lengthen {
		rand := uint(0) // random bits to lengthen URL
		urlEncode = (address << bitShift) | rand
		nChars = charEncodeLong
	}
	urlEncode = urlEncode | iShard // include database shard
	encoded, err := encodeURL(urlEncode, nChars)
	if err != nil {
		return "", err
	}
	if len(encoded) != numChar {
		return "", errors.New("Encoded URL wrong length")
	}
	return encoded, nil
}

var grayList = map[string]struct{}{"box.com":{}, "dropbox.com":{}, "googlemaps.com":{}}
// Sensitive domains to be encoded with longer shortened URL
// https://gobyexample.com/url-parsing
func urlGrayListed(longURL string) (bool, error) {
	// fmt.Println(grayList)
	u, err := url.Parse(longURL)
	if err != nil {
		fmt.Println("Failed parse")
		return false, nil // use default shortening length
	}
	host := u.Host
	if u.Port() != "" {
		host, _, _ = net.SplitHostPort(host)
	}
	// Strip "www" subdomains
	if (len(host) > len("www")) && host[:3] == "www" {
		host = host[4:]
	}
	fmt.Println("host:", host)
	// Check if domain is gray-listed
	if _, inList := grayList[host]; inList {
		fmt.Println("Found")
		return true, nil // use extended shortening length
	}
	fmt.Println("Not found")
	return false, nil
}

func bin2char() []byte {
	binChar := make([]byte, numChar)
	for i:=0; i < numLett; i++ {
		binChar[i] = 'a' + byte(i)
		binChar[i + numLett] = 'A' + byte(i)
		if i < numDig {
			binChar[i + 2 * numLett] = '0' + byte(i)
		}
	}
	return binChar
}
var binChar = bin2char()

// Convert binary value to 
func encodeURL(address uint, nChars int) (string, error) {
	encoded := ""
	// Convert 'address' to base numChar;
	//   convert each digit to character representation
	for i := 0; i < nChars; i++  {
		charVal := address % numChar
		address = address / numChar
		encoded = string(binChar[charVal]) + encoded
	}
	// return "", errors.New("URL encode failed")
	return encoded, nil
}

func testGrayListed() {
 	url := "postgres://user:pass@host.com:5432/path?k=v#f"
	// grayListed, err := sensitiveURL(url)
	fmt.Println(urlGrayListed(url))
	url = "https://www.dropbox.com/filename"
	// grayListed, err := urlGrayListed(url)
	fmt.Println()
	fmt.Println(urlGrayListed(url))
}

// Invert address encoding process for generating test vectors
func invertEncode(encoded string) uint {
	if _, ok := charBin['a']; !ok {
		invertCharBin()
	}
	addr := uint(0)
	for i := len(encoded) - 1; i >= 0; i-- {
		addr = numChar * addr + charBin[encoded[i]]
		fmt.Print(charBin[encoded[i]], ", ")
	}
	fmt.Println('\n', addr)
	return addr
}

// Invert binChar mapping for generating test vectors
var charBin = make(map[byte]uint)
func invertCharBin() {
	for i, b := range binChar {
		charBin[b] = uint(i)
	}
}
