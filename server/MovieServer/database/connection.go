package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client = Connect()

func Connect() *mongo.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("env fail read")
	}
	MONGODB := os.Getenv("MONGODB_URI")
	if MONGODB == "" {
		log.Println("empty MongDB")
	}
	clientOptions := options.Client().ApplyURI(MONGODB)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("connect fail:", err)
	}
	return client
}

func OpenCollection(collectionName string) *mongo.Collection {
	DATABASENAME := os.Getenv("MONGODB_DATABASE")
	collection := Client.Database(DATABASENAME).Collection(collectionName)
	return collection
}
