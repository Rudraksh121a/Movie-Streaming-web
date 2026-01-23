package routes

import (
	controller "github.com/Rudraksh121a/Movie-Streaming-web/controllers"
	"github.com/Rudraksh121a/Movie-Streaming-web/middleware"
	"github.com/gin-gonic/gin"
)

func SetupProtectedRoutes(router *gin.Engine) {
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/movie/:imdb_id", controller.GetMovie())
		protected.POST("/addmovie", controller.AddMovie())
	}
}
