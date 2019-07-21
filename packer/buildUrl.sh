#!/bin/sh -x

go build -o RequestAddr RequestAddr.go addr.go db.go dbAddr.go dbDrop.go dbUrl.go encode.go network.go
go build -o RequestShorten RequestShorten.go addr.go db.go dbAddr.go dbDROP.go dbUrl.go encode.go network.go
go build -o RequestExpand RequestExpand.go addr.go db.go dbAddr.go dbUrl.go dbDrop.go encode.go network.go
