package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Rudraksh121a/Movie-Streaming-web/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("Warning: unable to find .env file")
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")

	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		log.Println("Allowed Origin: http://localhost:5173")
	}

	config := cors.Config{}
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	//config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour

	router.Use(cors.New(config))
	router.Use(gin.Logger())

	routes.SetupProtectedRoutes(router)
	routes.SetupUnProtectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start Server", err)
	}
}
