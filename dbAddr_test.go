// dbAddr_test.go
// go test dbAddr_test.go dbAddr.go addr.go encode.go genAddr.go -args 'passwd

package main

import (
	"math/rand"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

const FullUrl = "http://Full.Url"
const Addr, RandExt = uint64(0xaaaa), 0xcccc
const Shard = 3

func TestAddr(t *testing.T) {
	passwd := os.Args[1]
	db, err := OpenDB(passwd)
	assert.Equal(t, err, nil)
	err = db.DropTable()
	assert.Equal(t, err, nil)
	err = db.CreateTable()
	assert.Equal(t, err, nil)

	rand.Seed(0) // pick non-random seed
	NrangeAlloc := Nrange >> NallocBits
	const len_test = 10
	addrArr := rand.Perm(NrangeAlloc)[:len_test]
	err = db.SaveAddrArr(addrArr)
	assert.Equal(t, err, nil)

	addrArrR, err := db.GetAddrArr()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(addrArrR), len_test)

	err = db.MarkAddrUsed(addrArrR[0])
	addrArrR1, err := db.GetAddrArr()
	assert.Equal(t, err, nil)
	assert.Equal(t, len(addrArrR1), len_test-1)
}
