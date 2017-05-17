package web

import (
	"encoding/json"
	"log"
	"net/http"

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
	// check for cookie, if there is a cookie validate to log the person in.

	// get the post parameters ...
	decoder := json.NewDecoder(r.Body)
	user := models.User{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Println("decode json: ", err)
	}
	log.Println("body:", user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

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
