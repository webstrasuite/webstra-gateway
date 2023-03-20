FROM golang:1.20.1

ADD . /go/src/github.com/webstrasuite/gateway-service

WORKDIR /go/src/github.com/webstrasuite/gateway-service

RUN go install

ENTRYPOINT /go/bin/gateway

EXPOSE 3000