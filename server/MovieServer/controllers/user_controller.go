package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/database"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/models"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection("users")
var validate = validator.New()

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
		defer cancel()
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "User register fail"})
			return
		}
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "validate fail"})
			return
		}
		hashPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "hash fail"})
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"err": "fail query user email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"err": "email has been use"})
			return
		}
		user.UserID = bson.NewObjectID().Hex()
		user.CreateAt = time.Now()
		user.UpdateAt = time.Now()
		user.Password = hashPassword
		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "insert fail"})
			return
		}
		c.JSON(http.StatusCreated, result)
	}
}
