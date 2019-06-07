// RequestShorten
// go run RequestShorten.go addr.go dbUrl.go encode.go 'passwd

// "localhost:8086/create" (working WSL)
// "localhost:8086/create/?source=&url=http://FullURL"

package main

import (
	"log"
	"os"
)

func password() string {
	if len(os.Args) <= 1 {
		log.Fatal("Supply DB password")
	}
	return os.Args[1]
}
