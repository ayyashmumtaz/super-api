package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // Ambil JWT Secret dari environment variable

// GenerateJWT menghasilkan token JWT untuk user
func GenerateJWT(username string) (string, error) {
	claims := jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour).Unix(), // 1 jam
		Subject:   username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

// Validasi token JWT
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
