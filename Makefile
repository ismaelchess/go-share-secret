FORCE:

dev: lint 
	CompileDaemon -directory=./server -log-prefix=false -build="go build" -command="./server/server" -exclude-dir=.git

run:
	go run ./server/

lint:
	golangci-lint run --timeout 5m

test:
	go test -v -count 1 ./server

testcover:
	go test -cover $(shell go list ./server)
