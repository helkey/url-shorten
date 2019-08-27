alias buildAddr='CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o addr/ReqAddr RequestAddr.go addr.go db.go dbAddr.go dbDrop.go dbUrl.go encode.go network.go param.go'
alias buildShort='CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o short/ReqShort RequestShorten.go addr.go db.go dbAddr.go dbDROP.go dbUrl.go encode.go network.go param.go'
alias buildExpand='CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o expand/ReqExpand RequestExpand.go addr.go db.go dbAddr.go dbUrl.go dbDrop.go encode.go network.go param.go'

