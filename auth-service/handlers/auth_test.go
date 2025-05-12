package handlers_test

import (
	"auth-service/handlers"
	"auth-service/models"
	"auth-service/utils"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/register", handlers.Register)
	r.POST("/login", handlers.Login)
	r.POST("/logout", handlers.Logout)
	return r
}

func TestRegisterHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	handlers.OverrideDB(db) // You should create a function to override the DB in handlers
	defer db.Close()

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("Test User", "testuser", "test@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	user := models.User{
		Name:     "Test User",
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(user)
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	handlers.OverrideDB(db)
	defer db.Close()

	hashed, _ := utils.HashPassword("password123")
	mock.ExpectQuery(`SELECT id, username, email, password FROM users WHERE username = \$1`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(1, "testuser", "test@example.com", hashed))

	body := `{"username": "testuser", "password": "password123"}`
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestLoginHandler_InvalidPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock, _ := sqlmock.New()
	handlers.OverrideDB(db)
	defer db.Close()

	hashed, _ := utils.HashPassword("correctpassword")
	mock.ExpectQuery(`SELECT id, username, email, password FROM users WHERE username = \$1`).
		WithArgs("testuser").
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password"}).
			AddRow(1, "testuser", "test@example.com", hashed))

	body := `{"username": "testuser", "password": "wrongpassword"}`
	r := setupRouter()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
}
