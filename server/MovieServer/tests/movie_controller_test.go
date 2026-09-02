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

func setupMovieRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	collection := TestDB.Collection("movies")
	mongoCollection := database.NewMongoCollection(collection)
	movieRepo := repository.NewMovieRepository(mongoCollection)
	controllers.SetMovieRepository(movieRepo)

	router.GET("/movies", controllers.GetMovies())
	router.GET("/movies/:imdb_id", controllers.GetMovie())
	router.POST("/movies", controllers.AddMovie())

	return router
}

func TestGetMovies_Empty(t *testing.T) {
	router := setupMovieRouter()
	SetupTestCollection(t, "movies")

	req, _ := http.NewRequest("GET", "/movies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var movies []models.Movie
	if err := json.Unmarshal(w.Body.Bytes(), &movies); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	if len(movies) != 0 {
		t.Errorf("Expected 0 movies, got %d", len(movies))
	}
}

func TestGetMovies_WithData(t *testing.T) {
	router := setupMovieRouter()
	collection := SetupTestCollection(t, "movies")

	testMovies := []interface{}{
		models.Movie{
			ImdbID:      "tt1234567",
			Title:       "Test Movie 1",
			PosterPath:  "https://example.com/poster1.jpg",
			YouTubeID:   "youtube123",
			AdminReview: "Great movie",
			Ranking:     models.Ranking{RankingValue: 5, RankingName: "Excellent"},
			Genre:       []models.Genre{{GenreID: 1, GenreName: "Action"}},
		},
		models.Movie{
			ImdbID:      "tt7654321",
			Title:       "Test Movie 2",
			PosterPath:  "https://example.com/poster2.jpg",
			YouTubeID:   "youtube456",
			AdminReview: "Good movie",
			Ranking:     models.Ranking{RankingValue: 4, RankingName: "Good"},
			Genre:       []models.Genre{{GenreID: 2, GenreName: "Comedy"}},
		},
	}
	InsertTestData(t, collection, testMovies)

	req, _ := http.NewRequest("GET", "/movies", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var movies []models.Movie
	if err := json.Unmarshal(w.Body.Bytes(), &movies); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	if len(movies) != 2 {
		t.Errorf("Expected 2 movies, got %d", len(movies))
	}
}

func TestGetMovie_Found(t *testing.T) {
	router := setupMovieRouter()
	collection := SetupTestCollection(t, "movies")

	testMovie := models.Movie{
		ImdbID:      "tt1234567",
		Title:       "Test Movie",
		PosterPath:  "https://example.com/poster.jpg",
		YouTubeID:   "youtube123",
		AdminReview: "Great movie",
		Ranking:     models.Ranking{RankingValue: 5, RankingName: "Excellent"},
		Genre:       []models.Genre{{GenreID: 1, GenreName: "Action"}},
	}
	InsertTestData(t, collection, []interface{}{testMovie})

	req, _ := http.NewRequest("GET", "/movies/tt1234567", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var movie models.Movie
	if err := json.Unmarshal(w.Body.Bytes(), &movie); err != nil {
		t.Fatal("Failed to unmarshal response:", err)
	}

	if movie.ImdbID != "tt1234567" {
		t.Errorf("Expected imdb_id tt1234567, got %s", movie.ImdbID)
	}
}

func TestGetMovie_NotFound(t *testing.T) {
	router := setupMovieRouter()
	SetupTestCollection(t, "movies")

	req, _ := http.NewRequest("GET", "/movies/nonexistent", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestGetMovie_EmptyID(t *testing.T) {
	router := setupMovieRouter()

	req, _ := http.NewRequest("GET", "/movies/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected status 301, got %d", w.Code)
	}
}

func TestAddMovie_Success(t *testing.T) {
	router := setupMovieRouter()
	collection := SetupTestCollection(t, "movies")

	testMovie := models.Movie{
		ImdbID:      "tt9999999",
		Title:       "New Test Movie",
		PosterPath:  "https://example.com/newposter.jpg",
		YouTubeID:   "youtube999",
		AdminReview: "New movie",
		Ranking:     models.Ranking{RankingValue: 4, RankingName: "Good"},
		Genre:       []models.Genre{{GenreID: 3, GenreName: "Drama"}},
	}

	jsonData, _ := json.Marshal(testMovie)
	req, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", w.Code)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	count, err := collection.CountDocuments(ctx, bson.M{"imdb_id": "tt9999999"})
	if err != nil {
		t.Fatal("Failed to count documents:", err)
	}
	if count != 1 {
		t.Errorf("Expected 1 movie, got %d", count)
	}
}

func TestAddMovie_InvalidJSON(t *testing.T) {
	router := setupMovieRouter()
	SetupTestCollection(t, "movies")

	invalidJSON := []byte(`{"invalid": json}`)
	req, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestAddMovie_ValidationFailure(t *testing.T) {
	router := setupMovieRouter()
	SetupTestCollection(t, "movies")

	invalidMovie := models.Movie{ImdbID: "tt8888888"}

	jsonData, _ := json.Marshal(invalidMovie)
	req, _ := http.NewRequest("POST", "/movies", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}