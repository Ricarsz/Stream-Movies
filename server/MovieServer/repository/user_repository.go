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

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var user models.User
	result, err := r.collection.FindOne(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	if err := result.Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
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

func (r *UserRepository) UpdateTokens(ctx context.Context, userID, token, refreshToken string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.collection.UpdateOne(ctx, bson.M{"user_id": userID}, bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"update_at":     time.Now(),
		},
	})
	return err
}
