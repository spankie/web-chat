package web

import (
	"net/http"

	"log"

	"encoding/json"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/spankie/web-chat/config"
	"github.com/spankie/web-chat/messages"
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

	// search the DB for the username sent
	db := config.Get().DB
	for _, user := range db {
		if friend.Username == user.Username && friend.Username != claims["Name"] {
			// send the user the id and username of the friend...
			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(struct {
				ID       int    `json:"id"`
				Username string `json:"username"`
			}{
				ID:       user.ID,
				Username: user.Username,
			})
			if err != nil {
				log.Println("Json err:", err)
				return
			}
			// i think i will have to add the user to friend list of the present(active) user
			// TODO:: find a way to reconsile or authenticate the two users wen sending messages
			return
		}

	}
	w.WriteHeader(http.StatusOK)
	messages.SendError(w, messages.UserNotFound)
	return
}
