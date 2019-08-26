// network.go

package main

import "os"

const PortAddr = ":8088" // 12.0.0.1 (IPv6 ::1)
const PortShorten = ":8088"
const PortExpand = ":8088" // "127.0.0.1:8090"

const UrlAddrServer = "34.83.95.20" + PortAddr

// ADDR DB endpoint
//  GCP database port number does not seem configurable
//  codelabs.developers.google.com/codelabs/connecting-to-cloud-sql/index.html
const hostAddr = "35.233.227.58"
const portAddr = 5432 

// URL DB shard endpoints (for 'Nshard' shards)
const db_url0 = hostAddr // share a database for initial testing
const db_url1 = ""
var hostsDbShard = []string{db_url0, db_url1}
const portUrl = portAddr

func getShortenInstance() string {
	// Requires setting each instance env variable differently
	instance := os.Getenv("INSTANCE")
	if instance == "" {
		instance = "?"
	}
	// AWS: to retrieve instance info from metadata of running instance
	//   $ curl http://169.254.169.254/latest/meta-data/
	return instance
}
