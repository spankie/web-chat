package models

// User contains individuals details
type User struct {
	ID       int    `json="id"`
	Username string `json:"username"`
	Password string `json:"password"`
}