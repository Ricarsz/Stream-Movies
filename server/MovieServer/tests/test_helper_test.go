package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var TestDB *mongo.Database

func init() {
	// Set test environment variable before any package initialization
	os.Setenv("GO_TEST", "1")
}

func TestMain(m *testing.M) {
	// Set test environment variable to prevent database package from connecting
	os.Setenv("GO_TEST", "1")
	log.Println("TestMain: Starting test setup")

	// Load test environment - try multiple paths
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Warning: Could not load .env file from parent directory, trying current directory")
		if err := godotenv.Load(".env"); err != nil {
			log.Println("Warning: Could not load .env file, using defaults")
		}
	}

	// Use development database (user said it's okay to test with it)
	mongodbURI := os.Getenv("MONGODB_URI")
	if mongodbURI == "" {
		mongodbURI = "mongodb://localhost:27017"
	}

	mongodbDatabase := os.Getenv("MONGODB_DATABASE")
	if mongodbDatabase == "" {
		mongodbDatabase = "stream-movies"
	}

	log.Printf("TestMain: Connecting to MongoDB at %s, database %s", mongodbURI, mongodbDatabase)

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongodbURI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Failed to ping MongoDB:", err)
	}

	TestDB = client.Database(mongodbDatabase)
	log.Println("TestMain: Database connection established")

	// Run tests
	code := m.Run()

	// Cleanup - disconnect only, don't drop the database
	if err := client.Disconnect(context.Background()); err != nil {
		log.Println("Warning: Failed to disconnect from MongoDB:", err)
	}

	log.Println("TestMain: Tests completed")
	os.Exit(code)
}

func SetupTestCollection(t *testing.T, collectionName string) *mongo.Collection {
	collection := TestDB.Collection(collectionName)

	// Clean up collection before test
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := collection.Drop(ctx); err != nil {
		t.Fatal("Failed to drop collection:", err)
	}

	return collection
}

func InsertTestData(t *testing.T, collection *mongo.Collection, documents []interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if len(documents) > 0 {
		_, err := collection.InsertMany(ctx, documents)
		if err != nil {
			t.Fatal("Failed to insert test data:", err)
		}
	}
}

func CountDocuments(t *testing.T, collection *mongo.Collection, filter bson.M) int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		t.Fatal("Failed to count documents:", err)
	}
	return count
}