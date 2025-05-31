package store

import (
	"context"

	"github.com/PraneGIT/devmatcher/internal/domain"
)

type UserStore interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
