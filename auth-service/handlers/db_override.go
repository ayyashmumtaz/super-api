package handlers

import (
	"auth-service/database"
	"database/sql"
)

func OverrideDB(mockDB *sql.DB) {
	database.DB = mockDB
}
