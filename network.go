// network.go

package main

import "os"

const UrlAddrServer = ":8088" // 12.0.0.1 (IPv6 ::1)
const UrlShorten = ":8088"
const UrlExpand = ":8088" // "127.0.0.1:8090"

// const hostAddr = "localhost" // "127.0.0.1"
const hostAddr = "terraform-20190730164638380000000002.cbnrinnowc9a.us-west-1.rds.amazonaws.com"
const portAddr = 5433

// URL DB shard IP addresses ('Nshard' shards)
const host0 = "127.0.0.1"
var hostsDbShard = []string{host0, host0} // , "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1"}
const portUrl = portAddr

func getShortenInstance() string {
	// Requires setting each instance env variable differently
	instance := os.Getenv("INSTANCE")
	if instance == "" {
		instance = "?"
	}
	// OR: to retrieve instance info from metadata of running AWS instance:
	//   $ curl http://169.254.169.254/latest/meta-data/
	return instance
}
