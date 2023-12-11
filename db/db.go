package database

import (
	"database/sql"
	"fmt"
	"log"

	"goproj2/models"
	"os"

	"github.com/joho/godotenv"
)

func ConnectDb() (*sql.DB, error) {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("failed to load dot env %s", err)
		return nil, fmt.Errorf("failed to load dot env: %v", err)
	}

	connectionStr := os.Getenv("DB_CONNECTION_STRING")

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	fmt.Println("Connected to the database")

	return db, nil
}

func CloseDb(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Closed the database connection")
}

func GetUserByID(db *sql.DB, userID int) (*models.User, error) {
	query := "SELECT * FROM users WHERE id = $1"
	row := db.QueryRow(query, userID)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		log.Fatal(err)
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	return user, nil
}
