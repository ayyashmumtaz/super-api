package main

import (
	"log"
	"net/http"
	"user-service/database"
	"user-service/handlers"
	"user-service/middleware"
	redisdb "user-service/redis"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using env from container")
	}
}

func main() {
	database.InitDB()
	redisdb.InitRedis()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/ping-db", func(c *gin.Context) {
		if err := database.DB.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"db": "not connected"})
		} else {
			c.JSON(http.StatusOK, gin.H{"db": "connected"})
		}
	})

	protected := r.Group("/api/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/users", handlers.GetUserData)
		protected.GET("/users/:id", handlers.GetUserByID)

	}

	r.Run(":8081")
}
