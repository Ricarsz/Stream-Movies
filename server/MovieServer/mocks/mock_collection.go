package mocks

import (
	"context"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/database"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MockCollection struct {
	FindFunc           func(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	FindOneFunc        func(ctx context.Context, filter interface{}) (*mongo.SingleResult, error)
	InsertOneFunc      func(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	CountDocumentsFunc func(ctx context.Context, filter interface{}) (int64, error)
}

var _ database.CollectionInterface = (*MockCollection)(nil)

func (m *MockCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	if m.FindFunc != nil {
		return m.FindFunc(ctx, filter)
	}
	return nil, nil
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}) (*mongo.SingleResult, error) {
	if m.FindOneFunc != nil {
		return m.FindOneFunc(ctx, filter)
	}
	return nil, nil
}

func (m *MockCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	if m.InsertOneFunc != nil {
		return m.InsertOneFunc(ctx, document)
	}
	return nil, nil
}

func (m *MockCollection) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	if m.CountDocumentsFunc != nil {
		return m.CountDocumentsFunc(ctx, filter)
	}
	return 0, nil
}
