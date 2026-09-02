package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Client *mongo.Client

func init() {
	// Only connect if not in test environment
	// Note: This check happens at package initialization time
	// For tests, we need to set GO_TEST before running tests
	if os.Getenv("GO_TEST") == "" {
		Client = Connect()
	}
}

func Connect() *mongo.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("env fail read")
	}
	MONGODB := os.Getenv("MONGODB_URI")
	if MONGODB == "" {
		log.Println("empty MongDB")
		return nil
	}
	clientOptions := options.Client().ApplyURI(MONGODB)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("connect fail:", err)
		return nil
	}
	return client
}

func OpenCollection(collectionName string) *mongo.Collection {
	if Client == nil {
		// In test environment, this should not be called
		// Return nil or panic with a helpful message
		if os.Getenv("GO_TEST") != "" {
			log.Fatal("OpenCollection called in test environment. Use test database directly.")
		}
		log.Fatal("Database client not initialized.")
	}
	DATABASENAME := os.Getenv("MONGODB_DATABASE")
	collection := Client.Database(DATABASENAME).Collection(collectionName)
	return collection
}
