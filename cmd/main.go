package main

import (
	"github.com/webstrasuite/webstra-gateway/gateway"
	"github.com/webstrasuite/webstra-gateway/proxy"
)

const webPort = ":3000"

func main() {
	// proxy := proxy.NewKubernetes("webstrasuite")
	proxy := proxy.NewLocal()
	gateway := gateway.New(webPort, proxy)

	gateway.RegisterRoutes()

	gateway.Start()
}
