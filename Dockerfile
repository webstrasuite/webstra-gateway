FROM golang:1.20.1

ADD . /go/src/github.com/webstrasuite/webstra-gateway

WORKDIR /go/src/github.com/webstrasuite/webstra-gateway

RUN go install

ENTRYPOINT /go/bin/gateway

EXPOSE 3000