package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

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

func GetUserByCreds(db *sql.DB, userName string, userPass string) (*models.User, error) {
	query := "SELECT * FROM users WHERE username = $1 AND password = $2"
	row := db.QueryRow(query, userName, userPass)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		} else {
			log.Fatal(err)
			return nil, fmt.Errorf("failed to fetch user: %v", err)
		}
	}

	return user, nil
}

func UserSignup(db *sql.DB, userName string, userPass string) (*models.User, error) {
	_, err := GetUserByCreds(db, userName, userPass)
	if err != nil {
		if err.Error() == "user not found" {
			addQuery := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
			row := db.QueryRow(addQuery, userName, userPass)
			user := &models.User{}

			if err := row.Scan(&user.ID); err != nil {
				log.Fatal(err)
				return nil, fmt.Errorf("failed to insert user: %v", err)
			}

			user.Username = userName
			user.Password = userPass

			return user, nil
		}

		log.Fatal(err)
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	return nil, fmt.Errorf("user already exists")
}

func FindServer(db *sql.DB, serverName string) (*models.Server, error) {

	query := "SELECT * FROM servers WHERE name = $1"
	row := db.QueryRow(query, serverName)

	var server models.Server
	err := row.Scan(&server.ID, &server.Name, &server.Members)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Server not found")
		} else {
			log.Fatal(err)
			return nil, fmt.Errorf("failed to fetch server: %v", err)
		}

	}

	return &server, nil
}

func CreateServer(db *sql.DB, serverName string, user *models.User) (*models.Server, error) {

	existingServer, err := FindServer(db, serverName)

	if err != nil {
		if strings.Contains(err.Error(), "Server not found") {

			_, err := db.Exec("INSERT INTO servers (name, members) VALUES ($1, 1)", serverName)
			if err != nil {
				log.Fatal(err)
				log.Printf("Failed to create server: %v", err)
				return nil, fmt.Errorf("failed to create server %v", err)
			}

			createdServer, err := FindServer(db, serverName)
			if err != nil {
				log.Printf("failed to fetch server after creation %v", err)
				return nil, fmt.Errorf("failed to fetch server after creation %v", err)
			}

			user.Servers = append(user.Servers, *createdServer)

			newServer := &models.Server{
				ID:      createdServer.ID,
				Name:    createdServer.Name,
				Members: createdServer.Members,
			}

			return newServer, nil
		}

		log.Printf("failed to fetch server: %v", err)
		return nil, fmt.Errorf("failed to fetch server: %v", err)
	}

	//otherwise return the existing server
	return existingServer, nil
}
