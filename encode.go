// Encode

package main

import (
	"errors"
	"fmt"
	// "github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"net/url"
	"testing"
)

// Number of chars for addr field
//   (may be increased in stages if service becomes popular and larger URL space is needed)
const charEncode = 10 // number chars for encoded address
const charEncodeLong = 14 // number chars for long encoded address
const charExtend = charEncodeLong - charEncode // additional chars for higher security

// Character encoding
const nDig = 10 // number of digits
const nLett = 26 // number of letters
const nChar = 2 * nLett + nDig // number of characters used for encoding


func main() {
	test()
}

// Encode ULR string by with counter and database shard
func EncodeURL(longURL string, address, iShard uint64, nAddr uint) (string, error) {
	urlEncode := address | iShard // include database shard
	encoded, err := encodeAddr(urlEncode, charEncode)
	if err != nil {
		return "", err
	}
	if len(encoded) != charEncode {
		return "", errors.New("Encoded URL wrong length")
	}

	// Check for gray-listed domains
	lengthen, err := urlGrayListed(longURL)
	if (err != nil) || !lengthen {
		return encoded, nil
	}
	sRand := randString(charExtend)
	return sRand + encoded, nil
}

// Generate rand string of encoded characters of specified length
func randString(strLen int) string {
	sRand := ""
	for i := 0; i < strLen; i++ {
		sRand = sRand + string(numChar[rand.Intn(nChar)])
	}
	return sRand
}

var grayList = map[string]struct{}{"box.com":{}, "dropbox.com":{}, "googlemaps.com":{}}
// Sensitive domains to be encoded with longer shortened URL
// https://gobyexample.com/url-parsing
func urlGrayListed(longURL string) (bool, error) {
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

func num2char() []byte {
	numChar := make([]byte, nChar)
	for i:=0; i < nLett; i++ {
		if i < nDig {
			numChar[i] = '0' + byte(i)
		}
		numChar[i + nDig] = 'a' + byte(i)
		numChar[i + nDig + nLett] = 'A' + byte(i)

	}
	return numChar
}
var numChar = num2char()

// Convert binary value to 
func encodeAddr(address uint64, nChars int) (string, error) {
	encoded := ""
	// Convert 'address' to base numChar;
	//   convert each digit to character representation
	for i := 0; i < nChars; i++  {
		charVal := address % nChar
		address = address / nChar
		// Construct 'encoded' string from right to left
		encoded = string(numChar[charVal]) + encoded
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
func invertEncode(encoded string) uint64 {
	addr := uint64(0)
	fmt.Print("invertEncode ", encoded, ": " )
	for _, char := range encoded {
		// char := encoded[i]
		addr = nChar * addr + uint64(charNum[byte(char)])
		fmt.Print(charNum[byte(char)], ", ")
	}
	fmt.Println(";  ", addr, (uint64(1)<<63)/addr)
	return addr
}

// Invert numChar mapping for generating test vectors
func invertCharNum() map[byte]uint {
	charNum := map[byte]uint{}
	for i, b := range numChar {
		charNum[b] = uint(i)
	}
	return charNum
}
var charNum = invertCharNum()


func TestEncode(t *testing.T) {
	const maxCharEncode = 10 // max characters for uint64 
	if charEncode > maxCharEncode {
		t.Errorf("charEncode Exceeded max chars for uint64 representation")
	}
}

func test() {
	en1 := "ABCabs0123"
	inv := invertEncode(en1)
	encoded, _ := encodeAddr(inv, len(en1))
	fmt.Println(encoded)
}
