package main

import (
	"fmt"
	"log"

	database "goproj2/db"

	_ "github.com/lib/pq"
)

func main() {
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal(err)
		return

	}
	defer database.CloseDb(db)

	// userID := 1
	username := "test5"
	password := "1abc"

	user, err := database.GetUserByCreds(db, username, password)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("User ID: %d, Username: %s, Password: %s\n", user.ID, user.Username, user.Password)
}
