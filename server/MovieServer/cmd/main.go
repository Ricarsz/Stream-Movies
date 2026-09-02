package main

import (
	"fmt"
	"time"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// CORS — allow React dev server
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Public routes
	router.POST("/users/register", controllers.RegisterUser())
	router.POST("/users/login", controllers.LoginUser())

	// Public movie reads
	router.GET("/movies", controllers.GetMovies())
	router.GET("/movies/:imdb_id", controllers.GetMovie())

	// Protected routes (require valid JWT)
	auth := router.Group("/")
	auth.Use(controllers.AuthMiddleware())
	{
		auth.POST("/movies", controllers.AddMovie())
		auth.PUT("/movies/:imdb_id", controllers.UpdateMovie())
		auth.DELETE("/movies/:imdb_id", controllers.DeleteMovie())
	}

	if err := router.Run(":8080"); err != nil {
		fmt.Println("err", err)
	}
}
