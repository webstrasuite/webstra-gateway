package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
	})

	router.Any("/api/*path", Gateway)

	router.Run(":3000")
}
