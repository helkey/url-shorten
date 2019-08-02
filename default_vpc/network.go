// network.go

package main

import "os"

const PortAddr = ":8088" // 12.0.0.1 (IPv6 ::1)
const PortShorten = ":8088"
const PortExpand = ":8088" // "127.0.0.1:8090"

const UrlAddrServer = "13.56.13.41" + PortAddr

// const UrlShorten = "12.0.0.1" + PortShorten
// const UrlExpand = "127.0.0.1" + PortExpand

// const hostAddr = "localhost" // "127.0.0.1"
const hostAddr = "terraform-20190731170657244900000002.cbnrinnowc9a.us-west-1.rds.amazonaws.com"
const portAddr = 5433

// URL DB shard endpoints (for 'Nshard' shards)
const db_url0 = "terraform-20190801224232423400000002.cbnrinnowc9a.us-west-1.rds.amazonaws.com" // 127.0.0.1"
const db_url1 = "terraform-20190801224232422900000001.cbnrinnowc9a.us-west-1.rds.amazonaws.com"
var hostsDbShard = []string{db_url0, db_url1}
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
