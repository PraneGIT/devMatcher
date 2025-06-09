package api

import (
	"github.com/gin-gonic/gin"

	"github.com/PraneGIT/devmatcher/internal/api/handlers"
	"github.com/PraneGIT/devmatcher/internal/api/middleware"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// User profile routes (protected)
		user := api.Group("/user")
		user.Use(middleware.Auth())
		{
			user.GET("/profile", handlers.GetProfile)
			user.PUT("/profile", handlers.UpdateProfile)
			user.GET("/preferences", handlers.GetPreferences)
			user.PUT("/preferences", handlers.UpdatePreferences)
		}

		// Discovery (swiping) routes (protected)
		discovery := api.Group("/discovery")
		// discovery.Use(middleware.Auth())
		{
			discovery.GET("/profiles", handlers.GetProfilesForDiscovery)
		}

		// Interaction routes
		interactions := api.Group("/interactions")
		{
			interactions.POST("/swipe", handlers.RecordSwipe)
		}

		api.GET("/ws", handlers.WebSocketHandler)
	}

	return r
}
