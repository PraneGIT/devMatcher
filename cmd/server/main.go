package main

import (
    "fmt"
    "log"

	"github.com/PraneGIT/devmatcher/internal/config"
	"github.com/PraneGIT/devmatcher/internal/store/mongodb"
	"github.com/PraneGIT/devmatcher/internal/api"
)

func main() {
	config.LoadConfig()
	mongodb.InitMongo()
    port := config.AppConfig.Port

    router := api.SetupRouter()

    fmt.Println("Server running on port", port)
    if err := router.Run(":" + port); err != nil {
        log.Fatal(err)
    }
}
