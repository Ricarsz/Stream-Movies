package repository

import (
	"context"
	"time"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/database"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UserRepository struct {
	collection database.CollectionInterface
}

func NewUserRepository(collection database.CollectionInterface) *UserRepository {
	return &UserRepository{collection: collection}
}

func (r *UserRepository) CountByEmail(ctx context.Context, email string) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	return r.collection.CountDocuments(ctx, bson.M{"email": email})
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}