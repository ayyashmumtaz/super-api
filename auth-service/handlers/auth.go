package handlers

import (
	"auth-service/database"
	"auth-service/models"
	redisdb "auth-service/redis"
	"auth-service/utils"
	"database/sql"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var DB *sql.DB

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

func Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No token provided"})
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse token
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	// Hitung waktu kadaluarsa token
	expiration := time.Until(claims.ExpiresAt.Time)
	if expiration <= 0 {
		expiration = time.Minute * 5 // fallback kalau exp di masa lalu
	}

	// Simpan ke Redis sebagai blacklist
	if err := redisdb.BlacklistToken(tokenStr, expiration); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to blacklist token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
