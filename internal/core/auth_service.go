package core

import (
	"context"
	"errors"
	"time"

	"github.com/PraneGIT/devmatcher/internal/domain"
	"github.com/PraneGIT/devmatcher/internal/store"
	"github.com/PraneGIT/devmatcher/internal/util"
)

type AuthService struct {
	UserStore store.UserStore
	JWTSecret string
}

func NewAuthService(userStore store.UserStore, jwtSecret string) *AuthService {
	return &AuthService{UserStore: userStore, JWTSecret: jwtSecret}
}

func (s *AuthService) RegisterUser(ctx context.Context, name, email, password string) (*domain.User, error) {
	existing, err := s.UserStore.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("email already registered")
	}
	hash, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: hash,
	}
	err = s.UserStore.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *AuthService) LoginUser(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.UserStore.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid credentials")
	}
	if !util.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}

func (s *AuthService) GenerateTokens(user *domain.User) (accessToken, refreshToken string, err error) {
	accessToken, err = util.GenerateJWT(user.ID, user.Email, s.JWTSecret, 24*time.Hour)
	if err != nil {
		return
	}
	refreshToken, err = util.GenerateJWT(user.ID, user.Email, s.JWTSecret, 7*24*time.Hour)
	return
}

func (s *AuthService) ParseToken(token string) (*util.TokenClaims, error) {
	return util.ParseJWT(token, s.JWTSecret)
}
