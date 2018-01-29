FROM golang:1.9

MAINTAINER Chris Cox <chriscox@abandonedfactory.net>

ENV PATH="${GOPATH}/bin:${PATH}"
ENV GOBIN="${GOPATH}/bin"

# retrieve tools and move them so they don't disappear when /go/bin
# volume is mounted.
RUN go get -u github.com/golang/dep/cmd/dep && \
	   	   	   go get -u github.com/stretchr/testify && \
			   go get -u gopkg.in/alecthomas/gometalinter.v1 && \
			   gometalinter.v1 -i -u && \
			   mkdir -p /go/toolbin && \
			   mv /go/bin/* /usr/local/bin
