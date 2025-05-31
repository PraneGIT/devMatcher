package handlers

import (
	"context"
	"net/http"

	"github.com/PraneGIT/devmatcher/internal/api/dto"
	"github.com/PraneGIT/devmatcher/internal/config"
	"github.com/PraneGIT/devmatcher/internal/core"
	"github.com/PraneGIT/devmatcher/internal/store/mongodb"
	"github.com/gin-gonic/gin"
)

var authService *core.AuthService

func getAuthService() *core.AuthService {
	if authService == nil {
		client := mongodb.Client
		db := client.Database("devmatcher")
		userStore := mongodb.NewUserMongoStore(db)
		authService = core.NewAuthService(userStore, config.AppConfig.JWTSecret)
	}
	return authService
}

func Register(c *gin.Context) {
	var req dto.UserRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service := getAuthService()
	user, err := service.RegisterUser(context.Background(), req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}
	resp := dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       user.ID,
		Name:         user.Name,
		Email:        user.Email,
	}
	c.JSON(http.StatusCreated, resp)
}

func Login(c *gin.Context) {
	var req dto.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service := getAuthService()
	user, err := service.LoginUser(context.Background(), req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	access, refresh, err := service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}
	resp := dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       user.ID,
		Name:         user.Name,
		Email:        user.Email,
	}
	c.JSON(http.StatusOK, resp)
}

func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	service := getAuthService()
	claims, err := service.ParseToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	user, err := service.UserStore.GetByEmail(context.Background(), claims.Email)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	access, refresh, err := service.GenerateTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}
	resp := dto.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		UserID:       user.ID,
		Name:         user.Name,
		Email:        user.Email,
	}
	c.JSON(http.StatusOK, resp)
}
