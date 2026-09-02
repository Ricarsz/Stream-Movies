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
			c.JSON(http.StatusBadRequest, gin.H{"error": "movieID is Empty"})
			return
		}
		movie, err := movieRepo.FindByImdbID(c.Request.Context(), movieID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movieID fetch fail"})
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "add failed"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation fail", "details": err.Error()})
			return
		}
		insertedID, err := movieRepo.Create(c.Request.Context(), &movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "insert fail"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"inserted_id": insertedID})
	}
}

func UpdateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movieID is Empty"})
			return
		}
		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validation fail", "details": err.Error()})
			return
		}
		count, err := movieRepo.Update(c.Request.Context(), movieID, &movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update fail"})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "movie updated", "modified_count": count})
	}
}

func DeleteMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "movieID is Empty"})
			return
		}
		count, err := movieRepo.Delete(c.Request.Context(), movieID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "delete fail"})
			return
		}
		if count == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "movie deleted", "deleted_count": count})
	}
}
