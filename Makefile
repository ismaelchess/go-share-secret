FORCE:

run:
	go run ./server/

dev:
	CompileDaemon -directory=./server -log-prefix=false -build="go build" -command="./server/server" -exclude-dir=.git -exclude-dir=./integration
