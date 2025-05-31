package handlers

import "github.com/gin-gonic/gin"

func WebSocketHandler(c *gin.Context) {
	// TODO: Implement WebSocket connection handler
	c.JSON(501, gin.H{"error": "not implemented"})
}
