FROM golang:1.15
RUN go get -v github.com/githubnemo/CompileDaemon
ENTRYPOINT CompileDaemon -directory=./server -log-prefix=false -build="go build" -command="./server/server" -exclude-dir=.git