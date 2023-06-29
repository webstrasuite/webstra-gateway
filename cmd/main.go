package main

import (
	"log"

	"github.com/webstrasuite/webstra-gateway/gateway"
	"github.com/webstrasuite/webstra-gateway/proxy"
)

const (
	listenAddr      = ":3000"
	authServiceAddr = ":3001"
)

func main() {
	// proxy := proxy.NewKubernetes("webstrasuite")
	proxy := proxy.NewLocal()

	gateway, err := gateway.New(listenAddr, authServiceAddr, proxy)
	if err != nil {
		log.Fatal(err)
	}

	gateway.RegisterRoutes()

	gateway.Start()
}
