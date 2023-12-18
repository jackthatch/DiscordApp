package models

type User struct {
	ID       int
	Username string
	Password string
	Servers  []Server
}

type Server struct {
	ID      int
	Name    string
	Members int //number of members
}
