package main

const webPort = ":3000"

func main() {
	router := NewRouter(webPort)

	router.RegisterRoutes()

	router.Start()
}
