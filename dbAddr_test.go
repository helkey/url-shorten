// dbAddr_test.go
// go test dbAddr_test.go addr.go db.go dbAddr.go dbDROP.go dbUrl.go encode.go -args 'passwd
// COULD DESTRY PRODUCTION DATABASE!!!
// DONT run this in same location as Prod DB

package main

import (
	"fmt"
	"math/rand"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func init() {
}

func TestDbaddr(t *testing.T) {
	rand.Seed(0) // pick non-random seed

	dB, err := OpenAddrDB(dbPassword())
	assert.Equal(t, nil, err)
	err = dB.DropTable("addrs")
	fmt.Println(err)
	assert.Equal(t, nil, err)
	err = dB.CreateAddrTable()
	assert.Equal(t, nil, err)

	addr1, err := dB.GetRandAddr()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(0x1f5b0412), addr1)

	addr2, err := dB.GetRandAddr()
	assert.Equal(t, nil, err)
	assert.Equal(t, uint64(0x23de7767), addr2)

	count, err := dB.NumRowsAddr(addr1)
	assert.Equal(t, nil, err)
	assert.Equal(t, count, 1)
}
