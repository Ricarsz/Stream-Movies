package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/database"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/models"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/repository"
	"github.com/Ricarse/goMovies/server/GoMoviesServer/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var userRepo *repository.UserRepository
var validate = validator.New()

func init() {
	if os.Getenv("GO_TEST") == "" {
		userRepo = repository.NewUserRepository(database.NewMongoCollection(database.OpenCollection("users")))
	}
}

func SetUserRepository(repo *repository.UserRepository) {
	userRepo = repo
}

func RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {
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
			return
		}
		count, err := userRepo.CountByEmail(c.Request.Context(), user.Email)
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
		insertedID, err := userRepo.Create(c.Request.Context(), &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "insert fail"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"inserted_id": insertedID})
	}
}