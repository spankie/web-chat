package web

import (
	"log"
	"net/http"

	"fmt"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
	"github.com/spankie/web-chat/config"
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
	log.Println("cookie: ", cookie.Value)

	conf := config.Get()

	// generate a user value from the passed cookie using jwt

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() == jwt.SigningMethodRS256.Alg() {
			pub, err := jwt.ParseRSAPrivateKeyFromPEM(conf.PrivateKey)
			if err != nil {
				log.Println("Could not create pub")
				return pub, err
			}
			log.Println("created pub")
			return pub, nil
		}
		log.Println("Unexpected signing method::", token.Header["alg"])
		return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
	})

	// Do your thing here...
	if err != nil {
		// could not parse the token, might be invalid...
		log.Println("could not parse the token, might be invalid.::", err)
		http.Redirect(w, r, "/", 302)
		return
	}

	params := r.Context().Value("params").(httprouter.Params)
	user := params.ByName("user")

	// If all goes well, continue to serve the page...
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid && claims["Name"] == user {
		log.Println("User Active :", user)
		// Should probably add the user to active users...
		http.ServeFile(w, r, "web/templates/user.html")
		return
	}
	// if the token is invalid, ok is not true or username does not match the page visited, redirect to login
	log.Println("the token is invalid, ok is not true or username does not match the page visited")
	http.Redirect(w, r, "/", 302)
	return
}
