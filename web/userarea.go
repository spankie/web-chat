package web

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// UserArea handles request for the userarea.
func UserArea(w http.ResponseWriter, r *http.Request) {
	// check if the user is authorized...
	// check for cookie.
	cookie, err := r.Cookie("deewebchat")
	if err != nil {
		log.Println("Could not get cookie:", err)
		http.Redirect(w, r, "/", 302)
		return
	}
	log.Println("cookie: ", cookie)
	// generate a user value from the passed cookie using jwt

	params := r.Context().Value("params").(httprouter.Params)
	user := params.ByName("user")

	log.Println("User Active :", user)
	http.ServeFile(w, r, "web/templates/user.html")
}
