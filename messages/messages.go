package messages

import (
	"encoding/json"
	"log"
	"net/http"
)

// TODO :: Collate all the error messages
var (
	InvalidNamePass = "Invalid Username or Password."
	UsernameTaken   = "Username is not available."
	NotLoggedIn     = "You are not logged in."
	UserNotFound    = "User not found"
	ImproperRequest = "Improper Request"
	FriendExists    = "You are already friends"
	FriendAdded     = "Friend Added"
	NoFriends       = "You have no Friends"
)

// UserResponse contains data to be sent to the user.
type UserResponse struct {
	Status  string `json:"status"`
	Cookie  string `json:"cookie"`
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SendError send a specified error message to the User.
func SendError(w http.ResponseWriter, e string) {
	err := json.NewEncoder(w).Encode(UserResponse{
		Status:  "error",
		Cookie:  "",
		Error:   e,
		Message: "",
	})
	if err != nil {
		log.Println("Json Error: ", err)
	}
}
