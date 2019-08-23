// Encode

package main

import (
	"errors"
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

// Test with 2 database shards to minimize resources
const NshardBits = 1 // number of bits used to encode database shard number
const Nshard = 2     // (1 << NshardBits) number of key-val database shards>
// const NshardBits = 3 // number of bits used to encode database shard number
// const Nshard = 8     // (1 << NshardBits) number of key-val database shards>

// Character encoding
const Ndig = 10              // number of digits
const Nlett = 26             // number of letters
const Nchar = 2*Nlett + Ndig // number of characters used for encoding

// Round down max rand integer to avoid char conversion overflow when adding shard number
var MaxRand int = ((pow(Nchar, NcharR) - 1) / Nshard)
var MaxRandLong int = ((pow(Nchar, NcharRLong) - 1) / Nshard)

var Nrange = pow(Nchar, NcharA) >> NaddrBit

func init() {
	rand.Seed(time.Now().UnixNano()) // pick random seed
}


// Encode ULR string with address, random address, and database shard
func EncodeURL(fullUrl string, addr uint64, shard int) (string, int, int, error) {
	isGrayList := UrlGrayList(fullUrl) // Check for gray-listed domains
	nChar := NcharR
	maxR := MaxRand
	if isGrayList {
		nChar = NcharRLong
		maxR = MaxRandLong
	}
	randExt := rand.Intn(maxR)
	shortUrl, err := encode(addr, randExt, shard, nChar)
	return shortUrl, randExt, nChar, err
}

// Encode ULR string with address, random address, and database shard
func encode(addr uint64, randExt, shard, nChar int) (string, error) {
	encodeA, err := encodeAddr(addr, NcharA)
	// fmt.Println("encode addr:", addr, "shard:", shard)
	if err != nil {
		return "", err
	}
	if len(encodeA) != NcharA {
		return "", errors.New("Encoded address is wrong length")
	}


	// String extension with rand number & shard ID
	// random extension; before conversion to char
	randShard := uint64((randExt << NshardBits) | shard)
	// fmt.Printf("randExt:%d %b  randShard:%b \n", randExt, randExt, randShard)
	encodeR, err := encodeAddr(randShard, nChar)
	// fmt.Println("encodeR:", encodeR, "err:", err)
	if err != nil {
		return "", err
	}
	if len(encodeR) != nChar {
		return "", errors.New("Encoded random extension wrong length")
	}
	shortUrl := encodeR + encodeA
	// fmt.Println("shortURL:", shortURL, "encodeR:", encodeR, "encodeA:", encodeA)
	return shortUrl, nil
}

func DecodeURL(shortUrl string) (uint64, int, int) {
	// fmt.Println("shortUrl:", shortUrl)
	lenExt := len(shortUrl) - NcharA
	// split shortUrl into extension and address
	encodeR, encodeA := shortUrl[:lenExt], shortUrl[lenExt:]
	decodeA := decode(encodeA)
	decodeRS := int(decode(encodeR))
	decodeR := decodeRS >> NshardBits // random value
	shard := decodeRS & (Nshard-1) // database shard
	// fmt.Println(encodeR, encodeA, decodeR, shard)
	return decodeA, decodeR, shard
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
func UrlGrayList(longUrl string) bool {
	u, err := url.Parse(longUrl)
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
	if _, inList := grayList[host]; inList {
		// fmt.Println("Found in gray-list")
		return true // use extended shortening length
	}
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

// Convert binary value to URL string
func encodeAddr(address uint64, nChars int) (string, error) {
	encoded := ""
	// Convert 'address' to numChar;
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
	en, shard:= "ABCabs012", 0
	s, _, _ := EncodeURL("https://goog.com", decode(en), shard)
	// assert.Equal(en, s[len(s)-len(en):])
	i := len(s) - len(en)
	decodeRS := decode(s[:i])
	shard := decodeRS & Nshard-1)
	fmt.Println(en, s[i:], shard, shard)
	fmt.Println()

	shard = 7
	s, _, _ = EncodeURL("https://dropbox.com", decode(en), shard)
	i = len(s) - len(en)
	decodeRS = decode(s[:i])
	shardExtr = uint32(decodeRS & uint64(Nshard-1))
	fmt.Println(en, s[i:], shard, shardExtr)
} */
