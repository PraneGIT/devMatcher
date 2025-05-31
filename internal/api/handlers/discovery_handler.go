package handlers

import "github.com/gin-gonic/gin"

func GetProfilesForDiscovery(c *gin.Context) {
	// TODO: Implement discovery profiles
	c.JSON(501, gin.H{"error": "not implemented"})
}
