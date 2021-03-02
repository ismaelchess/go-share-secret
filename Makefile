FORCE:

run:
	go run ./server/

dev: lint 
	CompileDaemon -directory=./server -log-prefix=false -build="go build" -command="./server/server" -exclude-dir=.git -exclude-dir=./integration

lint:
	golangci-lint run --timeout 5m
