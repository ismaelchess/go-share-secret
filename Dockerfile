FROM golang:1.15
ADD . /go/src/github.com/ismaelchess/go-share-secret
WORKDIR /go/src/github.com/ismaelchess/go-share-secret
RUN make install build
CMD [ "go", "version" ]