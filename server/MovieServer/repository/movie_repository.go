package repository

import (
	"context"
	"time"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/database"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MovieRepository struct {
	collection database.CollectionInterface
}

func NewMovieRepository(collection database.CollectionInterface) *MovieRepository {
	return &MovieRepository{collection: collection}
}

func (r *MovieRepository) FindAll(ctx context.Context) ([]models.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var movies []models.Movie
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &movies); err != nil {
		return nil, err
	}
	return movies, nil
}

func (r *MovieRepository) FindByImdbID(ctx context.Context, imdbID string) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var movie models.Movie
	result, err := r.collection.FindOne(ctx, bson.M{"imdb_id": imdbID})
	if err != nil {
		return nil, err
	}
	if err := result.Decode(&movie); err != nil {
		return nil, err
	}
	return &movie, nil
}

func (r *MovieRepository) Create(ctx context.Context, movie *models.Movie) (interface{}, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, movie)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}