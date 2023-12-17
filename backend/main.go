package main

import (
	"fmt"
	"log"
	"net/http"

	database "goproj2/db"
	"goproj2/websocket"

	_ "github.com/lib/pq"
)

func main() {

	http.HandleFunc("/webs", websocket.Handler)

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/submit", submitHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "static/index.html"
	fmt.Printf("Serving file: %s\n", filePath)
	http.ServeFile(w, r, filePath)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Printf("Submitted form: Username=%s, Password:%s", username, password)

	w.Write([]byte("Form submitted successfully!"))

	callDb(username, password)
}

func callDb(username string, password string) {
	db, err := database.ConnectDb()
	if err != nil {
		log.Fatal(err)
		return

	}
	defer database.CloseDb(db)

	user, err := database.UserSignup(db, username, password)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Printf("User ID: %d, Username: %s, Password: %s\n", user.ID, user.Username, user.Password)
}
