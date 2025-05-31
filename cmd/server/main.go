package main

import (
    "fmt"
    "log"

    "github.com/gin-gonic/gin"
	"github.com/PraneGIT/devmatcher/internal/config"
	"github.com/PraneGIT/devmatcher/internal/store/mongodb"
)

func main() {
	config.LoadConfig()
	mongodb.InitMongo()
    port := config.AppConfig.Port


    router := gin.Default()

    // TODO: Add middleware & routes
    router.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{"message": "pong"})
    })

    fmt.Println("Server running on port", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}
