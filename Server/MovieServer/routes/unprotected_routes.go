package routes

import (
	controller "github.com/Rudraksh121a/Movie-Streaming-web/controllers"
	"github.com/gin-gonic/gin"
)

func SetupUnProtectedRoutes(router *gin.Engine) {
	router.GET("/movies", controller.GetMovies())
	router.POST("/login", controller.LoginUser())
	router.POST("/register", controller.RegisterUser())
}
