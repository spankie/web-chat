package web

import (
	"net/http"

	"log"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
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
