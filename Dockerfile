FROM golang:1.17.5

ADD . /go/src/github.com/webstraservices/gateway

WORKDIR /go/src/github.com/webstraservices/gateway

RUN go install

ENTRYPOINT /go/bin/gateway

EXPOSE 3000