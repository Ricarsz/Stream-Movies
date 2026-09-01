package main

import (
	"fmt"

	"github.com/Ricarse/goMovies/server/GoMoviesServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/movies", controllers.GetMovies())
	if err := router.Run(":8080"); err != nil {
		fmt.Println("err", err)
	}
}
