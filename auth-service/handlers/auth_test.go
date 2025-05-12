package handlers_test

import (
	"auth-service/database"
	"auth-service/handlers"
	"auth-service/models"
	"auth-service/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestRegisterHandler(t *testing.T) {
	// Mock database connection
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	// Set the DB in handlers package to the mocked DB
	database.DB = db

	// Prepare hashed password
	password := "test1234"
	// hashed password is not used, so we remove its declaration
	_, _ = utils.HashPassword(password)

	// Mock the behavior of the database for the query
	mock.ExpectQuery("INSERT INTO users").
		WithArgs("Test", "testuser", "test@example.com", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Set up the router
	r := setupRouter()
	r.POST("/register", handlers.Register)

	// Create a user object for the request
	user := models.User{
		Name:     "Test",
		Username: "testuser",
		Email:    "test@example.com",
		Password: password,
	}
	body, _ := json.Marshal(user)
	fmt.Println(string(body))

	// Create a new HTTP POST request
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Record the response
	resp := httptest.NewRecorder()

	// Send the request
	r.ServeHTTP(resp, req)

	// Check if the response is as expected
	assert.Equal(t, http.StatusCreated, resp.Code)
}
