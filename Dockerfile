FROM golang:1.20.1

ADD . /go/src/github.com/webstraservices/gateway-service

WORKDIR /go/src/github.com/webstraservices/gateway-service

RUN go install

ENTRYPOINT /go/bin/gateway

EXPOSE 3000