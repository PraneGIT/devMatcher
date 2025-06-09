package mongodb

import (
	"context"

	"github.com/PraneGIT/devmatcher/internal/domain"
	"github.com/PraneGIT/devmatcher/internal/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoStore struct {
	Coll *mongo.Collection
}

func NewUserMongoStore(db *mongo.Database) store.UserStore {
	return &UserMongoStore{Coll: db.Collection("users")}
}

func (s *UserMongoStore) Create(ctx context.Context, user *domain.User) error {
	res, err := s.Coll.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		user.ID = oid.Hex()
	}
	return nil
}

func (s *UserMongoStore) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := s.Coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserMongoStore) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.User, error) {
	var user domain.User
	err := s.Coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserMongoStore) Update(ctx context.Context, user *domain.User) error {
	objID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		return err
	}
	_, err = s.Coll.UpdateOne(ctx, bson.M{"_id": objID}, bson.M{"$set": bson.M{"name": user.Name}})
	return err
}
