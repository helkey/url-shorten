// Encode

package main

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	// "github.com/gotestyourself/gotest.tools/assert"
)


func TestAddr(t *testing.T) {
	chAddr, baseIncr, nAddrs := Test()
	const sleepMs = 10 // avoid race condition between channels
	const nIter = 30
	for i:=0; i<nIter; i++ {
		addr := <-chAddr
		addrExpect := baseIncr * (uint64(i / nAddrs) + 1) + uint64(i % nAddrs)
		assert.Equal(t, addr, addrExpect)
		time.Sleep(sleepMs * time.Millisecond)
	}
}




