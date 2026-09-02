package controllers

import (
	"net/http"
	"os"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/database"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/models"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/repository"
	"github.com/gin-gonic/gin"
)

var movieRepo *repository.MovieRepository

func init() {
	if os.Getenv("GO_TEST") == "" {
		movieRepo = repository.NewMovieRepository(database.NewMongoCollection(database.OpenCollection("movies")))
	}
}

func SetMovieRepository(repo *repository.MovieRepository) {
	movieRepo = repo
}

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		movies, err := movieRepo.FindAll(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail fetch movies"})
			return
		}
		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"err": "movieID is Empty"})
			return
		}
		movie, err := movieRepo.FindByImdbID(c.Request.Context(), movieID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "movieID fetch fail"})
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "add failed"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "validation fail", "detals": err.Error()})
			return
		}
		insertedID, err := movieRepo.Create(c.Request.Context(), &movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "insert fail"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"inserted_id": insertedID})
	}
}