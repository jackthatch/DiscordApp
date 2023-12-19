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
	http.HandleFunc("/submit-login", submitLoginHandler)
	http.HandleFunc("/submit-signup", submitSignupHandler)

	fs := http.FileServer(http.Dir("dist"))
	http.Handle("/dist/", http.StripPrefix("/dist/", fs))

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("error starting server: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "dist/index.html"
	fmt.Printf("Serving file: %s\n", filePath)
	http.ServeFile(w, r, filePath)
}

func submitLoginHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Printf("Submitted form: Username=%s, Password:%s\n", username, password)

	db, err := database.ConnectDb()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to connect to the database: %v", err), http.StatusInternalServerError)
		return
	}
	defer database.CloseDb(db)

	user, err := database.UserLogin(db, username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to sign up user: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("User ID: %d, Username: %s, Password: %s\n", user.ID, user.Username, user.Password)

	w.Write([]byte("Form submitted successfully!"))
}

func submitSignupHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to parse form: %v", err), http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Printf("Submitted form: Username=%s, Password:%s\n", username, password)

	db, err := database.ConnectDb()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to connect to the database: %v", err), http.StatusInternalServerError)
		return
	}
	defer database.CloseDb(db)

	user, err := database.UserSignup(db, username, password)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to sign up user: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Printf("User ID: %d, Username: %s, Password: %s\n", user.ID, user.Username, user.Password)

	w.Write([]byte("Form submitted successfully!"))
}
