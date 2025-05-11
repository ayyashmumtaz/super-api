package middleware

import (
	"net/http"
	"os"
	"strings"
	redisdb "user-service/redis"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
			return
		}

		// Ambil token dari header Authorization
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Cek apakah token sudah diblacklist
		isBlacklisted, err := redisdb.IsTokenBlacklisted(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to check token blacklist"})
			return
		}
		if isBlacklisted {
			// Tambahkan log di sini untuk debugging
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has been invalidated"})
			return
		}

		// ✅ Lanjut Parse Token jika tidak diblacklist
		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		// Cek dan ambil claims dari token
		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			return
		}

		// Set username ke context untuk digunakan di handler berikutnya
		c.Set("username", claims.Subject)
		c.Next()
	}
}
