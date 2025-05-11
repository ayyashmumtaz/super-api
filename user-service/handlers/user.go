package handlers

import (
	"database/sql"
	"net/http"
	"user-service/database"
	"user-service/models"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func GetUserData(c *gin.Context) {
	// Ambil username dari context yang diset di middleware
	usernameRaw, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	username := usernameRaw.(string)

	// Ambil data user dari DB berdasarkan username
	var u models.User
	err := database.DB.QueryRow("SELECT id, name, username, email FROM users WHERE username = $1", username).
		Scan(&u.ID, &u.Name, &u.Username, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return data user
	c.JSON(http.StatusOK, u)
}

func GetUserByID(c *gin.Context) {
	// Ambil username dari context yang diset di middleware
	usernameRaw, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	username := usernameRaw.(string)

	// Ambil id dari parameter
	id := c.Param("id")
	var u models.User
	err := database.DB.QueryRow("SELECT id, username, name, email FROM users WHERE id = $1 AND username = $2", id, username).
		Scan(&u.ID, &u.Username, &u.Name, &u.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// Return data user
	c.JSON(http.StatusOK, u)
}
