package main

import (
	"github.com/gin-gonic/gin"
	"github.com/webstraservices/gateway/gateway"
)

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Unable to load environment files - Err: %s", err)
	// }

	// connectionString := os.Getenv("DB_CONNECTION_STRING")

	// db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	// if err != nil {
	// 	log.Fatalf("Unable to open database instance")
	// }

	// db.AutoMigrate(models.User{})

	router := gin.Default()
	router.GET("/healthz", func(c *gin.Context) {
		c.Writer.WriteHeader(200)
	})

	router.Any("/api/*path", gateway.Gateway)

	router.Run(":3000")
}
