// dbAddr_test.go
// go test dbAddr_test.go dbAddr.go addr.go encode.go -args 'passwd

package main

import (
	"math/rand"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(0) // pick non-random seed
}

const FullUrl = "http://Full.Url"
const Addr, RandExt = uint64(0xaaaa), 0xcccc
const Shard = 3

func TestAddr(t *testing.T) {
	passwd := os.Args[1]
	dB, err := OpenDB(passwd)
	assert.Equal(t, err, nil)
	err = dB.DropTable()
	assert.Equal(t, err, nil)
	err = dB.CreateTable()
	assert.Equal(t, err, nil)

	addr1, err := dB.GetRandAddr()
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(0x1f5b0412), addr1)

	addr2, err := dB.GetRandAddr()
	assert.Equal(t, err, nil)
	assert.Equal(t, uint64(0x23de7767), addr2)

	count, err := dB.NumAddrRows(addr1)
	assert.Equal(t, err, nil)
	assert.Equal(t, count, 1)
}
