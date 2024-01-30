
dev: 
	nodemon --exec go run --tags dynamic $(shell pwd)/cmd/main.go --signal SIGTERM
