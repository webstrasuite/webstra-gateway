package main

const webPort = ":3000"

func main() {
	router := NewRouter(webPort, "webstrasuite")

	router.RegisterRoutes()

	router.Start()
}
