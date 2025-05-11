package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET")) // Ambil JWT Secret dari environment variable

// GenerateJWT menghasilkan token JWT untuk user
func GenerateJWT(username string) (string, error) {
	// Set klaim token
	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // expired 1 hari
	}

	// Buat token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Tanda tangani token
	return token.SignedString(jwtKey)
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
