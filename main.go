package main

import "github.com/webstrasuite/webstra-gateway/proxy"

const webPort = ":3000"

func main() {
	// gateway := NewKubernetesProxy("webstrasuite")
	gateway := proxy.NewLocal()
	router := NewRouter(webPort, gateway)

	router.RegisterRoutes()

	router.Start()
}
