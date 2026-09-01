package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/database"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var movieCollection *mongo.Collection = database.OpenCollection("movies")
var validate = validator.New()

func GetMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		var movies []models.Movie
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail fetch movies"})
			return
		}
		defer cursor.Close(ctx)
		if err := cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, movies)
	}
}

func GetMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		movieID := c.Param("imdb_id")
		if movieID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"err": "movieID is Empty"})
			return
		}
		var movie models.Movie
		if err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieID}).Decode(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "movieID fetch fail"})
			return
		}
		c.JSON(http.StatusOK, movie)
	}
}

func AddMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "add failed"})
			return
		}
		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "validation fail", "detals": err.Error()})
			return
		}
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "insert fail"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}
