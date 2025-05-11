package main

import (
	"auth-service/database"
	"auth-service/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using container env")
	}

	database.InitDB()

	r := gin.Default()
	r.POST("/auth/register", handlers.Register)
	r.POST("/auth/login", handlers.Login)
	r.Run(":8080")
}
