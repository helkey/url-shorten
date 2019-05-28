// Encode

package main

import (
	"errors"
	"fmt"
	"math/rand"
	"net"
	"net/url"
	"time"
)

var grayList = map[string]struct{}{"box.com": {}, "dropbox.com": {}, "googlemaps.com": {}}

// Number of chars for addr field
//   (may be increased in stages if service becomes popular and larger URL space is needed)
const NcharA = 9     // number chars for encoded address
const NcharR = 3     // number additional random chars
const NcharRLong = 5 // additional random chars for higher security
const NshardBits = 3 // number of bits used to encode database shard number
const Nshard = 8     // (1 << NshardBits) number of key-val database shards>

// Character encoding
const Ndig = 10              // number of digits
const Nlett = 26             // number of letters
const Nchar = 2*Nlett + Ndig // number of characters used for encoding

// Round down max rand integer to avoid char conversion overflow when adding shard number
var MaxRand int = ((pow(Nchar, NcharR) - 1) / Nshard)
var MaxRandLong int = ((pow(Nchar, NcharRLong) - 1) / Nshard)

func init() {
	rand.Seed(time.Now().UnixNano()) // pick random seed
}

//func main() {
// TestEncode()
//}

// Encode ULR string with base address, random address, and database shard
func EncodeURL(fullURL string, baseAddr uint64, iShard uint32) (string, uint32, error) {
	encodeA, err := encodeAddr(baseAddr, NcharA)
	fmt.Println("ENCODEURL baseAddr:", baseAddr, "iShard:", iShard)
	if err != nil {
		return "", 0, err
	}
	if len(encodeA) != NcharA {
		return "", 0, errors.New("Encoded base address wrong length")
	}

	// Check for gray-listed domains
	lengthen := urlGrayListed(fullURL)
	charR := NcharR
	maxR := MaxRand
	if lengthen {
		charR = NcharRLong
		maxR = MaxRandLong
	}

	// String extension with rand number & shard ID
	// random extension; before conversion to char
	randExt := uint64(rand.Intn(maxR))
	randShard := (randExt << NshardBits) | uint64(iShard)
	// fmt.Printf("randExt:%d %b  randShard:%b \n", randExt, randExt, randShard)
	encodeR, err := encodeAddr(randShard, charR)
	// fmt.Println("encodeR:", encodeR, "err:", err)
	if err != nil {
		return "", 0, err
	}
	if len(encodeR) != charR {
		return "", 0, errors.New("Encoded random extension wrong length")
	}
	shortURL := encodeR + encodeA
	// fmt.Println("shortURL:", shortURL, "encodeR:", encodeR, "encodeA:", encodeA)
	return shortURL, uint32(randExt), nil
}

func DecodeURL(shortURL string) (uint64, uint32, uint32) {
	// fmt.Println("shortURL:", shortURL)
	lenExt := len(shortURL) - NcharA
	// split shortURL into extension and base address
	encodeR, encodeA := shortURL[:lenExt], shortURL[lenExt:]
	decodeRS := decode(encodeR)
	decodeR := uint32(decodeRS >> NshardBits)     // random value
	iShard := uint32(decodeRS & uint64(Nshard-1)) // database shard
	decodeA := decode(encodeA)
	// fmt.Println(encodeR, encodeA, decodeR, iShard)
	return decodeA, decodeR, iShard
}

// Generate rand string of encoded characters of specified length
func randString(strLen int) string {
	sRand := ""
	for i := 0; i < strLen; i++ {
		sRand = sRand + string(numChar[rand.Intn(Nchar)])
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
	numChar := make([]byte, Nchar)
	for i := 0; i < Nlett; i++ {
		if i < Ndig {
			numChar[i] = '0' + byte(i)
		}
		numChar[i+Ndig] = 'a' + byte(i)
		numChar[i+Ndig+Nlett] = 'A' + byte(i)

	}
	return numChar
}

var numChar = num2char()

// Convert binary value to
func encodeAddr(address uint64, nChars int) (string, error) {
	encoded := ""
	// Convert 'address' to base numChar;
	//   convert each digit to character representation
	for i := 0; i < nChars; i++ {
		charVal := address % Nchar
		address = address / Nchar
		// Construct 'encoded' string from right to left
		encoded = string(numChar[charVal]) + encoded
	}
	// return "", errors.New("URL encode failed")
	return encoded, nil
}

// Invert address encoding process for generating test vectors
func decode(encoded string) uint64 {
	addr := uint64(0)
	// fmt.Print("decode:", encoded, "; ")
	for _, char := range encoded {
		// char := encoded[i]
		addr = Nchar*addr + uint64(charNum[byte(char)])
		// fmt.Print(charNum[byte(char)], ", ")
	}
	// fmt.Println(";  ", addr, (uint64(1)<<63)/addr)
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

// Test encoding/decoding functions
//func  TestEncode(t *testing.T) {
/*
func  TestEncode() {
	en, iShard:= "ABCabs012", uint32(0)
	s, _, _ := EncodeURL("https://goog.com", decode(en), iShard)
	// assert.Equal(en, s[len(s)-len(en):])
	i := len(s) - len(en)
	decodeRS := decode(s[:i])
	shard := uint32(decodeRS & uint64(Nshard-1))
	fmt.Println(en, s[i:], iShard, shard)
	fmt.Println()

	iShard = 7
	s, _, _ = EncodeURL("https://dropbox.com", decode(en), iShard)
	i = len(s) - len(en)
	decodeRS = decode(s[:i])
	shard = uint32(decodeRS & uint64(Nshard-1))
	fmt.Println(en, s[i:], iShard, shard)
} */
