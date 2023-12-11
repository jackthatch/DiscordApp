package main

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	//"goproj2/models"
	database "goproj2/db"
)

func main() {
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal(err)
		return

	}
	defer database.CloseDb(db)

	userID := 1
	user, err := database.GetUserByID(db, userID)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("User ID: %d, Username: %s, Password: %s\n", user.ID, user.Username, user.Password)
}
