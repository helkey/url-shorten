// network.go

package main

import "os"

const UrlShorten = "127.0.0.1:8086"    // 12.0.0.1 (IPv6 ::1)
const UrlAddrServer = "127.0.0.1:8088" // 12.0.0.1 (IPv6 ::1)

const UrlExpand = "127.0.0.1:8090"     // 12.0.0.1 (IPv6 ::1)

const hostAddr = "localhost" // "127.0.0.1"
const portAddr = 5433

// URL DB shard IP addresses ('Nshard' shards)
const host0 = "terraform-20190718002900625800000001.cbeqcb3c0hcc.us-west-2.rds.amazonaws.com"
var hostUrl = []string{host0, host0, "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1", "127.0.0.1"}
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
