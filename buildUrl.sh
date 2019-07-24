alias buildAddr='go build -o ReqAddr -ldflags "-X db.db_password=$TF_VAR_db_password" RequestAddr.go addr.go db.go dbAddr.go dbDrop.go dbUrl.go encode.go network.go'

alias buildShort='go build -o ReqShort -ldflags "-X db.db_password=$TF_VAR_db_password" RequestShorten.go addr.go db.go dbAddr.go dbDROP.go dbUrl.go encode.go network.go'
alias buildExpand='go build -o ReqExpad -ldflags "-X db.db_password=$TF_VAR_db_password" RequestExpand.go addr.go db.go dbAddr.go dbUrl.go dbDrop.go encode.go network.go'

