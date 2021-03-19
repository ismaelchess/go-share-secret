FROM golang:1.15
RUN go get -v github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -directory=. -log-prefix=false -build="go build -o service" -command=./service -exclude-dir=.git -exclude-dir=./ui