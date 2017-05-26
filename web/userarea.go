package web

import (
	"net/http"

	"log"

	"encoding/json"

	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/spankie/web-chat/config"
	"github.com/spankie/web-chat/messages"
	"github.com/spankie/web-chat/models"
)

// UserArea handles request for the userarea.
func UserArea(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	params := ctx.Value("params").(httprouter.Params)
	user := params.ByName("user")

	claims := ctx.Value("Claims")
	if claims == nil {
		// log attempt to access unauthorized page...
		log.Println("No claims. redirecting to /")
		http.Redirect(w, r, "/", 302)
		return
	}

	claimsName := claims.(jwt.MapClaims)["Name"].(string)

	if claimsName != user {
		log.Println("username and page user mismatch... redirecting to /")
		http.Redirect(w, r, "/", 302)
		return
	}
	// log.Println("User Active :", user)
	// if the token is invalid, ok is not true or username does not match the page visited, redirect to login
	// log.Println("the token is invalid, ok is not true or username does not match the page visited")
	// http.Redirect(w, r, "/", 302)
	// return

	http.ServeFile(w, r, "web/templates/user.html")
	return
}

// SearchFriend searches for a friend with their username...
func SearchFriend(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// You cant search for friends if you are not logged in...
	ctx := r.Context()
	claims, ok := ctx.Value("Claims").(jwt.MapClaims)
	if claims == nil || !ok {
		// log attempt to access unauthorized page...
		log.Println("No claims. sending error")
		w.WriteHeader(http.StatusOK)
		messages.SendError(w, messages.NotLoggedIn)
		return
	}

	rbody := r.Body
	log.Println("the body: ", rbody)
	// get the post parameters ...
	decoder := json.NewDecoder(rbody)
	friend := struct {
		Username string `json:"username"`
	}{}
	err := decoder.Decode(&friend)
	if err != nil {
		log.Println("decode json: ", err)
		messages.SendError(w, messages.ImproperRequest)
		return
	}
	log.Println("body friend:", friend)

	claimsUser := claims["User"].(map[string]interface{})
	claimID := int(claimsUser["ID"].(float64))
	log.Println("claimID:", claimID)

	// search the DB for the username sent
	conf := config.Get()
	for _, user := range conf.DB {
		if friend.Username == user.Username && friend.Username != claims["Name"] {
			if ContainsInt(conf.Friends[claimID], user.ID) {
				continue
				// tell the user the friend has been added already...
			}

			// send the user the id and username of the friend...
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(models.Friend{
				ID:       user.ID,
				Username: user.Username,
			})
			if err != nil {
				log.Println("Json err:", err)
				return
			}
			// i think i will have to add the user to friend list of the present(active) user
			// TODO::: write trivial function to check if it contains user.ID
			// conf.Friends[claimID] = append(conf.Friends[claimID], user.ID)
			// log.Println("conf.Friends:", conf.Friends)
			// TODO:: find a way to reconsile or authenticate the two users wen sending messages
			return
		}

	}
	w.WriteHeader(http.StatusOK)
	messages.SendError(w, messages.UserNotFound)
	return
}

// AddFriend adds a friend to the users friends list...
func AddFriend(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value("Claims").(jwt.MapClaims)
	if claims == nil || !ok {
		// log attempt to access unauthorized page...
		log.Println("No claims. sending error")
		w.WriteHeader(http.StatusOK)
		messages.SendError(w, messages.NotLoggedIn)
		return
	}

	claimsUser := claims["User"].(map[string]interface{})
	claimID := int(claimsUser["ID"].(float64))
	log.Println("claimID:", claimID)

	conf := config.Get()

	params := ctx.Value("params").(httprouter.Params)
	friendID, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		log.Println("Invalid Friend ID")
		w.WriteHeader(http.StatusOK)
		messages.SendError(w, messages.ImproperRequest)
		return
	}

	if ContainsInt(conf.Friends[claimID], friendID) {
		log.Println("Already friends with this user.")
		log.Println("myfriends:", conf.Friends[claimID])
		w.WriteHeader(http.StatusOK)
		messages.SendError(w, messages.FriendExists)
		return
		// tell the user the friend has been added already...
	}

	// now add the friend to the friends list...
	conf.Friends[claimID] = append(conf.Friends[claimID], friendID)
	log.Println("conf.Friends:", conf.Friends)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(messages.UserResponse{
		Status:  "ok",
		Cookie:  "",
		Error:   "",
		Message: messages.FriendAdded,
	})
	if err != nil {
		log.Println("Json Error: ", err)
	}

}

// GetFriends gets list of all friends connected to.
func GetFriends(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	claims, ok := ctx.Value("Claims").(jwt.MapClaims)
	if claims == nil || !ok {
		// log attempt to access unauthorized page...
		log.Println("No claims. sending error")
		w.WriteHeader(http.StatusOK)
		messages.SendError(w, messages.NotLoggedIn)
		return
	}

	claimsUser := claims["User"].(map[string]interface{})
	claimID := int(claimsUser["ID"].(float64))
	log.Println("claimID:", claimID)

	conf := config.Get()
	db := conf.DB
	var friendsList []models.Friend

	for _, v := range conf.Friends[claimID] {
		user := db[v]
		afriend := models.Friend{
			ID:       user.ID,
			Username: user.Username,
		}
		friendsList = append(friendsList, afriend)
	}
	if len(friendsList) < 1 {
		log.Println("No friends")
		w.WriteHeader(http.StatusOK)
		messages.SendError(w, messages.NoFriends)
		return
	}
	err := json.NewEncoder(w).Encode(friendsList)
	if err != nil {
		log.Println("JSON ERROR:", err)
	}
}
