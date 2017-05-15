package messages

// TODO :: Collate all the error messages
var (
	InvalidNamePass = "Invalid Username or Password"
	UsernameTaken   = "Username is not available"
)

// UserResponse contains data to be sent to the user.
type UserResponse struct {
	Status string `json:"status"`
	Cookie string `json:"cookie"`
	Error  string `json:"error"`
}
