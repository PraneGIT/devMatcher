package store

import (
	"context"

	"github.com/PraneGIT/devmatcher/internal/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserStore interface {
	Create(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}
