package web

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/spankie/web-chat/config"
	"github.com/spankie/web-chat/messages"
	"github.com/spankie/web-chat/models"
)

// Signup handles adding new users
func Signup(w http.ResponseWriter, r *http.Request) {
	// get the post parameters ...
	decoder := json.NewDecoder(r.Body)
	user := models.User{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Println("decode json: ", err)
		messages.SendError(w, messages.ImproperRequest)
		return
	}
	log.Println("body:", user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// get DB
	db := config.Get().DB
	dbSize := len(db)

	// check if there is a user with the same username
	for _, v := range db {
		if v.Username == user.Username {
			// there is a user already with that username...return error
			err = json.NewEncoder(w).Encode(messages.UserResponse{
				Status: "error",
				Cookie: "",
				Error:  messages.UsernameTaken,
			})
			if err != nil {
				log.Println("Json Error: ", err)
			}
			return
		}
	}
	// add new user to DB
	user.ID = dbSize + 1
	dbSize++
	db[dbSize] = user

	// if the user is authenticated, send him a cookie, he/she has been a good boy/girl
	tokenstring, err := GenerateJWT(user)
	if err != nil {
		log.Println("jwt err:", err)
	}

	err = json.NewEncoder(w).Encode(messages.UserResponse{
		Status: "ok",
		Cookie: tokenstring,
		Error:  "",
	})
	if err != nil {
		log.Println("Json Error: ", err)
	}
}

// Login handles login requests
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// check for cookie, if there is a cookie validate to log the person in.

	// get the post parameters ...
	decoder := json.NewDecoder(r.Body)
	user := models.User{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Println("decode json: ", err)
		messages.SendError(w, messages.ImproperRequest)
		return
	}
	log.Println("body:", user)

	// userresponse := models.UserResponse{}
	// authenticate the user here ...

	// get DB
	db := config.Get().DB
	// user.ID = 1
	// db[1] = user
	dbSize := len(db)
	log.Println("len(db):", dbSize)
	// check if there is a user with the same username and password...
	for _, v := range db {
		if v.Username == user.Username && v.Password == user.Password {
			// if the user is authenticated, send him a cookie, he/she has been a good boy/girl
			tokenstring, err := GenerateJWT(user)
			if err != nil {
				log.Println("jwt err:", err)
			}

			err = json.NewEncoder(w).Encode(messages.UserResponse{
				Status: "ok",
				Cookie: tokenstring,
				Error:  "",
			})
			if err != nil {
				log.Println("Json Error: ", err)
			}
			return
		}
	}
	// there is a user already with that username...return error
	err = json.NewEncoder(w).Encode(messages.UserResponse{
		Status: "error",
		Cookie: "",
		Error:  messages.InvalidNamePass,
	})
	if err != nil {
		log.Println("Json Error: ", err)
	}
	return
}

// AuthHandler authencticates each request...
func AuthHandler(next http.Handler) http.Handler {
	conf := config.Get()
	fn := func(w http.ResponseWriter, r *http.Request) {
		// check if the user is authorized...
		// check for cookie.
		cookie, err := r.Cookie("deewebchat")
		if err != nil {
			log.Println("Could not get cookie:", err)
			next.ServeHTTP(w, r)
			// http.Redirect(w, r, "/", 302)
			return
		}
		// log.Println("cookie: ", cookie.Value)

		// generate a user value from the passed cookie using jwt

		token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
			if token.Method.Alg() == jwt.SigningMethodRS256.Alg() {
				pub, err := jwt.ParseRSAPublicKeyFromPEM(conf.CertKey)
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
		// Check if there was any error wen parsing the Cookie token
		if err != nil {
			// could not parse the token, might be invalid...
			log.Println("could not parse the token, might be invalid.::", err)
			next.ServeHTTP(w, r)
			return
		}

		// If all goes well, continue to serve the page...
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			/*
					First check if this user is in the DB, incase he is accessing
				   	the page with a valid cookie but he has been removed from DB
			*/
			user := claims["User"].(map[string]interface{})
			db := conf.DB
			log.Printf("%T : %v", user, user)
			id := int(user["ID"].(float64))
			if db[id].Username == user["username"].(string) {
				// set the context for the request so the user will be available to the next handler
				ctx := context.WithValue(r.Context(), "Claims", claims)
				r = r.WithContext(ctx)
				next.ServeHTTP(w, r)
				return
			}
			// the user is not in the DB
			next.ServeHTTP(w, r)
			return
		}

		// if there was an unforseen error, just proceed...
		next.ServeHTTP(w, r)
		return
	}

	return http.HandlerFunc(fn)
}
