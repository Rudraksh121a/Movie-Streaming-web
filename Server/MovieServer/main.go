package main

import (
	"fmt"

	"github.com/Rudraksh121a/Movie-Streaming-web/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	routes.SetupProtectedRoutes(router)
	routes.SetupUnProtectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start Server", err)
	}
}
