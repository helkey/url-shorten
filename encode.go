// Encode

package main

import (
	"errors"
	"fmt"
	// "github.com/stretchr/testify/assert"
	"math/rand"
	"net"
	"net/url" // url.Parse
	"testing"
	"time"
)

var grayList = map[string]struct{}{"box.com":{}, "dropbox.com":{}, "googlemaps.com":{}}

// Number of chars for addr field
//   (may be increased in stages if service becomes popular and larger URL space is needed)
const charAddr = 10 // number chars for encoded address
const charRand = 2 // number additional random chars
const charRandLong = 4 // additional random chars for higher security
const nShard = 8 // number of key-val database shards

// Character encoding
const nDig = 10 // number of digits
const nLett = 26 // number of letters
const nChar = 2 * nLett + nDig // number of characters used for encoding

// Round down max rand integer to avoid overflow when ORing with shard number
var maxRand = ((pow(nChar, charRand)- 1) / nShard) * nShard
var maxRandLong = ((pow(nChar, charRandLong)- 1) / nShard) * nShard


func main() {
	rand.Seed(time.Now().UnixNano()) // pick random seed
	// test()
	en := "ABCabs0123"
	s, _ := EncodeURL("https://goog.com", invertEncode(en), 1)
	fmt.Println(s)
	s, _ = EncodeURL("https://dropbox.com", invertEncode(en), 1)
	fmt.Println(s)
}

// Encode ULR string by with counter and database shard
func EncodeURL(longURL string, address, iShard uint64) (string, error) {
	encodeA, err := encodeAddr(address, charAddr)
	if err != nil {
		return "", err
	}
	if len(encodeA) != charAddr {
		return "", errors.New("Encoded base address wrong length")
	}

	// Check for gray-listed domains
	lengthen := urlGrayListed(longURL)
	charR := charRand

	maxR := maxRand
	if lengthen {
		charR = charRandLong
		maxR = maxRandLong
	}

	// String extension with rand number & shard ID
	fmt.Println(uint64(rand.Intn(maxR)))
	fmt.Println(nShard, ^(uint64(nShard)), ^(uint64(0)))
	randExt := (uint64(rand.Intn(maxR)) & ^(uint64(nShard))) | iShard
	encodeR, err := encodeAddr(randExt, charR)
	if err != nil {
		return "", err
	}
	if len(encodeR) != charR {
		return "", errors.New("Encoded random extension wrong length")
	}
	return encodeR + encodeA, nil
}

// Generate rand string of encoded characters of specified length
func randString(strLen int) string {
	sRand := ""
	for i := 0; i < strLen; i++ {
		sRand = sRand + string(numChar[rand.Intn(nChar)])
	}
	return sRand
}

// Sensitive domains to be encoded with longer shortened URL
// https://gobyexample.com/url-parsing
func urlGrayListed(longURL string) bool {
	u, err := url.Parse(longURL)
	if err != nil {
		// Failed parse - assume domain not gray-listed
		return false
	}
	host := u.Host
	if u.Port() != "" {
		host, _, _ = net.SplitHostPort(host)
	}
	// Strip "www" subdomains
	if (len(host) > len("www")) && host[:3] == "www" {
		host = host[4:]
	}
	// Check if domain is gray-listed
	// fmt.Println("host:", host)
	if _, inList := grayList[host]; inList {
		// fmt.Println("Found in gray-list")
		return true // use extended shortening length
	}
	// fmt.Println("Domain not gray-listed")
	return false
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


// Integer power: compute a**b
// Donald Knuth, The Art of Computer Programming, Volume 2, Section 4.6.3
func pow(a, b int) int {
        pow := 1
        for b > 0 {
                if (b & 1) != 0 {
                        pow *= a
                }
                b >>= 1
                a *= a
        }
        return pow
}


func TestEncode(t *testing.T) {
	const maxCharEncode = 10 // max characters for uint64 representation
	if charAddr> maxCharEncode {
		t.Errorf("charEncode Exceeded max chars for uint64 representation")
	}
}

func test() {
	en := "ABCabs0123"
	inv := invertEncode(en)
	encoded, _ := encodeAddr(inv, len(en))
	fmt.Println(encoded)
}
