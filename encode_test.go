// encode_test
// go test encode_test.go addr.go encode.go 

package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

// Test numeric limits; likely to break if NcharA increased too much
func TestRange(t *testing.T) {
	const maxCharEncode = 10 // max characters for uint64 representation
	if NcharA > maxCharEncode {
		t.Errorf("NcharA too large - exceeds  uint64 representation")
	}

	// Address server using https
	maxStr := strings.Repeat("Z", NcharA) // max encoded base address
	maxAddr := decode(maxStr)             // integer representation
	s := strconv.Itoa(int(maxAddr))
	b := []byte(s) // addr server converts to []byte
	maxAddr2, _ := strconv.Atoi(string(b))
	assert.Equal(t, maxAddr, uint64(maxAddr2))
	maxS2, _ := encodeAddr(uint64(maxAddr2), NcharA)
	assert.Equal(t, maxStr, maxS2)
}

// Test encoding/decoding functions
func TestEncode(t *testing.T) {
	en, iShard := "ABCabs012", 0
	s, _, _, _ := EncodeURL("https://goog.com", decode(en), iShard)
	i := len(s) - len(en)
	decodeRS := decode(s[:i])
	shard := int(decodeRS & (Nshard-1))
	assert.Equal(t, en, s[i:])
	assert.Equal(t, iShard, shard)
	return

	iShard = 7
	s, _, _, _ = EncodeURL("https://dropbox.com", decode(en), iShard)
	i = len(s) - len(en)
	decodeRS = decode(s[:i])
	shard = int(decodeRS & (Nshard-1))

	encoded := "ABCabs0123"
	decoded := decode(encoded)
	encoded2, _ := encodeAddr(decoded, len(en))
	assert.Equal(t, decode, 495548099420723299)
	assert.Equal(t, encoded2, "ABCabs0123")

	decodeA, decodeR, iShard := DecodeURL("oxABCabs0123") // randSlice=1521
	assert.Equal(t, decodeR, 190)
	assert.Equal(t, decodeA, 495548099420723299)
	assert.Equal(t, iShard, 1)

	decodeA, decodeR, iShard = DecodeURL("ZG8xABCabs0123") // randSlice=14699985
	assert.Equal(t, decodeR, 1837498)
	assert.Equal(t, decodeA, 495548099420723299)
	assert.Equal(t, iShard, 1)
}

func TestGraylisted(t *testing.T) {
	url := "postgres://user:pass@host.com:5432/path?k=v#f"
	// fmt.Println("graylisted:", urlGrayListed(url))
	assert.False(t, UrlGrayList(url))
	url = "https://www.dropbox.com/filename"
	// fmt.Println(urlGrayListed(url)) // assert true
	assert.True(t, UrlGrayList(url))
}
