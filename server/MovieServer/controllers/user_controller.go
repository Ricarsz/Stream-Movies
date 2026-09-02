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
			c.JSON(http.StatusBadRequest, gin.H{"error": "User register fail"})
			return
		}
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validate fail", "details": err.Error()})
			return
		}
		hashPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "hash fail"})
			return
		}
		count, err := userRepo.CountByEmail(c.Request.Context(), user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "fail query user email"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email has been used"})
			return
		}
		user.UserID = bson.NewObjectID().Hex()
		user.CreateAt = time.Now()
		user.UpdateAt = time.Now()
		user.Password = hashPassword

		accessToken, refreshToken, err := utils.GenerateToken(user.UserID, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation fail"})
			return
		}
		user.Token = accessToken
		user.RefreshToken = refreshToken

		insertedID, err := userRepo.Create(c.Request.Context(), &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "insert fail"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"inserted_id":   insertedID,
			"token":         accessToken,
			"refresh_token": refreshToken,
		})
	}
}

func LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Email    string `json:"email" validate:"required,email"`
			Password string `json:"password" validate:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}
		if err := validate.Struct(input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validate fail", "details": err.Error()})
			return
		}

		user, err := userRepo.FindByEmail(c.Request.Context(), input.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		if err := utils.VerifyPassword(user.Password, input.Password); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
			return
		}

		accessToken, refreshToken, err := utils.GenerateToken(user.UserID, user.Email, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "token generation fail"})
			return
		}

		_ = userRepo.UpdateTokens(c.Request.Context(), user.UserID, accessToken, refreshToken)

		c.JSON(http.StatusOK, gin.H{
			"token":         accessToken,
			"refresh_token": refreshToken,
			"user_id":       user.UserID,
			"email":         user.Email,
			"role":          user.Role,
			"first_name":    user.FirstName,
			"last_name":     user.LastName,
		})
	}
}
