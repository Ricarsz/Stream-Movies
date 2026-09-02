package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/controllers"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/database"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/models"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/repository"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func setupUserRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	collection := TestDB.Collection("users")
	mongoCollection := database.NewMongoCollection(collection)
	userRepo := repository.NewUserRepository(mongoCollection)
	controllers.SetUserRepository(userRepo)

	router.POST("/users/register", controllers.RegisterUser())

	return router
}

func TestRegisterUser_Success(t *testing.T) {
	router := setupUserRouter()
	collection := SetupTestCollection(t, "users")

	testUser := models.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password123",
		Role:      "USER",
		FavouriteGenres: []models.Genre{
			{GenreID: 1, GenreName: "Action"},
		},
	}

	jsonData, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, bson.M{"email": "john.doe@example.com"})
	if err != nil {
		t.Fatal("Failed to count documents:", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 user, got %d", count)
	}
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	router := setupUserRouter()
	collection := SetupTestCollection(t, "users")

	firstUser := models.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password123",
		Role:      "USER",
		FavouriteGenres: []models.Genre{
			{GenreID: 1, GenreName: "Action"},
		},
	}

	jsonData, _ := json.Marshal(firstUser)
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("First registration failed with status %d", w.Code)
	}

	req, _ = http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for duplicate email, got %d", w.Code)
	}

	_ = collection
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	router := setupUserRouter()
	SetupTestCollection(t, "users")

	invalidJSON := []byte(`{"invalid": json}`)
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRegisterUser_ValidationFailure(t *testing.T) {
	router := setupUserRouter()
	SetupTestCollection(t, "users")

	invalidUser := models.User{
		FirstName: "J",
		LastName:  "D",
		Email:     "invalid-email",
		Password:  "123",
		Role:      "INVALID",
	}

	jsonData, _ := json.Marshal(invalidUser)
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestRegisterUser_PasswordHashing(t *testing.T) {
	router := setupUserRouter()
	collection := SetupTestCollection(t, "users")

	testUser := models.User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "test.user@example.com",
		Password:  "mypassword123",
		Role:      "USER",
		FavouriteGenres: []models.Genre{
			{GenreID: 1, GenreName: "Action"},
		},
	}

	jsonData, _ := json.Marshal(testUser)
	req, _ := http.NewRequest("POST", "/users/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": "test.user@example.com"}).Decode(&user)
	if err != nil {
		t.Fatal("Failed to find user:", err)
	}

	if user.Password == "mypassword123" {
		t.Error("Password should be hashed, not stored in plain text")
	}

	if len(user.Password) < 10 {
		t.Error("Hashed password seems too short")
	}
}