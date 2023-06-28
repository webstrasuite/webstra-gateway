package main

const webPort = ":3000"

func main() {
	// gateway := NewKubernetesProxy("webstrasuite")
	gateway := NewLocalProxy()
	router := NewRouter(webPort, gateway)

	router.RegisterRoutes()

	router.Start()
}
