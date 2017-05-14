package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/spankie/web-chat/models"
)

// Login handles login requests
func Login(w http.ResponseWriter, r *http.Request) {
	// get the post parameters ...
	decoder := json.NewDecoder(r.Body)
	user := models.User{}
	err := decoder.Decode(&user)
	if err != nil {
		log.Println("decode json: ", err)
	}
	log.Println("body:", user)
	// authenticate the user here ...

	// if the user is authenticated, send him a cookie, he/she has been a good boy/girl
	tokenstring, err := GenerateJWT(user)
	if err != nil {
		log.Println("jwt err:", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
		Cookie string `json:"cookie"`
	}{
		Status: "ok",
		Cookie: tokenstring,
	})
	if err != nil {
		log.Println("Json Error: ", err)
	}
}
