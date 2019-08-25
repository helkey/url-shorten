// network.go

package main

import "os"

const PortAddr = ":8088" // 12.0.0.1 (IPv6 ::1)
const PortShorten = ":8088"
const PortExpand = ":8088" // "127.0.0.1:8090"


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
