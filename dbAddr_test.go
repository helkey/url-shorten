// dbAddr_test.go
// go test dbAddr_test.go encode.go -args 'passwd'

package main

import (
	"os"
	"testing"
	
	"github.com/stretchr/testify/assert"
	_ "github.com/lib/pq"
)


const FullUrl = "http://Full.Url"
const Addr, RandExt = uint64(0xaaaa), 0xcccc
const Shard = 3
func TestSaveurl(t *testing.T) {
	const tableName = "url"

	passwd := os.Args[1]
	db, err := OpenDB(passwd)
	err = db.DropTable(addrTable)
	err = db.CreateTable(addrTable)
	rand.Seed(0) // pick random seed
	
	/*
	// encodeA, _ := encodeAddr(addr, NcharA)
	// randShard := (RandExt << NshardBits) | Shard
	// encodeR, _ := encodeAddr(randShard, charR)
	const isGrayList = false
	shortUrl, _ := encode(isGrayList, Addr, RandExt, Shard, NcharR)
	fmt.Println("TestSaveurl shortURL: ", shortUrl)
	
	dbS, err := NewDBconn(Shard)
	if err != nil {
		return
	}
	dbS.DropTable(tableName)
	dbS.CreateTable(tableName)
	if err != nil {
		return
	}

	// assert.Equal(t, os.Args[1], "passwd")
	passwd := os.Args[1]
	err = dbS.SaveUrl(FullUrl, Addr, RandExt, Shard, dbS.passwd)
	if err != nil {
		return
	}

	fullUrlR, randExtR, err := dbS.ReadUrlDB(Addr, Shard, dbS.passwd)
	fmt.Println(fullUrlR, randExtR, err)
	return */
}
