package main

import (
	"auth-service/database"
	"auth-service/handlers"
	redisdb "auth-service/redis"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using container env")
	}

	database.InitDB()
	redisdb.InitRedis()

	r := gin.Default()
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)
	r.GET("/auth/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	r.POST("/auth/logout", handlers.Logout)
	r.Run(":8080")
}
