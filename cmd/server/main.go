package main

import (
    "fmt"
    "log"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
)

func main() {
    if err := godotenv.Load(); err != nil {
        log.Fatal("Error loading .env file")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

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
