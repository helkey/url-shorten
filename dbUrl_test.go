// dbUrl_test.go
// go test dbUrl_test.go dbUrl.go dbAddr.go addr.go encode.go -args 'passwd'

package main

import (
	"testing"
	// "github.com/stretchr/testify/assert"
)

func TestSaveurl(t *testing.T) {
	const tableName = "url"

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

		err = dbS.SaveUrl(FullUrl, Addr, RandExt, Shard, dbS.passwd)
		if err != nil {
			return
		}

		fullUrlR, randExtR, err := dbS.ReadUrlDB(Addr, Shard, dbS.passwd)
		fmt.Println(fullUrlR, randExtR, err)
		return */
}
