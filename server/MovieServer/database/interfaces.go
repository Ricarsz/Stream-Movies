package database

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CollectionInterface interface {
	Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error)
	FindOne(ctx context.Context, filter interface{}) (*mongo.SingleResult, error)
	InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error)
	CountDocuments(ctx context.Context, filter interface{}) (int64, error)
}

type MongoCollection struct {
	collection *mongo.Collection
}

func NewMongoCollection(collection *mongo.Collection) *MongoCollection {
	return &MongoCollection{collection: collection}
}

func (mc *MongoCollection) Find(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return mc.collection.Find(ctx, filter)
}

func (mc *MongoCollection) FindOne(ctx context.Context, filter interface{}) (*mongo.SingleResult, error) {
	return mc.collection.FindOne(ctx, filter), nil
}

func (mc *MongoCollection) InsertOne(ctx context.Context, document interface{}) (*mongo.InsertOneResult, error) {
	return mc.collection.InsertOne(ctx, document)
}

func (mc *MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}) (*mongo.UpdateResult, error) {
	return mc.collection.UpdateOne(ctx, filter, update)
}

func (mc *MongoCollection) DeleteOne(ctx context.Context, filter interface{}) (*mongo.DeleteResult, error) {
	return mc.collection.DeleteOne(ctx, filter)
}

func (mc *MongoCollection) CountDocuments(ctx context.Context, filter interface{}) (int64, error) {
	return mc.collection.CountDocuments(ctx, filter)
}
