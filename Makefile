
dev: 
	nodemon --exec go run --tags dynamic $(shell pwd)/main.go --signal SIGTERM
