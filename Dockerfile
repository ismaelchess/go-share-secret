FROM golang:1.15
ARG NETRC
ADD . /go/src/github.com/ismaelchess/go-share-secret
WORKDIR /go/src/github.com/ismaelchess/go-share-secret
RUN echo "${NETRC}" > /root/.netrc
RUN make install build
CMD [ "go", "version" ]