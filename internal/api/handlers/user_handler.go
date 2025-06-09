package handlers

import (
	"context"
	"net/http"

	"github.com/PraneGIT/devmatcher/internal/store/mongodb"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getUserStore() *mongodb.UserMongoStore {
	client := mongodb.Client
	db := client.Database("devmatcher")
	return mongodb.NewUserMongoStore(db).(*mongodb.UserMongoStore)
}

func GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID") // userID set by auth middleware
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	store := getUserStore()
	user, err := store.GetByID(context.Background(), objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID") // userID set by auth middleware
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	objID, err := primitive.ObjectIDFromHex(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	store := getUserStore()
	user, err := store.GetByID(context.Background(), objID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	err = store.Update(context.Background(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func GetPreferences(c *gin.Context) {
	// TODO: Implement get user preferences
	c.JSON(501, gin.H{"error": "not implemented"})
}

func UpdatePreferences(c *gin.Context) {
	// TODO: Implement update user preferences
	c.JSON(501, gin.H{"error": "not implemented"})
}
