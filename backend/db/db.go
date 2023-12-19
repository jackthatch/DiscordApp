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

func CheckUsernameExists(db *sql.DB, userName string) (bool, error) { //return true if username is in database
	query := "SELECT COUNT (*) FROM users WHERE username = $1"
	row := db.QueryRow(query, userName)

	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false, fmt.Errorf("failed to check user exists: %v", err)
	}

	return count > 0, nil
}

func CheckUserPassword(db *sql.DB, userName string, userPass string) (bool, error) { //return true if password is correct
	query := "SELECT COUNT(*) FROM users WHERE username = $1 AND password = $2"
	row := db.QueryRow(query, userName, userPass)

	var count int
	err := row.Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false, fmt.Errorf("Password for user: %v, is incorrect", userName)
	}

	return count > 0, nil
}

func UserLogin(db *sql.DB, userName string, userPass string) (*models.User, error) {

	userExists, err := CheckUsernameExists(db, userName)

	if err != nil {
		return nil, fmt.Errorf("error checking if user exists: %v", err)
	}

	if userExists {
		validPassword, err := CheckUserPassword(db, userName, userPass)
		if err != nil {
			return nil, fmt.Errorf("error checking password :%v", err)
		}

		if !validPassword {
			return nil, fmt.Errorf("incorrect password")
		}

		user, err := GetUserByCreds(db, userName, userPass)
		if err != nil {
			return nil, fmt.Errorf("error retrieving user details: %v", err)
		}

		loggedUser := &models.User{}

		loggedUser.ID = user.ID
		loggedUser.Username = user.Username
		loggedUser.Password = user.Password

		return loggedUser, nil
	}
	return nil, fmt.Errorf("username not found in the datbase")
}

func UserSignup(db *sql.DB, userName string, userPass string) (*models.User, error) {
	userExists, err := CheckUsernameExists(db, userName)

	if err != nil {
		return nil, fmt.Errorf("error checking if user exists")
	}

	if !userExists {
		query := "INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id"
		row := db.QueryRow(query, userName, userPass)
		user := &models.User{}

		if err := row.Scan(&user.ID); err != nil {
			log.Fatal(err)
			return nil, fmt.Errorf("failed to insert user: %v", err)
		}

		user.Username = userName
		user.Password = userPass

		return user, nil
	}

	return nil, fmt.Errorf("username already exists in database")
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

func FindOrCreateServer(db *sql.DB, serverName string, user *models.User) (*models.Server, error) {
	existingServer, err := FindServer(db, serverName)

	if err != nil {
		if strings.Contains(err.Error(), "Server not found") {
			createdServer, err := CreateServer(db, serverName, user)
			if err != nil {
				return nil, fmt.Errorf("failed to create server %v", err)
			}
			return createdServer, nil
		}

		return nil, fmt.Errorf("failed to find server %v", err)
	}

	return existingServer, nil
}
