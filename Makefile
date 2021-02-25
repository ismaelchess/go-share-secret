FORCE:

run:
	go run ./server/

dev:
	CompileDaemon -directory=. -log-prefix=false -build="go run ./server/" -command="." -exclude-dir=.git -exclude-dir=./integration