package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("No .env file found, using env from container")
	}
}

func getDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)
}

func main() {
	fmt.Println("DSN:", getDSN())

	db, err := sql.Open("postgres", getDSN())
	if err != nil {
		log.Fatal("Failed to connect DB:", err)
	}
	defer db.Close()

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/ping-db", func(c *gin.Context) {
		if err := db.Ping(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"db": "not connected"})
		} else {
			c.JSON(http.StatusOK, gin.H{"db": "connected"})
		}
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		var u User
		err := db.QueryRow("SELECT id,username, name, email FROM users WHERE id = $1", id).Scan(&u.ID, &u.Username, &u.Name, &u.Email)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}
		c.JSON(http.StatusOK, u)
	})

	r.GET("/users", func(c *gin.Context) {
		rows, err := db.Query("SELECT id, name, username, email, password FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()

		var users []User
		for rows.Next() {
			var u User
			if err := rows.Scan(&u.ID, &u.Name, &u.Username, &u.Email, &u.Password); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			users = append(users, u)
		}
		c.JSON(http.StatusOK, users)
	})

	r.POST("/users", func(c *gin.Context) {
		var u User
		if err := c.ShouldBindJSON(&u); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", u.Name, u.Email).Scan(&u.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, u)
	})

	r.Run(":8081")
}
