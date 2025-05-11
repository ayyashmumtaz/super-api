package handlers

import (
	"auth-service/database"
	"auth-service/models"
	"auth-service/utils"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	hashedPassword, err := utils.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	err = database.DB.QueryRow(
		"INSERT INTO users (name, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id",
		newUser.Name, newUser.Username, newUser.Email, hashedPassword,
	).Scan(&newUser.ID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newUser.Password = "" // Don't return the password in the response
	c.JSON(http.StatusCreated, newUser)
}

func Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var dbUser models.User
	err := database.DB.QueryRow(
		"SELECT id, username, email, password FROM users WHERE username = $1",
		user.Username,
	).Scan(&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	if !utils.CheckPasswordHash(user.Password, dbUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful!",
		"token":   token,
	})

}
